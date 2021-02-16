package main

// S is a definition for test
type S struct {
	Int int
}

func main() {
	s := S{Int: 1}
	s.Int += 1
}
