// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
)

func isMazeFullyReachable(maze [][]int, exits [][2]int) bool {
	rows := len(maze)
	cols := len(maze[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	// Directions for moving up, down, left, right
	directions := [][2]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	}

	// BFS function
	bfs := func(start [2]int) {
		queue := [][2]int{start}
		visited[start[0]][start[1]] = true

		for len(queue) > 0 {
			x, y := queue[0][0], queue[0][1]
			queue = queue[1:]

			for _, dir := range directions {
				nx, ny := x+dir[0], y+dir[1]
				if nx >= 0 && nx < rows && ny >= 0 && ny < cols && !visited[nx][ny] && maze[nx][ny] == 0 {
					visited[nx][ny] = true
					queue = append(queue, [2]int{nx, ny})
				}
			}
		}
	}

	// Run BFS from each exit
	for _, exit := range exits {
		if !visited[exit[0]][exit[1]] {
			bfs(exit)
		}
	}

	// Check if all cells are reachable
	for i := range visited {
		for j := range visited[i] {
			if maze[i][j] == 0 && !visited[i][j] {
				return false
			}
		}
	}
	return true
}

func main() {
	// Example maze (0 = path, 1 = wall)
	maze := [][]int{
		{0, 1, 0, 0, 0},
		{0, 1, 1, 1, 0},
		{0, 0, 0, 1, 0},
		{0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0},
	}

	// Exits (assuming they are the coordinates of the exits)
	exits := [][2]int{{0, 0}, {4, 4}}

	if isMazeFullyReachable(maze, exits) {
		fmt.Println("The maze is fully reachable from the exits.")
	} else {
		fmt.Println("Not all cells in the maze are reachable from the exits.")
	}
}
