package main

import (
	// "encoding/json"
	"fmt"
	"os"

	messagestats "github.com/drbh/imessage-stats/go-text-statistics"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {

	argsWithoutProg := os.Args[1:]

	var command string
	if len(argsWithoutProg) > 0 {
		command = argsWithoutProg[0]

	} else {
		os.Exit(0)
	}

	validCommands := []string{"all", "counts", "number"}
	if !contains(validCommands, command) {
		os.Exit(1)
	}

	// fmt.Println("Running", command)

	if command == "all" {
		allstats := messagestats.GetFullProfileStatsFullDatabase()
		fmt.Println(string(allstats))
	}

	if command == "counts" {

		allstats := messagestats.GetStringCountsFullDatabase()
		fmt.Println(string(allstats))
	}

	if command == "number" {
		number := argsWithoutProg[1]
		allstats := messagestats.GetFullProfileStats(number)
		fmt.Println(string(allstats))
	}

}
