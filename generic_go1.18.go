//go:build go1.18
// +build go1.18

package usecase

import (
	"context"
	"fmt"
)

// ErrInvalidType is returned on port type assertion error.
const ErrInvalidType = sentinelError("invalid type")

// IOInteractorOf is an IOInteractor with parametrized input/output types.
type IOInteractorOf[i, o any] struct {
	IOInteractor

	InteractFunc func(ctx context.Context, input i, output *o) error
}

// Invoke calls interact function in a type-safe way.
func (ioi IOInteractorOf[i, o]) Invoke(ctx context.Context, input i, output *o) error {
	return ioi.InteractFunc(ctx, input, output)
}

// NewInteractor creates generic use case interactor with input and output ports.
//
// It pre-fills name and title with caller function.
// Input is passed by value, while output is passed by pointer to be mutable.
func NewInteractor[i, o any](interact func(ctx context.Context, input i, output *o) error, options ...func(i *IOInteractor)) IOInteractorOf[i, o] {
	u := IOInteractorOf[i, o]{}
	u.Input = *new(i)
	u.Output = new(o)
	u.InteractFunc = interact
	u.Interactor = Interact(func(ctx context.Context, input, output any) error {
		inp, ok := input.(i)
		if !ok {
			return fmt.Errorf("%w of input: %T, expected: %T", ErrInvalidType, input, u.Input)
		}

		out, ok := output.(*o)
		if !ok {
			return fmt.Errorf("%w f output: %T, expected: %T", ErrInvalidType, output, u.Output)
		}

		return interact(ctx, inp, out)
	})

	u.name, u.title = callerFunc()
	u.name = filterName(u.name)

	for _, o := range options {
		o(&u.IOInteractor)
	}

	return u
}
