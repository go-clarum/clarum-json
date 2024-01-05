package comparator

import (
	"github.com/go-clarum/clarum-json/internal"
	"github.com/go-clarum/clarum-json/recorder"
	"log/slog"
)

type Builder struct {
	options
}

// NewComparator is the builder initiator. Always use the builder to create a [Comparator]
// as this will set the default options.
func NewComparator() *Builder {
	return &Builder{
		options{
			strictObjectCheck: true,
			// pathsToIgnore:     []string{},
			logger:   slog.Default(),
			recorder: internal.NewNoopRecorder(),
		},
	}
}

// StrictObjectCheck determines if the [Comparator] will do a strict check on object fields.
// If set to 'true', the following checks will be done:
// - actual JSON has the same number of fields
// - actual JSON has extra unexpected fields
//
// Default is 'true'.
func (builder *Builder) StrictObjectCheck(check bool) *Builder {
	builder.strictObjectCheck = check
	return builder
}

// PathsToIgnore is a list of json paths that the comparator will ignore during validation.
// Default is empty.
/*
func (builder *Builder) PathsToIgnore(paths ...string) *Builder {
	builder.pathsToIgnore = append(builder.pathsToIgnore, paths...)
	return builder
}
*/

func (builder *Builder) Logger(logger *slog.Logger) *Builder {
	builder.logger = logger
	return builder
}

// Recorder to be used.
// Default is the NoopRecorder.
func (builder *Builder) Recorder(recorder recorder.Recorder) *Builder {
	builder.recorder = recorder
	return builder
}

func (builder *Builder) Build() *Comparator {
	return &Comparator{builder.options}
}
