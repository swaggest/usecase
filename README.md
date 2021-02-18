# Use Case Interactor

[![Build Status](https://github.com/swaggest/usecase/workflows/test-unit/badge.svg)](https://github.com/swaggest/usecase/actions?query=branch%3Amaster+workflow%3Atest-unit)
[![Coverage Status](https://codecov.io/gh/swaggest/usecase/branch/master/graph/badge.svg)](https://codecov.io/gh/swaggest/usecase)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/swaggest/usecase)
![Code lines](https://sloc.xyz/github/swaggest/usecase/?category=code)
![Comments](https://sloc.xyz/github/swaggest/usecase/?category=comments)

This module defines generalized contract of *Use Case Interactor* to enable 
[The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) 
in Go application.

![Clean Architecture](https://blog.cleancoder.com/uncle-bob/images/2012-08-13-the-clean-architecture/CleanArchitecture.jpg)

This abstraction is intended for use with automated transport layer, for example see [`REST`](https://github.com/swaggest/rest). 

## Usage

```go
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
```

Then use configured use case interactor with transport/documentation/etc adapter.

For example with [REST](https://github.com/swaggest/rest/blob/v0.1.18/_examples/basic/main.go#L95-L96) router:
```go
// Add use case handler to router.
r.Method(http.MethodPost, "/double/{param1}", nethttp.NewHandler(u))
```