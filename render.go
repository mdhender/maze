// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package maze

import (
	"bytes"
	"github.com/fogleman/gg"
	"os"
)

func (r *Rectangle) RenderPNG(path string, scale int) error {
	return r.g.toPNG(path, scale)
}

func (r *Rectangle) RenderText(path string) error {
	return r.g.toText(path)
}

// toPNG renders the grid as a PNG image file.
// each cell is scaled and a gutter is added to the final image.
func (g *grid) toPNG(path string, scale int) error {
	// calculate the gutter
	gutter := scale / 2
	if gutter < 5 {
		gutter = 5
	}

	// set the width and height of the image, assuming cells are scaled
	// and including room for the gutter
	width, height := g.width*scale+gutter*2, g.height*scale+gutter*2

	dc := gg.NewContext(width, height)

	// set the background of the image to white
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// the offset will be half the scale and allows for the gutter
	offset := scale/2 + gutter
	for x := 0; x < g.width; x++ {
		// derive the center x value of the cell in the image, assuming cells are 10x10
		cx := x*scale + offset
		for y := 0; y < g.height; y++ {
			// c is the cell that we're adding to the image
			c := g.cells[y][x]

			// derive the center y value of the cell in the image
			cy := y*scale + offset

			// derive values for the four corners of the cell
			type point struct {
				x, y float64
			}
			nw := point{x: float64(cx - scale/2), y: float64(cy - scale/2)}
			ne := point{x: float64(cx + scale/2), y: float64(cy - scale/2)}
			sw := point{x: float64(cx - scale/2), y: float64(cy + scale/2)}
			se := point{x: float64(cx + scale/2), y: float64(cy + scale/2)}

			// draw walls as black lines
			dc.SetRGB(0, 0, 0)

			// make the walls 3 pixels wide
			dc.SetLineWidth(3)

			// if there is a wall blocking the path north, draw a line from NW to NE corners.
			if c.walls.north {
				dc.DrawLine(nw.x, nw.y, ne.x, ne.y)
				dc.Stroke()
			}
			// if there is a wal blocking the path east, draw a line from the NE to SE corners.
			if c.walls.east {
				dc.DrawLine(ne.x, ne.y, se.x, se.y)
				dc.Stroke()
			}
			// if there is a wall blocking the path south, draw a line from SE to SW corners.
			if c.walls.south {
				dc.DrawLine(se.x, se.y, sw.x, sw.y)
				dc.Stroke()
			}
			// if there is a wall blocking the path west, draw a line from the SW to NW corners.
			if c.walls.west {
				dc.DrawLine(sw.x, sw.y, nw.x, nw.y)
				dc.Stroke()
			}
		}
	}

	// save the image as a PNG file
	err := dc.SavePNG(path)
	if err != nil {
		return err
	}

	return nil
}

// toText renders the grid using IBM box glyphs
func (g *grid) toText(path string) error {
	// define constants for the edges of the maze
	north, east, south, west := 0, g.width-1, g.height-1, 0

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

	return os.WriteFile(path, b.Bytes(), 0644)
}
