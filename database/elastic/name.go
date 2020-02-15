package elastic

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
)

func structIndex(v interface{}) (string, error) {

	typo := reflect.TypeOf(v)
	if typo == nil {
		return "", fmt.Errorf("value is nil")
	}

	if typo.Kind() == reflect.Ptr {
		typo = typo.Elem()
	}

	if typo.Kind() != reflect.Struct {
		return "", fmt.Errorf("value must be struct or pointer of struct, got %v", typo)
	}

	return gorm.ToColumnName(typo.Name()), nil
}

func sliceIndex(v interface{}) (string, error) {

	typo := reflect.TypeOf(v)
	if typo == nil {
		return "", fmt.Errorf("nil value")
	}

	if typo.Kind() != reflect.Slice {
		return "", fmt.Errorf("value must be slice, got %v", typo)
	}
	typo = typo.Elem()

	if typo.Kind() == reflect.Ptr {
		typo = typo.Elem()
	}

	if typo.Kind() != reflect.Struct {
		return "", fmt.Errorf("value's element must be struct or pointer of struct, got %v", typo)
	}

	return gorm.ToColumnName(typo.Name()), nil
}
