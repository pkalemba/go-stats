package stats_test

import (
	"fmt"
	"github.com/pkalemba/go-stats"
)

func Dupa() {
	fmt.Println("dupa")
	s := stats.Stats{}
	s.Start()
	// Output: dupa
}
