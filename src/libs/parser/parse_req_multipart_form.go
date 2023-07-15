package parser

import (
	"reflect"
	"sim-puskesmas/src/libs/validator"

	"github.com/gofiber/fiber/v2"
)

func ParseReqMultipartForm[T any](c *fiber.Ctx, req T) error {
	parsedMultipart, err := c.MultipartForm()
	if err != nil {
		logger.Error(err)
		return err
	}

	reqValue := reflect.ValueOf(req).Elem()
	reqType := reqValue.Type()

	for k, v := range parsedMultipart.Value {
		decode(&k, &v, &reqType, &reqValue)
	}
	for k, v := range parsedMultipart.File {
		decode(&k, &v, &reqType, &reqValue)
	}

	if err := validator.Struct(req); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func decode[T any](k *string, v *[]T, reqType *reflect.Type, reqValue *reflect.Value) {
	index := -1
	for i := 0; i < (*reqType).NumField(); i++ {
		fieldType := (*reqType).Field(i)
		if *k == fieldType.Tag.Get("form") {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}

	fieldValue := reqValue.Field(index)
	fieldType := (*reqType).Field(index)
	if reflect.ValueOf((*v)[0]).Kind() == reflect.Pointer {
		fieldValue.Set(reflect.ValueOf((*v)[0]).Convert(fieldType.Type))
	} else {
		fieldValue.Set(reflect.ValueOf(&(*v)[0]).Convert(fieldType.Type))
	}
}
