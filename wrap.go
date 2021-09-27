package usecase

import (
	"context"
	"reflect"
)

// Middleware creates decorated use case interactor.
type Middleware[i interface{}, o interface{}] interface {
	Wrap(interactor Interactor[i, o]) Interactor[i, o]
}

// ErrorCatcher is a use case middleware that collects non empty errors.
type ErrorCatcher[i interface{}, o interface{}] func(ctx context.Context, input i, err error)

// Wrap implements Middleware.
func (e ErrorCatcher[i, o]) Wrap(u Interactor[i, o]) Interactor[i, o] {
	return &wrappedInteractor[i, o]{
		Interactor: Interact[i, o](func(ctx context.Context, input i, output o) error {
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
type MiddlewareFunc[i interface{}, o interface{}] func(next Interactor[i, o]) Interactor[i, o]

// Wrap decorates use case interactor.
func (mwf MiddlewareFunc[i, o]) Wrap(interactor Interactor[i, o]) Interactor[i, o] {
	return mwf(interactor)
}

// Wrap decorates Interactor with Middlewares.
//
// Having arguments i, mw1, mw2 the order of invocation is: mw1, mw2, i, mw2, mw1.
// Middleware mw1 can find behaviors of mw2 with As, but not vice versa.
func Wrap[i interface{}, o interface{}](interactor Interactor[i, o], mw ...Middleware[i, o]) Interactor[i, o] {
	for j := len(mw) - 1; j >= 0; j-- {
		w := mw[j].Wrap(interactor)
		if w != nil {
			interactor = &wrappedInteractor[i, o]{
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
func As[i interface{}, o interface{}](interactor Interactor[i, o], target interface{}) bool {
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
		wrap, isWrap := interactor.(*wrappedInteractor[i, o])

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

type wrappedInteractor[i interface{}, o interface{}] struct {
	Interactor[i, o]
	wrapped Interactor[i, o]
}
