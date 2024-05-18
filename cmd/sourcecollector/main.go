package main

import (
	sourcecollector "github.com/hitesh22rana/sourcecollector/pkg"
)

func main() {
	sc, err := sourcecollector.NewSourceCollector("/Users/hiteshrana/Work/sourcecollector", "/Users/hiteshrana/Work/sourcecollector/output.txt")
	if err != nil {
		panic(err)
	}

	if err := sc.Save(); err != nil {
		panic(err)
	}
}
