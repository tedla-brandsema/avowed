package avowed

import (
	"cmp"
)

type Validator[T cmp.Ordered] interface {
	Validate(val T) (ok bool, err error)
}

type ValidatorFunc[T cmp.Ordered] func(val T) (ok bool, err error)

func (p ValidatorFunc[T]) Validate(val T) (ok bool, err error) {
	return p(val)
}

type ValidatedValue[T cmp.Ordered] struct {
	value     T
	Validator Validator[T]
}

func (v *ValidatedValue[T]) Set(val T) error {
	if ok, err := v.Validator.Validate(val); !ok {
		return err
	}
	v.value = val

	return nil
}

func (v *ValidatedValue[T]) Get() T {
	return v.value
}
