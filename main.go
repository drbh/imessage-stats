package main

import (
	// "encoding/json"
	"fmt"
	"os"

	messagestats "github.com/drbh/imessage-stats/go-text-statistics"
)

func main() {

	argsWithoutProg := os.Args[1:]

	var number string
	if len(argsWithoutProg) > 0 {
		number = argsWithoutProg[0]
	} else {
		os.Exit(0)
	}

	if number == "--all" {
		allstats := messagestats.GetFullProfileStatsFullDatabase()
		fmt.Println(string(allstats))
	}

	if number == "counts" {

		allstats := messagestats.GetStringCountsFullDatabase()
		fmt.Println(string(allstats))
	} else {
		allstats := messagestats.GetFullProfileStats(number)
		fmt.Println(string(allstats))
	}

}
