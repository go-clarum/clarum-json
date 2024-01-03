package comparator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/goclarum/clarum/json/internal/path"
	"github.com/goclarum/clarum/json/recorder"
	"log/slog"
	"reflect"
	"strconv"
)

const ignoreFlag = "@ignore@"

type options struct {
	strictObjectCheck bool
	pathsToIgnore     []string
	logger            *slog.Logger
	recorder          recorder.Recorder
}

// Comparator used for comparing JSON structures. It returns detailed errors about how the compared structures do not match.
// The errors may also contain JSON paths that point to the fields discovered.
//
// Always create a comparator using the [NewComparator] builder.
//
// The comparator does not fail fast and goes over the entire structure for validation.
// It is goroutine safe.
type Comparator struct {
	options
}

func (comparator *Comparator) Compare(expected []byte, actual []byte) (string, error) {
	var expectedJsonObject, actualJsonObject any
	comparator.logger.Debug(fmt.Sprintf("json comparator - comparing [%s] to [%s]", expected, actual))

	expectedJsonObject, err1 := unmarshalJson(expected)
	if err1 != nil {
		return "", err1
	}

	actualJsonObject, err2 := unmarshalJson(actual)
	if err2 != nil {
		return "", err2
	}

	typeOfExpected := reflect.TypeOf(expectedJsonObject)
	typeOfActual := reflect.TypeOf(actualJsonObject)

	var compareErrors []error

	if typeOfExpected.Kind() != typeOfActual.Kind() {
		compareErrors = append(compareErrors,
			errors.New(fmt.Sprintf("root object mismatch - expected [%s] but found [%s]",
				convertToJsonType(typeOfExpected), convertToJsonType(typeOfActual))))
	} else if typeOfExpected.Kind() == reflect.Map {
		compareErrors = comparator.compareJsonMaps("$",
			expectedJsonObject.(map[string]any), actualJsonObject.(map[string]any),
			"", compareErrors)
	} else if typeOfExpected.Kind() == reflect.Slice {
		compareErrors = comparator.compareSlices("$",
			expectedJsonObject.([]interface{}), actualJsonObject.([]interface{}),
			"", compareErrors)
	}

	if len(compareErrors) > 0 {
		comparator.logger.Debug(fmt.Sprintf("json comparator - JSON structures do not match"))
	} else {
		comparator.logger.Debug(fmt.Sprintf("json comparator - JSON structures match"))
	}

	return comparator.recorder.GetLog(), errors.Join(compareErrors...)
}

func (comparator *Comparator) compareJsonMaps(parentPath string, expected map[string]any, actual map[string]any,
	logIndent string, compareErrors []error) []error {
	currIndent := logIndent + "  "

	compareErrors = handleFieldsCheck(parentPath, expected, actual, comparator.strictObjectCheck,
		comparator.recorder, logIndent, compareErrors)

	for key, expectedValue := range expected {
		if actualValue, exists := actual[key]; exists {
			comparator.recorder.AppendFieldName(currIndent, key)

			expectedValueType := reflect.TypeOf(expectedValue)
			actualValueType := reflect.TypeOf(actualValue)

			ignoreValueValidation := ignoreValue(expectedValueType, expectedValue)
			if ignoreValueValidation {
				comparator.recorder.AppendIgnoreField(currIndent, parentPath)
				continue
			}

			if expectedValueType.Kind() != actualValueType.Kind() {
				compareErrors = handleTypeMismatch(path.GetObjectChildPath(parentPath, key),
					expectedValueType, actualValueType, comparator.recorder, compareErrors)
			} else {
				// we only consider JSON Kinds, since the Unmarshal already parsed & checked them
				switch actualValueType.Kind() {
				case reflect.String:
					expectedString := expectedValue.(string)
					actualString := actualValue.(string)

					compareErrors = compareValue(path.GetObjectChildPath(parentPath, key),
						expectedString != actualString,
						expectedString, actualString, comparator.recorder, logIndent, compareErrors)
				case reflect.Float64:
					compareErrors = compareValue(path.GetObjectChildPath(parentPath, key),
						expectedValue.(float64) != actualValue.(float64),
						formatFloat(expectedValue), formatFloat(actualValue), comparator.recorder, logIndent, compareErrors)
				case reflect.Bool:
					expectedBool := expectedValue.(bool)
					actualBool := actualValue.(bool)

					compareErrors = compareValue(path.GetObjectChildPath(parentPath, key),
						expectedBool != actualBool,
						strconv.FormatBool(expectedBool), strconv.FormatBool(actualBool), comparator.recorder, logIndent, compareErrors)
				case reflect.Slice:
					compareErrors = comparator.compareSlices(path.GetObjectChildPath(parentPath, key),
						expectedValue.([]interface{}), actualValue.([]interface{}),
						currIndent, compareErrors)
				case reflect.Map:
					compareErrors = comparator.compareJsonMaps(path.GetObjectChildPath(parentPath, key),
						expectedValue.(map[string]any), actualValue.(map[string]any),
						currIndent, compareErrors)
				}
			}
		} else {
			compareErrors = handleMissingField(path.GetObjectChildPath(parentPath, key),
				key, currIndent, comparator.recorder, compareErrors)
		}
	}

	if comparator.strictObjectCheck {
		compareErrors = handleUnexpectedFields(parentPath, expected, actual, comparator.recorder,
			currIndent, compareErrors)
	}

	comparator.recorder.AppendEndObject(logIndent, parentPath)
	return compareErrors
}

