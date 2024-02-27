// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package maze implements a maze generator using Wilson's algorithm
package maze

import (
	"math/rand"
)

type Rectangle struct {
	g *grid
}

func RectangleMaze(height, width int) (*Rectangle, error) {
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

	return &Rectangle{g: g}, nil
}

func SquareMaze(height int) (*Rectangle, error) {
	return RectangleMaze(height, height)
}
