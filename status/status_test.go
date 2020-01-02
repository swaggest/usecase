package status_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/usecase/status"
)

func TestCode_Status(t *testing.T) {
	assert.Equal(t, "FAILED_PRECONDITION", status.FailedPrecondition.String())
	assert.Equal(t, status.FailedPrecondition, status.FailedPrecondition.Status())
}
