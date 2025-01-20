package validator

import (
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	Validate *validator.Validate

	validateSyncOne sync.Once
)

func init() {
	validateSyncOne.Do(func() {
		Validate = validator.New()

		Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			fieldName := field.Tag.Get("json")
			if fieldName != "" {
				return fieldName
			}

			return field.Tag.Get("name")
		})
	})
}
