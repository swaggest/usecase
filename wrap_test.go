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
		withInput       usecase.HasInputPort
		withOutput      usecase.HasOutputPort
	)

	f := func(next usecase.Interactor, name string) usecase.Interactor {
		return usecase.Interact(func(ctx context.Context, input, output interface{}) error {
			invocationOrder = append(invocationOrder, name+" start")
			err := next.Interact(ctx, input, output)
			invocationOrder = append(invocationOrder, name+" end")

			return err
		})
	}

	mw1 := usecase.MiddlewareFunc(func(next usecase.Interactor) usecase.Interactor {
		assert.False(t, usecase.As(next, &withInput))
		assert.True(t, usecase.As(next, &withOutput))

		u := struct {
			usecase.HasInputPort
			usecase.Interactor
		}{
			HasInputPort: usecase.WithInput{Input: new(string)},
			Interactor:   f(next, "mw1"),
		}

		return u
	})

	mw2 := usecase.MiddlewareFunc(func(next usecase.Interactor) usecase.Interactor {
		assert.False(t, usecase.As(next, &withInput))
		assert.False(t, usecase.As(next, &withOutput))

		u := struct {
			usecase.HasOutputPort
			usecase.Interactor
		}{
			HasOutputPort: usecase.WithOutput{Output: new(string)},
			Interactor:    f(next, "mw2"),
		}

		return u
	})

	i := usecase.Wrap(usecase.Interact(func(ctx context.Context, input, output interface{}) error {
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
		usecase.Interactor
		usecase.HasInputPort
		usecase.HasOutputPort
	}{
		Interactor: usecase.Interact(func(ctx context.Context, input, output interface{}) error {
			o, ok := output.(*Response)
			assert.True(t, ok)

			o.Name = "Jane"

			return nil
		}),
		HasOutputPort: usecase.WithOutput{
			Output: Response{},
		},
	}

	var (
		withOutput usecase.HasOutputPort
		withInput  usecase.HasInputPort
	)

	assert.True(t, usecase.As(u, &withInput))
	assert.True(t, usecase.As(u, &withOutput))
}

func TestAs_panics(t *testing.T) {
	u := struct {
		usecase.Interactor
		usecase.Info
	}{}

	// target cannot be nil.
	assert.Panics(t, func() {
		usecase.As(u, nil)
	})

	// target must be a non-nil pointer.
	assert.Panics(t, func() {
		usecase.As(u, 123)
	})

	// *target must be interface.
	assert.Panics(t, func() {
		usecase.As(u, &usecase.Info{})
	})
}

func TestErrorCatcher_Wrap(t *testing.T) {
	u := usecase.NewIOI(nil, nil, func(ctx context.Context, input, output interface{}) error {
		return errors.New("failed")
	})

	called := false
	uw := usecase.Wrap(u, usecase.ErrorCatcher(func(err error) {
		called = true
		assert.EqualError(t, err, "failed")
	}))

	assert.EqualError(t, uw.Interact(context.Background(), nil, nil), "failed")
	assert.True(t, called)
}
