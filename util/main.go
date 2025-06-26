package util

import (
	"math"
	"math/rand/v2"
	"slices"

	"golang.org/x/exp/constraints"
)

func XYinRect(x, y, rx, ry, rw, rh float64) bool {
	return x >= rx && x <= rx+rw && y >= ry && y <= ry+rh
}

func Distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow((x1-x2), 2) + math.Pow((y1-y2), 2))
}

func TaxiDistance(x1, y1, x2, y2 int) int {
	return IntAbs(x1-x2) + IntAbs(y1-y2)
}

// maybe check for 0?
func UnitDirection(x1, y1, x2, y2 float64) (float64, float64) {
	d := Distance(x1, y1, x2, y2)
	return (x2 - x1) / d, (y2 - y1) / d
}

func IntAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func IntSign(x int) int {
	if x < 0 {
		return -1
	} else if x == 0 {
		return 0
	}
	return 1
}

func Clamp[T constraints.Ordered](x, a, b T) T {
	return max(min(x, b), a)
}

func NearestUnitDirectionVector(x1, y1, x2, y2 int) (int, int) {
	dx := x2 - x1
	dy := y2 - y1

	// not sure if we actually need this case
	if dx == 0 && dy == 0 {
		return 0, 0
	}
	// prefer dx to dy, may change
	if IntAbs(dx) >= IntAbs(dy) {
		return IntSign(dx), 0
	} else {
		return 0, IntSign(dy)
	}
}

func TaxicabDistance(x1, y1, x2, y2 int) int {
	return IntAbs(x1-x2) + IntAbs(y1-y2)
}

func ExpandArray[T any](arr []T, times int) []T {
	out := make([]T, 0)
	for _, elem := range arr {
		for i := 0; i < times; i++ {
			out = append(out, elem)
		}
	}
	return out
}

func FilteredChoice[T any](s []T, f func(T) bool, rand *rand.Rand) T {
	if len(s) == 0 {
		panic("Cannot make choice from 0 length slice")
	}
	var ret T
	seen := 0
	for _, o := range s {
		if f(o) {
			seen += 1
			if rand.IntN(seen) == 0 {
				ret = o
			}
		}
	}
	return ret
}

func Choice[T any](s []T, filter func(T) bool, rand *rand.Rand) T {
	if len(s) == 0 {
		panic("Cannot make choice from 0 length slice")
	}
	var ret T
	seen := 0
	for _, o := range s {
		if filter(o) {
			seen += 1
			if rand.IntN(seen) == 0 {
				ret = o
			}
		}
	}
	return ret
}

func MinItem[T any](s []T, f func(t T) int) T {
	val := math.MaxInt
	var best T
	if len(s) == 0 {
		panic("Cannot select from 0 length slice")
	}
	for _, item := range s {
		v := f(item)
		if v < val {
			best = item
			val = v
		}
	}
	return best
}

func FakePerutations[T any](s []T, rand *rand.Rand) [][]T {
	var out [][]T
	for range 1 {
		new := make([]T, len(s))
		copy(new, s)
		rand.Shuffle(len(new), func(i, j int) {
			new[i], new[j] = new[j], new[i]
		})
		out = append(out, new)
	}
	return out
}

func Permutations[T any](s []T) [][]T {
	if len(s) == 0 {
		return [][]T{{}}
	}
	var out [][]T
	for i, v := range s {
		rest := make([]T, len(s)-1)
		copy(rest, s[:i])
		copy(rest[len(rest):], s[i+1:])

		for _, p := range Permutations(rest) {
			out = append(out, append([]T{v}, p...))
		}
	}
	return out
}

func Union[T comparable](a, b []T) []T {
	m := make(map[T]struct{})
	for _, v := range a {
		m[v] = struct{}{}
	}
	for _, v := range b {
		m[v] = struct{}{}
	}
	out := make([]T, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func StringsUnion(input ...string) string {
	if len(input) == 0 {
		return ""
	}
	out := []byte(input[0])
	for i := 1; i < len(input); i++ {
		out = Union(out, []byte(input[i]))
	}
	slices.Sort(out)
	return string(out)
}

type LocRect struct {
	X, Y, W, H float64
}

func (l LocRect) Contains(x, y float64) bool {
	return XYinRect(x, y, l.X, l.Y, l.W, l.H)
}
