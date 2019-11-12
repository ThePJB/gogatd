package main

import (
	"fmt"
	"math"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	GAMEXRES = 900
	GAMEYRES = 900
	UIH      = 300
	GRIDW    = 15
	GRIDH    = 15
)

type Context struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	atlas    Atlas

	xres, yres int32

	spawnidx, goalidx int32
	grid              []Cell
	gridw, gridh      int32
	cellw, cellh      int32

	enemies []Enemy

	lives int32
}

type Cell struct {
	cellType  CellType
	towerType TowerType
	pathDir   [2]int32
	// to come stuff about state like attack cooldown,
	// ability cooldown etc
}

type Enemy struct {
	enemyType EnemyType
	position  vec2f
	velocity  vec2f

	w, h      int32
	animstage float64
	animmax   float64
	speedBase float64 // pixels per second

	alive bool
}

type EnemyType int32

const (
	Slime EnemyType = 0
)

func (e Enemy) rect() *sdl.Rect {
	return &sdl.Rect{int32(e.position[0]) - e.w/2, int32(e.position[1]) - e.h/2, e.w, e.h}
}

type CellType int32

const (
	Path       CellType = 0
	Buildable  CellType = 1
	TowerSkull CellType = 2
	Portal     CellType = 3
	Orb        CellType = 4
	Wall       CellType = 5
	WallTop    CellType = 6
)

type TowerType int32

const (
	None TowerType = 0

	Skull TowerType = 3
	Fire  TowerType = 4
)

var context Context = Context{}

