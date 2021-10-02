//go:build go1.18
// +build go1.18

package usecase_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/usecase"
)

func TestNewInteractor(t *testing.T) {
	u := usecase.NewInteractor(func(ctx context.Context, input int, output *string) error {
		*output = strconv.Itoa(input)

		return nil
	})

	u.SetDescription("Foo.")

	ctx := context.Background()

	var out string
	assert.NoError(t, u.Interact(ctx, 123, &out))
	assert.Equal(t, "123", out)

	out = ""
	assert.NoError(t, u.Invoke(ctx, 123, &out))
	assert.Equal(t, "123", out)
}
