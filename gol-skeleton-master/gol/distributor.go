package gol

import (
	"fmt"

	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

func calcNextState(p Params, world [][]byte) [][]byte {
	worldUpdate := make([][]byte, p.ImageHeight)
	for i := range worldUpdate {
		worldUpdate[i] = make([]byte, p.ImageWidth)
	}

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			sum := int(world[(y+p.ImageHeight-1)%p.ImageHeight][(x+p.ImageWidth-1)%p.ImageWidth])/255 + int(world[(y+p.ImageHeight-1)%p.ImageHeight][(x+p.ImageWidth)%p.ImageWidth])/255 + int(world[(y+p.ImageHeight-1)%p.ImageHeight][(x+p.ImageWidth+1)%p.ImageWidth])/255 +
				int(world[(y+p.ImageHeight)%p.ImageHeight][(x+p.ImageWidth-1)%p.ImageWidth])/255 + int(world[(y+p.ImageHeight)%p.ImageHeight][(x+p.ImageWidth+1)%p.ImageWidth])/255 +
				int(world[(y+p.ImageHeight+1)%p.ImageHeight][(x+p.ImageWidth-1)%p.ImageWidth])/255 + int(world[(y+p.ImageHeight+1)%p.ImageHeight][(x+p.ImageWidth)%p.ImageWidth])/255 + int(world[(y+p.ImageHeight+1)%p.ImageHeight][(x+p.ImageWidth+1)%p.ImageWidth])/255
			//fmt.Println(world)
			if world[y][x] == 255 {
				if sum < 2 {
					worldUpdate[y][x] = 0
				} else if sum == 2 || sum == 3 {
					worldUpdate[y][x] = 255
				} else {
					worldUpdate[y][x] = 0
				}
			} else {
				if sum == 3 {
					worldUpdate[y][x] = 255
				} else {
					worldUpdate[y][x] = 0
				}
			}
		}
	}
	return worldUpdate
}

func calculateAliveCells(p Params, world [][]byte) []util.Cell {
	var liveCells []util.Cell
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == 255 {
				liveCells = append(liveCells, util.Cell{X: x, Y: y})
			}
		}
	}
	return liveCells
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	// TODO: Create a 2D slice to store the world.

	var receivedData []uint8
	//input := make(chan int)
	imageFilename := "16x16"

	var world [][]uint8
	//var worldUpdate [][]uint8
	go func() {
		c.ioCommand <- 1
		c.ioFilename <- imageFilename
	}()

	for input := range c.ioInput {
		receivedData = append(receivedData, input)
	}

	turn := 0

	c.events <- StateChange{turn, Executing}

	world = make([][]byte, p.ImageHeight)
	for i := range world {
		world[i] = make([]byte, p.ImageWidth)
	}

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			world[y][x] = receivedData[y*p.ImageWidth+x]
		}
	}
	// TODO: Execute all turns of the Game of Life.
	for i := 0; i < p.Turns; i++ {
		world = calcNextState(p, world)
	}

	finalState := calculateAliveCells(p, world)
	fmt.Println(finalState)
	// TODO: Report the final state using FinalTurnCompleteEvent.
	c.events <- FinalTurnComplete{p.Turns, finalState}
	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle
	c.events <- StateChange{turn, Quitting}
	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
