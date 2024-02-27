// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package main implements a command line application to generate mazes using Wilson's algorithm
package main

import (
	"bytes"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	started := time.Now()
	height, width := 25, 25
	if g, err := run(height, width); err != nil {
		log.Fatal(err)
	} else if err := os.WriteFile("wilson.txt", g.toText(), 0644); err != nil {
		log.Fatal(err)
	}
	log.Printf("maze: %5d x %5d in %v\n", height, width, time.Now().Sub(started))
}

func run(height, width int) (*grid, error) {
	g := createGrid(height, width)

	// create a stack containing all the cells in the grid in a random order
	stack := g.allCells()
	rand.Shuffle(len(stack), func(i, j int) {
		stack[i], stack[j] = stack[j], stack[i]
	})

	// randomly add a cell to the maze.
	// since the stack contains all cells in a random order, we can just pop the first cell from it
	// and mark it as in.
	stack[0].in = true
	stack = stack[1:]

	// while the stack is not empty, pop a cell.
	// perform a random walk from that cell, stopping only when we encounter a cell that is already in the maze.
	// for every cell that we visit, we record the direction that we exited so that we'll be able to retrace our path.
	for len(stack) != 0 {
		// pick a cell at random from the stack.
		// since the stack is randomly shuffled before we start, we can just pop the first cell.
		from := stack[0]
		if from == nil {
			panic("assert(from != nil)")
		}
		stack = stack[1:]

		// clear the walk pointers for this iteration
		g.clearWalk()

		// randomly walk until we find a cell that is already in the maze
		for to := from; !to.in; {
			// pick a neighboring cell at random
			to.to = to.randomNeighbor()
			// and move to it
			to = to.to
		}

		// retrace the walk, removing walls as needed, until we find a cell that is in the maze
		for !from.in {
			to := from.to
			// remove the wall between the from and to cells
			if from.neighbors.north == to {
				from.walls.north = false
				to.walls.south = false
			} else if from.neighbors.east == to {
				from.walls.east = false
				to.walls.west = false
			} else if from.neighbors.south == to {
				from.walls.south = false
				to.walls.north = false
			} else if from.neighbors.west == to {
				from.walls.west = false
				to.walls.east = false
			}
			// the cell is now in the maze, so mark it
			from.in = true
			// walk to the next cell
			from = from.to
		}
	}

	return g, nil
}

type grid struct {
	height, width int
	cells         [][]*cell
}

type cell struct {
	row, col int
	// todo: implement hex grid
	neighbors struct {
		north *cell
		east  *cell
		south *cell
		west  *cell
	}
	walls struct {
		north bool
		east  bool
		south bool
		west  bool
	}
	// set of all neighbors
	neighborhood []*cell
	// in is set to true if the cell has been added to the maze
	in bool
	// to points the last cell visited in the walk
	to *cell
}

func createGrid(height, width int) *grid {
	g := &grid{height: height, width: width, cells: make([][]*cell, height)}

	// allocate memory for all the cells in the grid
	for row := 0; row < height; row++ {
		g.cells[row] = make([]*cell, width)
		for col := 0; col < width; col++ {
			c := &cell{row: row, col: col}
			c.walls.north = true
			c.walls.east = true
			c.walls.south = true
			c.walls.west = true
			g.cells[row][col] = c
		}
	}

	// link neighboring cells
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			c := g.cells[row][col]
			if row > 0 {
				neighbor := g.cells[row-1][col]
				c.neighbors.north = neighbor
				c.neighborhood = append(c.neighborhood, neighbor)
			}
			if col < width-1 {
				neighbor := g.cells[row][col+1]
				c.neighbors.east = neighbor
				c.neighborhood = append(c.neighborhood, neighbor)
			}
			if row < height-1 {
				neighbor := g.cells[row+1][col]
				c.neighbors.south = neighbor
				c.neighborhood = append(c.neighborhood, neighbor)
			}
			if col > 0 {
				neighbor := g.cells[row][col-1]
				c.neighbors.west = neighbor
				c.neighborhood = append(c.neighborhood, neighbor)
			}
		}
	}

	return g
}

