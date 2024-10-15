package gol

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
	//fmt.Println(worldUpdate)
	return worldUpdate
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	// TODO: Create a 2D slice to store the world.
	//initialWorld := readPgmImage(params, inputFilename)
	turn := 0
	c.events <- StateChange{turn, Executing}

	// TODO: Execute all turns of the Game of Life.
	for i := 0; i < p.Turns; i++ {
		calcNextState(p)
	}

	// TODO: Report the final state using FinalTurnCompleteEvent.

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
