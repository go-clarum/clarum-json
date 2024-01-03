package comparator

import (
	"github.com/goclarum/clarum/json/internal"
	"testing"
)

func TestComparatorBuilderDefaults(t *testing.T) {
	comparator := NewComparator().Build()

	if !comparator.strictObjectCheck {
		t.Error("default StrictObjectCheck must be true")
	}
	if len(comparator.pathsToIgnore) != 0 {
		t.Error("default PathsToIgnore is empty")
	}
	if comparator.logger == nil {
		t.Error("default Logger must not be nil")
	}
	if _, isNoopRecorder := comparator.recorder.(*internal.NoopRecorder); !isNoopRecorder {
		t.Error("default Recorder must be NoopRecorder")
	}
}
