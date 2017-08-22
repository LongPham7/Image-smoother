package main

import (
	"fmt"
	"math/rand"
)

var height int
var width int

var image [][]bool
var auxiliaryImage [][]bool

// initialize initializes the image to be smoothed. 
func initialize(seed int64) {
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

// printImage displays the input image. 
func printImage(image [][]bool) {
	for r := 0; r != len(image); r++ {
		printRow(image, r)
	}
}

// printRow displays the r-th row of the input image. The indices
// of rows are zero-based. 
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

// system executes a concurrent system consisting of one coordinator
// and the given number of workers. 
func system(WORKERS int) {
	var blockLowHeight int = height / WORKERS
	blockHighHeight := blockLowHeight + 1
	remainder := height % WORKERS
	base := remainder * blockHighHeight

	fromW := make(chan bool)
	toW := make(chan bool)

	for i := 0; i != remainder; i++ {
		go worker(i*blockHighHeight, blockHighHeight, toW, fromW)
	}

	for i := 0; i != WORKERS - remainder; i++ {
		go worker(base + i*blockLowHeight, blockLowHeight, toW, fromW)
	}
	coordinator(WORKERS, fromW, toW)
}

// coordinator executes a coordinator process. 
func coordinator(WORKERS int, in, out chan bool) {
	noChange := false
	input := false
	for !noChange {
		noChange = true
		for i := 0; i != WORKERS; i++ {
			input = <-in
			noChange = input && noChange
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

// worker executes a worker process. 
func worker(start int, length int, in, out chan bool) {
	proceed := true
	for proceed {
		out <- smoothBlock(start, length, in, out)
		proceed = <-in
	}
}

// smoothBlock smoothes the specified block of an input image. 
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

// calculateMajority returns the majority of the neighbours of the specified cell in
// an input image. 
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
	var seed int64
	fmt.Println("Please enter an integer to be used as a seed to generate random numbers:")
	fmt.Scan(&seed)

	fmt.Println("Please specify the height of an image:")
	fmt.Scan(&height)

	fmt.Println("Please specify the width of an image:")
	fmt.Scan(&width)

	initialize(seed)
	fmt.Println("The original image:")
	printImage(image)
	fmt.Println("")
	system(5)
	fmt.Println("The new image:")
	printImage(image)
}
