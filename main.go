package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	GAMEXRES  = 900
	GAMEYRES  = 900
	UIH       = 200
	GRIDW     = 15
	GRIDH     = 15
	GRID_SZ_X = GAMEXRES / GRIDW
	GRID_SZ_Y = GAMEYRES / GRIDH
)

const (
	INTER_WAVE_TIME  = 3.0
	INTER_ENEMY_TIME = 1.0
)

const (
	DESIRED_ENEMIES         = 10
	ENEMY_STRENGTH_PER_WAVE = 20
)

const (
	SELECTIVE_PRESSURE = 0.5
)

const (
	FFWD_SPEED = 10
)

type GameState int

const (
	BETWEEN_WAVE GameState = iota
	IN_WAVE
)

type Context struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	atlas    []*sdl.Texture
	chunks   []*mix.Chunk

	xres, yres int32

	spawnidx, goalidx int32
	grid              []Cell
	gridw, gridh      int32
	cellw, cellh      int32
	path              []PathSegment

	parentGeneration []Enemy
	enemies          []Enemy

	lives int32

	selectedEnemy int
	selectedTower int32

	simTime float64
	tweens  []Tween
	events  []Event
	paused  bool

	waveNumber         int
	enemyStrength      float64
	state              GameState
	stateChangeTimeAcc float64

	placingTower TowerType
	money        int
}

var context Context = Context{}

var keymap []sdl.Keycode

