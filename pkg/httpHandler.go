package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"slices"

	"github.com/olbrichattila/gofra/pkg/app/gofraerror"
	"github.com/olbrichattila/gofra/pkg/app/logger"
	"github.com/olbrichattila/gofra/pkg/app/request"
	"github.com/olbrichattila/gofra/pkg/app/router"
	"github.com/olbrichattila/gofra/pkg/app/session"
	"github.com/olbrichattila/gofra/pkg/app/validator"
	"github.com/olbrichattila/gofra/pkg/app/view"
	internalconfig "github.com/olbrichattila/gofra/pkg/internal-config"
)

type hTTPHandler struct {
	app             *App
	routes          []router.ControllerAction
	customValidator validator.Validator
	session         session.Sessioner
	requester       request.Requester
	logger          logger.Logger
}

func (h *hTTPHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			stackTrace := string(buf[:n])
			h.logCritical(fmt.Sprintf("error: %v\nStack Trace:\n%s\n", err, stackTrace))
			fmt.Printf("error: %v\nStack Trace:\n%s\n", err, stackTrace)

			h.renderGofraError(w, err.(error))
		}
	}()

	h.initRoutes()
	h.initValidator()
	h.initSession(w, r)
	h.initLogger()

	h.app.di.Set("http.ResponseWriter", w)
	if h.runMiddlewares(h.app.conf.Middlewares()) {
		// If middleware want to terminate process
		return
	}

	if !h.renderActionIfRouteFind(w, r) {
		http.NotFound(w, r)
	}
}

func (h *hTTPHandler) renderActionIfRouteFind(w http.ResponseWriter, r *http.Request) bool {
	for _, action := range h.routes {
		match, routePars := h.app.router.Match(action.Path, r.RequestURI)

		if match {
			if action.IsStatic {
				// Serve static file
				baseStaticPath := action.StaticPath
				if !strings.HasPrefix(baseStaticPath, "/") {
					baseStaticPath = "/" + baseStaticPath
				}

				if fileToServe, ok := routePars["*"]; ok {
					http.ServeFile(w, r, string(http.Dir("."+baseStaticPath+fileToServe)))
					return true
				}

				http.ServeFile(w, r, string(http.Dir("."+baseStaticPath)))
				return true
			}

			if !h.isInRequestTypes(action.RequestType, r.Method) {
				continue
			}

			if h.requester != nil {
				h.requester.SetRouteParameters(routePars)
			}

			if h.runMiddlewares(action.Middlewares) {
				// If middleware want to terminate process
				return true
			}

			if h.runValidator(w, r, action.ValidationRules) {
				// redirected, stop execution
				return true
			}

			h.loadRouteViewAutoLoads(action.ViewAutoLoads)

			// Crete controller from struct if provided
			if action.Controller != nil {
				result, err := h.resolveControllerActionFromStruct(action, r)
				if err != nil {
					h.renderGofraError(w, err)
					return true
				}

				return h.renderControllerResult(result, w)
			}

			bodyAsStruct, err := h.mapRouteParamsIfResolvable(action.Path, action.Fn, r)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return true
			}

			// This is the main controller call
			result, err := h.app.di.Call(action.Fn, bodyAsStruct...)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return true
			}

			return h.renderControllerResult(result, w)
		}
	}

	return false
}

func (h *hTTPHandler) isInRequestTypes(requestTypes []string, requestType string) bool {
	return slices.Contains(requestTypes, requestType)
}

func (h *hTTPHandler) resolveControllerActionFromStruct(action router.ControllerAction, r *http.Request) ([]reflect.Value, error) {
	// Crete controller from struct if provided
	controllerStruct := action.Controller()
	if controllerStruct == nil {
		return nil, fmt.Errorf("controller is nil")
	}

	// TODO check if it is a struct
	val := reflect.ValueOf(controllerStruct)
	method := val.MethodByName(action.ActionName)
	if !method.IsValid() {
		return nil, fmt.Errorf("action does not exists: %s", action.ActionName)
	}

	// Call before action with DI
	beforeMethod := val.MethodByName("Before")
	if beforeMethod.IsValid() {

		bodyAsStruct, err := h.mapRouteParamsIfResolvable(action.Path, beforeMethod.Interface(), r)
		if err != nil {
			return nil, err
		}

		beforeResult, err := h.app.di.Call(beforeMethod.Interface(), bodyAsStruct...)
		if err != nil {
			return nil, err
		}

		if len(beforeResult) == 1 && beforeResult[0].Interface() != nil && beforeResult[0].Interface().(error) != nil {
			return nil, beforeResult[0].Interface().(error)
		}

	}

	bodyAsStruct, err := h.mapRouteParamsIfResolvable(action.Path, method.Interface(), r)
	if err != nil {
		return nil, err
	}

	result, err := h.app.di.Call(method.Interface(), bodyAsStruct...)
	if err != nil {
		return nil, err
	}

	// Call after method if exists
	afterMethod := val.MethodByName("After")
	if afterMethod.IsValid() {
		bodyAsStruct, err := h.mapRouteParamsIfResolvable(action.Path, afterMethod.Interface(), r)
		if err != nil {
			return nil, err
		}

		afterResult, err := h.app.di.Call(afterMethod.Interface(), bodyAsStruct...)
		if err != nil {
			return nil, err
		}

		if len(afterResult) == 1 && afterResult[0].Interface() != nil && afterResult[0].Interface().(error) != nil {
			return nil, afterResult[0].Interface().(error)
		}
	}

	return result, nil
}

