package option

import (
	"fmt"
)

func HelloError(fail bool) error {
	if fail {
		return fmt.Errorf("fail")
	}

	return nil
}

func ExampleMustZero() {
	MustZero(HelloError(false))
}

func CalcSum(a, b int, fail bool) (int, error) {
	if fail {
		return 0, fmt.Errorf("fail")
	}

	return a + b, nil
}

func ExampleMust() {
	// Output:
	// 3
	sum := Must(CalcSum(1, 2, false))
	fmt.Println(sum)
}

func Close() error {
	return fmt.Errorf("close failed")
}

func ExampleTry() {
	fn := func() (err error) {
		defer Try(Close, &err)

		return nil
	}

	// Output:
	// close failed
	fmt.Println(fn())
}