func main() {
	keymap = []sdl.Keycode{
		sdl.K_ESCAPE,
		sdl.K_q,
		sdl.K_w,
		sdl.K_e,
		sdl.K_r,
		sdl.K_t,
		sdl.K_y,
		sdl.K_u,
	}
	rand.Seed(time.Now().UnixNano())
	context.xres = GAMEXRES
	context.yres = GAMEYRES
	context.gridw = GRIDW
	context.gridh = GRIDH
	context.cellw = GAMEXRES / GRIDW
	context.cellh = GAMEYRES / GRIDH
	context.lives = 20
	context.selectedEnemy = -1
	context.selectedTower = -1
	context.enemyStrength = 100
	context.waveNumber = 0
	context.state = 0
	context.paused = false
	context.money = 20

	context.stateChangeTimeAcc = 0
	initSDL()
	defer teardownSDL()

	loadTextures()
	loadChunks()
	initTowerProps()

	context.grid, context.spawnidx, context.goalidx = makeGrid()

	// initial chromosomes
	for i := 0; i < DESIRED_ENEMIES; i++ {
		context.parentGeneration = append(context.parentGeneration, makeEnemy(0, uniformChromosome()))
	}

	doffwd := false

	running := true
	tStart := time.Now().UnixNano()
	tCurrentStart := float64(tStart) / 1000000000
	var tLastStart float64
	var CursorX int32
	var CursorY int32
	var TileX int32
	var TileY int32
	var HoverTile int32
	for running {
		tStart = time.Now().UnixNano()
		tLastStart = tCurrentStart
		tCurrentStart = float64(tStart) / 1000000000

		dt := tCurrentStart - tLastStart
		if doffwd {
			dt *= 5
		}

		// handle input
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
						gx := t.X / context.cellw
						gy := t.Y / context.cellh
						clickedCellIdx := gy*context.gridw + gx

						if context.placingTower != None && context.grid[clickedCellIdx].cellType == Buildable && context.grid[clickedCellIdx].tower.towerType == None {
							context.selectedTower = -1
							if towerProperties[context.placingTower].cost <= context.money {
								context.money -= towerProperties[context.placingTower].cost
								context.grid[clickedCellIdx].tower = makeTower(context.placingTower)
								context.placingTower = None
							} else {
								// not enough money
								fmt.Println("insufficient money", context.money)
							}
						} else if context.grid[clickedCellIdx].tower.towerType != None {
							context.selectedTower = clickedCellIdx
							context.selectedEnemy = -1
						} else {
							context.selectedTower = -1
							// see if we clicked an enemy
							for i := range context.enemies {
								if !context.enemies[i].alive {
									continue
								}
								r := context.enemies[i].rect()
								if clickpt.InRect(r) {
									context.selectedEnemy = i
									context.selectedTower = -1
									continue OUTER
								}
							}
							context.selectedEnemy = -1
						}
					} else {
						// UI LMB event
						fmt.Println("ui clicc")
					}
				} else if t.Button == sdl.BUTTON_RIGHT {
					context.placingTower = None
				}
			case *sdl.KeyboardEvent:
				if t.Keysym.Sym == sdl.K_RIGHTBRACKET {
					if t.State == sdl.PRESSED {
						doffwd = true
					} else {
						doffwd = false
					}
				}
				if t.State == sdl.PRESSED {
					for i, k := range keymap {
						if t.Keysym.Sym == k {
							context.selectedTower = -1
							context.selectedEnemy = -1
							if context.placingTower == TowerType(i) {
								context.placingTower = None
							} else {
								context.placingTower = TowerType(i)
							}
						}
					}
				}
				if t.Keysym.Sym == sdl.K_SPACE && t.State == sdl.PRESSED {
					context.paused = !context.paused
				}
			case *sdl.MouseMotionEvent:
				CursorX = t.X
				CursorY = t.Y
				TileX = t.X / context.cellw
				TileY = t.Y / context.cellh
				HoverTile = TileY*context.gridw + TileX
			}
		}
		if !context.paused {
			context.simTime += dt

			for i := range context.events {
				if !context.events[i].done && context.simTime > context.events[i].when {
					context.events[i].done = true
					context.events[i].action()
				}
			}

			// Game progression stuff:
			context.stateChangeTimeAcc += dt
			switch context.state {
			case BETWEEN_WAVE:
				// does nothing in this state

				// transition:
				if context.stateChangeTimeAcc > INTER_WAVE_TIME {
					context.waveNumber += 1
					context.enemyStrength += ENEMY_STRENGTH_PER_WAVE
					waveAnnounce(context.waveNumber, context.simTime) // doesnt happen
					context.stateChangeTimeAcc = -2
					context.chunks[CHUNK_WAVE].Play(-1, 0)
					context.state = IN_WAVE
				}
			case IN_WAVE:
				// spawns enemy if timer < inter enemy time and number of enemies less than desired number
				if len(context.enemies) < DESIRED_ENEMIES && context.stateChangeTimeAcc > INTER_ENEMY_TIME {
					context.stateChangeTimeAcc = 0
					context.enemies = append(context.enemies, makeEnemy(context.enemyStrength, doTournamentSelection().mutate()))
				}

				// check if all enemies are dead and done splatting (if I refactored splatting into a vfx it wouldnt be necessary to check here)
				// also unselect any selection if enemy is removed
				anyLivingEnemies := false
				for i := range context.enemies {
					if context.enemies[i].alive || context.enemies[i].splatTime > 0 {
						anyLivingEnemies = true
					}
				}

				if !anyLivingEnemies && len(context.enemies) == DESIRED_ENEMIES {
					// initiate transition to next wave
					fmt.Println("wave end")
					// parent gen, enemies confused?

					context.selectedEnemy = -1
					context.parentGeneration = context.enemies

					f, err := os.OpenFile("fit.log", os.O_CREATE|os.O_APPEND|os.O_RDONLY, 0600)
					if err != nil {
						fmt.Println("error writing fit log")
					}
					for i := range context.parentGeneration {
						context.parentGeneration[i].measuredFitness = fitness(i)
					}

					fmt.Fprintf(f, "Wave %d pre sort:\t\t\t", context.waveNumber+1)
					for i := range context.parentGeneration {
						fmt.Fprintf(f, " %f ", fitness(i))
					}
					fmt.Fprintf(f, "\n")

					sort.Slice(context.parentGeneration, func(a, b int) bool {
						return context.parentGeneration[a].measuredFitness > context.parentGeneration[b].measuredFitness
					})

					fmt.Fprintf(f, "Wave %d post sort pre cull:\t", context.waveNumber+1)
					for i := range context.parentGeneration {
						fmt.Fprintf(f, " %f ", context.parentGeneration[i].measuredFitness)
					}
					fmt.Fprintf(f, "\n")
					context.parentGeneration = context.parentGeneration[:int(float64(len(context.parentGeneration))*SELECTIVE_PRESSURE)]
					fmt.Fprintf(f, "Wave %d post cull:\t\t\t", context.waveNumber+1)
					for i := range context.parentGeneration {
						fmt.Fprintf(f, " %f ", context.parentGeneration[i].measuredFitness)
					}
					fmt.Fprintf(f, "\n")
					context.enemies = []Enemy{}
					context.state = BETWEEN_WAVE
					context.stateChangeTimeAcc = 0
				}
			}

			// update
			updateEnemies(dt)

			sort.Slice(context.enemies, func(i, j int) bool {
				return context.enemies[i].distance > context.enemies[j].distance
			})

			for i := range context.grid {
				if context.grid[i].tower.towerType != None {
					context.grid[i].tower.cooldown -= dt
					tryAttack(i)
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

		// draw selected ui
		//background
		context.renderer.CopyEx(context.atlas[TEX_PATH], nil, &sdl.Rect{0, GAMEYRES, GAMEXRES, UIH}, 0, nil, sdl.FLIP_NONE)

		// a stateful cursor thing would actually probably be quite good
		if context.selectedEnemy != -1 {
			drawSelectedEnemy()
		}
		if context.selectedTower != -1 {
			drawSelectedTower()
		}

		for i := range context.tweens {
			if context.tweens[i].from < context.simTime && context.tweens[i].to > context.simTime {
				t := (context.simTime - context.tweens[i].from) / (context.tweens[i].to - context.tweens[i].from)
				context.tweens[i].update(t)
			}
		}

		// if in tower place mode, draw
		if context.placingTower != None {
			indicateRange(HoverTile, towerProperties[context.placingTower].attack.dist)
			drawTowerInfo(context.placingTower)
			var i int32
			for i = 0; i < GRIDW; i++ {
				context.renderer.FillRect(&sdl.Rect{-1 + i*GAMEXRES/GRIDW, 0, 2, GAMEYRES})
			}
			for i = 0; i < GRIDH; i++ {
				context.renderer.FillRect(&sdl.Rect{0, -1 + i*GAMEYRES/GRIDH, GAMEXRES, 2})
			}
			// draw ghost tower
			if CursorY < GAMEYRES && CursorX < GAMEXRES {
				t := context.atlas[towerProperties[context.placingTower].texture]
				t.SetAlphaMod(128)
				context.renderer.CopyEx(t, nil, getTileRect(HoverTile), 0.0, nil, sdl.FLIP_NONE)
				t.SetAlphaMod(255)
			}
		}

		if context.selectedTower == -1 && context.selectedEnemy == -1 && context.placingTower == None {
			// draw normal UI card
			drawTowerBtn("q", TOWER_SKULL, 0, 0)
			drawTowerBtn("w", TOWER_LASER, 1, 0)
			drawTowerBtn("e", TOWER_FIRE, 2, 0)
			drawTowerBtn("r", TOWER_LIGHTNING, 3, 0)
			drawTowerBtn("t", TOWER_ARROW, 4, 0)
			drawTowerBtn("y", TOWER_BLACKSMITH, 5, 0)
			drawTowerBtn("u", TOWER_TREBUCHET, 6, 0)
		}

		if context.paused {
			context.renderer.SetDrawColor(200, 200, 150, 80)
			context.renderer.FillRect(&sdl.Rect{0, 0, GAMEXRES, GAMEYRES})
			w := int32(7)
			scale := int32(8)
			s := "paused"
			est_w := int32(len(s)) * w * scale
			y := int32(500)
			drawText(GAMEXRES/2-est_w/2, y, s, scale)
			drawText(70, y+w*scale+20, "space to pause/unpause", 2)
			drawText(70, y+w*scale+w*2+20+10, "] to fast forward", 2)
		}

		// draw some text
		drawText(10, 10, fmt.Sprintf("%.0f FPS", 1/dt), 2)
		drawText(10, 30, fmt.Sprintf("%d Lives", context.lives), 2)
		context.renderer.CopyEx(context.atlas[TEX_CASH], nil, &sdl.Rect{10, 50, 20, 20}, 0, nil, sdl.FLIP_NONE)
		drawText(30, 50, fmt.Sprintf("%d", context.money), 2)

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

	if err := mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		panic(err)
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
	context.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	music, err := mix.LoadMUS("assets/Apero Hour.wav")
	if err != nil {
		panic(err)
	}
	music.FadeIn(999999, 1000)
	mix.VolumeMusic(32)
}

func teardownSDL() {
	fmt.Print("Tearing down SDL...")
	for i := 0; i < int(NUM_TEXTURES); i++ {
		context.atlas[i].Destroy()
	}
	// could clean up chunks but idk how

	context.window.Destroy()
	context.renderer.Destroy()
	img.Quit()
	mix.CloseAudio()
	mix.Quit()
	sdl.Quit()
	fmt.Println("Done")
}
