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
	UIH      = 200
	GRIDW    = 15
	GRIDH    = 15
)

type Context struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	atlas    []*sdl.Texture

	xres, yres int32

	spawnidx, goalidx int32
	grid              []Cell
	gridw, gridh      int32
	cellw, cellh      int32

	enemies []Enemy

	wave  int
	lives int32

	beams       []Beam
	projectiles []Projectile

	selectedEnemy int

	simTime    float64
	eventQueue []DoLater
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
	Path CellType = iota
	Buildable
	Portal
	Orb
	Wall
	WallTop
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
	context.selectedEnemy = -1

	initSDL()
	defer teardownSDL() // defer means we should tear down correctly in the event of a panic

	// catch signals to ensure we tear down correctly if someone ctrl Cs

	loadTextures()
	initTowerProps()

	context.eventQueue = append(context.eventQueue, DoLater{
		from: 1,
		to:   5,
		update: func(t float64) {
			var y int32 = 100
			if t < 0.1 {
				tn := t * 10
				tn = slowStop4(tn)
				y = int32(0.5 + tn*100)
			}
			if t > 0.6 {
				tn := (t - 0.6) * 2.5
				tn = slowStart4(tn)
				y = int32(0.5 + (1-tn)*100)
			}
			var alpha uint8 = 255
			if t > 0.8 {
				tn := (t - 0.8) * 5
				alpha = uint8(0.5 + 255.0*(1-tn))
				//fmt.Println(t, tn, alpha)
			}
			context.atlas[TEX_FONT].SetAlphaMod(alpha)
			w := int32(7)
			scale := int32(8)
			s := fmt.Sprintf("Wave 1")
			est_w := int32(len(s)) * w * scale
			drawText(GAMEXRES/2-est_w/2, y, s, scale)
			context.atlas[TEX_FONT].SetAlphaMod(255)
		},
	})

	context.grid, context.spawnidx, context.goalidx = makeGrid()

	tspawn := 1.0
	tsinceSpawn := 0.0

	running := true
	tStart := time.Now().UnixNano()
	tCurrentStart := float64(tStart) / 1000000000
	var tLastStart float64
	points := 100.0
	for running {
		tStart = time.Now().UnixNano()
		tLastStart = tCurrentStart
		tCurrentStart = float64(tStart) / 1000000000

		dt := tCurrentStart - tLastStart

		tsinceSpawn += dt
		if tsinceSpawn >= tspawn {
			tsinceSpawn = 0
			spawnEnemy(points, uniformChromosome())
			points += 1
		}

	OUTER:
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.MouseButtonEvent:
				if t.Button == sdl.BUTTON_LEFT && t.State == 1 {
					if t.Y < GAMEYRES {
						clickpt := sdl.Point{t.X, t.Y}
						for i := range context.enemies {
							if !context.enemies[i].alive {
								continue
							}
							r := context.enemies[i].rect()
							if clickpt.InRect(r) {
								fmt.Println("you clicked enemy", context.enemies[i])
								context.selectedEnemy = i
								continue OUTER
							}
						}
						context.selectedEnemy = -1

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
			case *sdl.KeyboardEvent:

			}
		}
		// update
		updateEnemies(dt)

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
								damage(j, props.damage, props.damageType)
							} else if props.attackType == ATTACK_PROJECTILE {
								target := context.enemies[j].position
								makeProjectile(getTileCenter(int32(i)), target, props.attackTexture, 600, func() {
									// hope closures work how i think. args vs closing over. args means we provide it at the time? or later idk
									for k := range context.enemies {
										if dist(context.enemies[k].position, target) < 100 {
											damage(j, props.damage, props.damageType)
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
			context.renderer.CopyEx(context.atlas[TEX_OFFSET_TILES+TextureID(cell.cellType)], nil, toRect, 0.0, nil, sdl.FLIP_NONE)
			if cell.tower.towerType != None {
				drawTower(cell.tower, toRect)
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

		// draw selected ui
		// a stateful cursor thing would actually probably be quite good
		if context.selectedEnemy != -1 {
			drawSelectedEnemy()
		}

		context.simTime += dt
		for i := range context.eventQueue {
			if context.eventQueue[i].from < context.simTime && context.eventQueue[i].to > context.simTime {
				t := (context.simTime - context.eventQueue[i].from) / (context.eventQueue[i].to - context.eventQueue[i].from)
				context.eventQueue[i].update(t)
			}
		}

		// draw some text
		drawText(10, 10, fmt.Sprintf("%.0f FPS", 1/dt), 2)
		drawText(10, 30, fmt.Sprintf("%d Lives", context.lives), 2)
		context.renderer.Present()
		tnow := time.Now().UnixNano()
		currdt := tnow - tStart
		c := 1000000000/60 - currdt
		if c > 0 {
			time.Sleep(time.Nanosecond * time.Duration(c))
		}
	}
}

func initSDL() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	if imgflags := img.Init(img.INIT_PNG); imgflags != img.INIT_PNG {
		panic("failed to init png loading")
	}

	window, err := sdl.CreateWindow("GOGATD v0.1", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		GAMEXRES, GAMEYRES+UIH, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	context.window = window

	renderer, err := sdl.CreateRenderer(context.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	context.renderer = renderer
}

func teardownSDL() {
	fmt.Print("Tearing down SDL...")
	for i := 0; i < int(NUM_TEXTURES); i++ {
		context.atlas[i].Destroy()
	}
	context.window.Destroy()
	context.renderer.Destroy()
	img.Quit()
	sdl.Quit()
	fmt.Println("Done")
}
