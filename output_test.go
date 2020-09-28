package usecase_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/usecase"
)

func TestOutputWithEmbeddedWriter_SetWriter(t *testing.T) {
	w := bytes.NewBuffer(nil)
	o := usecase.OutputWithEmbeddedWriter{}
	o.SetWriter(w)

	_, err := o.Write([]byte("hello"))
	require.NoError(t, err)

	assert.Equal(t, "hello", w.String())
}

func TestOutputWithNoContent_NoContent(t *testing.T) {
	o := usecase.OutputWithNoContent{}
	assert.True(t, o.NoContent())
	o.SetNoContent(false)
	assert.False(t, o.NoContent())
}
