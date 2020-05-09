package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func TestInteract_Interact(t *testing.T) {
	err := errors.New("failed")
	inp := new(string)
	out := new(int)

	i := usecase.Interact(func(ctx context.Context, input, output interface{}) error {
		assert.Equal(t, inp, input)
		assert.Equal(t, out, output)

		return err
	})
	assert.Equal(t, err, i.Interact(context.Background(), inp, out))
}

func TestInfo(t *testing.T) {
	i := usecase.Info{}
	i.SetDescription("Description")
	i.SetTitle("Title")
	i.SetTags("tag1", "tag2")
	i.SetIsDeprecated(true)
	i.SetExpectedErrors(usecase.Error{StatusCode: status.InvalidArgument})

	assert.Equal(t, "Description", i.Description())
	assert.Equal(t, "Title", i.Title())
	assert.Equal(t, []string{"tag1", "tag2"}, i.Tags())
	assert.Equal(t, true, i.IsDeprecated())
	assert.Equal(t, []error{usecase.Error{StatusCode: status.InvalidArgument}}, i.ExpectedErrors())
}
