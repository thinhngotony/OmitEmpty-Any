package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func OmitEmptyFields(data interface{}) interface{} {
	switch t := data.(type) {
	case map[string]interface{}:
		return omitEmptyMap(t)
	case []interface{}:
		return omitEmptySlice(t)
	case interface{}:
		return omitEmptyInterface(t)
	default:
		return data
	}
}

func omitEmptyMap(val map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range val {
		newValue := OmitEmptyFields(value) // Recursively call for nested structures
		if !isEmptyValue(newValue) {
			result[strings.ToLower(key)] = newValue
		}
	}

	return result
}

func omitEmptySlice(val []interface{}) []interface{} {
	result := []interface{}{}

	for _, v := range val {
		if !isEmptyValue(v) {
			result = append(result, OmitEmptyFields(v))
		}
	}

	return result
}

func omitEmptyInterface(val interface{}) interface{} {
	valReflect := reflect.ValueOf(val)
	if valReflect.Kind() == reflect.Ptr {
		valReflect = valReflect.Elem()
	}

	switch valReflect.Kind() {
	case reflect.Struct:
		return omitEmptyFieldsStruct(valReflect)
	case reflect.Map:
		return omitEmptyMap(valReflect.Interface().(map[string]interface{}))
	default:
		return val
	}
}

func omitEmptyFieldsStruct(valReflect reflect.Value) interface{} {
	result := make(map[string]interface{})

	for i := 0; i < valReflect.NumField(); i++ {
		field := valReflect.Type().Field(i)
		fieldValue := valReflect.Field(i)

		newValue := OmitEmptyFields(fieldValue.Interface()) // Recursively call for nested structures

		if !isEmptyValue(newValue) {
			result[strings.ToLower(field.Name)] = newValue
		}
	}

	return result
}

// Function to check if a value is considered empty (modify as needed)
func isEmptyValue(v interface{}) bool {
	// Expand this logic to handle additional types within the interface if needed
	switch t := v.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(t) == "" // Use type assertion for string type
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return t == 0
	case bool:
		return !t
	case map[string]interface{}:
		return len(t) == 0
	default:
		return false
	}
}

type TestStruct struct {
	Name     string
	Age      int
	Email    string
	Nullable *string
	EmptyMap map[string]interface{}
}

func test1() {
	// Example usage with a struct
	data := TestStruct{
		Name:     "John",
		Age:      30,
		Email:    "   ", // Whitespace string
		Nullable: nil,
		EmptyMap: map[string]interface{}{},
	}

	// Process the struct with OmitEmptyFields
	processedData := OmitEmptyFields(data)

	// Print the resulting data (without empty fields)
	fmt.Println(InterfaceToString(processedData))
}

func InterfaceToString(data interface{}) string {
	manifestJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", data)
	}
	return string(manifestJson)
}

func test2() {
	// Example usage with your interface structure (including whitespace strings)
	data := map[string]interface{}{
		"addons": []interface{}{
			map[string]interface{}{
				"api_url": "   ", // Whitespace string
				"chart": map[string]interface{}{
					"api_version": " ", // Whitespace string
					"deleted":     "",
					"deployed_at": "",
					"name":        "",
					"namespace":   "default",
					"release":     "mysql-release-name",
					"resources":   map[string]interface{}{},
					"revision":    "0",
					"status":      "",
					"version":     "",
				},
				"error":       "",
				"error_code":  "",
				"name":        "mysql",
				"status_code": 0,
			},
		},
		"adminConfig":     "",
		"ffk8s_version":   "",
		"kubeadmConfig":   "",
		"kubekey_version": "",
		"network":         nil,
		"spec":            "",
		"version":         "",
	}

	// Process the interface with OmitEmptyFields
	processedData := OmitEmptyFields(data)

	// Print the resulting data (without empty fields)
	fmt.Println(InterfaceToString(processedData))
}

func main() {
	test1()
	test2()
}
