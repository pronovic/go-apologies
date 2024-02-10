package main

import (
	"flag"
	"fmt"
	"github.com/pronovic/go-apologies/model"
)

func main() {
	players, mode := parseArgs()
	fmt.Printf("Running with %d players in %s\n", players, mode.Value())
}

func parseArgs() (int, model.GameMode) {
	players := flag.Int("players", 2, "number of players")
	adult := flag.Bool("adult", false, "run in adult mode")

	flag.Parse()

	mode := model.StandardMode
	if *adult {
		mode = model.AdultMode
	}

	return *players, mode
}