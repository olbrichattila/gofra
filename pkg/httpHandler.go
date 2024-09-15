package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/olbrichattila/gofra/pkg/app/request"
	"github.com/olbrichattila/gofra/pkg/app/router"
	"github.com/olbrichattila/gofra/pkg/app/session"
	"github.com/olbrichattila/gofra/pkg/app/validator"
	internalconfig "github.com/olbrichattila/gofra/pkg/internal-config"
)

type hTTPHandler struct {
	app             *App
	routes          []router.ControllerAction
	customValidator validator.Validator
	session         session.Sessioner
	requester       request.Requester
}

func (h *hTTPHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.initRoutes()
	h.initValidator()
	h.initSession(w, r)

	h.app.di.Set("http.ResponseWriter", w)
	h.runMiddlewares(h.app.conf.Middlewares())

	if !h.renderActionIfRouteFind(w, r) {
		http.NotFound(w, r)
	}
}

func (h *hTTPHandler) renderActionIfRouteFind(w http.ResponseWriter, r *http.Request) bool {
	for _, action := range h.routes {
		match, routePars := h.app.router.Match(action.Path, r.RequestURI)

		if match {
			if action.RequestType != r.Method {
				continue
			}

			if h.requester != nil {
				h.requester.SetRouteParameters(routePars)
			}

			if h.runMiddlewares(action.Middlewares) {
				return true
			}

			if h.runValidator(w, r, action.ValidationRules) {
				// redirected, stop execution
				return true
			}

			// This is the main controller call
			result, err := h.app.di.Call(action.Fn)
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

func (h *hTTPHandler) runMiddlewares(middlewares []interface{}) bool {
	for _, middleware := range middlewares {

		res, err := h.app.di.Call(middleware)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
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
		errorInterface := reflect.TypeOf((*error)(nil)).Elem()
		if result[1].Type().Implements(errorInterface) {
			// Use type assertion to get the error
			if err, ok := result[1].Interface().(error); ok {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return true
			}
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

func (h *hTTPHandler) getRequestFromDi() request.Requester {
	dep, err := h.app.di.GetDependency("olbrichattila.gofra.pkg.app.request.Requester")
	if err == nil {
		if req, ok := dep.(request.Requester); ok {
			return req
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
	}

	return nil
}

func (h *hTTPHandler) getSessionerFromDi() session.Sessioner {
	dep, err := h.app.di.GetDependency("olbrichattila.gofra.pkg.app.session.Sessioner")
	if err == nil {
		if req, ok := dep.(session.Sessioner); ok {
			return req
		}
	}

	return nil
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
