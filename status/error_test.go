package status_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/usecase/status"
)

func TestWrap(t *testing.T) {
	err := status.Wrap(errors.New("failed"), status.AlreadyExists)
	assert.EqualError(t, err, "already exists: failed")
	assert.True(t, errors.Is(err, status.AlreadyExists))
	assert.False(t, errors.Is(err, status.NotFound))
	assert.EqualError(t, err.(interface{ Unwrap() error }).Unwrap(), "failed")
	assert.Equal(t, status.AlreadyExists, err.(interface{ Status() status.Code }).Status())
}

func TestCode_Error(t *testing.T) {
	assert.Equal(t, "deadline exceeded", status.DeadlineExceeded.Error())
}

func TestWithDescription(t *testing.T) {
	err := status.WithDescription(status.AlreadyExists, "This is a description.")
	assert.EqualError(t, err, "already exists")
	assert.True(t, errors.Is(err, status.AlreadyExists))
	assert.False(t, errors.Is(err, status.NotFound))

	d, ok := err.(interface {
		Description() string
	})
	assert.True(t, ok)
	assert.Equal(t, "This is a description.", d.Description())
}
