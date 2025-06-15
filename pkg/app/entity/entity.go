package entity

// TODO add Count and Exist (width where)
// Auto detect SQL DB flavour and decide to use ? or $1 for binding
// Return or Alter entity to save the ID after Save (Last Insert ID)
// Can I combine with SQL builder as well?

import (
	"errors"
	"reflect"
	"strings"

	"github.com/olbrichattila/gofra/pkg/app/db"
)

const (
	tagJSONFieldName = "json"
	tagFieldName     = "fieldName"
	tagTableName     = "tableName"
)

type fieldDef struct {
	name  string
	value any
}

func Nullable[T any](v T) *T {
	return &v
}

func ById[T any](db db.DBer, id int64) (T, error) {
	var entity T
	_, tableName, _, err := parseEntity(entity)
	if err != nil {
		return entity, err
	}

	row, err := db.QueryOne("SELECT * FROM `"+tableName+"` WHERE `id` = ?", id)
	if err != nil {
		return entity, err
	}

	if err := mapToStruct(row, &entity); err != nil {
		return entity, err
	}

	return entity, nil
}

// All is an alias of ByWhere without parameters
func All[T any](db db.DBer) ([]T, error) {
	return ByWhere[T](db, "")
}

func ByWhere[T any](db db.DBer, where string, params ...any) ([]T, error) {
	var entity T
	_, tableName, _, err := parseEntity(entity)
	if err != nil {
		return nil, err
	}

	sql := "SELECT * FROM `" + tableName + "` " + where
	rows := db.QueryAll(sql, params...)
	if err := db.GetLastError(); err != nil {
		return nil, err
	}

	var result []T
	for row := range rows {
		var item T
		if err := mapToStruct(row, &item); err == nil {
			result = append(result, item)
		}
	}
	return result, nil
}

func Count[T any](db db.DBer, where string, params ...any) (int64, error) {
	var entity T
	_, tableName, _, err := parseEntity(entity)
	if err != nil {
		return 0, err
	}

	sql := "SELECT count(*) as cnt FROM `" + tableName + "` " + where
	row, err := db.QueryOne(sql, params...)
	if err != nil {
		return 0, err
	}

	return row["cnt"].(int64), nil
}

func Exists[T any](db db.DBer, where string, params ...any) (bool, error) {
	var entity T
	_, tableName, _, err := parseEntity(entity)
	if err != nil {
		return false, err
	}

	sql := "SELECT * FROM `" + tableName + "` " + where + " LIMIT 1"
	rows := db.QueryAll(sql, params...)
	if err := db.GetLastError(); err != nil {
		return false, err
	}

	result := false
	for range rows {
		result = true
	}
	return result, nil
}

func Delete(db db.DBer, entity any) error {
	_, tableName, id, err := parseEntity(entity)
	if err != nil {
		return err
	}

	_, err = db.Execute("DELETE FROM `"+tableName+"` WHERE `id` = ?", id)
	return err
}

func DeleteById[T any](db db.DBer, id int64) error {
	var entity T
	_, tableName, _, err := parseEntity(entity)
	if err != nil {
		return err
	}

	_, err = db.Execute("DELETE FROM `"+tableName+"` WHERE `id` = ?", id)
	return err
}

func DeleteWhere(db db.DBer, entity any, where string, pars ...any) error {
	_, tableName, _, err := parseEntity(entity)
	if err != nil {
		return err
	}

	_, err = db.Execute("DELETE FROM `"+tableName+"` "+where, pars...)
	return err
}

func Save(db db.DBer, entity any) error {
	fields, tableName, id, err := parseEntity(entity)
	if err != nil {
		return err
	}

	var (
		sql  string
		args []any
	)
	if id == 0 {
		sql, args, err = toInsertSQL(fields, tableName)
	} else {
		sql, args, err = toUpdateSQL(id, fields, tableName)
	}

	if err != nil {
		return err
	}

	insertId, err := db.Execute(sql, args...)
	if id == 0 && err == nil {
		return updateIdIfPossible(insertId, entity)
	}

	return err
}

