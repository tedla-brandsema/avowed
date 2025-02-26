package avowed

import (
	"cmp"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const tagID = "val"

func ValidateStruct(data interface{}) (ok bool, err error) {
	val := reflect.ValueOf(data)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if tag, ok := field.Tag.Lookup(tagID); ok {
			i := val.FieldByName(field.Name).Interface()
			switch v := i.(type) {
			case string:
				if ok, err := stringValidators(tag, field.Name, v); !ok {
					return false, err
				}
			case int:
				if ok, err := intValidators(tag, field.Name, v); !ok {
					return false, err
				}
			}
		}

	}
	return true, nil
}
func intValidators(tag string, name string, value int) (ok bool, err error) {
	vals, err := vals(tag, name)
	if err != nil {
		return false, err
	}
	var v Validator[int]
	switch id := strings.TrimSpace(vals[0]); id {
	case "range":
		min, max, err := rangeFinder(vals[1:])
		if err != nil {
			return false, err
		}
		v = &IntRangeValidator{
			Min: min,
			Max: max,
		}
	case "pos":
		v = &NonNegativeIntValidator{}
	case "neg":
		v = &NonPositiveIntValidator{}
	default:
		return false, fmt.Errorf("unknown validator %q  for field %q", id, name)
	}
	return fieldValidate(name, value, v)
}

const (
	sizeKey = "size"
)

func stringValidators(tag string, name string, value string) (ok bool, err error) {
	vals, err := vals(tag, name)
	if err != nil {
		return false, err
	}
	var v Validator[string]
	switch id := strings.TrimSpace(vals[0]); id {
	case "length":
		min, max, err := rangeFinder(vals[1:])
		if err != nil {
			return false, err
		}
		v = &LengthRangeValidator{
			Min: min,
			Max: max,
		}
	case "min":
		size, err := intParam(sizeKey, vals[1:])
		if err != nil {
			return false, err
		}
		v = &MinLengthValidator{
			Size: size,
		}
	case "max":
		size, err := intParam(sizeKey, vals[1:])
		if err != nil {
			return false, err
		}
		v = &MaxLengthValidator{
			Size: size,
		}
	case "regex":
		pattern, err := patternParam(vals[1:])
		if err != nil {
			return false, err
		}
		v = &RegexValidator{
			Pattern: pattern,
		}
	case "alphanum":
		v = &AlphaNumericValidator{}
	case "ipv4":
		v = &IPv4Validator{}
	case "ipv6":
		v = &IPv6Validator{}
	case "mac":
		v = &MACAddressValidator{}
	case "json":
		v = &JSONValidator{}
	case "xml":
		v = &XMLValidator{}
	case "url":
		v = &UrlValidator{}
	case "email":
		v = &EmailValidator{}
	case "!empty":
		v = &NonEmptyStringValidator{}
	default:
		return false, fmt.Errorf("unknown validator %q  for field %q", id, name)
	}
	return fieldValidate(name, value, v)
}

func fieldValidate[T cmp.Ordered](name string, value T, v Validator[T]) (ok bool, err error) {
	if ok, err := v.Validate(value); !ok {
		return false, fmt.Errorf("error validating field %q: %v", name, err)
	}
	return true, nil
}

func vals(tag, name string) ([]string, error) {
	vals := strings.Split(tag, ",")
	if len(vals) == 0 {
		return nil, fmt.Errorf("missing validator for field %q", name)
	}
	return vals, nil
}

const (
	minKey = "min"
	maxKey = "max"
)

func rangeFinder(params []string) (min int, max int, err error) {
	if len(params) != 2 {
		return 0, 0, fmt.Errorf("expected 2 parameters (%s, %s), found: %v", minKey, maxKey, params)
	}
	for _, pair := range params {
		k, v, err := kv(pair)
		if err != nil {
			return 0, 0, err
		}
		if k == minKey {
			min, err = strconv.Atoi(v)
			if err != nil {
				return 0, 0, fmt.Errorf("invalid value %q for parameter %q", v, k)
			}
			continue
		}

		if k == maxKey {
			max, err = strconv.Atoi(v)
			if err != nil {
				return 0, 0, fmt.Errorf("invalid value %q for parameter %q", v, k)
			}
			continue
		}
	}
	return min, max, nil
}

const patternKey = "pattern"

func patternParam(params []string) (pattern *regexp.Regexp, err error) {
	v, err := stringParam(patternKey, params)
	if err != nil {
		return nil, err
	}
	pattern, err = regexp.Compile(v)
	if err != nil {
		return nil, err
	}
	return pattern, nil
}

func intParam(key string, params []string) (int, error) {
	v, err := stringParam(key, params)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid value %q for parameter %q", v, key)
	}
	return i, nil
}

func stringParam(key string, params []string) (val string, err error) {
	if len(params) != 1 || params[0] == "" {
		return "", fmt.Errorf("expected 1 parameter (%s), found: %v", key, params)
	}
	k, v, err := kv(params[0])
	if err != nil {
		return "", err
	}
	if k != key {
		return "", fmt.Errorf("expected parameter %q, found: %q", key, k)
	}
	return v, nil
}

func kv(pair string) (k string, v string, err error) {
	kv := strings.Split(pair, "=")
	if len(kv) == 2 {
		k = strings.TrimSpace(kv[0])
		v = strings.TrimSpace(kv[1])
		if k != "" && v != "" {
			return k, v, nil
		}
	}
	return "", "", fmt.Errorf("malformed key value pair %q, expected format is \"key=value\"", pair)
}
