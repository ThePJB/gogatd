package main

import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

type EnemyType int32

const (
	Slime EnemyType = 0
)

type Enemy struct {
	enemyType EnemyType
	position  vec2f
	velocity  vec2f

	w, h      int32
	animstage float64
	animmax   float64
	speedBase float64 // pixels per second

	res [3]float64

	r, g, b float64

	alive     bool
	splatTime float64
	hpmax     float64
	hp        float64
}

// another thing is like the global score weights
// how does normalizing it work, it needs to be fair
// actually just weight how many points each one translates into

// theres a lot of curvey relationships going on here e.g. diminishing returns and how much
// think about it later

// 0-1 weights
type Chromosome struct {
	speed      float64
	resistance [3]float64 // idx off damageType
	health     float64
}

func uniformChromosome() Chromosome {
	c := Chromosome{
		speed: rand.Float64(),
		resistance: [3]float64{
			rand.Float64(),
			rand.Float64(),
			rand.Float64(),
		},
		health: rand.Float64(),
	}

	sum := c.speed + c.health + c.resistance[0] + c.resistance[1] + c.resistance[2]
	c.speed /= sum
	c.resistance[0] /= sum
	c.resistance[1] /= sum
	c.resistance[2] /= sum
	c.health /= sum

	return c
}

// appearance generate from chromosome
// (phenotypes are represented visually)
// this is all multiplicative with points right now. mayeb you would want it to be additive idk

func spawnEnemy(points float64, c Chromosome) {
	speed := 10 + points*c.speed // * some coefficient
	hp := 10 + points*c.health

	r0 := c.resistance[0]
	r1 := c.resistance[1]
	r2 := c.resistance[2]

	//r0 := logisticFunction(0.25 * points * c.resistance[0])
	//r1 := logisticFunction(0.25 * points * c.resistance[1])
	//r2 := logisticFunction(0.25 * points * c.resistance[2])

	newEnemy := Enemy{
		enemyType: Slime,
		w:         int32(2 * hp),
		h:         int32(2 * hp),
		animstage: 0,
		animmax:   20 / speed,
		speedBase: speed,
		res:       [3]float64{r0, r1, r2},
		r:         1 - (r1+r2)/2,
		g:         1 - (r0+r2)/2,
		b:         1 - (r0+r1)/2,
		position:  getTileCenter(context.spawnidx),
		velocity:  vecMulScalar(asF64(context.grid[context.spawnidx].pathDir), speed),
		alive:     true,
		hpmax:     hp,
		hp:        hp,
	}

	context.enemies = append(context.enemies, newEnemy)
}

func (e Enemy) rect() *sdl.Rect {
	return &sdl.Rect{int32(e.position[0]) - e.w/2, int32(e.position[1]) - e.h/2, e.w, e.h}
}

func updateEnemies(dt float64) {
	for i, enemy := range context.enemies {
		if !enemy.alive {
			context.enemies[i].splatTime -= dt
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
}

func drawEnemies() {
	for _, enemy := range context.enemies {
		if enemy.alive {
			context.atlas.dude.SetColorMod(uint8(255*enemy.r), uint8(255*enemy.g), uint8(255*enemy.b))
			if enemy.animstage > enemy.animmax/2 {
				context.renderer.CopyEx(context.atlas.dude, nil, enemy.rect(), 0.0, nil, sdl.FLIP_HORIZONTAL)
			} else {
				context.renderer.CopyEx(context.atlas.dude, nil, enemy.rect(), 0.0, nil, sdl.FLIP_NONE)
			}
			if enemy.hp < enemy.hpmax {
				r := enemy.rect()
				drawHPBar(r.X, r.Y-5, r.W, int32(float64(r.W)*(enemy.hp/enemy.hpmax)))
			}
			context.atlas.dude.SetColorMod(255, 255, 255)
		} else if enemy.splatTime > 0 {
			context.atlas.splat.SetColorMod(uint8(255*enemy.r), uint8(255*enemy.g), uint8(255*enemy.b))
			if enemy.animstage > enemy.animmax/2 {
				context.renderer.CopyEx(context.atlas.splat, nil, enemy.rect(), 0.0, nil, sdl.FLIP_HORIZONTAL)
			} else {
				context.renderer.CopyEx(context.atlas.splat, nil, enemy.rect(), 0.0, nil, sdl.FLIP_NONE)
			}
			context.atlas.splat.SetColorMod(255, 255, 255)

		}
	}
}

func drawEnemyToRect(i int, r *sdl.Rect) {
	enemy := context.enemies[i]
	if enemy.alive {
		context.atlas.dude.SetColorMod(uint8(255*enemy.r), uint8(255*enemy.g), uint8(255*enemy.b))
		if enemy.animstage > enemy.animmax/2 {
			context.renderer.CopyEx(context.atlas.dude, nil, r, 0.0, nil, sdl.FLIP_HORIZONTAL)
		} else {
			context.renderer.CopyEx(context.atlas.dude, nil, r, 0.0, nil, sdl.FLIP_NONE)
		}
		if enemy.hp < enemy.hpmax {
			drawHPBar(r.X, r.Y-5, r.W, int32(float64(r.W)*(enemy.hp/enemy.hpmax)))
		}
		context.atlas.dude.SetColorMod(255, 255, 255)
	} else /* if enemy.splatTime > 0 */ {
		context.atlas.splat.SetColorMod(uint8(255*enemy.r), uint8(255*enemy.g), uint8(255*enemy.b))
		if enemy.animstage > enemy.animmax/2 {
			context.renderer.CopyEx(context.atlas.splat, nil, r, 0.0, nil, sdl.FLIP_HORIZONTAL)
		} else {
			context.renderer.CopyEx(context.atlas.splat, nil, r, 0.0, nil, sdl.FLIP_NONE)
		}
		context.atlas.splat.SetColorMod(255, 255, 255)

	}
}

func drawHPBar(x, y, w, g int32) {
	var h int32 = 6
	context.renderer.CopyEx(context.atlas.bargreen, nil, &sdl.Rect{x, y, g, h}, 0.0, nil, sdl.FLIP_NONE)
	context.renderer.CopyEx(context.atlas.barred, nil, &sdl.Rect{x + g, y, w - g, h}, 0.0, nil, sdl.FLIP_NONE)
	context.renderer.CopyEx(context.atlas.barend, nil, &sdl.Rect{x - 1, y, 2, h}, 0.0, nil, sdl.FLIP_NONE)
	context.renderer.CopyEx(context.atlas.barend, nil, &sdl.Rect{x + w - 1, y, 2, h}, 0.0, nil, sdl.FLIP_HORIZONTAL)
}

const SPLAT_TIME = 0.25

func killEnemy(i int) {
	context.enemies[i].alive = false
	context.enemies[i].splatTime = SPLAT_TIME
}
