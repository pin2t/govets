// const and var declarations, none of them grouped, wherever Go allows them.
package clean

const a = 1
const b = 2

var c = 3
var d = 4

type (
	one int
	two string
)

func run(x int) {
	var e = x
	var f = 5
	const g = 6
	const h = 7
	if x > 0 {
		var i = 8
		var j = 9
		const k = 10
		const l = 11
		_, _, _, _ = i, j, k, l
	}
	_ = e
	_ = f
	_ = g
	_ = h
}

var closure = func(x int) {
	var m = x
	var n = 12
	const o = 13
	const p = 14
	if x > 0 {
		var q = 15
		const r = 16
		_, _ = q, r
	}
	_, _, _, _ = m, n, o, p
}

var _, _ = one(0), two("")
