package short

func f() (int, error) {
	return 0, nil
}

func g(ok bool, ch chan int) {
	a := 1        // want "short variable declaration: declare it with the var keyword"
	b, err := f() // want "short variable declaration"
	_, c := f()   // want "short variable declaration"
	x, y := a, b  // want "short variable declaration"
	if ok {
		err := error(nil) // want "short variable declaration"
		_ = err
	}
	if v := a; v > 0 {
		w := v // want "short variable declaration"
		_ = w
	}
	for i := 0; i < a; i++ {
		d := i // want "short variable declaration"
		_ = d
	}
	switch {
	case ok:
		e := b // want "short variable declaration"
		_ = e
	}
	select {
	case v := <-ch:
		w := v // want "short variable declaration"
		_ = w
	}
	var h = func() {
		i := c // want "short variable declaration"
		_ = i
	}
	h()
	_, _, _ = err, x, y
}
