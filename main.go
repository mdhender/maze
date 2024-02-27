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
	height, width := 22, 45
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
	// define constants for the edges of the maze
	north, east, south, west := 0, g.width-1, g.height-1, 0

	// randomly assign an entrance and exit to the maze.
	// entrances and exits will be on the western and eastern sides of the maze.
	theGate := g.width / 6
	// the entrance will be on the western third of the northern edge of the maze.
	entranceRow, entranceCol := north, west
	entranceCol = west + rand.Intn(theGate)
	// the exit will be on the easter third of the southern edge of the maze.
	exitRow, exitCol := south, east
	exitCol = east - rand.Intn(theGate)
	// set the flags on the entrance and exit cells
	g.cells[entranceRow][entranceCol].walls.north = false
	g.cells[exitRow][exitCol].walls.south = false

	// allocate memory for the maze, which we're representing as runes
	maze := make([][]rune, g.height*2+1)
	for row := 0; row < len(maze); row++ {
		maze[row] = make([]rune, g.width*2+1)
		for n := range maze[row] {
			maze[row][n] = '+'
		}
	}

	// now add the walls based on each cell's attributes
	for row := north; row <= south; row++ {
		for col := west; col <= east; col++ {
			c := g.cells[row][col]

			// derive the coordinates of the center of the cell in the maze array
			cRow, cCol := row*2+1, col*2+1

			// define flags for edges, rows, and columns
			isNorthEdge, isSouthEdge := row == north, row == south
			isWestEdge, isEastEdge := col == west, col == east

			var glyph rune

			// set the corners of the cell to the correct IBM box glyph
			// start with the northwest corner of the cell
			if isNorthEdge && isWestEdge {
				glyph = '╔'
			} else if isNorthEdge {
				glyph = '╦'
			} else if isWestEdge {
				glyph = '╠'
			} else {
				glyph = '╬'
			}
			maze[cRow-1][cCol-1] = glyph
			// set the northern edge of the cell
			if c.walls.north {
				glyph = '═'
			} else {
				glyph = ' '
			}
			maze[cRow-1][cCol] = glyph
			// set the northeast corner of the cell
			if isNorthEdge && isEastEdge {
				glyph = '╗'
			} else if isNorthEdge {
				glyph = '╦'
			} else if isEastEdge {
				glyph = '╣'
			} else {
				glyph = '╬'
			}
			maze[cRow-1][cCol+1] = glyph
			// set the eastern edge of the cell
			if c.walls.east {
				glyph = '║'
			} else {
				glyph = ' '
			}
			maze[cRow][cCol+1] = glyph
			// set the southeast corner of the cell
			if isSouthEdge && isEastEdge {
				glyph = '╝'
			} else if isSouthEdge {
				glyph = '╩'
			} else if isEastEdge {
				glyph = '╣'
			} else {
				glyph = '╬'
			}
			maze[cRow+1][cCol+1] = glyph
			// set the southern edge of the cell
			if c.walls.south {
				glyph = '═'
			} else {
				glyph = ' '
			}
			maze[cRow+1][cCol] = glyph
			// set the southwest corner of the cell
			if isSouthEdge && isWestEdge {
				glyph = '╚'
			} else if isSouthEdge {
				glyph = '╩'
			} else if isWestEdge {
				glyph = '╠'
			} else {
				glyph = '╬'
			}
			maze[cRow+1][cCol-1] = glyph
			// set the western edge of the cell
			if c.walls.west {
				glyph = '║'
			} else {
				glyph = ' '
			}
			maze[cRow][cCol-1] = glyph
			// always set the center of the cell to a space
			maze[cRow][cCol] = ' '
		}
	}

	// convert the runes in the maze to a slice of bytes
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
