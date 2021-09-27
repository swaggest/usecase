package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/usecase"
)

func TestWrap(t *testing.T) {
	var (
		invocationOrder []string
		withInput       usecase.HasInputPort[*string]
		withOutput      usecase.HasOutputPort[*string]
	)

	f := func(next usecase.Interactor[*string, *string], name string) usecase.Interactor[*string, *string] {
		return usecase.Interact[*string, *string](func(ctx context.Context, input *string, output *string) error {
			invocationOrder = append(invocationOrder, name+" start")
			err := next.Interact(ctx, input, output)
			invocationOrder = append(invocationOrder, name+" end")

			return err
		})
	}

	mw1 := usecase.MiddlewareFunc[*string, *string](func(next usecase.Interactor[*string, *string]) usecase.Interactor[*string, *string] {
		assert.False(t, usecase.As(next, &withInput))
		assert.True(t, usecase.As(next, &withOutput))

		u := struct {
			usecase.HasInputPort[*string]
			usecase.Interactor[*string, *string]
		}{
			HasInputPort: usecase.WithInput[*string]{Input: new(string)},
			Interactor:   f(next, "mw1"),
		}

		return u
	})

	mw2 := usecase.MiddlewareFunc[*string, *string](func(next usecase.Interactor[*string, *string]) usecase.Interactor[*string, *string] {
		assert.False(t, usecase.As(next, &withInput))
		assert.False(t, usecase.As(next, &withOutput))

		u := struct {
			usecase.HasOutputPort[*string]
			usecase.Interactor[*string, *string]
		}{
			HasOutputPort: usecase.WithOutput[*string]{Output: new(string)},
			Interactor:    f(next, "mw2"),
		}

		return u
	})

	i := usecase.Wrap[*string, *string](usecase.Interact[*string, *string](func(ctx context.Context, input *string, output *string) error {
		invocationOrder = append(invocationOrder, "interaction")

		return nil
	}), mw1, mw2)
	err := i.Interact(context.Background(), nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, []string{
		"mw1 start", "mw2 start", "interaction", "mw2 end", "mw1 end",
	}, invocationOrder)
	assert.True(t, usecase.As(i, &withInput))
	assert.True(t, usecase.As(i, &withOutput))
}

func TestAs(t *testing.T) {
	type Response struct {
		Name string `json:"name"`
	}

	u := struct {
		usecase.Interactor[interface{}, *Response]
		usecase.HasInputPort[interface{}]
		usecase.HasOutputPort[*Response]
	}{
		Interactor: usecase.Interact[interface{}, *Response](func(ctx context.Context, input interface{}, output *Response) error {
			o := output

			o.Name = "Jane"

			return nil
		}),
		HasOutputPort: usecase.WithOutput[*Response]{
			Output: &Response{},
		},
	}

	var (
		withOutput usecase.HasOutputPort[*Response]
		withInput  usecase.HasInputPort[interface{}]
	)

	assert.True(t, usecase.As[interface{}, *Response](u, &withInput))
	assert.True(t, usecase.As[interface{}, *Response](u, &withOutput))
}

func TestAs_panics(t *testing.T) {
	u := struct {
		usecase.Interactor[interface{}, interface{}]
		usecase.Info
	}{}

	// target cannot be nil.
	assert.Panics(t, func() {
		usecase.As[interface{}, interface{}](u, nil)
	})

	// target must be a non-nil pointer.
	assert.Panics(t, func() {
		usecase.As[interface{}, interface{}](u, 123)
	})

	// *target must be interface.
	assert.Panics(t, func() {
		usecase.As[interface{}, interface{}](u, &usecase.Info{})
	})
}

func TestErrorCatcher_Wrap(t *testing.T) {
	u := usecase.NewIOI(nil, nil, func(ctx context.Context, input, output interface{}) error {
		return errors.New("failed")
	})

	called := false
	uw := usecase.Wrap[interface{}, interface{}](u, usecase.ErrorCatcher[interface{}, interface{}](func(ctx context.Context, input interface{}, err error) {
		called = true
		assert.EqualError(t, err, "failed")
	}))

	assert.EqualError(t, uw.Interact(context.Background(), nil, nil), "failed")
	assert.True(t, called)
}
