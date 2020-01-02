package usecase

import "io"

// OutputWithWriter defines output with streaming writer.
type OutputWithWriter interface {
	SetWriter(w io.Writer)
}

// OutputWithEmbeddedWriter implements streaming use case output.
type OutputWithEmbeddedWriter struct {
	io.Writer
}

// SetWriter implements OutputWithWriter.
func (o *OutputWithEmbeddedWriter) SetWriter(w io.Writer) {
	o.Writer = w
}

// OutputWithNoContent is embeddable structure to provide conditional output discard state.
type OutputWithNoContent struct {
	enabled bool
}

// SetNoContent controls output discard state.
func (o *OutputWithNoContent) SetNoContent(enabled bool) {
	o.enabled = enabled
}

// NoContent returns output discard state.
func (o OutputWithNoContent) NoContent() bool {
	return o.enabled
}
