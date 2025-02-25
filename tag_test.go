package avowed

import (
	"testing"
)

type MyStruct struct {
	Number int    `val:"range,min=4,max=6"`
	Word   string `val:"length,min=4,max=6"`
	Must   string `val:"alphnum"`
}

func TestValidateStruct(t *testing.T) {
	if ok, err := ValidateStruct(MyStruct{
		Word: "Pluk",
		Must: "hello6#",
	}); !ok {
		t.Error(err)
		return
	}

	t.Log("struct is valid!")
}
