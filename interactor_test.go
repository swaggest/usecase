package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

//func TestInteract_Interact(t *testing.T) {
//	err := errors.New("failed")
//	inp := new(string)
//	out := new(int)
//
//	i := usecase.Interact(func(ctx context.Context, input, output interface{}) error {
//		assert.Equal(t, inp, input)
//		assert.Equal(t, out, output)
//
//		return err
//	})
//	assert.Equal(t, err, i.Interact(context.Background(), inp, out))
//}

func TestInfo(t *testing.T) {
	i := usecase.Info{}
	i.SetName("name")
	i.SetDescription("Description")
	i.SetTitle("Title")
	i.SetTags("tag1", "tag2")
	i.SetIsDeprecated(true)
	i.SetExpectedErrors(usecase.Error{StatusCode: status.InvalidArgument})

	assert.Equal(t, "name", i.Name())
	assert.Equal(t, "Description", i.Description())
	assert.Equal(t, "Title", i.Title())
	assert.Equal(t, []string{"tag1", "tag2"}, i.Tags())
	assert.Equal(t, true, i.IsDeprecated())
	assert.Equal(t, []error{usecase.Error{StatusCode: status.InvalidArgument}}, i.ExpectedErrors())
}

type Foo struct{}

func (f *Foo) Bar() usecase.IOInteractor[interface{}, interface{}] {
	return usecase.NewIOI(nil, nil, func(ctx context.Context, input, output interface{}) error {
		return nil
	})
}

func (f Foo) Baz() usecase.IOInteractor[interface{}, interface{}] {
	return usecase.NewIOI[interface{}, interface{}](nil, nil, func(ctx context.Context, input, output interface{}) error {
		return nil
	})
}

func TestNewIOI(t *testing.T) {
	u1 := usecase.NewIOI[*string, *int](new(string), new(int), func(ctx context.Context, input *string, output *int) error {
		return nil
	})

	assert.Equal(t, "swaggest/usecase_test.TestNewIOI", u1.Name())
	assert.Equal(t, "Test New IOI", u1.Title())

	u2 := (&Foo{}).Bar()
	assert.Equal(t, "swaggest/usecase_test.(*Foo).Bar", u2.Name())
	assert.Equal(t, "Foo Bar", u2.Title())

	u3 := Foo{}.Baz()
	assert.Equal(t, "swaggest/usecase_test.Foo.Baz", u3.Name())
	assert.Equal(t, "Foo Baz", u3.Title())
}
