package validator

import (
    "regexp"
    "strconv"
)

func IsBool(value string) (err error) {
    _, err = strconv.ParseBool(value)
    if err != nil {
        return err
    }

    return
}


func IsInt(value string) (err error) {
    _, err = strconv.Atoi(value)
    if err != nil {
        return err
    }

    return
}
func IsNonzeroInt(value string) (err error) {
    i, err := strconv.Atoi(value)
    if err != nil {
        return err
    }

    if i == 0 {
        return err
    }

    return
}


func IsFloat64(value string) (err error) {
    _, err = strconv.ParseFloat(value, 64)
    if err != nil {
        return err
    }

    return
}
func IsNonzeroFloat64(value string) (err error) {
    i, err := strconv.ParseFloat(value, 64)
    if err != nil {
        return err
    }

    if i == 0.0 {
        return err
    }

    return
}


func IsNonzeroString(value string) (err error) {
    if value == "" {
        return err
    }

    return
}


func IsAlphanumeric(value string) (err error) {
    re := regexp.MustCompile("^[A-Za-z0-9]*$")
    if re.MatchString(value) == false {
        return err
    }

    return nil
}
