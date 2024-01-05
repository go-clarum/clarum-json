package internal

import (
	"github.com/go-clarum/clarum-json/recorder"
	"reflect"
)

type NoopRecorder struct {
}

func NewNoopRecorder() recorder.Recorder {
	return &NoopRecorder{}
}

func (recorder *NoopRecorder) GetLog() string {
	return ""
}

func (recorder *NoopRecorder) AppendFieldName(indent string, fieldName string) recorder.Recorder {
	return recorder
}

func (recorder *NoopRecorder) AppendIgnoreField(indent string, jsonPath string) recorder.Recorder {
	return recorder
}

func (recorder *NoopRecorder) AppendValue(indent string, path string, value any, kind reflect.Kind) recorder.Recorder {
	return recorder
}

func (recorder *NoopRecorder) AppendValidationErrorSignal(message string) recorder.Recorder {
	return recorder
}

func (recorder *NoopRecorder) AppendMissingFieldErrorSignal(indent string, path string) recorder.Recorder {
	return recorder
}

func (recorder *NoopRecorder) AppendStartObject(indent string, path string) recorder.Recorder {
	return recorder
}

func (recorder *NoopRecorder) AppendEndObject(indent string, path string) recorder.Recorder {
	return recorder
}

func (recorder *NoopRecorder) AppendStartArray(indent string, path string) recorder.Recorder {
	return recorder
}

func (recorder *NoopRecorder) AppendEndArray(indent string, path string) recorder.Recorder {
	return recorder
}

func (recorder *NoopRecorder) AppendNewLine() recorder.Recorder {
	return recorder
}
