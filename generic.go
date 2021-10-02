//go:build go1.18

package usecase

import (
	"context"
	"fmt"
)

// NewInteractor creates generic use case interactor with input and output ports.
//
// It pre-fills name and title with caller function.
func NewInteractor[i, o interface{}](interact func(ctx context.Context, input *i, output *o) error) IOInteractor {
	u := IOInteractor{}
	u.Input = new(i)
	u.Output = new(o)
	u.Interactor = Interact(func(ctx context.Context, input, output interface{}) error {
		inp, ok := input.(*i)
		if !ok {
			return fmt.Errorf("invalid input type received: %T, %T expected", input, u.Input)
		}

		out, ok := output.(*o)
		if !ok {
			return fmt.Errorf("invalid output type received: %T, %T expected", output, u.Output)
		}

		return interact(ctx, inp, out)
	})

	u.name, u.title = callerFunc()
	u.name = filterName(u.name)

	return u
}
