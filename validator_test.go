package validator

import (
    "errors"
    "net/url"
    "strconv"
    "testing"
)

func TestEmptyForm(t *testing.T) {
    t.Parallel()
    form := NewForm()
    if form.Validated() != true {
        t.Error("Empty form should validate")
    }

    fields := form.Fields()
    if len(fields) != 0 {
        t.Error("Empty form contained fields")
    }
}


func TestFormAddField(t *testing.T) {
    t.Parallel()
    field1 := NewBaseField("field")

    form := NewForm()
    form.AddField(&field1)
    fields := form.Fields()
    if len(fields) != 1 {
        t.Fatal("Form should have 1 field, has " + strconv.Itoa(len(fields)))
    }

    // Test validating a form with a field with no validators
    if form.Validated() == false {
        t.Error("Form has no validators, should validate")
    }

    field2 := NewBaseField("field2")
    field2.Require()
    field2.SetValues([]string{"a"})
    field2.AddValidator(func(value string) error {
        if value != "a" {
            return errors.New("Value is not 'a'")
        }
        return nil
    })
    form.AddField(&field2)

    if form.Validated() == false {
        t.Error("Form should validate")
    }

    field3 := NewBaseField("field3")
    field3.Require()
    field3.SetValues([]string{"b"})
    field3.AddValidator(func(value string) error {
        if value != "a" {
            return errors.New("Value is not 'a'")
        }
        return nil
    })
    form.AddField(&field3)

    if form.Validated() == true {
        t.Error("Form with failing validator should not validate")
    }
}


func TestFormSubmission(t *testing.T) {
    form := NewForm()
    field := NewBaseField("field")
    field.Required()
    field.AddValidator(func(value string) error {
        if value != "a" {
            return errors.New("Value != 'a'")
        }
        return nil
    })

    form.AddField(&field)

    values := url.Values{}
    values.Set("field", "a")
    form.SetValues(values)
    valid := form.Validated()
    if valid == false {
        t.Error("Form should validate")
    }

    // Test the opposite case, just to be sure the validators are being run
    values.Set("field", "b")
    form.SetValues(values)
    valid = form.Validated()
    if valid == true {
        t.Error("Form shouldn't validate")
    }
}


func TestFieldAddValidator(t *testing.T) {
    t.Parallel()
    field := new(BaseField)
    field.AddValidator(func(value string) error {
        return nil
    })

    validators := field.Validators()
    if len(validators) != 1 {
        t.Error("Should have 1 validator, have " + strconv.Itoa(len(validators)))
    }
}


func TestFieldAddValues(t *testing.T) {
    t.Parallel()
    field := NewBaseField("field")
    values := field.Values()
    if len(values) != 0 {
        t.Fatal("Empty field has values")
    }

    values = []string{"a", "b"}
    field.SetValues(values)
    values = field.Values()
    if len(values) != 2 {
        t.Fatal("Field should have 2 values, has " + strconv.Itoa(len(values)))
    }
}


func TestFieldValidation(t *testing.T) {
    t.Parallel()
    field := NewBaseField("field")
    err := field.Validate()
    if err != nil {
        t.Fatal("Field with no validators should validate")
    }

    field.AddValidator(func(value string) error {
        return errors.New("Dummy validation error")
    })

    err = field.Validate()
    if err != nil {
        t.Fatal("Non-required field should validate")
    }

    field.Require()

    err = field.Validate()
    if err == nil {
        t.Fatal("Dummy failing validator should cause validation to fail")
    }
}

func TestFieldMultiError(t *testing.T) {
    t.Parallel()
    field := NewBaseField("field")
    field.Require()

    err := field.Validate()
    me, ok := err.(MultiError)
    if ok != true {
        t.Fatal("Could not convert validation error to MultiError")
    }

    errors := me.Errors()
    if len(errors) != 1 {
        t.Fatal("MultiError should contain 1 error, contains " + strconv.Itoa(len(errors)))
    }

    if errors[0] == nil {
        t.Error("Error should not be nil")
    }
}


func TestRequiredField(t *testing.T) {
    t.Parallel()
    field := NewBaseField("field")

    required := field.Required()
    if required == true {
        t.Fatal("Field should not be required")
    }

    field.Require()
    required = field.Required()
    if required == false {
        t.Fatal("Field should be required")
    }

    err := field.Validate()
    if err == nil {
        t.Error("Required field with no value should not validate")
    }

    field = NewBaseField("field")  // New field because the old one has errors attached to it
    field.Require()
    values := []string{"a", "b"}
    field.SetValues(values)
    err = field.Validate()
    if err != nil {
        t.Error("Required field with values should validate")
    }
}
