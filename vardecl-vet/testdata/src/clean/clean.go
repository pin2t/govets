// Variables declared with var, and := only where var cannot be written.
package clean

import "errors"

func f() (int, error) {
	return 0, errors.New("f")
}

func declared() {
	var a = 1
	var b, err = f()
	var c int
	var (
		d = 2
		e string
	)
	_, _, _, _, _, _ = a, b, err, c, d, e
}

func grammar(s []int, m map[string]int, ch chan int, x any) {
	for i := 0; i < len(s); i++ {
	}
	for i := range s {
		_ = i
	}
	for k, v := range m {
		_, _ = k, v
	}
	if v, ok := m["k"]; ok {
		_ = v
	} else if w := len(s); w > 0 {
		_ = w
	}
	switch n := len(s); n {
	case 0:
	}
	switch t := x.(type) {
	case int:
		_ = t
	}
	switch y := x; t := y.(type) {
	case string:
		_ = t
	}
	select {
	case v := <-ch:
		_ = v
	case v, ok := <-ch:
		_, _ = v, ok
	default:
	}
}

func redeclared() error {
	var err error
	a, err := f()
	_ = a
	return err
}
