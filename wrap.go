package usecase

import (
	"context"
	"reflect"
)

// Middleware creates decorated use case interactor.
type Middleware interface {
	Wrap(interactor Interactor) Interactor
}

// ErrorCatcher is a use case middleware that collects non empty errors.
type ErrorCatcher func(ctx context.Context, input interface{}, err error)

// Wrap implements Middleware.
func (e ErrorCatcher) Wrap(u Interactor) Interactor {
	return &wrappedInteractor{
		Interactor: Interact(func(ctx context.Context, input, output interface{}) error {
			err := u.Interact(ctx, input, output)
			if err != nil {
				e(ctx, input, err)
			}

			return err
		}),
		wrapped: u,
	}
}

// MiddlewareFunc makes Middleware from function.
type MiddlewareFunc func(next Interactor) Interactor

// Wrap decorates use case interactor.
func (mwf MiddlewareFunc) Wrap(interactor Interactor) Interactor {
	return mwf(interactor)
}

// Wrap decorates Interactor with Middlewares.
//
// Having arguments i, mw1, mw2 the order of invocation is: mw1, mw2, i, mw2, mw1.
// Middleware mw1 can find behaviors of mw2 with As, but not vice versa.
func Wrap(interactor Interactor, mw ...Middleware) Interactor {
	for i := len(mw) - 1; i >= 0; i-- {
		w := mw[i].Wrap(interactor)
		if w != nil {
			interactor = &wrappedInteractor{
				Interactor: w,
				wrapped:    interactor,
			}
		}
	}

	return interactor
}

// As finds the first Interactor in Interactor's chain that matches target, and if so, sets
// target to that Interactor value and returns true.
//
// An Interactor matches target if the Interactor's concrete value is assignable to the value
// pointed to by target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// Interactor, or to any interface type.
func As(interactor Interactor, target interface{}) bool {
	if interactor == nil {
		return false
	}

	if target == nil {
		panic("target cannot be nil")
	}

	val := reflect.ValueOf(target)
	typ := val.Type()

	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("target must be a non-nil pointer")
	}

	if e := typ.Elem(); e.Kind() != reflect.Interface {
		panic("*target must be interface")
	}

	targetType := typ.Elem()

	for {
		wrap, isWrap := interactor.(*wrappedInteractor)

		if isWrap {
			interactor = wrap.Interactor
		}

		if reflect.TypeOf(interactor).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(interactor))

			return true
		}

		if !isWrap {
			break
		}

		interactor = wrap.wrapped
	}

	return false
}

type wrappedInteractor struct {
	Interactor
	wrapped Interactor
}
