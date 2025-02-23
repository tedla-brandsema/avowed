package avowed

import "errors"

func Assert(eval bool, msg string) {
	if eval {
		return
	}
	panic(errors.New(msg))
}
