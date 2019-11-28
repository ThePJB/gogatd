package main

import "fmt"

type PathSegment struct {
	start vec2f
	end   vec2f
}

var TOTAL_DISTANCE = 0

func makeGrid() ([]Cell, int32, int32) {
	grid := []Cell{}
	levelString := `
		ttttttttttttPtt
		twwwwwwwwwwwpwt
		t...........p.t
		tpppppppppppp.t
		tp............t
		tp............t
		tppppppppppp..t
		t..........p..t
		t..........p..t
		tppppp.....p..t
		tp...p.....p..t
		tp...p.....p..t
		tp...p.....p..t
		tp...ppppppp..t
		tOttttttttttttt`

	var PortalIdx int32 = -1
	var OrbIdx int32 = -1
	for _, c := range levelString {
		if c == '.' {
			grid = append(grid, Cell{cellType: Buildable})
		} else if c == 'p' {
			grid = append(grid, Cell{cellType: Path})
		} else if c == 't' {
			grid = append(grid, Cell{cellType: WallTop})
		} else if c == 'w' {
			grid = append(grid, Cell{cellType: Wall})
		} else if c == 'P' {
			grid = append(grid, Cell{cellType: Portal})
			if PortalIdx != -1 {
				panic("duplicate portals")
			}
			PortalIdx = int32(len(grid) - 1)
		} else if c == 'O' {
			grid = append(grid, Cell{cellType: Orb})
			if OrbIdx != -1 {
				panic("duplicate orbs")
			}
			OrbIdx = int32(len(grid) - 1)
			fmt.Println("setting orb to", len(grid), OrbIdx)
		} else {
			// its a comment or whitespace
		}
	}

	if OrbIdx == -1 || PortalIdx == -1 {
		panic("either portal or orb wasn't set")
	}

	// It pathfinding time
	// just do something greedy, start at one end and find the neighbour that hasnt already been found
	currentIndex := int32(PortalIdx)
OUTER:
	for currentIndex != OrbIdx {
		println("\n\nbegin loop, cidx =", currentIndex)
		neighbours := GetNeighbours(GRIDW, GRIDH, currentIndex)
		for _, n := range neighbours {
			fmt.Println(grid[n])
			x, y := idxToCoords(n)
			fmt.Println("inspecting", x, y)
			if grid[n].cellType == Path || grid[n].cellType == Orb {
				fmt.Println("found candidate")
				if grid[n].pathDir[0] == 0 && grid[n].pathDir[1] == 0 {
					gx := currentIndex % GRIDW
					gy := currentIndex / GRIDW
					nx := n % GRIDW
					ny := n / GRIDW
					// wrong
					grid[n].pathDir[0] = gx - nx
					grid[n].pathDir[1] = gy - ny
					start := getTileCenter(currentIndex)
					currentIndex = n
					end := getTileCenter(n)
					context.path = append(context.path, PathSegment{start, end})
					TOTAL_DISTANCE += (GAMEXRES / GRIDW)
					println("found pathy neighbour at", n, "setting currentIdx to that")
					continue OUTER
				} else {
					fmt.Println("failed path")
					fmt.Println(grid[n])
				}
			}
		}
		x, y := idxToCoords(currentIndex)
		fmt.Println("At: ", x, y, "Neighbours:", neighbours)
		fmt.Println(idxToCoords(currentIndex))
		for _, n := range neighbours {
			x, y := idxToCoords(n)
			fmt.Println("neigh index:", x, y, "neigh cell type:", grid[n].cellType)
		}
		panic("couldnt find a path neighbour")
	}

	fmt.Println(context.path)
	fmt.Println(len(context.path))

	return grid, PortalIdx, OrbIdx
}

func idxToCoords(idx int32) (int32, int32) {
	return idx % GRIDW, idx / GRIDW
}

func GetNeighbours(gridw, gridh, idx int32) []int32 {
	// 4 connected
	gridlen := gridw * gridh
	ret := []int32{}
	if idx+gridw <= gridlen {
		ret = append(ret, idx+gridw)
	}
	if idx-gridw >= 0 {
		ret = append(ret, idx-gridw)
	}
	if gridlen%(idx+1) != 0 {
		ret = append(ret, idx+1)
	}
	if gridlen%(idx) != 0 {
		ret = append(ret, idx-1)
	}
	return ret
}
