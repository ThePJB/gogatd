package main

import (
	"fmt"
	"math/rand"
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

const (
	INTER_WAVE_TIME  = 5.0
	INTER_ENEMY_TIME = 1.0
)

const (
	DESIRED_ENEMIES         = 20
	ENEMY_STRENGTH_PER_WAVE = 40
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

	xres, yres int32

	spawnidx, goalidx int32
	grid              []Cell
	gridw, gridh      int32
	cellw, cellh      int32

	parentGeneration []Enemy
	enemies          []Enemy

	lives int32

	beams       []Beam
	projectiles []Projectile

	selectedEnemy int
	selectedTower int32

	simTime float64
	tweens  []Tween
	events  []Event

	waveNumber         int
	enemyStrength      float64
	state              GameState
	stateChangeTimeAcc float64

	placingTower TowerType
	money        int
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

var keymap []sdl.Keycode

func main() {
	keymap = []sdl.Keycode{
		sdl.K_ESCAPE,
		sdl.K_q,
		sdl.K_w,
		sdl.K_e,
		sdl.K_r,
		sdl.K_t,
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
	context.enemyStrength = 60
	context.waveNumber = 0
	context.state = 0

	context.money = 20

	context.stateChangeTimeAcc = 0
	initSDL()
	defer teardownSDL() // defer means we should tear down correctly in the event of a panic

	// catch signals to ensure we tear down correctly if someone ctrl Cs

	loadTextures()
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
		context.simTime += dt

		for i := range context.events {
			if !context.events[i].done && context.simTime > context.events[i].when {
				context.events[i].done = true
				context.events[i].action()
			}
		}

		//fmt.Println(context.stateChangeTimeAcc, context.state, len(context.enemies), len(context.parentGeneration))

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
				context.stateChangeTimeAcc = 99999999999          // spawn enemy immediately
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
				context.selectedEnemy = -1
				context.parentGeneration = context.enemies
				context.enemies = []Enemy{}
				context.state = BETWEEN_WAVE
				context.stateChangeTimeAcc = 0
			}
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

						if context.placingTower != None && context.grid[clickedCellIdx].cellType == Buildable {
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
						} else {
							context.selectedTower = -1
							// see if we clicked an enemy
							for i := range context.enemies {
								if !context.enemies[i].alive {
									continue
								}
								r := context.enemies[i].rect()
								if clickpt.InRect(r) {
									fmt.Println("you clicked enemy", context.enemies[i])
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
							if context.placingTower == TowerType(i) {
								context.placingTower = None
							} else {
								context.placingTower = TowerType(i)
							}
						}
					}
				}
			case *sdl.MouseMotionEvent:
				CursorX = t.X
				CursorY = t.Y
				TileX = t.X / context.cellw
				TileY = t.Y / context.cellh
				HoverTile = TileY*context.gridw + TileX
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
								damage(i, j)
							} else if props.attackType == ATTACK_PROJECTILE {
								target := context.enemies[j].position
								makeProjectileAoE(getTileCenter(int32(i)), target, props.attackTexture, 600, 100, i)
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
			context.renderer.SetDrawColor(255, 255, 255, 64)
			context.renderer.FillRect(&sdl.Rect{0, 0, GAMEXRES, GAMEYRES})
			context.renderer.SetDrawColor(200, 200, 200, 255)
			var i int32
			for i = 0; i < GRIDW; i++ {
				context.renderer.FillRect(&sdl.Rect{-1 + i*GAMEXRES/GRIDW, 0, 2, GAMEYRES})
			}
			for i = 0; i < GRIDH; i++ {
				context.renderer.FillRect(&sdl.Rect{0, -1 + i*GAMEYRES/GRIDH, GAMEXRES, 2})
			}
			// draw ghost tower
			if CursorY < GAMEYRES && CursorX < GAMEXRES {
				towerProperties[context.placingTower].texture.SetAlphaMod(128)
				context.renderer.CopyEx(towerProperties[context.placingTower].texture, nil, getTileRect(HoverTile), 0.0, nil, sdl.FLIP_NONE)
				towerProperties[context.placingTower].texture.SetAlphaMod(255)
			}

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
