package gol

// Params provides the details of how to run the Game of Life and which image to load.
type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

func calculateNextState(p Params, world [][]byte) [][]byte {
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

func calculateAliveCells(p Params, world [][]byte) []cell {
	var liveCells []cell
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == 255 {
				liveCells = append(liveCells, cell{x: x, y: y})
			}
		}
	}
	return liveCells
}

// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run(p Params, events chan<- Event, keyPresses <-chan rune) {

	//	TODO: Put the missing channels in here.

	ioCommand := make(chan ioCommand)
	ioIdle := make(chan bool)

	ioChannels := ioChannels{
		command:  ioCommand,
		idle:     ioIdle,
		filename: nil,
		output:   nil,
		input:    nil,
	}
	go startIo(p, ioChannels)

	distributorChannels := distributorChannels{
		events:     events,
		ioCommand:  ioCommand,
		ioIdle:     ioIdle,
		ioFilename: nil,
		ioOutput:   nil,
		ioInput:    nil,
	}
	distributor(p, distributorChannels)
}
