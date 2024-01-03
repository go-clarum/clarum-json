package recorder

import (
	"reflect"
)

// Recorder is an interface that allows to create a human-readable output from the comparator's validation errors.
// The [Comparator] calls the methods below at different points during the validation.
// The use of a Recorder is entirely optional.
type Recorder interface {
	AppendFieldName(indent string, fieldName string) Recorder
	AppendIgnoreField(indent string, jsonPath string) Recorder
	AppendValue(indent string, path string, value any, kind reflect.Kind) Recorder
	AppendValidationErrorSignal(message string) Recorder
	AppendMissingFieldErrorSignal(indent string, path string) Recorder
	AppendStartObject(indent string, path string) Recorder
	AppendEndObject(indent string, path string) Recorder
	AppendStartArray(indent string, path string) Recorder
	AppendEndArray(indent string, path string) Recorder
	AppendNewLine() Recorder
	GetLog() string
}
