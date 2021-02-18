package usecase_test

import (
	"context"
	"fmt"
	"log"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func ExampleNewIOI() {
	// Configure use case interactor in application layer.
	type myInput struct {
		Param1 int    `path:"param1" description:"Parameter in resource path." multipleOf:"2"`
		Param2 string `json:"param2" description:"Parameter in resource body."`
	}

	type myOutput struct {
		Value1 int    `json:"value1"`
		Value2 string `json:"value2"`
	}

	u := usecase.NewIOI(new(myInput), new(myOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*myInput)
			out = output.(*myOutput)
		)

		if in.Param1%2 != 0 {
			return status.InvalidArgument
		}

		// Do something to set output based on input.
		out.Value1 = in.Param1 + in.Param1
		out.Value2 = in.Param2 + in.Param2

		return nil
	})

	// Additional properties can be configured for purposes of automated documentation.
	u.SetTitle("Doubler")
	u.SetDescription("Doubler doubles parameter values.")
	u.SetTags("transformation")
	u.SetExpectedErrors(status.InvalidArgument)
	u.SetIsDeprecated(true)

	// The code below illustrates transport side.

	// At transport layer, input and out ports are to be examined and populated using reflection.
	// For example request body could be json unmarshaled, or request parameters can be mapped.
	input := new(myInput)
	// input := reflect.New(reflect.TypeOf(u.InputPort()))
	input.Param1 = 1234
	input.Param2 = "abc"

	output := new(myOutput)
	// output := reflect.New(reflect.TypeOf(u.OutputPort()))

	// When input is prepared and output is initialized, transport should invoke interaction.
	err := u.Interact(context.TODO(), input, output)
	if err != nil {
		log.Fatal(err)
	}

	// And make some use of prepared output.
	fmt.Printf("%+v\n", output)

	// Output:
	// &{Value1:2468 Value2:abcabc}
}
