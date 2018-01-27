package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
)

var count = flag.Int("count", 0, "count of lines to generate")
var length = flag.Int("length", 0, "length of each line")

func main() {
	flag.Parse()
	if *count == 0 || *length == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	for i := 0; i < *count; i++ {
		line := ""
		for j := 0; j < *length; j++ {
			line += string(rune('a') + rune(rand.Intn(26)))
		}
		fmt.Printf("%s\n", line)
	}
}
