package avowed

import (
	"cmp"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
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

type LengthRangeValidator struct {
	Min int
	Max int
}

func (v *LengthRangeValidator) Validate(val string) (ok bool, err error) {
	l := len(val)
	if l < v.Min || l > v.Max {
		return false, fmt.Errorf("length %d is not in range [%d, %d]", l, v.Min, v.Max)
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

type AlphaNumericValidator struct{}

func (v *AlphaNumericValidator) Validate(val string) (ok bool, err error) {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, val)
	if err != nil {
		return false, err
	}
	if !matched {
		return false, fmt.Errorf("value %q is not alphanumeric", val)
	}
	return true, nil
}

type MACAddressValidator struct{}

func (v *MACAddressValidator) Validate(val string) (ok bool, err error) {
	_, err = net.ParseMAC(val)
	if err != nil {
		return false, fmt.Errorf("invalid MAC address %q: %v", val, err)
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

type IPv4Validator struct{}

func (v *IPv4Validator) Validate(val string) (ok bool, err error) {
	ip := net.ParseIP(val)
	if ip == nil || ip.To4() == nil {
		return false, fmt.Errorf("invalid IPv4 address %q", val)
	}
	return true, nil
}

type IPv6Validator struct{}

func (v *IPv6Validator) Validate(val string) (ok bool, err error) {
	ip := net.ParseIP(val)
	if ip == nil || ip.To4() != nil {
		return false, fmt.Errorf("invalid IPv6 address %q", val)
	}
	return true, nil
}

type XMLValidator struct{}

func (v *XMLValidator) Validate(val string) (ok bool, err error) {
	decoder := xml.NewDecoder(strings.NewReader(val))
	var hasElement bool

	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return false, fmt.Errorf("XML parsing error: %w", err)
		}

		if _, ok := tok.(xml.StartElement); ok { // atleast one tag
			hasElement = true
		}
	}

	if !hasElement {
		return false, fmt.Errorf("XML document must contain at least one element")
	}

	return true, nil
}

type JSONValidator struct{}

func (v *JSONValidator) Validate(val string) (ok bool, err error) {
	if !json.Valid([]byte(val)) {
		return false, fmt.Errorf("invalid JSON")
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
