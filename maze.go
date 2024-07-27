// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package maze implements a maze generator using Wilson's algorithm
package maze

import (
	"log"
	"math/rand"
	"time"
)

type Rectangle struct {
	g        *grid
	entrance *cell
	exit     *cell
	solved   bool
}

func RectangleMaze(height, width int, solve bool) (*Rectangle, error) {
	g := createGrid(height, width)

	// create a stack containing all the cells in the grid in a random order
	var stack []*cell
	stack = g.allCells()
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
	// the exit will be on the eastern third of the southern edge of the maze.
	exitRow, exitCol := south, east
	exitCol = east - rand.Intn(theGate)
	// set the flags on the entrance and exit cells
	entrance := g.cells[entranceRow][entranceCol]
	entrance.entrance = true
	entrance.walls.north = false
	exit := g.cells[exitRow][exitCol]
	exit.exit = true
	exit.walls.south = false

	if solve {
		started := time.Now()
		log.Printf("maze: solving maze\n")

		// clear the walk pointers for this search
		g.clearWalk()

		// solve the maze using depth-first search
		stack = []*cell{entrance}
		entrance.visited = true
		for !stack[len(stack)-1].isExit() {
			// pop current cell off top of stack
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			//log.Printf("maze: depth %6d current %4d %4d\n", len(stack), current.row, current.col)

			// optimization - if neighbor is the exit, push it and quit searching
			if current.southIsOpen() {
				if neighbor := current.neighbors.south; neighbor.isExit() {
					neighbor.visited = true
					neighbor.to = current
					stack = append(stack, neighbor)
					break
				}
			}

			// push all neighbors that haven't yet been visited on to the stack
			if current.northIsOpen() {
				if neighbor := current.neighbors.north; !neighbor.hasBeenVisited() {
					neighbor.visited = true
					neighbor.to = current
					stack = append(stack, neighbor)
				}
			}
			if current.eastIsOpen() {
				if neighbor := current.neighbors.east; !neighbor.hasBeenVisited() {
					neighbor.visited = true
					neighbor.to = current
					stack = append(stack, neighbor)
				}
			}
			if current.southIsOpen() {
				if neighbor := current.neighbors.south; !neighbor.hasBeenVisited() {
					neighbor.visited = true
					neighbor.to = current
					stack = append(stack, neighbor)
				}
			}
			if current.westIsOpen() {
				if neighbor := current.neighbors.west; !neighbor.hasBeenVisited() {
					neighbor.visited = true
					neighbor.to = current
					stack = append(stack, neighbor)
				}
			}
		}
		log.Printf("maze: solved  %5d x %5d maze in %v\n", g.height, g.width, time.Now().Sub(started))

		// flag each cell that is on the path between the entrance and the exit
		for c := exit; c != nil; c = c.to {
			c.onPath = true
		}
	}

	return &Rectangle{
		g:        g,
		entrance: entrance,
		exit:     exit,
	}, nil
}

func (r *Rectangle) Solve() {
	if r.solved {
		return
	}
	started := time.Now()
	log.Printf("maze: solving maze\n")

	// clear the walk pointers for this search
	r.g.clearWalk()

	// solve the maze using depth-first search
	stack := []*cell{r.entrance}
	r.entrance.visited = true
	for !stack[len(stack)-1].isExit() {
		// pop current cell off top of stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		//log.Printf("maze: depth %6d current %4d %4d\n", len(stack), current.row, current.col)

		// optimization - if neighbor is the exit, push it and quit searching
		if current.southIsOpen() {
			if neighbor := current.neighbors.south; neighbor.isExit() {
				neighbor.visited = true
				neighbor.to = current
				stack = append(stack, neighbor)
				break
			}
		}

		// push all neighbors that haven't yet been visited on to the stack
		if current.northIsOpen() {
			if neighbor := current.neighbors.north; !neighbor.hasBeenVisited() {
				neighbor.visited = true
				neighbor.to = current
				stack = append(stack, neighbor)
			}
		}
		if current.eastIsOpen() {
			if neighbor := current.neighbors.east; !neighbor.hasBeenVisited() {
				neighbor.visited = true
				neighbor.to = current
				stack = append(stack, neighbor)
			}
		}
		if current.southIsOpen() {
			if neighbor := current.neighbors.south; !neighbor.hasBeenVisited() {
				neighbor.visited = true
				neighbor.to = current
				stack = append(stack, neighbor)
			}
		}
		if current.westIsOpen() {
			if neighbor := current.neighbors.west; !neighbor.hasBeenVisited() {
				neighbor.visited = true
				neighbor.to = current
				stack = append(stack, neighbor)
			}
		}
	}
	log.Printf("maze: solved  %5d x %5d maze in %v\n", r.g.height, r.g.width, time.Now().Sub(started))

	// flag each cell that is on the path between the entrance and the exit
	for c := r.exit; c != nil; c = c.to {
		c.onPath = true
	}

	r.solved = true
}

func SquareMaze(height int, solve bool) (*Rectangle, error) {
	return RectangleMaze(height, height, solve)
}
