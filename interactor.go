package usecase

import (
	"context"
	"path"
	"runtime"
	"strings"
)

// Interactor orchestrates the flow of data to and from the entities,
// and direct those entities to use their enterprise
// wide business rules to achieve the goals of the use case.
type Interactor interface {
	// Interact sets output port value with regards to input port value or fails.
	Interact(ctx context.Context, input, output interface{}) error
}

// Interact makes use case interactor from function.
type Interact func(ctx context.Context, input, output interface{}) error

// Interact implements Interactor.
func (i Interact) Interact(ctx context.Context, input, output interface{}) error {
	return i(ctx, input, output)
}

// HasInputPort declares input port.
type HasInputPort interface {
	// InputPort returns sample of input value, e.g. new(MyInput).
	InputPort() interface{}
}

// WithInput is an embeddable implementation of HasInputPort.
type WithInput struct {
	Input interface{}
}

// InputPort implements HasInputPort.
func (wi WithInput) InputPort() interface{} {
	return wi.Input
}

// HasOutputPort declares output port.
type HasOutputPort interface {
	// OutputPort returns sample of output value, e.g. new(MyOutput).
	OutputPort() interface{}
}

// WithOutput is an embeddable implementation of HasOutputPort.
type WithOutput struct {
	Output interface{}
}

// OutputPort implements HasOutputPort.
func (wi WithOutput) OutputPort() interface{} {
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
type IOInteractor struct {
	Interactor
	Info
	WithInput
	WithOutput
}

// NewIOI creates use case interactor with input, output and interact action function.
//
// It pre-fills name and title with caller function.
func NewIOI(input, output interface{}, interact Interact) IOInteractor {
	u := IOInteractor{}
	u.Input = input
	u.Output = output
	u.Interactor = interact

	u.name, u.title = callerFunc()

	return u
}

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
		title = parts[1]
	}

	return pathName, title
}
