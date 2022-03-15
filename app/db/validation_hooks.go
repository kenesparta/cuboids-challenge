package db

import (
	"cuboid-challenge/app/models"
	"reflect"

	"gorm.io/gorm"
)

func registerDBValidationsHooks(db *gorm.DB) {
	callback := db.Callback()
	if callback.Create().Get("validations:validate") == nil {
		if err := callback.Create().Before("gorm:before_create").Register("validations:validate", validateModel); err != nil {
			panic(err)
		}
	}

	if callback.Update().Get("validations:validate") == nil {
		if err := callback.Update().Before("gorm:before_update").Register("validations:validate", validateModel); err != nil {
			panic(err)
		}
	}
}

func validateModel(dbs *gorm.DB) {
	if dbs.Error != nil || dbs.Statement.Schema == nil || dbs.Statement.SkipHooks {
		return
	}

	record := getRecord(dbs)
	if record == nil {
		return
	}

	if ok, valErr := models.Validate(record); !ok {
		_ = dbs.AddError(valErr)
	}
}

func getRecord(db *gorm.DB) interface{} {
	record := db.Statement.ReflectValue.Interface()
	val := reflect.ValueOf(record)

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	return record
}
