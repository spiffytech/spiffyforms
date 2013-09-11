package validator

import (
    "errors"
    "strconv"
)

type Form struct {
    fields map[string]Field
}

func (form *Form) AddField(field Field) {
    form.fields[field.Name()] = field
}
func NewForm() Form {
    form := Form{}
    form.fields = make(map[string]Field)
    return form
}
func (form *Form) Fields() map[string]Field {
    return form.fields
}

func (form *Form) SetValues(submission map[string][]string) {
    for name, field := range form.fields {
        values, ok := submission[name]
        if ok == false {
            continue
        }
        field.SetValues(values)
    }
}

type MultiError struct {
    errors []error
}
func NewMultiError(errors []error) (me MultiError) {
    me.errors = errors
    return me
}
func (me MultiError) Error() string {
    if len(me.errors) == 0 {
        return ""
    }

    ret := me.errors[0].Error()
    if len(me.errors) > 1 {
        ret += " (and " + strconv.Itoa(len(me.errors)-1) + "more"
    }

    return ret
}
func (me *MultiError) Errors() []error {
    return me.errors
}

// Signature for a validator. Accepts the HTML field's value, and returns an error if the field does not validate.
type ValidationFunc func(string) error

type Field interface {
    AddValidator(ValidationFunc)
    Name() string
    Validators() []ValidationFunc
    Validate() error
    AddError(error)
    Error() string
    Errors() error
    Values() []string
    SetValues([]string)
    Require()
    Required() bool
}

type HTMLField interface {
    ErrorClass() string
    ErrorHTML() string
}


// Returns a new BaseField
func NewBaseField(name string) BaseField {
    field := BaseField{}
    field.name = name
    return field
}

// Implements core HTML field logic. May be extended for specialized functionality.
type BaseField struct {
    name string
    validators []ValidationFunc
    errors []error
    values []string
    required bool
}


// Returns the field's name
func (field *BaseField) Name() string {
    return field.name
}


// Add a validator to the field
func (field *BaseField) AddValidator(f ValidationFunc) {
    field.validators = append(field.validators, f)
}


// Used to indicate the field didn't validate
func (field *BaseField) AddError(err error) {
    field.errors = append(field.errors, err)
}


// Returns the field's values
func (field *BaseField) Values() []string {
    return field.values
}


// Sets the submitted values for the field.
func (field *BaseField) SetValues(values []string) {
    field.values = values
}
func (field *BaseField) Error() string {
    return NewMultiError(field.errors).Error()
}
func (field *BaseField) Errors() error {
    if len(field.errors) == 0 {
        return nil
    }
    return NewMultiError(field.errors)
}


// Mark the field as required
func (field *BaseField) Require() {
    field.required = true
}


// Indicates whether the field is required to be filled in
func (field *BaseField) Required() bool {
    return field.required
}


// Returns all validators associated with the field
func (field *BaseField) Validators() []ValidationFunc {
    return field.validators
}


// Validate an individual form field
func (field *BaseField) Validate() error {
    if field.Required() && len(field.Values()) == 0 {
        field.AddError(errors.New("Field is required"))
    }

    for _, validator := range field.Validators() {
        for _, value := range field.Values() {
            err := validator(value)
            if err != nil {
                field.AddError(err)
            }
        }
    }

    return field.Errors()
}


// Validates a form's fields. Returns true if all fields pass all validators.
func (form *Form) Validate() (ok bool) {
    ok = true
    if len(form.fields) == 0 {
        return
    }

    for _, field := range form.fields {
        if field.Validate() != nil {
            ok = false
        }
    }

    return
}
