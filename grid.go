// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package maze

// grid contains all the cells in the maze.
type grid struct {
	height int
	width  int
	cells  [][]*cell
}

// createGrid creates a new rectangular grid with the given height and width.
func createGrid(height, width int) *grid {
	g := &grid{
		height: height,
		width:  width,
		cells:  make([][]*cell, height),
	}

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

	// link neighboring cells. required to implement the random walk.
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			// c is the current cell
			c := g.cells[row][col]

			// link northern neighbor
			if row > 0 {
				neighbor := g.cells[row-1][col]
				c.neighbors.north = neighbor
				c.neighborhood = append(c.neighborhood, neighbor)
			}
			// link eastern neighbor
			if col < width-1 {
				neighbor := g.cells[row][col+1]
				c.neighbors.east = neighbor
				c.neighborhood = append(c.neighborhood, neighbor)
			}
			// link southern neighbor
			if row < height-1 {
				neighbor := g.cells[row+1][col]
				c.neighbors.south = neighbor
				c.neighborhood = append(c.neighborhood, neighbor)
			}
			// link western neighbor
			if col > 0 {
				neighbor := g.cells[row][col-1]
				c.neighbors.west = neighbor
				c.neighborhood = append(c.neighborhood, neighbor)
			}
		}
	}

	return g
}

// allCells returns a new slice containing all the cells in the grid.
func (g *grid) allCells() []*cell {
	var cells []*cell
	for row := 0; row < g.height; row++ {
		for col := 0; col < g.width; col++ {
			cells = append(cells, g.cells[row][col])
		}
	}
	return cells
}

// clearWalk resets the cells in the grid to ready it for another walk.
// it loops through all cells and updates the `to` field to nil.
func (g *grid) clearWalk() {
	for row := 0; row < g.height; row++ {
		for col := 0; col < g.width; col++ {
			g.cells[row][col].to = nil
		}
	}
}