// Arrays in json are represented as slices of type interface because they can contain anything.
// Each item in the slice can be of any valid JSON type.
func (comparator *Comparator) compareSlices(parentPath string, expected []interface{}, actual []interface{},
	currIndent string, compareErrors []error) []error {
	comparator.recorder.AppendStartArray(currIndent, parentPath)

	expectedLen := len(expected)
	actualLen := len(actual)
	if expectedLen != actualLen {
		comparator.recorder.AppendValidationErrorSignal(fmt.Sprintf("size mismatch - expected [%d]", expectedLen)).
			AppendEndArray(currIndent, parentPath)
		return append(compareErrors,
			errors.New(fmt.Sprintf("[%s] - array size mismatch - expected [%d] but received [%d]", parentPath, expectedLen, actualLen)))
	} else {
		comparator.recorder.AppendNewLine()
	}

	valIdent := currIndent + "  "
	for i, expectedValue := range expected {
		expectedValueType := reflect.TypeOf(expectedValue)
		actualValue := actual[i]
		actualValueType := reflect.TypeOf(actualValue)
		jsonPathArray := path.GetArrayIndexPath(parentPath, i)

		ignoreValueValidation := ignoreValue(expectedValueType, expectedValue)
		if ignoreValueValidation {
			comparator.recorder.AppendIgnoreField(valIdent, jsonPathArray)
			continue
		}

		if expectedValueType.Kind() != actualValueType.Kind() {
			comparator.recorder.AppendValue(valIdent, jsonPathArray, actualValue, actualValueType.Kind())
			baseErrorMessage := fmt.Sprintf("value type mismatch - expected [%s] but found [%s]",
				convertToJsonType(expectedValueType), convertToJsonType(actualValueType))

			compareErrors = append(compareErrors, errors.New(fmt.Sprintf("[%s] - %s", jsonPathArray,
				baseErrorMessage)))
			comparator.recorder.AppendValidationErrorSignal(baseErrorMessage)
		} else {
			switch actualValueType.Kind() {
			case reflect.String:
				expectedString := expectedValue.(string)
				actualString := actualValue.(string)

				compareErrors = compareValue(jsonPathArray,
					expectedString != actualString,
					expectedString, actualString, comparator.recorder, valIdent, compareErrors)
			case reflect.Float64:
				compareErrors = compareValue(jsonPathArray,
					expectedValue.(float64) != actualValue.(float64),
					formatFloat(expectedValue), formatFloat(actualValue), comparator.recorder, valIdent, compareErrors)
			case reflect.Bool:
				expectedBool := expectedValue.(bool)
				actualBool := actualValue.(bool)

				compareErrors = compareValue(jsonPathArray,
					expectedBool != actualBool,
					strconv.FormatBool(expectedBool), strconv.FormatBool(actualBool), comparator.recorder, valIdent, compareErrors)
			case reflect.Slice:
				compareErrors = comparator.compareSlices(jsonPathArray,
					expectedValue.([]interface{}), actualValue.([]interface{}),
					valIdent, compareErrors)
			case reflect.Map:
				compareErrors = comparator.compareJsonMaps(jsonPathArray,
					expectedValue.(map[string]any), actualValue.(map[string]any),
					valIdent, compareErrors)
			}
		}
	}
	comparator.recorder.AppendEndArray(currIndent, parentPath)
	return compareErrors
}