func parseEntity(entity any) ([]fieldDef, string, int64, error) {
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(entity)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, "", 0, errors.New("parseEntity: not a struct or pointer to struct")
	}

	tableName := strings.ToLower(t.Name())
	fields := make([]fieldDef, 0, t.NumField())
	var id int64

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue // skip unexported
		}

		if tn := field.Tag.Get(tagTableName); tn != "" {
			tableName = tn
			continue
		}

		name := field.Name
		tagName := field.Tag.Get(tagJSONFieldName)
		if tagName != "" {
			name = tagName
		}

		tagName = field.Tag.Get(tagFieldName)
		if tagName != "" {
			name = tagName
		}

		val := v.Field(i)
		if name == "id" && val.Kind() == reflect.Int64 {
			id = val.Int()
		}

		if val.Kind() == reflect.Ptr && val.IsNil() {
			fields = append(fields, fieldDef{name: name})
			continue
		}

		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		fields = append(fields, fieldDef{name: name, value: val.Interface()})
	}

	return fields, tableName, id, nil
}

func updateIdIfPossible(newId int64, entity any) error {
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(entity)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return errors.New("parseEntity: not a struct or pointer to struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue // skip unexported
		}

		name := field.Name
		tagName := field.Tag.Get(tagJSONFieldName)
		if tagName != "" {
			name = tagName
		}

		tagName = field.Tag.Get(tagFieldName)
		if tagName != "" {
			name = tagName
		}

		val := v.Field(i)
		if name == "id" && val.Kind() == reflect.Int64 {
			val.SetInt(newId)
		}
	}

	return nil
}

func toInsertSQL(fields []fieldDef, table string) (string, []any, error) {
	var cols []string
	var placeholders []string
	var args []any

	for _, f := range fields {
		if f.name == "id" {
			continue
		}
		cols = append(cols, "`"+f.name+"`")
		placeholders = append(placeholders, "?")
		args = append(args, f.value)
	}

	sql := "INSERT INTO `" + table + "` (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(placeholders, ",") + ")"
	return sql, args, nil
}

func toUpdateSQL(id int64, fields []fieldDef, table string) (string, []any, error) {
	var sets []string
	var args []any

	for _, f := range fields {
		if f.name == "id" {
			continue
		}
		sets = append(sets, "`"+f.name+"`=?")
		args = append(args, f.value)
	}

	sql := "UPDATE `" + table + "` SET " + strings.Join(sets, ",") + " WHERE `id` = ?"
	args = append(args, id)
	return sql, args, nil
}

func mapToStruct(data map[string]any, dest any) error {
	ptrVal := reflect.ValueOf(dest)
	if ptrVal.Kind() != reflect.Ptr || ptrVal.Elem().Kind() != reflect.Struct {
		return errors.New("mapToStruct: destination must be pointer to struct")
	}

	val := ptrVal.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" || field.Tag.Get(tagTableName) != "" {
			continue
		}

		key := field.Name
		tagName := field.Tag.Get(tagJSONFieldName)
		if tagName != "" {
			key = tagName
		}

		tagName = field.Tag.Get(tagFieldName)
		if tagName != "" {
			key = tagName
		}

		mapVal, ok := data[key]
		if !ok || mapVal == nil {
			continue
		}

		fieldVal := val.Field(i)
		mv := reflect.ValueOf(mapVal)

		if fieldVal.Kind() == reflect.Ptr {
			elem := reflect.New(fieldVal.Type().Elem())
			if mv.Type().AssignableTo(elem.Elem().Type()) {
				elem.Elem().Set(mv)
			} else if mv.Type().ConvertibleTo(elem.Elem().Type()) {
				elem.Elem().Set(mv.Convert(elem.Elem().Type()))
			} else {
				continue
			}
			fieldVal.Set(elem)
		} else {
			if mv.Type().AssignableTo(fieldVal.Type()) {
				fieldVal.Set(mv)
			} else if mv.Type().ConvertibleTo(fieldVal.Type()) {
				fieldVal.Set(mv.Convert(fieldVal.Type()))
			}
		}
	}

	return nil
}