// this function tries to assign parameters to the controller from the route by name or render body to a struct, or form to struct ot map[string]any
// if the parameter type hint is struct. If int and the provided value is not int it returns false
// Note these parameters must be at the beginning of the parameter list
func (h *hTTPHandler) mapRouteParamsIfResolvable(route string, fn any, r *http.Request) ([]any, error) {
	bodyAsStruct := []any{}
	var bodyBytes []byte
	fnType := reflect.TypeOf(fn)

	parIndex := 0
	for i := range fnType.NumIn() {
		paramType := fnType.In(i)
		paramValueAsString, err := h.requester.URLParByIndex(route, parIndex)
		if err != nil {
			return bodyAsStruct, err
		}

		if paramType.Kind() == reflect.String {
			parIndex++
			if h.requester != nil {
				bodyAsStruct = append(bodyAsStruct, paramValueAsString)
			} else {
				return bodyAsStruct, fmt.Errorf("requested not set when trying to resolve string parameter")
			}
		} else if paramType.Kind() == reflect.Int {
			parIndex++
			if h.requester != nil {
				if val, err := strconv.Atoi(paramValueAsString); err == nil {
					bodyAsStruct = append(bodyAsStruct, val)
				} else {
					return bodyAsStruct, fmt.Errorf("expecting integer parameter, got string `%s`", paramValueAsString)
				}
			} else {
				return bodyAsStruct, fmt.Errorf("requested not set when trying to resolve int parameter")
			}
		} else if paramType.Kind() == reflect.Struct || paramType.Kind() == reflect.Map {
			if h.isFormRequest(r) {
				if paramType.Kind() == reflect.Struct {
					err = h.marshalFormDataToBodyAsStruct(paramType, r, &bodyAsStruct)
				} else {
					err = h.marshalFormDataToBodyAsMap(paramType, r, &bodyAsStruct)
				}
			} else {
				if bodyBytes == nil {
					bodyBytes, err = io.ReadAll(r.Body)
					if err != nil {
						return bodyAsStruct, fmt.Errorf("failed to read body: %w", err)
					}
				}
				err = h.marshalToBodyAsStruct(bodyBytes, paramType, r, &bodyAsStruct)
			}
			if err != nil {
				return bodyAsStruct, fmt.Errorf("cannot parse request body %w", err)
			}
		} else {
			return bodyAsStruct, nil
		}
	}

	return bodyAsStruct, nil
}

func (h *hTTPHandler) initRoutes() {
	h.routes = h.app.conf.Routes()
}

func (h *hTTPHandler) initValidator() {
	h.customValidator = h.getValidatorFromDi()
	if h.customValidator != nil {
		h.customValidator.SetRules(h.app.validationRuleFuncs)
		h.customValidator.SetRules(internalconfig.ValidatorRules)
	}
}

func (h *hTTPHandler) initSession(w http.ResponseWriter, r *http.Request) {
	h.session = h.getSessionerFromDi()
	if h.session != nil {
		h.session.Init(w, r)
	}
	h.requester = h.getRequestFromDi()
	if h.requester != nil {
		h.requester.SetRequest(r)
	}
}

func (h *hTTPHandler) initLogger() {
	h.logger = h.getLoggerFromDi()
}

func (h *hTTPHandler) runMiddlewares(middlewares []interface{}) bool {
	for _, middleware := range middlewares {
		res, err := h.app.di.Call(middleware)
		if err != nil {
			fmt.Println(err.Error())
			// Log instead
			h.logger.Error(err.Error())
		}

		if len(res) > 0 && res[0].Kind() == reflect.Bool {
			if !res[0].Bool() {
				return true
			}
		}
	}

	return false
}