func mkmaze(height, width int) [][]rune {
	// create the maze
	maze := make([][]rune, height*2+1)
	for row := range maze {
		maze[row] = make([]rune, width*2+1)
	}
	north, east, south, west := 0, width*2, height*2, 0
	for row := north; row <= south; row++ {
		isNorthEdge, isSouthEdge := row == north, row == south
		isRow := row%2 == 0
		for col := west; col <= east; col++ {
			isWestEdge, isEastEdge := col == west, col == east
			isCol := col%2 == 0
			if isNorthEdge && isWestEdge {
				maze[row][col] = '╔'
			} else if isNorthEdge && isEastEdge {
				maze[row][col] = '╗'
			} else if isNorthEdge && isCol {
				maze[row][col] = '╦'
			} else if isNorthEdge {
				maze[row][col] = '═'
			} else if isSouthEdge && isWestEdge {
				maze[row][col] = '╚'
			} else if isSouthEdge && isEastEdge {
				maze[row][col] = '╝'
			} else if isSouthEdge && isCol {
				maze[row][col] = '╩'
			} else if isSouthEdge {
				maze[row][col] = '═'
			} else if isWestEdge && isRow {
				maze[row][col] = '╠'
			} else if isWestEdge {
				maze[row][col] = '║'
			} else if isEastEdge && isRow {
				maze[row][col] = '╣'
			} else if isEastEdge {
				maze[row][col] = '║'
			} else if isRow && isCol {
				maze[row][col] = '╬'
			} else if isRow {
				maze[row][col] = '═'
			} else if isCol {
				maze[row][col] = '║'
			} else {
				maze[row][col] = ' '
			}
		}
	}

	// create an entrance and an exit
	entranceRow, exitRow := north, south

	// entrances and exits will be on the western and eastern thirds of the maze.
	thirdWidth := width / 3
	if thirdWidth == 0 {
		panic("maze too small to generate entrance or exit")
	}
	log.Printf("width is %5d\n", thirdWidth)

	// Choose an entrance column within the western (first) third of the maze.
	// The entrance should be a column that is an even number (a path column, not a wall) and
	// is not on the west edge of the maze.
	// We start from 2 (not 0, to not be on the western edge) and end at thirdWidth (not including).
	entranceCol := west + rand.Intn(thirdWidth)*2 + 1

	// Choose an exit column within the eastern (last) third of the maze.
	// The exit should be a column that is an even number (a path, not a wall) and
	// is not on the east edge of the maze.
	exitCol := east - (rand.Intn(thirdWidth)*2 + 1)

	// punch the entrance and exit
	maze[entranceRow][entranceCol] = ' '
	maze[exitRow][exitCol] = ' '

	return maze
}

func (g *grid) allCells() []*cell {
	var cells []*cell
	for row := 0; row < g.height; row++ {
		for col := 0; col < g.width; col++ {
			cells = append(cells, g.cells[row][col])
		}
	}
	return cells
}

func (g *grid) clearWalk() {
	for row := 0; row < g.height; row++ {
		for col := 0; col < g.width; col++ {
			g.cells[row][col].to = nil
		}
	}
}

func (g *grid) toText() []byte {
	// create the maze
	maze := mkmaze(g.height, g.width)

	// punch out walls
	for row := 0; row < g.height; row++ {
		crow := row*2 + 1
		for col := 0; col < g.width; col++ {
			c := g.cells[row][col]
			ccol := col*2 + 1
			if !c.walls.north {
				maze[crow-1][ccol] = ' '
			}
			if !c.walls.east {
				maze[crow][ccol+1] = ' '
			}
			if !c.walls.south {
				maze[crow+1][ccol] = ' '
			}
			if !c.walls.west {
				maze[crow][ccol-1] = ' '
			}
		}
	}

	// print out the maze
	b := &bytes.Buffer{}
	for _, line := range maze {
		for _, ch := range line {
			b.WriteRune(ch)
		}
		b.WriteByte('\n')
	}
	b.WriteByte('\n')

	return b.Bytes()
}

// randomNeighbor returns a neighboring cell at random.
// if the cell is on an edge, the set won't include the walls.
func (c *cell) randomNeighbor() *cell {
	// pick a random direction
	direction := rand.Intn(len(c.neighborhood))
	rn := c.neighborhood[direction]
	if rn == nil {
		panic("assert(rn != nil)")
	}
	return rn
}
