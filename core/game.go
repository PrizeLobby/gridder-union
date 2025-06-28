package core

import (
	"crypto/sha256"
	"math/rand/v2"
	"strings"
	"time"

	"slices"

	"github.com/prizelobby/union-gridder/util"
)

const NUM_SETS = 9

type Game struct {
	Rand     *rand.Rand
	Sets     []string
	Targets  [6]string
	Matches  [6][]bool
	Slots    [9]string
	Extras   [9][]bool
	Solution []string
	Solved   bool
}

func (g *Game) Reset() {
	g.Solved = false
	var ALPHABET = "ABCDEFGHI"
	var letters = []byte(ALPHABET)

	var found = false

	for !found {
		var sets = []string{}
		for len(sets) < NUM_SETS {
			var setSize = g.Rand.IntN(2) + 2
			set := make([]byte, 0, setSize)
			for j := 0; j < setSize; j++ {
				set = append(set, util.Choice(letters, func(c byte) bool { return !slices.Contains(set, c) }, g.Rand))
			}
			slices.Sort(set)
			if !slices.Contains(sets, string(set)) {
				sets = append(sets, string(set))
			}
		}
		slices.Sort(sets)

		permutations := make([][]string, 0, 72)
		permutations = append(permutations, sets)
		for i := 0; i < 9; i++ {
			for j := i + 1; j < 9; j++ {
				sets2 := slices.Clone(sets)
				sets2[i], sets2[j] = sets2[j], sets2[i]
				permutations = append(permutations, sets2)
			}
		}
		seen := make(map[string][]int)
		for i, p := range permutations {
			unions := []string{util.StringsUnion(p[0], p[1], p[2]),
				util.StringsUnion(p[3], p[4], p[5]),
				util.StringsUnion(p[6], p[7], p[8]),
				util.StringsUnion(p[0], p[3], p[6]),
				util.StringsUnion(p[1], p[4], p[7]),
				util.StringsUnion(p[2], p[5], p[8]),
			}
			var u_string = ""
			for _, u := range unions {
				u_string += string(u) + ","
			}
			if s, ok := seen[u_string]; ok {
				seen[u_string] = []int{s[0], s[1] + 1}
			} else {
				seen[u_string] = []int{i, 1}
			}
		}

		for u_string, s := range seen {
			if s[1] == 1 {
				//fmt.Printf("Unique union %s found in permutation %d\n", u_string, s[0])
				perm := permutations[s[0]]
				permStr := ""
				for _, p := range perm {
					permStr += string(p) + " "
				}
				//fmt.Printf("Permutation: %s\n", permStr)
				g.Solution = perm
				p := slices.Clone(perm)
				slices.Sort(p)
				g.Sets = p

				unions := strings.Split(u_string, ",")
				targets := [6]string{}
				matches := [6][]bool{}
				for i, u := range unions {
					if u == "" {
						continue
					}
					targets[i] = u
					matches[i] = make([]bool, len(u))
				}
				g.Targets = targets
				g.Matches = matches
				found = true
				break
			}
		}
	}

	g.Extras = [9][]bool{}
	g.Slots = [9]string{}
}

func (g *Game) SetSlot(index int, set string) {
	g.Slots[index] = set
	t := []string{util.StringsUnion(g.Slots[0], g.Slots[1], g.Slots[2]),
		util.StringsUnion(g.Slots[3], g.Slots[4], g.Slots[5]),
		util.StringsUnion(g.Slots[6], g.Slots[7], g.Slots[8]),
		util.StringsUnion(g.Slots[0], g.Slots[3], g.Slots[6]),
		util.StringsUnion(g.Slots[1], g.Slots[4], g.Slots[7]),
		util.StringsUnion(g.Slots[2], g.Slots[5], g.Slots[8])}

	solved := true
	g.Extras[index] = make([]bool, len(set))
	for i, b := range []byte(set) {
		g.Extras[index][i] = !slices.Contains([]byte(g.Targets[index/3]), b) || !slices.Contains([]byte(g.Targets[index%3+3]), b)
	}

	for j := range 6 {
		if t[j] != g.Targets[j] {
			solved = false
		}
		for i, b := range g.Targets[j] {
			if strings.Contains(t[j], string(b)) {
				g.Matches[j][i] = true
			} else {
				g.Matches[j][i] = false
			}
		}
	}

	g.Solved = solved
}

func NewGame() *Game {
	seed := time.Now().String()
	sum := sha256.Sum256([]byte(seed))
	return &Game{
		Rand: rand.New(rand.NewChaCha8(sum)),
	}
}

func NewGameSeeded(seed string) *Game {
	sum := sha256.Sum256([]byte(seed))
	return &Game{
		Rand: rand.New(rand.NewChaCha8(sum)),
	}
}