func (h *hTTPHandler) runValidator(w http.ResponseWriter, r *http.Request, validationRules string) bool {
	// Route validator logic
	if validationRules != "" {
		genericErrors := make(validator.ValidationErrors)
		funcErrors := make(validator.ValidationErrors)
		isValid := true

		if rule, ok := h.app.validationRules[validationRules]; ok {
			if h.customValidator != nil {
				allRequests := h.requester.AllFlat()
				if rule.Rules != nil {

					ok, errors, _ := h.customValidator.Validate(allRequests, rule.Rules)
					if !ok {
						genericErrors = errors
						isValid = false
					}
				}

				if rule.CustomRule != nil {
					if customFuncErrors, ok := rule.CustomRule(allRequests); !ok {
						funcErrors = customFuncErrors
						isValid = false
					}
				}

				if !isValid {
					if h.session != nil {
						combinedErrors := h.mergeValidationErrors(genericErrors, funcErrors)
						jSONError, err := json.Marshal(combinedErrors)
						if err == nil {
							h.session.Set("lastValidationError", string(jSONError))
						}

						requestJSON, err := json.Marshal(allRequests)
						if err == nil {
							h.session.Set("lastRequest", string(requestJSON))
						}
					}

					if rule.Redirect != "" {
						http.Redirect(w, r, rule.Redirect, http.StatusSeeOther)
						return true
					}
				}
			}
		}
	}

	return false
}

func (h *hTTPHandler) renderControllerResult(result []reflect.Value, w http.ResponseWriter) bool {
	if len(result) == 0 {
		// nothing to render
		return true
	}

	// If second parameter is error, and not nill return error
	if len(result) == 2 {
		if h.renderErrorIfNecessary(w, result[1]) {
			return true
		}
	}

	// If first parameter is string, render string
	if result[0].Kind() == reflect.String {
		w.Write([]byte(result[0].String()))
		return true
	}

	// If first parameter is a struct or map, render json
	if result[0].Kind() == reflect.Struct || result[0].Kind() == reflect.Map {
		h.renderJson(result[0].Interface(), w)
		return true
	}

	// Cannot be rendered
	return false
}

func (h *hTTPHandler) renderJson(jsonData interface{}, w http.ResponseWriter) {
	jsonRes, err := json.Marshal(jsonData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRes)
}

// TODO core repetition here, fix it
func (h *hTTPHandler) getRequestFromDi() request.Requester {
	dep, err := h.app.di.GetDependency("olbrichattila.gofra.pkg.app.request.Requester")
	if err == nil {
		if req, ok := dep.(request.Requester); ok {
			return req
		}

		if dep, ok := dep.(func() any); ok {
			resolvedDep := dep()
			if req, ok := resolvedDep.(request.Requester); ok {
				return req
			}
		}
	}

	return nil
}

func (h *hTTPHandler) getValidatorFromDi() validator.Validator {
	dep, err := h.app.di.GetDependency("olbrichattila.gofra.pkg.app.validator.Validator")
	if err == nil {
		if req, ok := dep.(validator.Validator); ok {
			return req
		}

		if dep, ok := dep.(func() any); ok {
			resolvedDep := dep()
			if req, ok := resolvedDep.(validator.Validator); ok {
				return req
			}
		}
	}

	return nil
}

func (h *hTTPHandler) getSessionerFromDi() session.Sessioner {
	dep, err := h.app.di.GetDependency("olbrichattila.gofra.pkg.app.session.Sessioner")
	if err == nil {
		if req, ok := dep.(session.Sessioner); ok {
			return req
		}

		if dep, ok := dep.(func() any); ok {
			resolvedDep := dep()
			if req, ok := resolvedDep.(session.Sessioner); ok {
				return req
			}
		}
	}

	return nil
}

func (h *hTTPHandler) getViewFromDi() view.Viewer {
	dep, err := h.app.di.GetDependency("olbrichattila.gofra.pkg.app.view.Viewer")
	if err == nil {
		if req, ok := dep.(view.Viewer); ok {
			return req
		}

		if dep, ok := dep.(func() any); ok {
			resolvedDep := dep()
			if req, ok := resolvedDep.(view.Viewer); ok {
				return req
			}
		}
	}

	return nil
}

func (h *hTTPHandler) getLoggerFromDi() logger.Logger {
	dep, err := h.app.di.GetDependency("olbrichattila.gofra.pkg.app.logger.Logger")
	if err == nil {
		if req, ok := dep.(logger.Logger); ok {
			return req
		}

		if dep, ok := dep.(func() any); ok {
			resolvedDep := dep()
			if req, ok := resolvedDep.(logger.Logger); ok {
				return req
			}
		}
	}

	return nil
}

func (h *hTTPHandler) logWarning(message string) {
	if h.logger != nil {
		h.logger.Warning(message)
	}
}

