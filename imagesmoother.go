package main

import (
	"fmt"
	"math/rand"
)

const height = 50
const width = 30

var image1 [height][width]bool
var image2 [height][width]bool

func initialize() {
	rand.Seed(10)
	for i := 0; i != height; i++ {
		for j := 0; j != width; j++ {
			if rand.Intn(50) > 14 {
				image1[i][j] = true
			}
		}
	}
}

func printImage() {
	for i := 0; i != height; i++ {
		printRow(i)
	}
}

func printRow(i int) {
	for j := 0; j != width; j++ {
		if image1[i][j] {
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
			image2[i][j] = calculateMajority1(i, j)
		}
	}

	out <- false
	<-in

	for i := start; i != start+length; i++ {
		for j := 0; j != width; j++ {
			p = calculateMajority2(i, j)
			if p != image1[i][j] {
				nochange = false
			}
			image1[i][j] = p
		}
	}
	return nochange
}

func calculateMajority1(r, c int) bool {
	total := 0
	count := 0
	for i := -1; i != 2; i++ {
		for j := -1; j != 2; j++ {
			if 0 <= r+i && r+i < height && 0 <= c+j && c+j < width {
				total++
				if image1[r+i][c+j] {
					count++
				}
			}
		}
	}
	total--
	if image1[r][c] {
		count--
	}

	return majority(count, total)
}

func calculateMajority2(r, c int) bool {
	total := 0
	count := 0
	for i := -1; i != 2; i++ {
		for j := -1; j != 2; j++ {
			if 0 <= r+i && r+i < height && 0 <= c+j && c+j < width {
				total++
				if image2[r+i][c+j] {
					count++
				}
			}
		}
	}
	total--
	if image2[r][c] {
		count--
	}

	return majority(count, total)
}

func main() {
	initialize()
	fmt.Println("The original image:")
	printImage()
	fmt.Println("")
	system(3)
	fmt.Println("The new image:")
	printImage()
}

func majority(x, y int) bool {
	if y%2 == 0 {
		return x > y / 2
	} else {
		return x > y / 2
	}
}
