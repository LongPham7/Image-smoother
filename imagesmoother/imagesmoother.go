package main

import (
	"fmt"
	"math/rand"
)

const height = 50
const width = 30
const seed = 42

var image [][]bool
var auxiliaryImage [][]bool

func initialize() {
	image = make([][]bool, height)
	auxiliaryImage = make([][]bool, height)
	rand.Seed(seed)
	for r := 0; r != height; r++ {
		image[r] = make([]bool, width)
		auxiliaryImage[r] = make([]bool, width)
		for c := 0; c != width; c++ {
			if rand.Intn(50) > 14 {
				image[r][c] = true
			}
		}
	}
}

func printImage(image [][]bool) {
	for r := 0; r != len(image); r++ {
		printRow(image, r)
	}
}

func printRow(image [][]bool, r int) {
	for c := 0; c != len(image[r]); c++ {
		if image[r][c] {
			fmt.Print("*")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}

func system(WORKERS int) {
	blockheight := height / WORKERS
	fromW := make(chan bool)
	toW := make(chan bool)
	for i := 0; i != WORKERS; i++ {
		go worker(i*blockheight, blockheight, toW, fromW)
	}
	coordinator(WORKERS, fromW, toW)
}

func coordinator(WORKERS int, in, out chan bool) {
	noChange := false
	for !noChange {
		noChange = true
		for i := 0; i != WORKERS; i++ {
			noChange = <-in && noChange // For some reason, noChnage && <-in does not work.
		}
		if !noChange {
			for i := 0; i != WORKERS; i++ {
				out <- true
			}
		}
	}
	for i := 0; i != WORKERS; i++ {
		out <- false
	}
}

func worker(start int, length int, in, out chan bool) {
	proceed := true
	for proceed {
		out <- smoothBlock(start, length, in, out)
		proceed = <-in
	}
}

func smoothBlock(start int, length int, in, out chan bool) bool {
	noChange := true
	majority := false

	for r := start; r != start+length; r++ {
		for c := 0; c != width; c++ {
			auxiliaryImage[r][c] = calculateMajority(image, r, c)
		}
	}

	out <- false
	<-in

	for r := start; r != start+length; r++ {
		for c := 0; c != width; c++ {
			majority = calculateMajority(auxiliaryImage, r, c)
			if majority != image[r][c] {
				noChange = false
			}
			image[r][c] = majority
		}
	}
	return noChange
}

func calculateMajority(image [][]bool, r, c int) bool {
	total := 0
	count := 0
	for i := -1; i != 2; i++ {
		for j := -1; j != 2; j++ {
			if 0 <= r+i && r+i < height && 0 <= c+j && c+j < width {
				total++
				if image[r+i][c+j] {
					count++
				}
			}
		}
	}
	total--
	if image[r][c] {
		count--
	}

	return count > total/2
}

func main() {
	initialize()
	fmt.Println("The original image:")
	printImage(image)
	fmt.Println("")
	system(5)
	fmt.Println("The new image:")
	printImage(image)
}
