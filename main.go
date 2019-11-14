package main

import (
	"fmt"
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

	beams       []Beam
	projectiles []Projectile
}

type Cell struct {
	cellType CellType
	tower    Tower
	pathDir  [2]int32

	// to come stuff about state like attack cooldown,
	// ability cooldown etc
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
	initTowerProps()

	context.grid, context.spawnidx, context.goalidx = makeGrid()

	tspawn := 1.0
	tsinceSpawn := 0.0

	running := true
	tStart := time.Now().UnixNano()
	tCurrentStart := float64(tStart) / 1000000000
	var tLastStart float64
	for running {
		tStart = time.Now().UnixNano()
		tLastStart = tCurrentStart
		tCurrentStart = float64(tStart) / 1000000000

		dt := tCurrentStart - tLastStart

		tsinceSpawn += dt
		if tsinceSpawn >= tspawn {
			tsinceSpawn = 0
			spawnEnemy()
		}

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

						context.grid[clickedCellIdx].tower = makeTower((context.grid[clickedCellIdx].tower.towerType + 1) % NUM_TOWERS)
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

		for i := range context.grid {
			if context.grid[i].tower.towerType != None {
				props := towerProperties[context.grid[i].tower.towerType]
				context.grid[i].tower.cooldown -= dt
				if context.grid[i].tower.cooldown <= 0 {
					// we can attack: look for targets
					// at the oment just pick 1st enemy
					for j := range context.enemies {
						if !context.enemies[j].alive {
							continue
						}
						if dist(context.enemies[j].position, getTileCenter(int32(i))) < props.attackRange {
							// found an enemy
							// probably factor into a damage function eventually that accounts for attack,res and handles death etc
							if props.attackType == ATTACK_BEAM {
								makeBeam(props.attackTexture, getTileCenter(int32(i)), context.enemies[j].position, 0.4)
								context.enemies[j].hp -= props.damage
								if context.enemies[j].hp <= 0 {
									context.enemies[j].alive = false
								}
							} else if props.attackType == ATTACK_PROJECTILE {
								target := context.enemies[j].position
								makeProjectile(getTileCenter(int32(i)), target, props.attackTexture, 1000, func() {
									// hope closures work how i think. args vs closing over. args means we provide it at the time? or later idk
									for k := range context.enemies {
										if dist(context.enemies[k].position, target) < 50 {
											context.enemies[j].hp -= props.damage
											if context.enemies[j].hp <= 0 {
												context.enemies[j].alive = false
											}
										}
									}
								})
							}
							context.grid[i].tower.cooldown = props.cooldown
							// you would play a sound or something too
							break // could have multishot
						}
					}
				}
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
				drawTower(cell.tower, toRect)

			case Wall:
				context.renderer.CopyEx(context.atlas.wall, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
			case WallTop:
				context.renderer.CopyEx(context.atlas.wallTop, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
			}
		}
		// draw enemies
		drawEnemies()

		for i := range context.beams {
			if context.beams[i].timeRemaining > 0 {
				context.beams[i].update(dt)
				context.beams[i].draw()
			}
		}
		for i := range context.projectiles {
			if context.projectiles[i].timeRemaining > 0 {
				context.projectiles[i].update(dt)
				context.projectiles[i].draw()
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
