package main

import (
	// "encoding/json"
	"fmt"

	messagestats "github.com/drbh/imessage-stats/go-text-statistics"
)

func main() {

	allstats := messagestats.GetFullProfileStats("+13478344775")
	fmt.Println(string(allstats))
}
