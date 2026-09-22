// Grouped const and var declarations everywhere: package level, functions,
// nested blocks and function literals.
package grouped

const (
	a = 1 // want "grouped const declaration: give each constant its own const keyword and line"
	b = 2 // want "grouped const declaration: give each constant its own const keyword and line"
)

var (
	c = 3 // want "grouped var declaration: give each variable its own var keyword and line"
	d = 4 // want "grouped var declaration: give each variable its own var keyword and line"
)

const (
	single = 5 // want "grouped const declaration: give each constant its own const keyword and line"
)

var (
	e, f = 6, 7 // want "grouped var declaration: give each variable its own var keyword and line"
	g    = 8    // want "grouped var declaration: give each variable its own var keyword and line"
)

func run() {
	const (
		h = 9  // want "grouped const declaration: give each constant its own const keyword and line"
		i = 10 // want "grouped const declaration: give each constant its own const keyword and line"
	)
	var (
		j = 11 // want "grouped var declaration: give each variable its own var keyword and line"
		k = 12 // want "grouped var declaration: give each variable its own var keyword and line"
	)
	if j > 0 {
		const (
			l = 13 // want "grouped const declaration: give each constant its own const keyword and line"
		)
		var (
			m = 14 // want "grouped var declaration: give each variable its own var keyword and line"
		)
		_, _ = l, m
	}
	for n := 0; n < 1; n++ {
		const (
			o = 15 // want "grouped const declaration: give each constant its own const keyword and line"
		)
		var (
			p = 16 // want "grouped var declaration: give each variable its own var keyword and line"
		)
		_, _ = o, p
	}
	switch {
	default:
		var (
			q = 17 // want "grouped var declaration: give each variable its own var keyword and line"
		)
		_ = q
	}
	_, _, _, _, _ = h, i, j, k, 0
}

var closure = func() {
	const (
		r = 18 // want "grouped const declaration: give each constant its own const keyword and line"
	)
	var (
		s = 19 // want "grouped var declaration: give each variable its own var keyword and line"
	)
	_, _ = r, s
}

var _ = closure
var _, _, _, _, _ = a, b, c, d, single
var _, _, _, _ = e, f, g, 0
