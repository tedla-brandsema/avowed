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

func validateIntField(directive string, args []string, val int) error {
	var v Validator[int]
	switch directive {
	case "range":
		keys := []string{"min", "max"}
		pairs, err := findIntPairs(keys, args)
		if err != nil {
			return err
		}
		v = &IntRangeValidator{
			Min: pairs[keys[0]],
			Max: pairs[keys[1]],
		}
	case "pos":
		v = &NonNegativeIntValidator{}
	case "neg":
		v = &NonPositiveIntValidator{}
	default:
		return fmt.Errorf("unknown validator %q", directive)
	}
	return fieldValidate(val, v)
}

func validateStrField(directive string, args []string, val string) error {
	var v Validator[string]
	switch directive {
	case "length":
		keys := []string{"min", "max"}
		pairs, err := findIntPairs(keys, args)
		if err != nil {
			return err
		}
		v = &LengthRangeValidator{
			Min: pairs[keys[0]],
			Max: pairs[keys[1]],
		}
	case "min":
		keys := []string{"size"}
		pairs, err := findIntPairs(keys, args)
		if err != nil {
			return err
		}
		v = &MinLengthValidator{
			Size: pairs[keys[0]],
		}
	case "max":
		keys := []string{"size"}
		pairs, err := findIntPairs(keys, args)
		if err != nil {
			return err
		}
		v = &MinLengthValidator{
			Size: pairs[keys[0]],
		}
	case "regex":
		keys := []string{"pattern"}
		pairs, err := findStringPairs(keys, args)
		if err != nil {
			return err
		}
		pattern, err := regexp.Compile(pairs[keys[0]])
		if err != nil {
			return err
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
		return fmt.Errorf("unknown validator %q", directive)
	}
	return fieldValidate(val, v)
}

func ValidateStruct(data interface{}) (bool, error) {
	var err error

	val := reflect.ValueOf(data)
	for n := 0; n < val.NumField(); n++ {
		field := val.Type().Field(n)
		if tag, ok := field.Tag.Lookup(tagID); ok {
			directive, args := splitTag(tag)
			i := val.FieldByName(field.Name).Interface()
			switch v := i.(type) {
			case string:
				err = validateStrField(directive, args, v)
			case int:
				err = validateIntField(directive, args, v)
			}
		}
		if err != nil {
			return false, fmt.Errorf("error validating field %q: %v", field.Name, err)
		}
	}
	return true, nil
}

func splitTag(tag string) (id string, args []string) {
	args = strings.Split(tag, ",")
	return strings.TrimSpace(args[0]), args[1:]
}

func fieldValidate[T cmp.Ordered](value T, v Validator[T]) error {
	if ok, err := v.Validate(value); !ok {
		return err
	}
	return nil
}

func findIntPairs(keys []string, params []string) (map[string]int, error) {
	intPairs := make(map[string]int)
	pairs, err := findStringPairs(keys, params)
	if err != nil {
		return nil, err
	}
	for k, v := range pairs {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value %q for parameter %q", v, k)
		}
		intPairs[k] = i
	}
	return intPairs, nil
}

func findStringPairs(keys []string, args []string) (map[string]string, error) {
	pairs, err := extractPairs(args)
	if err != nil {
		return nil, err
	}
	if len(keys) != len(pairs) {
		return nil, fmt.Errorf("expected %d parameter(s) (%s), found: %s", len(keys), keys, args)
	}

	for _, key := range keys {
		if _, ok := pairs[key]; !ok {
			return nil, fmt.Errorf("missing required parameter %q", keys)
		}
	}
	return pairs, nil
}

func extractPairs(args []string) (map[string]string, error) {
	pairs := make(map[string]string)
	for _, pair := range args {
		k, v, err := kv(pair)
		if err != nil {
			return nil, err
		}
		pairs[k] = v
	}
	return pairs, nil
}

func kv(pair string) (k string, v string, err error) {
	split := strings.Split(pair, "=")
	if len(split) == 2 {
		k = strings.TrimSpace(split[0])
		v = strings.TrimSpace(split[1])
		if k != "" && v != "" {
			return k, v, nil
		}
	}
	return "", "", fmt.Errorf("malformed key value pair %q, expected format is \"key=value\"", pair)
}
