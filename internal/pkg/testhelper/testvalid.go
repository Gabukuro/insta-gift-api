package testhelper

import (
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type (
	ValidationFunc func(value any) bool
	ValidationMap  map[string]ValidationFunc
)

func VerifyHTTPResponseBodyJSON(
	t *testing.T,
	expectedResponseBody string,
	responseBody io.ReadCloser,
	customValidations ValidationMap,
) bool {
	t.Helper()

	bodyBytes, err := io.ReadAll(responseBody)
	if err != nil {
		panic(err)
	}

	if !json.Valid([]byte(expectedResponseBody)) || !json.Valid(bodyBytes) {
		return assert.Equal(t, expectedResponseBody, string(bodyBytes))
	}

	var expected any

	err = json.Unmarshal([]byte(expectedResponseBody), &expected)
	if err != nil {
		panic(err)
	}

	expectedMapped := getMappedObject("", expected, nil)

	var body any

	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		panic(err)
	}

	bodyMapped := getMappedObject("", body, nil)

	validations := getDefaultValidations()
	for validationKey, validationFunc := range customValidations {
		validations[validationKey] = validationFunc
	}

	runValidations(expectedMapped, bodyMapped, validations)

	return assert.Equal(t, expectedMapped, bodyMapped)
}

func runValidations(expectedMapped, bodyMapped map[string]any, validations ValidationMap) {
	for expectedKey, expectedValue := range expectedMapped {
		for validationKey, validationFunc := range validations {
			if expectedValue == validationKey {
				if validationFunc(bodyMapped[expectedKey]) {
					delete(bodyMapped, expectedKey)
					delete(expectedMapped, expectedKey)
				}
			}
		}
	}
}

func getDefaultValidations() ValidationMap {
	return ValidationMap{
		"!uuid": validateUUID,
		"!date": validateDate,
	}
}

func validateUUID(value any) bool {
	valueString, ok := value.(string)
	if !ok {
		return false
	}

	if _, err := uuid.Parse(valueString); err != nil {
		return false
	}

	return true
}

func validateDate(value any) bool {
	valueString, ok := value.(string)
	if !ok {
		return false
	}

	if _, err := time.Parse(time.RFC3339, valueString); err != nil {
		return false
	}

	return true
}

func getMappedObject(prefix string, input any, outputRef map[string]any) map[string]any {
	if outputRef == nil {
		outputRef = make(map[string]any)
	}

	if input == nil {
		outputRef[prefix] = input
		return outputRef
	}

	switch reflect.TypeOf(input).Kind() {
	case reflect.Map:
		mapCase(prefix, input, outputRef)
	case reflect.Array, reflect.Slice:
		sliceCase(prefix, input, outputRef)
	case reflect.Pointer:
		input = reflect.ValueOf(input).Elem()
		getMappedObject(prefix, input, outputRef)
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Chan, reflect.Func, reflect.Interface, reflect.String, reflect.Struct, reflect.UnsafePointer:
		outputRef[prefix] = input
	case reflect.Invalid:
		panic(reflect.Invalid.String())
	default:
		panic(reflect.Invalid.String())
	}

	return outputRef
}

func mapCase(prefix string, input any, outputRef map[string]any) {
	valueMap, _ := input.(map[string]any)
	for key, value := range valueMap {
		getMappedObject(fmt.Sprintf("%s.%s", prefix, key), value, outputRef)
	}
}

func sliceCase(prefix string, input any, outputRef map[string]any) {
	value := reflect.ValueOf(input)

	if value.Len() == 0 {
		outputRef[prefix] = input
	}

	for index := 0; index < value.Len(); index++ {
		getMappedObject(fmt.Sprintf("%s.%d", prefix, index), value.Index(index).Interface(), outputRef)
	}
}
