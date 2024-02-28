// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package maze

import "math/rand"

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
	// entrance is set to true if the cell is an entrance
	entrance bool
	// exit is set to true if the cell is an exit
	exit bool
	// in is set to true if the cell has been added to the maze
	in bool
	// onPath is set if the cell is on the path between the entrance and the exit
	onPath bool
	// visited is set to true if the cell has been visited while trying to solve
	visited bool
	// to points the last cell visited in the walk
	to *cell
}

func (c *cell) hasBeenVisited() bool {
	return c != nil && c.visited
}

// isEntrance returns true if the cell is an entrance
func (c *cell) isEntrance() bool {
	return c != nil && c.entrance
}

// isExit returns true if the cell is either an exit
func (c *cell) isExit() bool {
	return c != nil && c.exit
}

// eastIsOpen returns true if the cell has a neighbor to the east and no wall between them
func (c *cell) eastIsOpen() bool {
	return c.neighbors.east != nil && !c.walls.east
}

// northIsOpen returns true if the cell has a neighbor to the north and no wall between them
func (c *cell) northIsOpen() bool {
	return c.neighbors.north != nil && !c.walls.north
}

// southIsOpen returns true if the cell has a neighbor to the south and no wall between them
func (c *cell) southIsOpen() bool {
	return c.neighbors.south != nil && !c.walls.south
}

// westIsOpen returns true if the cell has a neighbor to the west and no wall between them
func (c *cell) westIsOpen() bool {
	return c.neighbors.west != nil && !c.walls.west
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
