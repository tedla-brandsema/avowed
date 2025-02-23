package avowed

import (
	"regexp"
	"testing"
)

func TestIntRangeValidator(t *testing.T) {
	v := &IntRangeValidator{Min: 10, Max: 20}
	tests := []struct {
		input int
		ok    bool
	}{
		{15, true},
		{10, true},
		{20, true},
		{9, false},
		{21, false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("IntRangeValidator(%d): expected ok=%v, got ok=%v (err: %v)", tc.input, tc.ok, ok, err)
		}
	}
}

func TestNonNegativeIntValidator(t *testing.T) {
	v := &NonNegativeIntValidator{}
	tests := []struct {
		input int
		ok    bool
	}{
		{-1, false},
		{0, true},
		{1, true},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("NonNegativeIntValidator(%d): expected ok=%v, got ok=%v (err: %v)", tc.input, tc.ok, ok, err)
		}
	}
}

func TestNonPositiveIntValidator(t *testing.T) {
	v := &NonPositiveIntValidator{}
	tests := []struct {
		input int
		ok    bool
	}{
		{-1, true},
		{0, true},
		{1, false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("NonPositiveIntValidator(%d): expected ok=%v, got ok=%v (err: %v)", tc.input, tc.ok, ok, err)
		}
	}
}

func TestUrlValidator(t *testing.T) {
	v := &UrlValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"https://www.example.com", true},
		{"ftp://example.com", true}, // Acceptable as a valid URL scheme.
		{"invalid-url", false},
		{"", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("UrlValidator(%q): expected ok=%v, got ok=%v (err: %v)", tc.input, tc.ok, ok, err)
		}
	}
}

func TestEmailValidator(t *testing.T) {
	v := &EmailValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"user@example.com", true},
		{"invalid-email", false},
		{"", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("EmailValidator(%q): expected ok=%v, got ok=%v (err: %v)", tc.input, tc.ok, ok, err)
		}
	}
}

func TestNonEmptyStringValidator(t *testing.T) {
	v := &NonEmptyStringValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"hello", true},
		{"", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("NonEmptyStringValidator(%q): expected ok=%v, got ok=%v (err: %v)", tc.input, tc.ok, ok, err)
		}
	}
}

func TestMinLengthValidator(t *testing.T) {
	v := &MinLengthValidator{Min: 3}
	tests := []struct {
		input string
		ok    bool
	}{
		{"abc", true},
		{"abcd", true},
		{"ab", false},
		{"", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("MinLengthValidator(%q): expected ok=%v, got ok=%v (err: %v)", tc.input, tc.ok, ok, err)
		}
	}
}

func TestMaxLengthValidator(t *testing.T) {
	v := &MaxLengthValidator{Max: 3}
	tests := []struct {
		input string
		ok    bool
	}{
		{"abc", true},
		{"ab", true},
		{"abcd", false},
		{"", true},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("MaxLengthValidator(%q): expected ok=%v, got ok=%v (err: %v)", tc.input, tc.ok, ok, err)
		}
	}
}

func TestRegexValidator(t *testing.T) {
	pattern := regexp.MustCompile(`^\d+$`)
	v := &RegexValidator{Pattern: pattern}
	tests := []struct {
		input string
		ok    bool
	}{
		{"123", true},
		{"abc", false},
		{"123abc", false},
		{"", false},
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("RegexValidator(%q): expected ok=%v, got ok=%v (err: %v)", tc.input, tc.ok, ok, err)
		}
	}
}

func TestIpValidator(t *testing.T) {
	v := &IpValidator{}
	tests := []struct {
		input string
		ok    bool
	}{
		{"192.168.1.1", true},    // valid IPv4
		{"2001:0db8::1", true},   // valid IPv6
		{"invalid-ip", false},    // invalid IP
		{"123.456.789.0", false}, // invalid IPv4
	}
	for _, tc := range tests {
		ok, err := v.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("IpValidator(%q): expected ok=%v, got ok=%v (err: %v)", tc.input, tc.ok, ok, err)
		}
	}
}

func TestCompositeValidator_String(t *testing.T) {
	nonEmpty := &NonEmptyStringValidator{}
	minLength := &MinLengthValidator{Min: 3}
	composite := &CompositeValidator[string]{Validators: []Validator[string]{nonEmpty, minLength}}

	tests := []struct {
		input string
		ok    bool
	}{
		{"abc", true},
		{"ab", false}, // Fails min length
		{"", false},   // Fails non-empty
	}

	for _, tc := range tests {
		ok, err := composite.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("CompositeValidator (string) for input %q: expected ok=%v, got ok=%v (err: %v)", tc.input, tc.ok, ok, err)
		}
	}
}

func TestCompositeValidator_Int(t *testing.T) {
	nonNegative := &NonNegativeIntValidator{}
	rangeValidator := &IntRangeValidator{Min: 0, Max: 100}
	composite := &CompositeValidator[int]{Validators: []Validator[int]{nonNegative, rangeValidator}}

	tests := []struct {
		input int
		ok    bool
	}{
		{-5, false}, // Fails non-negative check
		{50, true},
		{150, false}, // Fails range check
	}

	for _, tc := range tests {
		ok, err := composite.Validate(tc.input)
		if ok != tc.ok {
			t.Errorf("CompositeValidator (int) for input %d: expected ok=%v, got ok=%v (err: %v)", tc.input, tc.ok, ok, err)
		}
	}
}
