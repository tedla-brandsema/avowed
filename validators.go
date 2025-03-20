package valex

import (
	"cmp"
	"encoding/json"
	"encoding/xml"
	"errors"
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
	Min int `param:"min"`
	Max int `param:"max"`
}

func (v *IntRangeValidator) Validate(val int) (ok bool, err error) {
	if val < v.Min || val > v.Max {
		return false, fmt.Errorf("value %d is out of range [%d, %d]", val, v.Min, v.Max)
	}
	return true, nil
}

func (v *IntRangeValidator) Name() string {
	return "range"
}

func (v *IntRangeValidator) Handle(val int) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type NonNegativeIntValidator struct{}

func (v *NonNegativeIntValidator) Validate(val int) (ok bool, err error) {
	if val < 0 {
		return false, fmt.Errorf("value %d is a negative integer", val)
	}
	return true, nil
}

func (v *NonNegativeIntValidator) Name() string {
	return "pos"
}

func (v *NonNegativeIntValidator) Handle(val int) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type NonPositiveIntValidator struct{}

func (v *NonPositiveIntValidator) Validate(val int) (ok bool, err error) {
	if val > 0 {
		return false, fmt.Errorf("value %d is a positive integer", val)
	}
	return true, nil
}

func (v *NonPositiveIntValidator) Name() string {
	return "neg"
}

func (v *NonPositiveIntValidator) Handle(val int) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type UrlValidator struct{}

func (v *UrlValidator) Validate(val string) (ok bool, err error) {
	_, err = url.ParseRequestURI(val)
	if err == nil {
		ok = true
	}
	return
}

func (v *UrlValidator) Name() string {
	return "url"
}

func (v *UrlValidator) Handle(val string) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type EmailValidator struct{}

func (v *EmailValidator) Validate(val string) (ok bool, err error) {
	_, err = mail.ParseAddress(val)
	if err == nil {
		ok = true
	}
	return
}

func (v *EmailValidator) Name() string {
	return "email"
}

func (v *EmailValidator) Handle(val string) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type NonEmptyStringValidator struct{}

func (v *NonEmptyStringValidator) Validate(val string) (ok bool, err error) {
	if val == "" {
		return false, fmt.Errorf("string is empty")
	}
	return true, nil
}

func (v *NonEmptyStringValidator) Name() string {
	return "!empty"
}

func (v *NonEmptyStringValidator) Handle(val string) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type MinLengthValidator struct {
	Size int `param:"size"`
}

func (v *MinLengthValidator) Validate(val string) (ok bool, err error) {
	if v.Size == 0 {
		return false, errors.New(`value of parameter "size" cannot be 0`)
	}
	if len(val) < v.Size {
		return false, fmt.Errorf("value %s exeeds minimum length %d", val, v.Size)
	}
	return true, nil
}

func (v *MinLengthValidator) Name() string {
	return "min"
}

func (v *MinLengthValidator) Handle(val string) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type MaxLengthValidator struct {
	Size int `param:"size"`
}

func (v *MaxLengthValidator) Validate(val string) (ok bool, err error) {
	if v.Size == 0 {
		return false, errors.New(`value of parameter "size" cannot be 0`)
	}
	if len(val) > v.Size {
		return false, fmt.Errorf("value %s exeeds maximum length %d", val, v.Size)
	}
	return true, nil
}

func (v *MaxLengthValidator) Name() string {
	return "max"
}

func (v *MaxLengthValidator) Handle(val string) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type LengthRangeValidator struct {
	Min int `param:"min"`
	Max int `param:"max"`
}

func (v *LengthRangeValidator) Validate(val string) (ok bool, err error) {
	l := len(val)
	if v.Min == 0 {
		return false, errors.New(`"min" value cannot be 0`)
	}
	if v.Max == 0 {
		return false, errors.New(`"max" value cannot be 0`)
	}
	if l < v.Min || l > v.Max {
		return false, fmt.Errorf("value %q with length %d is not in range [%d, %d]", val, l, v.Min, v.Max)
	}
	return true, nil
}

func (v *LengthRangeValidator) Name() string {
	return "len"
}

func (v *LengthRangeValidator) Handle(val string) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type RegexValidator struct {
	Pattern *regexp.Regexp
}

func (v *RegexValidator) Validate(val string) (ok bool, err error) {
	if !v.Pattern.MatchString(val) {
		return false, fmt.Errorf("value %q does not match pattern %q", val, v.Pattern.String())
	}
	return true, nil
}

// TODO: implement Directive for RegexValidator

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

func (v *AlphaNumericValidator) Name() string {
	return "alphanum"
}

func (v *AlphaNumericValidator) Handle(val string) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type MACAddressValidator struct{}

func (v *MACAddressValidator) Validate(val string) (ok bool, err error) {
	_, err = net.ParseMAC(val)
	if err != nil {
		return false, fmt.Errorf("invalid MAC address %q: %v", val, err)
	}
	return true, nil
}

func (v *MACAddressValidator) Name() string {
	return "mac"
}

func (v *MACAddressValidator) Handle(val string) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type IpValidator struct{}

func (v *IpValidator) Validate(val string) (ok bool, err error) {
	if ip := net.ParseIP(val); ip == nil {
		return false, fmt.Errorf("invalid IP address %q", val)
	}
	return true, nil
}

func (v *IpValidator) Name() string {
	return "ip"
}

func (v *IpValidator) Handle(val string) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type IPv4Validator struct{}

func (v *IPv4Validator) Validate(val string) (ok bool, err error) {
	ip := net.ParseIP(val)
	if ip == nil || ip.To4() == nil {
		return false, fmt.Errorf("invalid IPv4 address %q", val)
	}
	return true, nil
}

func (v *IPv4Validator) Name() string {
	return "ipv4"
}

func (v *IPv4Validator) Handle(val string) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type IPv6Validator struct{}

func (v *IPv6Validator) Validate(val string) (ok bool, err error) {
	ip := net.ParseIP(val)
	if ip == nil || ip.To4() != nil {
		return false, fmt.Errorf("invalid IPv6 address %q", val)
	}
	return true, nil
}

func (v *IPv6Validator) Name() string {
	return "ipv6"
}

func (v *IPv6Validator) Handle(val string) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
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

func (v *XMLValidator) Name() string {
	return "xml"
}

func (v *XMLValidator) Handle(val string) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
}

type JSONValidator struct{}

func (v *JSONValidator) Validate(val string) (ok bool, err error) {
	if !json.Valid([]byte(val)) {
		return false, fmt.Errorf("invalid JSON")
	}
	return true, nil
}

func (v *JSONValidator) Name() string {
	return "json"
}

func (v *JSONValidator) Handle(val string) error {
	if ok, err := v.Validate(val); !ok {
		return err
	}
	return nil
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
