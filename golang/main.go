package main

import "fmt"

type penguinPart int

const (
	adelieBottom penguinPart = iota
	adelieTop
	chinstrapBottom
	chinstrapTop
	emperorBottom
	emperorTop
	gentooBottom
	gentooTop
)

var partNames = [...]string{
	"Adelie Bottom",
	"Adelie Top",
	"Chinstrap Bottom",
	"Chinstrap Top",
	"Emperor Bottom",
	"Emperor Top",
	"Gentoo Bottom",
	"Gentoo Top",
}

const (
	north int = iota
	east
	south
	west
)

type card struct {
	id    int
	rot   int // 0, 1, 2, or 3
	parts [4]penguinPart
}

// 1/2 3/2 0/0
// 5/0 7/0 2/2
// 6/2 4/2 8/0
// 17000
// 1/2 3/2 0/0
// 5/0 6/0 2/2
// 7/2 4/2 8/0

var cards = []card{
	card{0, 0, [4]penguinPart{emperorTop, adelieBottom, chinstrapTop, chinstrapTop}},
	card{1, 0, [4]penguinPart{emperorTop, gentooTop, chinstrapTop, adelieBottom}},
	card{2, 0, [4]penguinPart{emperorTop, gentooTop, chinstrapBottom, adelieTop}},
	card{3, 0, [4]penguinPart{emperorTop, adelieTop, gentooBottom, chinstrapBottom}},
	card{4, 0, [4]penguinPart{emperorBottom, adelieTop, chinstrapBottom, gentooTop}},
	card{5, 0, [4]penguinPart{emperorBottom, adelieTop, chinstrapBottom, gentooBottom}},
	card{6, 0, [4]penguinPart{emperorBottom, gentooBottom, chinstrapTop, adelieBottom}},
	card{7, 0, [4]penguinPart{emperorBottom, gentooBottom, chinstrapTop, adelieBottom}}, // note: same as above
	card{8, 0, [4]penguinPart{emperorBottom, adelieBottom, gentooTop, gentooBottom}},
}

func (c *card) part(face int) penguinPart {
	return c.parts[(face+c.rot)%4]
}

func (c *card) name(face int) string {
	return partNames[c.part(face)]
}

func (c card) String() string {
	return fmt.Sprintf("{N: %s, E: %s, S: %s, W: %s}",
		c.name(north), c.name(east), c.name(south), c.name(west))
}

func dump(board []card) {
	for i, c := range board {
		fmt.Printf("%d/%d", c.id, c.rot)
		if (i+1)%3 == 0 {
			fmt.Println()
		} else {
			fmt.Print(" ")
		}
	}
}

func permutations(arr []card) [][]card {
	var helper func([]card, int)
	res := [][]card{}

	helper = func(arr []card, n int) {
		if n == 1 {
			tmp := make([]card, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					arr[i], arr[n-1] = arr[n-1], arr[i]
				} else {
					arr[0], arr[n-1] = arr[n-1], arr[0]
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func b4dig(n int) []int {
	r := make([]int, 9)
	for i := 0; i < 9; i++ {
		d := n % 4
		n = n / 4
		r[8-i] = d
	}
	return r
}

var dirCombs [][]int

func genDirs() {
	dirCombs = make([][]int, 262144)
	for x := 0; x < 262144; x++ {
		dirCombs[x] = b4dig(x)
	}
}

func applyDirs(board []card, n int) {
	r := dirCombs[n]
	for i := range r {
		board[i].rot = r[i]
	}
}

var neighbors = []map[int]int{
	// row 0 [0 1 2]
	{north: -1, east: 1, south: 3, west: -1},
	{north: -1, east: 2, south: 4, west: 0},
	{north: -1, east: -1, south: 5, west: 1},
	// row 1 [3 4 5]
	{north: 0, east: 4, south: 6, west: -1},
	{north: 1, east: 5, south: 7, west: 3},
	{north: 2, east: -1, south: 8, west: 4},
	// row 1 [6 7 8]
	{north: 3, east: 7, south: -1, west: -1},
	{north: 4, east: 8, south: -1, west: 6},
	{north: 5, east: -1, south: -1, west: 7},
}

var dirs = [...]int{north, east, south, west}

var pair = map[penguinPart]penguinPart{
	adelieBottom:    adelieTop,
	adelieTop:       adelieBottom,
	gentooBottom:    gentooTop,
	gentooTop:       gentooBottom,
	chinstrapBottom: chinstrapTop,
	chinstrapTop:    chinstrapBottom,
	emperorBottom:   emperorTop,
	emperorTop:      emperorBottom,
}

func matches(a, b penguinPart) bool {
	return pair[a] == b
}

func checkOne(idx int, board []card) bool {
	ns := neighbors[idx]
	c := board[idx]
	for _, d := range dirs {
		// get the neighbor in the d direction if there is one:
		if n := ns[d]; n > -1 {
			nc := board[n]
			od := (d + 2) % 4 // get opposite direction
			//fmt.Printf("%v %v %v %v\n", c, nc, d, od)
			if !matches(c.part(d), nc.part(od)) {
				return false
			}
		}
	}
	return true
}

func checkBoard(board []card) bool {
	// dump(board)
	// fmt.Println()
	for i := range board {
		if !checkOne(i, board) {
			return false
		}
	}
	return true
}

func checkAllBoardCombs(board []card) {
	for x := range dirCombs {
		applyDirs(board, x)
		if checkBoard(board) {
			dump(board)
		}
	}
}

func worker(n int, jobs <-chan []card, done chan<- bool) {
	fmt.Printf("start worker %d\n", n)
	i := 0
	for j := range jobs {
		checkAllBoardCombs(j)
		i++
		if i%100 == 0 {
			fmt.Printf("%d: processed %d\n", n, i)
		}
	}
	done <- true
}

func main() {
	genDirs()
	fmt.Printf("calculated %d direction sets\n", len(dirCombs))
	perms := permutations(cards)
	fmt.Printf("calculated %d permutations\n", len(perms))
	jobs := make(chan []card, 362880)
	done := make(chan bool)

	parallel := 11

	for w := 1; w <= parallel; w++ {
		go worker(w, jobs, done)
	}

	for _, p := range perms {
		jobs <- p
	}
	close(jobs)
	for a := 1; a <= parallel; a++ {
		<-done
	}

	fmt.Println("done")
}
