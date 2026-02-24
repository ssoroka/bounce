package main

import (
	"flag"

	"github.com/ssoroka/bounce/levels"
)

var level = flag.Int("level", 1, "level you want to run")

func main() {
	flag.Parse()
	switch *level {
	case 1:
		levels.Level1()
	}
}