// When describing types we have to consider that the users will think of JSON types when reading logs & error messages.
// This is why we have to translate Go types into JSON types.
//
// json.Unmarshal returns a map[string]interface{} with all the fields of the JSON object:
// - number is a reflect.Float64
// - string is a reflect.String
// - boolean is a reflect.Bool
// - array is a reflect.Slice
// - struct is a reflect.Map
func convertToJsonType(goType reflect.Type) string {
	switch goType.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Float64:
		return "number"
	case reflect.Map:
		return "object"
	case reflect.Slice:
		return "array"
	default:
		return goType.String()
	}
}

func handleFieldsCheck(pathParent string, expected map[string]any, actual map[string]any, strictObjectCheck bool,
	recorder recorder.Recorder, indent string, compareErrors []error) []error {
	if strictObjectCheck && len(expected) != len(actual) {
		recorder.AppendStartObject(indent, pathParent).
			AppendValidationErrorSignal("number of fields does not match")

		compareErrors = append(compareErrors,
			errors.New(fmt.Sprintf("[%s] - number of fields does not match", pathParent)))
	} else {
		recorder.AppendStartObject(indent, pathParent).AppendNewLine()
	}
	return compareErrors
}

func handleUnexpectedFields(pathParent string, expected map[string]any, actual map[string]any,
	recorder recorder.Recorder, indent string, compareErrors []error) []error {
	for key := range actual {
		if _, exists := expected[key]; !exists {
			recorder.AppendFieldName(indent, key).
				AppendValidationErrorSignal("unexpected field")

			compareErrors = append(compareErrors,
				errors.New(fmt.Sprintf("[%s] - unexpected field", path.GetObjectChildPath(pathParent, key))))
		}
	}

	return compareErrors
}

func handleTypeMismatch(path string, expectedValueType reflect.Type, actualValueType reflect.Type,
	recorder recorder.Recorder, compareErrors []error) []error {

	baseErrorMessage := fmt.Sprintf("type mismatch - expected [%s] but found [%s]",
		convertToJsonType(expectedValueType), convertToJsonType(actualValueType))

	compareErrors = append(compareErrors, errors.New(fmt.Sprintf("[%s] - %s", path, baseErrorMessage)))
	recorder.AppendValidationErrorSignal(baseErrorMessage)

	return compareErrors
}

func compareValue(path string, mismatch bool, expectedValue string, actualValue string, recorder recorder.Recorder,
	indent string, compareErrors []error) []error {
	recorder.AppendValue(indent, path, actualValue, reflect.String)

	if mismatch {
		compareErrors = append(compareErrors,
			errors.New(fmt.Sprintf("[%s] - value mismatch - expected [%s] but received [%s]", path, expectedValue, actualValue)))
		recorder.AppendValidationErrorSignal(fmt.Sprintf("value mismatch - expected [%s]", expectedValue))
	} else {
		recorder.AppendNewLine()
	}

	return compareErrors
}

func handleMissingField(path string, fieldName string, indent string, recorder recorder.Recorder, compareErrors []error) []error {
	compareErrors = append(compareErrors, errors.New(fmt.Sprintf("[%s] - field is missing", path)))
	recorder.AppendMissingFieldErrorSignal(indent, fieldName)

	return compareErrors
}

// We rely on json.Unmarshal to detect invalid json structures here.
func unmarshalJson(rawJson []byte) (any, error) {
	var result any
	if err := json.Unmarshal(rawJson, &result); err != nil {
		return nil, handleError("unable to parse JSON - error [%s] - from string [%s]", err, rawJson)
	}

	return result, nil
}

// We have to trim trailing zeroes from the parsed float64 number before logging them.
func formatFloat(expectedValue any) string {
	return strconv.FormatFloat(expectedValue.(float64), 'f', -1, 64)
}

func ignoreValue(valueType reflect.Type, value any) bool {
	if valueType.Kind() == reflect.String {
		return value == ignoreFlag
	}

	return false
}

func handleError(format string, a ...any) error {
	errorMessage := fmt.Sprintf(format, a...)
	return errors.New(errorMessage)
}
