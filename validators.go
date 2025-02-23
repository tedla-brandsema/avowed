package avowed

import (
	"cmp"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
)

type CmpRangeValidator[T cmp.Ordered] struct {
	Min T
	Max T
}

func (v *CmpRangeValidator[T]) Validate(val T) (ok bool, err error) {
	if cmp.Less(val, v.Min) || cmp.Less(v.Max, val) {
		return false, fmt.Errorf("value %v is out of range [%v, %v]", val, v.Min, v.Max)
	}
	return true, nil
}

type IntRangeValidator struct {
	Min int
	Max int
}

func (v *IntRangeValidator) Validate(val int) (ok bool, err error) {
	if val < v.Min || val > v.Max {
		return false, fmt.Errorf("value %d is out of range [%d, %d]", val, v.Min, v.Max)
	}
	return true, nil
}

type NonNegativeIntValidator struct{}

func (v *NonNegativeIntValidator) Validate(val int) (ok bool, err error) {
	if val < 0 {
		return false, fmt.Errorf("value %d is a negative integer", val)
	}
	return true, nil
}

type NonPositiveIntValidator struct{}

func (v *NonPositiveIntValidator) Validate(val int) (ok bool, err error) {
	if val > 0 {
		return false, fmt.Errorf("value %d is a positive integer", val)
	}
	return true, nil
}

type UrlValidator struct{}

func (v *UrlValidator) Validate(val string) (ok bool, err error) {
	_, err = url.ParseRequestURI(val)
	if err == nil {
		ok = true
	}
	return
}

type EmailValidator struct{}

func (v *EmailValidator) Validate(val string) (ok bool, err error) {
	_, err = mail.ParseAddress(val)
	if err == nil {
		ok = true
	}
	return
}

type NonEmptyStringValidator struct{}

func (v *NonEmptyStringValidator) Validate(val string) (ok bool, err error) {
	if val == "" {
		return false, fmt.Errorf("string is empty")
	}
	return true, nil
}

type MinLengthValidator struct {
	Min int
}

func (v *MinLengthValidator) Validate(val string) (ok bool, err error) {
	if len(val) < v.Min {
		return false, fmt.Errorf("value %s exeeds minimum length %d", val, v.Min)
	}
	return true, nil
}

type MaxLengthValidator struct {
	Max int
}

func (v *MaxLengthValidator) Validate(val string) (ok bool, err error) {
	if len(val) > v.Max {
		return false, fmt.Errorf("value %s exeeds maximum length %d", val, v.Max)
	}
	return true, nil
}

type RegexValidator struct {
	Pattern *regexp.Regexp
}

func (v *RegexValidator) Validate(val string) (ok bool, err error) {
	if !v.Pattern.MatchString(val) {
		return false, fmt.Errorf("value %s does not match pattern %s", val, v.Pattern.String())
	}
	return true, nil
}

type IpValidator struct{}

func (v *IpValidator) Validate(val string) (ok bool, err error) {
	if ip := net.ParseIP(val); ip == nil {
		return false, fmt.Errorf("invalid IP address %q", val)
	}
	return true, nil
}

type CompositeValidator[T cmp.Ordered] struct {
	Validators []Validator[T]
}

func (cv *CompositeValidator[T]) Validate(val T) (ok bool, err error) {
	for _, validator := range cv.Validators {
		if ok, err = validator.Validate(val); !ok {
			return false, err
		}
	}
	return true, nil
}
