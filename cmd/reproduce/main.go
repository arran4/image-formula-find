package main

import (
	"fmt"
	"image-formula-find/dna4"
)

func main() {
    dna := "ABQ" // X Y +
    rf, _, _ := dna4.ParseDNA(dna)
    fmt.Printf("Formula: %s\n", rf.String())
}