func (h *hTTPHandler) logInfo(message string) {
	if h.logger != nil {
		h.logger.Info(message)
	}
}

func (h *hTTPHandler) logCritical(message string) {
	if h.logger != nil {
		h.logger.Critical(message)
	}
}

func (h *hTTPHandler) logError(message string) {
	if h.logger != nil {
		h.logger.Error(message)
	}
}

func (h *hTTPHandler) mergeValidationErrors(errorSet1, errorSet2 validator.ValidationErrors) validator.ValidationErrors {
	result := make(validator.ValidationErrors)
	for key, value := range errorSet1 {
		if value == nil {
			result[key] = make([]string, 0)
			continue
		}

		result[key] = value
	}

	for key, value := range errorSet2 {
		subset, ok := result[key]
		if ok && value != nil {
			result[key] = append(subset, value...)
			continue
		}

		if value == nil {
			result[key] = make([]string, 0)
			continue
		}

		result[key] = value

	}
	return result
}

func (h *hTTPHandler) loadRouteViewAutoLoads(loads []string) {
	view := h.getViewFromDi()
	view.LoadTemplateParts(loads)
}

func (h *hTTPHandler) renderErrorIfNecessary(w http.ResponseWriter, reflectMethod reflect.Value) bool {
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if reflectMethod.Type().Implements(errorInterface) {
		// Use type assertion to get the errorf
		if err, ok := reflectMethod.Interface().(error); ok {
			h.renderGofraError(w, err)
			return true
		}
	}

	return false
}

func (h *hTTPHandler) renderGofraError(w http.ResponseWriter, err error) {
	// Handle specific error type use case
	if gofraError, ok := err.(*gofraerror.Error); ok {
		w.WriteHeader(gofraError.ResponseStatus)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(err.Error()))
}

func (h *hTTPHandler) isFormRequest(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "application/x-www-form-urlencoded")
}

func (h *hTTPHandler) marshalToBodyAsStruct(bodyBytes []byte, paramType reflect.Type, r *http.Request, bodyAsStruct *[]any) error {
	paramPtr := reflect.New(paramType)

	if err := json.Unmarshal(bodyBytes, paramPtr.Interface()); err != nil {
		return fmt.Errorf("decode into first struct failed: %w", err)
	}

	*bodyAsStruct = append(*bodyAsStruct, paramPtr.Elem().Interface())

	return nil
}

func (h *hTTPHandler) marshalFormDataToBodyAsStruct(paramType reflect.Type, r *http.Request, bodyAsStruct *[]any) error {
	paramPtr := reflect.New(paramType)
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("cannot parse Form")
	}

	typ := paramPtr.Elem().Type()
	val := paramPtr.Elem()
	for i := 0; i < typ.NumField(); i++ {
		fieldName := typ.Field(i)
		field := val.Field(i)
		if fieldName.PkgPath != "" {
			continue
		}

		key := fieldName.Tag.Get("json")
		if key == "" {
			key = fieldName.Name
		}

		formVal, ok := r.PostForm[key]
		if !ok || len(formVal) == 0 {
			continue
		}

		raw := strings.TrimSpace(formVal[0])

		// Resolve pointers
		fieldType := field.Type()
		isPtr := fieldType.Kind() == reflect.Ptr
		if isPtr {
			fieldType = fieldType.Elem()
		}

		var parsed reflect.Value
		switch fieldType.Kind() {
		case reflect.String:
			parsed = reflect.ValueOf(raw)
		case reflect.Int:
			if i, err := strconv.Atoi(raw); err == nil {
				parsed = reflect.ValueOf(i)
			} else {
				continue
			}
		case reflect.Int64:
			if i64, err := strconv.ParseInt(raw, 10, 64); err == nil {
				parsed = reflect.ValueOf(i64)
			} else {
				continue
			}
		case reflect.Bool:
			if b, err := strconv.ParseBool(raw); err == nil {
				parsed = reflect.ValueOf(b)
			} else {
				continue
			}
		default:
			continue // unsupported type
		}

		if isPtr {
			field.Set(reflect.New(fieldType))
			field.Elem().Set(parsed.Convert(fieldType))
		} else {
			field.Set(parsed.Convert(fieldType))
		}
	}

	*bodyAsStruct = append(*bodyAsStruct, paramPtr.Elem().Interface())

	return nil
}

func (h *hTTPHandler) marshalFormDataToBodyAsMap(paramType reflect.Type, r *http.Request, bodyAsStruct *[]any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	result := make(map[string]any)
	for key, values := range r.PostForm {
		if len(values) > 0 {
			result[key] = values[0] // only the first value
		}
	}

	*bodyAsStruct = append(*bodyAsStruct, result)

	return nil
}
