package valex

import (
	"github.com/tedla-brandsema/tagex"
)

const tagKey = "val"

var (
	tag tagex.Tag
)

func init() {
	tag = tagex.NewTag(tagKey)

	// Int directives
	tagex.RegisterDirective[int](&tag, &IntRangeValidator{})
	tagex.RegisterDirective[int](&tag, &NonNegativeIntValidator{})
	tagex.RegisterDirective[int](&tag, &NonPositiveIntValidator{})

	// String directives
	tagex.RegisterDirective[string](&tag, &UrlValidator{})
	tagex.RegisterDirective[string](&tag, &EmailValidator{})
	tagex.RegisterDirective[string](&tag, &NonEmptyStringValidator{})
	tagex.RegisterDirective[string](&tag, &MinLengthValidator{})
	tagex.RegisterDirective[string](&tag, &MaxLengthValidator{})
	tagex.RegisterDirective[string](&tag, &LengthRangeValidator{})
	tagex.RegisterDirective[string](&tag, &AlphaNumericValidator{})
	tagex.RegisterDirective[string](&tag, &MACAddressValidator{})
	tagex.RegisterDirective[string](&tag, &IpValidator{})
	tagex.RegisterDirective[string](&tag, &IPv4Validator{})
	tagex.RegisterDirective[string](&tag, &IPv6Validator{})
	tagex.RegisterDirective[string](&tag, &XMLValidator{})
	tagex.RegisterDirective[string](&tag, &JSONValidator{})
}

func ValidateStruct(data interface{}) (bool, error) {
	return tag.ProcessStruct(data)
}
