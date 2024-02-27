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
	// in is set to true if the cell has been added to the maze
	in bool
	// to points the last cell visited in the walk
	to *cell
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
