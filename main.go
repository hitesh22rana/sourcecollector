package main

import (
	"fmt"

	sourcecollector "github.com/hitesh22rana/sourcecollector/pkg"
)

func main() {
	repo, _ := sourcecollector.NewRepository("https://github.com/hitesh22rana/SoundScripter")
	fmt.Println(repo)
}
