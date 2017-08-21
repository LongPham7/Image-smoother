package main

import (
	"fmt"
	"math/rand"
)

const height = 50
const width = 30

var image [][]bool
var auxiliaryImage [][]bool

func initialize() {
	image = make([][]bool, height)
	auxiliaryImage = make([][]bool, height)
	rand.Seed(10)
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

func system(WORKER int) {
	blockheight := height / WORKER
	fromW := make(chan bool)
	toW := make(chan bool)
	for i := 0; i != WORKER; i++ {
		go worker(i*blockheight, blockheight, toW, fromW)
	}
	coordinate(WORKER, fromW, toW)
}

func coordinate(WORKER int, in, out chan bool) {
	var nochange bool = false
	for !nochange {
		nochange = true
		for i := 0; i != WORKER; i++ {
			input := <-in
			nochange = nochange && input
		}
		if !nochange {
			for i := 0; i != WORKER; i++ {
				out <- true
			}
		}
	}
	for i := 0; i != WORKER; i++ {
		out <- false
	}
}

func worker(start int, length int, in, out chan bool) {
	var proceed bool = true
	var change bool
	for proceed {
		change = smoothBlock(start, length, in, out)
		out <- change
		proceed = <-in
	}
}

func smoothBlock(start int, length int, in, out chan bool) bool {
	nochange := true
	p := false

	for i := start; i != start+length; i++ {
		for j := 0; j != width; j++ {
			auxiliaryImage[i][j] = calculateMajority(image, i, j)
		}
	}

	out <- false
	<-in

	for i := start; i != start+length; i++ {
		for j := 0; j != width; j++ {
			p = calculateMajority(auxiliaryImage, i, j)
			if p != image[i][j] {
				nochange = false
			}
			image[i][j] = p
		}
	}
	return nochange
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

	return count > total / 2
}

func main() {
	initialize()
	fmt.Println("The original image:")
	printImage(image)
	fmt.Println("")
	system(3)
	fmt.Println("The new image:")
	printImage(image)
}
