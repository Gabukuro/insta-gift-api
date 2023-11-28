package testhelper

import (
	"math/rand"
	"testing"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func VerifyError(t *testing.T, expectedString string, err error) bool {
	t.Helper()

	if expectedString == "" {
		if err == nil {
			return false
		}

		assert.Equal(t, expectedString, err.Error())
	}

	if expectedString != "" {
		if err == nil {
			assert.Equal(t, expectedString, "")
			return true
		}

		assert.Equal(t, expectedString, err.Error())
	}

	if err != nil {
		return true
	}

	return false
}

func VerifyHTTPResponseBody(t *testing.T, expectedResponseBody string, responseBody []byte) bool {
	t.Helper()

	if !json.Valid([]byte(expectedResponseBody)) {
		return assert.Equal(t, expectedResponseBody, string(responseBody))
	}

	expected := make(map[string]any)

	err := json.Unmarshal([]byte(expectedResponseBody), &expected)
	if err != nil {
		panic(err)
	}

	actual := make(map[string]any)

	err = json.Unmarshal(responseBody, &actual)
	if err != nil {
		panic(err)
	}

	for key, value := range expected {
		if value == "!uuid" {
			if id, ok := actual[key]; ok {
				_, err := uuid.Parse(id.(string))
				if err != nil {
					t.Errorf("%s: %v", key, err)
				}

				delete(expected, key)
				delete(actual, key)
			}
		}
	}

	return assert.Equal(t, expected, actual)
}

// RandStringBytes  generate random string from given string list with provided size.
func RandStringBytes(stringList string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = stringList[rand.Intn(len(stringList))]
	}

	return string(b)
}
