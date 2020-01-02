package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func TestError(t *testing.T) {
	e := usecase.Error{
		AppCode:    123,
		StatusCode: status.FailedPrecondition,
		Context: map[string]interface{}{
			"foo": "bar",
		},
	}
	assert.Equal(t, 123, e.AppErrCode())
	assert.Equal(t, e.Context, e.Fields())
	assert.Equal(t, status.FailedPrecondition, e.Status())
	assert.EqualError(t, e, "failed precondition")
	assert.EqualError(t, e.Unwrap(), "failed precondition")

	e.Value = errors.New("failed")
	assert.EqualError(t, e, "failed precondition: failed")
	assert.EqualError(t, e.Unwrap(), "failed precondition: failed")

	e.StatusCode = 0
	assert.EqualError(t, e, "failed")
	assert.EqualError(t, e.Unwrap(), "failed")
}