func main() {
	context.xres = GAMEXRES
	context.yres = GAMEYRES
	context.gridw = GRIDW
	context.gridh = GRIDH
	context.cellw = GAMEXRES / GRIDW
	context.cellh = GAMEYRES / GRIDH
	context.lives = 20

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	if imgflags := img.Init(img.INIT_PNG); imgflags != img.INIT_PNG {
		panic("failed to init png loading")
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("GOGATD v0.1", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		GAMEXRES, GAMEYRES+UIH, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	context.window = window
	defer context.window.Destroy()

	renderer, err := sdl.CreateRenderer(context.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	context.renderer = renderer
	context.atlas = loadAssets()

	context.grid, context.spawnidx, context.goalidx = makeGrid()

	spawnEnemy()

	running := true
	tStart := time.Now().UnixNano()
	tCurrentStart := float64(tStart) / 1000000000
	var tLastStart float64
	for running {
		tStart = time.Now().UnixNano()
		tLastStart = tCurrentStart
		tCurrentStart = float64(tStart) / 1000000000

		dt := tCurrentStart - tLastStart

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.MouseButtonEvent:
				if t.Button == sdl.BUTTON_LEFT && t.State == 1 {
					if t.Y < GAMEYRES {
						// game LMB event
						gx := t.X / context.cellw
						gy := t.Y / context.cellh

						fmt.Println("you clicked", gx, gy)
						clickedCellIdx := gy*context.gridw + gx

						fmt.Println("tower:", context.grid[clickedCellIdx].towerType)
						context.grid[clickedCellIdx].towerType = Skull
						fmt.Println("tower:", context.grid[clickedCellIdx].towerType)
					} else {
						// UI LMB event
						fmt.Println("ui clicc")
					}
				}
			}
		}
		// update
		for i, enemy := range context.enemies {
			if !enemy.alive {
				continue
			}
			// move
			context.enemies[i].position = vecAdd(enemy.position, vecMulScalar(enemy.velocity, dt))
			currentCell := context.grid[getTileFromPos(enemy.position)]
			if currentCell.cellType == Orb {
				context.lives--
				context.enemies[i].alive = false
			} else {
				// only set this near the center of the tile
				d := dist(getTileCenter(getTileFromPos(enemy.position)), enemy.position)
				if d < 5.0 {
					context.enemies[i].velocity = vecMulScalar(asF64(currentCell.pathDir), enemy.speedBase)
				}
			}

			// update animation
			context.enemies[i].animstage += dt
			if context.enemies[i].animstage > enemy.animmax {
				context.enemies[i].animstage -= enemy.animmax
			}
		}

		// render loop
		context.renderer.Clear()
		context.renderer.SetDrawColor(0, 0, 0, 255)
		context.renderer.FillRect(&sdl.Rect{0, 0, GAMEXRES, GAMEYRES})
		// draw grid
		for i, cell := range context.grid {
			toRect := &sdl.Rect{
				(int32(i) % context.gridw) * context.cellw,
				(int32(i) / context.gridw) * context.cellh,
				context.cellw, context.cellh,
			}
			switch cell.cellType {
			case Path:
				context.renderer.CopyEx(context.atlas.path, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
			case Portal:
				context.renderer.CopyEx(context.atlas.path, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
				context.renderer.CopyEx(context.atlas.portal, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
			case Orb:
				context.renderer.CopyEx(context.atlas.path, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
				context.renderer.CopyEx(context.atlas.orb, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
			case Buildable:
				context.renderer.CopyEx(context.atlas.buildable, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
				switch cell.towerType {
				case None:
					break
				case Skull:
					context.renderer.CopyEx(context.atlas.skull, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
				case Fire:
					context.renderer.CopyEx(context.atlas.fire, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
				}
			case Wall:
				context.renderer.CopyEx(context.atlas.wall, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
			case WallTop:
				context.renderer.CopyEx(context.atlas.wallTop, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
			}
		}
		// draw enemies
		for _, enemy := range context.enemies {
			if enemy.alive {
				if enemy.animstage > enemy.animmax/2 {
					context.renderer.CopyEx(context.atlas.dude, nil, enemy.rect(), 0.0, nil, sdl.FLIP_HORIZONTAL)
				} else {
					context.renderer.CopyEx(context.atlas.dude, nil, enemy.rect(), 0.0, nil, sdl.FLIP_NONE)
				}
			}
		}

		// draw some text
		drawText(10, 10, fmt.Sprintf("%.0f FPS", 1/dt), 2)
		drawText(GAMEXRES-120, 10, fmt.Sprintf("%d Lives", context.lives), 2)
		context.renderer.Present()
		tnow := time.Now().UnixNano()
		currdt := tnow - tStart
		c := 1000000000/60 - currdt
		if c > 0 {
			time.Sleep(time.Nanosecond * time.Duration(c))
		}
	}
}

func spawnEnemy() {
	speed := float64(150.0)

	newEnemy := Enemy{
		enemyType: Slime,
		w:         80,
		h:         80,
		animstage: 0,
		animmax:   0.8,
		speedBase: speed,
		position:  getTileCenter(context.spawnidx),
		velocity:  vecMulScalar(asF64(context.grid[context.spawnidx].pathDir), speed),
		alive:     true,
	}

	context.enemies = append(context.enemies, newEnemy)
}

func getTileCenter(idx int32) vec2f {
	return vec2f{(float64(idx%GRIDW) + 0.5) * float64(context.cellw), (float64(idx/GRIDH) + 0.5) * float64(context.cellh)}
}

func getTileFromPos(pos vec2f) int32 {
	ivec := asI32(pos)
	gx := ivec[0] / context.cellw
	gy := ivec[1] / context.cellh

	if gx < 0 {
		panics("Tile gotten from", pos, " is out of bounds x < 0", gx)
	} else if gx > context.gridw {
		panics("Tile gotten from", pos, " is out of bounds x > gridw", gx, context.gridw)
	} else if gy > context.gridh {
		panics("Tile gotten from", pos, " is out of bounds y > gridh", gy, context.gridh)
	} else if gy < 0 {
		panics("Tile gotten from", pos, " is out of bounds y < 0", gy)
	}

	return gy*context.gridw + gx
}

type vec2f [2]float64
type vec2i [2]int32

func vecMul(a, b [2]float64) [2]float64 {
	return [2]float64{a[0] * b[0], a[1] * b[1]}
}
func vecMulScalar(a vec2f, b float64) vec2f {
	return vec2f{a[0] * b, a[1] * b}
}
func vecAdd(a, b [2]float64) [2]float64 {
	return [2]float64{a[0] + b[0], a[1] + b[1]}
}
func asF64(a [2]int32) [2]float64 {
	return [2]float64{float64(a[0]), float64(a[1])}
}
func asI32(a vec2f) vec2i {
	return vec2i{int32(a[0]), int32(a[1])}
}
func dist(a, b vec2f) float64 {
	return math.Sqrt((a[0]-b[0])*(a[0]-b[0]) + (a[1]-b[1])*(a[1]-b[1]))
}

func panics(a ...interface{}) {
	panic(fmt.Sprint(a...))
}
