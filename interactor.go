package usecase

import (
	"context"
	"path"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Interactor orchestrates the flow of data to and from the entities,
// and direct those entities to use their enterprise
// wide business rules to achieve the goals of the use case.
type Interactor[i interface{}, o interface{}] interface {
	// Interact sets output port value with regards to input port value or fails.
	Interact(ctx context.Context, input i, output o) error
}

// Interact makes use case interactor from function.
type Interact[i interface{}, o interface{}] func(ctx context.Context, input i, output o) error

// Interact implements Interactor.
func (f Interact[i, o]) Interact(ctx context.Context, input i, output o) error {
	return f(ctx, input, output)
}

// HasInputPort declares input port.
type HasInputPort[i interface{}] interface {
	// InputPort returns sample of input value, e.g. new(MyInput).
	InputPort() i
}

// WithInput is an embeddable implementation of HasInputPort.
type WithInput[i interface{}] struct {
	Input i
}

// InputPort implements HasInputPort.
func (wi WithInput[i]) InputPort() i {
	return wi.Input
}

// HasOutputPort declares output port.
type HasOutputPort[o interface{}] interface {
	// OutputPort returns sample of output value, e.g. new(MyOutput).
	OutputPort() o
}

// WithOutput is an embeddable implementation of HasOutputPort.
type WithOutput[o interface{}] struct {
	Output o
}

// OutputPort implements HasOutputPort.
func (wi WithOutput[o]) OutputPort() o {
	return wi.Output
}

// HasTitle declares title.
type HasTitle interface {
	Title() string
}

// HasName declares title.
type HasName interface {
	Name() string
}

// HasDescription declares description.
type HasDescription interface {
	Description() string
}

// HasTags declares tags of use cases group.
type HasTags interface {
	Tags() []string
}

// HasExpectedErrors declares errors that are expected to cause use case failure.
type HasExpectedErrors interface {
	ExpectedErrors() []error
}

// HasIsDeprecated declares status of deprecation.
type HasIsDeprecated interface {
	IsDeprecated() bool
}

// Info exposes information about use case.
type Info struct {
	name           string
	title          string
	description    string
	tags           []string
	expectedErrors []error
	isDeprecated   bool
}

var (
	_ HasTags           = Info{}
	_ HasTitle          = Info{}
	_ HasName           = Info{}
	_ HasDescription    = Info{}
	_ HasIsDeprecated   = Info{}
	_ HasExpectedErrors = Info{}
)

// IsDeprecated implements HasIsDeprecated.
func (i Info) IsDeprecated() bool {
	return i.isDeprecated
}

// SetIsDeprecated sets status of deprecation.
func (i *Info) SetIsDeprecated(isDeprecated bool) {
	i.isDeprecated = isDeprecated
}

// ExpectedErrors implements HasExpectedErrors.
func (i Info) ExpectedErrors() []error {
	return i.expectedErrors
}

// SetExpectedErrors sets errors that are expected to cause use case failure.
func (i *Info) SetExpectedErrors(expectedErrors ...error) {
	i.expectedErrors = expectedErrors
}

// Tags implements HasTag.
func (i Info) Tags() []string {
	return i.tags
}

// SetTags sets tags of use cases group.
func (i *Info) SetTags(tags ...string) {
	i.tags = tags
}

// Description implements HasDescription.
func (i Info) Description() string {
	return i.description
}

// SetDescription sets use case description.
func (i *Info) SetDescription(description string) {
	i.description = description
}

// Title implements HasTitle.
func (i Info) Title() string {
	return i.title
}

// SetTitle sets use case title.
func (i *Info) SetTitle(title string) {
	i.title = title
}

// Name implements HasName.
func (i Info) Name() string {
	return i.name
}

// SetName sets use case title.
func (i *Info) SetName(name string) {
	i.name = name
}

// IOInteractor is an interactor with input and output.
type IOInteractor[i interface{}, o interface{}] struct {
	Interactor[i, o]
	Info
	WithInput[i]
	WithOutput[o]
}

// NewIOI creates use case interactor with input, output and interact action function.
//
// It pre-fills name and title with caller function.
func NewIOI[i interface{}, o interface{}](input i, output o, interact Interact[i, o]) IOInteractor[i, o] {
	u := IOInteractor[i, o]{}
	u.Input = input
	u.Output = output
	u.Interactor = interact

	u.name, u.title = callerFunc()

	u.name = strings.TrimPrefix(u.name, "internal/")
	u.name = strings.TrimPrefix(u.name, "usecase.")
	u.name = strings.TrimPrefix(u.name, "./main.")

	return u
}

var titleReplacer = strings.NewReplacer(
	"(", "",
	".", "",
	"*", "",
	")", "",
)

// callerFunc returns trimmed path and name of parent function.
func callerFunc() (string, string) {
	skipFrames := 2

	pc, _, _, ok := runtime.Caller(skipFrames)
	if !ok {
		return "", ""
	}

	f := runtime.FuncForPC(pc)

	pathName := path.Base(path.Dir(f.Name())) + "/" + path.Base(f.Name())
	title := path.Base(f.Name())

	parts := strings.SplitN(title, ".", 2)
	if len(parts) != 1 {
		title = parts[len(parts)-1]
		title = titleReplacer.Replace(title)
		title = splitCamelcase(title)
	}

	return pathName, title
}

// borrowed from https://pkg.go.dev/github.com/fatih/camelcase#Split to avoid external dependency.
func splitCamelcase(src string) string { // nolint:cyclop
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return src
	}

	var (
		entries []string
		runes   [][]rune
	)

	var class, lastClass int

	// split into fields based on class of unicode character
	for _, r := range src {
		switch {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}

		if class == lastClass && runes != nil {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}

		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}

	return strings.Join(entries, " ")
}
