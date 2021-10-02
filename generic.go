//go:build go1.18
// +build go1.18

package usecase

import (
	"context"
	"fmt"
)

type IOInteractorOf[i, o interface{}] struct {
	IOInteractor

	InteractFunc func(ctx context.Context, input i, output *o) error
}

func (ioi IOInteractorOf[i, o]) Invoke(ctx context.Context, input i, output *o) error {
	return ioi.InteractFunc(ctx, input, output)
}

// NewInteractor creates generic use case interactor with input and output ports.
//
// It pre-fills name and title with caller function.
// Input is passed by value, while output is passed by pointer to be mutable.
func NewInteractor[i, o interface{}](interact func(ctx context.Context, input i, output *o) error) IOInteractorOf[i, o] {
	u := IOInteractorOf[i, o]{}
	u.Input = *new(i)
	u.Output = new(o)
	u.InteractFunc = interact
	u.Interactor = Interact(func(ctx context.Context, input, output interface{}) error {
		inp, ok := input.(i)
		if !ok {
			return fmt.Errorf("invalid input type received: %T, expected: %T", input, u.Input)
		}

		out, ok := output.(*o)
		if !ok {
			return fmt.Errorf("invalid output type received: %T, expected: %T", output, u.Output)
		}

		return interact(ctx, inp, out)
	})

	u.name, u.title = callerFunc()
	u.name = filterName(u.name)

	return u
}
