package option_test

import (
	"encoding/json"
	"testing"

	"github.com/worldiety/option"
)

func TestWrap(t *testing.T) {
	type Car struct {
		Serial string
	}

	type Test struct {
		RentedCar option.Ptr[Car]
	}

	test := Test{}
	test.RentedCar = option.Pointer(&Car{Serial: "1234"})

	test2 := test
	// test.RentedCar.Unwrap().Serial="1234" compiler: unassignable
	tmp := test2.RentedCar.Unwrap()
	tmp.Serial = "2345" // this only changes the dereferenced copy of the original car.
	if test.RentedCar.Unwrap().Serial == tmp.Serial {
		t.Failed()
	}

	buf := option.Must(json.Marshal(test))
	t.Log(string(buf))
	var tmp2 option.Opt[Car]
	if err := json.Unmarshal(buf, &tmp2); err != nil {
		t.Fatal(err)
	}

	if tmp2.Unwrap().Serial != test.RentedCar.Unwrap().Serial {
		t.Failed()
	}
}
