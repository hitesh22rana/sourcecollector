package main

import (
	"fmt"

	sourcecollector "github.com/hitesh22rana/sourcecollector/pkg"
)

func main() {
	repo, err := sourcecollector.NewRepository("https://github.com/hitesh22rana/SoundScripter")
	if err != nil {
		fmt.Println(err)
	}

	metaData, err := repo.GetMetadata()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(metaData)
}
