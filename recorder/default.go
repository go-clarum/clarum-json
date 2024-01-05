package recorder

import (
	"fmt"
	"github.com/go-clarum/clarum-json/internal/path"
	"reflect"
	"strings"
)

// DefaultRecorder used in clarum validations.
// As this implementation uses the strings.Builder, it is not goroutine safe!
type DefaultRecorder struct {
	logResult strings.Builder
}

func NewDefaultRecorder() Recorder {
	return &DefaultRecorder{}
}

func (recorder *DefaultRecorder) GetLog() string {
	return recorder.logResult.String()
}

func (recorder *DefaultRecorder) AppendFieldName(indent string, fieldName string) Recorder {
	recorder.logResult.WriteString(fmt.Sprintf("%s\"%s\": ", indent, fieldName))
	return recorder
}

func (recorder *DefaultRecorder) AppendIgnoreField(indent string, jsonPath string) Recorder {
	childOfArray := path.IsChildOfArray(jsonPath)

	if childOfArray {
		recorder.logResult.WriteString(fmt.Sprintf("%s <-- ignoring field\n", indent))
	} else {
		recorder.logResult.WriteString(fmt.Sprintf("%s <-- ignoring field\n", ""))
	}

	return recorder
}

func (recorder *DefaultRecorder) AppendValue(indent string, jsonPath string, value any, kind reflect.Kind) Recorder {
	childOfArray := path.IsChildOfArray(jsonPath)

	var indentToSet string
	if childOfArray {
		indentToSet = indent
	} else {
		indentToSet = ""
	}

	if kind == reflect.Map {
		recorder.logResult.WriteString(fmt.Sprintf("%sobject,", indentToSet))
	} else if kind == reflect.Slice {
		recorder.logResult.WriteString(fmt.Sprintf("%sarray,", indentToSet))
	} else if kind != reflect.Invalid {
		recorder.logResult.WriteString(fmt.Sprintf("%s%v,", indentToSet, value))
	}
	return recorder
}

func (recorder *DefaultRecorder) AppendValidationErrorSignal(message string) Recorder {
	recorder.logResult.WriteString(fmt.Sprintf(" <-- %s\n", message))
	return recorder
}

func (recorder *DefaultRecorder) AppendMissingFieldErrorSignal(indent string, path string) Recorder {
	recorder.logResult.WriteString(fmt.Sprintf("%s X-- missing field [%s]\n", indent, path))
	return recorder
}

func (recorder *DefaultRecorder) AppendStartObject(indent string, jsonPath string) Recorder {
	childOfArray := path.IsChildOfArray(jsonPath)

	if childOfArray {
		recorder.logResult.WriteString(fmt.Sprintf("%s{", indent))
	} else {
		recorder.logResult.WriteString(fmt.Sprintf("%s{", ""))
	}
	return recorder
}

func (recorder *DefaultRecorder) AppendEndObject(indent string, jsonPath string) Recorder {
	root := path.IsRoot(jsonPath)

	if root {
		recorder.logResult.WriteString(fmt.Sprintf("%s}\n", ""))
	} else {
		recorder.logResult.WriteString(fmt.Sprintf("%s},\n", indent))
	}
	return recorder
}

func (recorder *DefaultRecorder) AppendStartArray(indent string, jsonPath string) Recorder {
	childOfArray := path.IsChildOfArray(jsonPath)

	if childOfArray {
		recorder.logResult.WriteString(fmt.Sprintf("%s[", indent))
	} else {
		recorder.logResult.WriteString(fmt.Sprintf("%s[", ""))
	}
	return recorder
}

func (recorder *DefaultRecorder) AppendEndArray(indent string, jsonPath string) Recorder {
	root := path.IsRoot(jsonPath)

	if root {
		recorder.logResult.WriteString(fmt.Sprintf("%s]\n", indent))
	} else {
		recorder.logResult.WriteString(fmt.Sprintf("%s],\n", indent))
	}
	return recorder
}

func (recorder *DefaultRecorder) AppendNewLine() Recorder {
	recorder.logResult.WriteString("\n")
	return recorder
}
