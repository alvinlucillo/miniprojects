package main

import (
	"fmt"
	"os"

	processor "alvinlucillo/toy-robot-challenge/internal/processor"
)

// main - entrypoint to the program
func main() {
	processor, err := processor.NewProcessor(os.Stdin, processor.SourceTypeStdin)
	if err != nil {
		fmt.Println(err)
		return
	}
	processor.Execute()
}
