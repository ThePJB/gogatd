package main

import (
	"math"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

type EnemyType int32

// fitness coeffs
// normalize == tile px (60) * num tiles (50?)

const (
	Slime EnemyType = 0
)

type Enemy struct {
	enemyType EnemyType
	position  vec2f
	velocity  vec2f
	distance  float64

	w, h      int32
	animstage float64
	animmax   float64
	speedBase float64 // pixels per second

	res [4]float64

	r, g, b float64

	alive     bool
	splatTime float64
	hpmax     float64
	hp        float64
	regen     float64

	chromosome      Chromosome
	measuredFitness float64
	fitness         FitnessMetric
}

type FitnessMetric struct {
	damageReduced float64
	distance      float64
	madeToEnd     float64
}

// another thing is like the global score weights
// how does normalizing it work, it needs to be fair
// actually just weight how many points each one translates into

// theres a lot of curvey relationships going on here e.g. diminishing returns and how much
// think about it later

// 0-1 weights
type Chromosome struct {
	speed      float64
	resistance [4]float64 // idx off damageType
	health     float64
	regen      float64
}

func scoreFitness(f FitnessMetric) float64 {
	//fmt.Println("f", f)
	//fitness := f.damageReduced*FIT_RES_COEFF + f.distance*FIT_DISTANCE_COEFF + f.madeToEnd*FIT_END_COEFF + FIT_DISTANCE_AND_RES_COEFF*f.distance*f.damageReduced
	//fmt.Println(fitness)
	//10000 * fitness
	return f.distance * FIT_DISTANCE_COEFF
}

func uniformChromosome() Chromosome {
	c := Chromosome{
		speed: rand.Float64(),
		resistance: [4]float64{
			rand.Float64(),
			rand.Float64(),
			rand.Float64(),
			rand.Float64(),
		},
		health: rand.Float64(),
		regen:  rand.Float64(),
	}

	return c.norm()
}

func (c Chromosome) norm() Chromosome {
	sum := c.speed + c.health + c.resistance[0] + c.resistance[1] + c.resistance[2] + c.resistance[3]
	c.speed /= sum
	c.resistance[0] /= sum
	c.resistance[1] /= sum
	c.resistance[2] /= sum
	c.resistance[3] /= sum
	c.health /= sum
	c.regen /= sum

	return c
}

// appearance generate from chromosome
// (phenotypes are represented visually)
// this is all multiplicative with points right now. mayeb you would want it to be additive idk

func makeEnemy(points float64, c Chromosome) Enemy {
	speed := 20 + points*c.speed*c.speed*2 // * some coefficient
	hp := 10 + points*c.health*0.5

	r0 := c.resistance[0]
	r1 := c.resistance[1]
	r2 := c.resistance[2]
	r3 := c.resistance[3]

	//r0 := logisticFunction(0.25 * points * c.resistance[0])
	//r1 := logisticFunction(0.25 * points * c.resistance[1])
	//r2 := logisticFunction(0.25 * points * c.resistance[2])

	newEnemy := Enemy{
		enemyType: Slime,
		w:         int32(math.Sqrt(200 * hp)),
		h:         int32(math.Sqrt(200 * hp)),
		animstage: 0,
		animmax:   20 / speed,
		speedBase: speed,
		res:       [4]float64{r0, r1, r2, r3},
		r:         1 - (r1+r2)/2,
		g:         1 - (r0+r2)/2,
		b:         1 - (r0+r1)/2,
		position:  getTileCenter(context.spawnidx),
		velocity:  vecMulScalar(asF64(context.grid[context.spawnidx].pathDir), speed),
		alive:     true,
		hpmax:     hp,
		hp:        hp,
		regen:     points * c.regen * 0.01,

		chromosome:      c,
		measuredFitness: 0,
		fitness:         FitnessMetric{},
	}
	return newEnemy
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
		// health regen
		context.enemies[i].hp += context.enemies[i].regen * dt
		if context.enemies[i].hp > context.enemies[i].hpmax {
			context.enemies[i].hp = context.enemies[i].hpmax
		}

		// move
		context.enemies[i].distance += enemy.speedBase * dt
		d := context.enemies[i].distance
		for j := range context.path {
			d -= float64(GRID_SZ_X)
			if d <= float64(GRID_SZ_X) {
				// we belong to this segment
				context.enemies[i].position = interp(context.path[j].start, context.path[j].end, d/float64(GRID_SZ_X))
				break
			}
		}

		//context.enemies[i].position = vecAdd(enemy.position, vecMulScalar(enemy.velocity, dt))
		context.enemies[i].fitness.distance += enemy.speedBase * dt

		currentCell := context.grid[getTileFromPos(enemy.position)]
		if currentCell.cellType == Orb {
			context.lives--
			context.enemies[i].fitness.madeToEnd = 1
			killEnemy(i)
			context.chunks[CHUNK_LEAK].Play(-1, 0)
		} else {
			/*
				selectedTile := getTileFromPos(enemy.position)
				towardsCenter := vecSub(getTileCenter(selectedTile), enemy.position)
				x := cross2d(rot90(enemy.velocity.unit()), towardsCenter)

				// only set this near the center of the tile
				if x > -10.0 {
					context.enemies[i].velocity = vecMulScalar(asF64(currentCell.pathDir), enemy.speedBase)
				}
			*/
		}

		// update animation
		context.enemies[i].animstage += dt
		if context.enemies[i].animstage > enemy.animmax {
			context.enemies[i].animstage -= enemy.animmax
		}
	}
}

func drawEnemies() {
	for i := len(context.enemies) - 1; i >= 0; i-- {
		enemy := context.enemies[i]
		if enemy.alive {
			context.atlas[TEX_DUDE].SetColorMod(uint8(255*enemy.r), uint8(255*enemy.g), uint8(255*enemy.b))
			if enemy.animstage > enemy.animmax/2 {
				context.renderer.CopyEx(context.atlas[TEX_DUDE], nil, enemy.rect(), 0.0, nil, sdl.FLIP_HORIZONTAL)
			} else {
				context.renderer.CopyEx(context.atlas[TEX_DUDE], nil, enemy.rect(), 0.0, nil, sdl.FLIP_NONE)
			}
			if enemy.hp < enemy.hpmax {
				r := enemy.rect()
				drawHPBar(r.X, r.Y-5, r.W, int32(float64(r.W)*(enemy.hp/enemy.hpmax)))
			}
			context.atlas[TEX_DUDE].SetColorMod(255, 255, 255)
		} else if enemy.splatTime > 0 {
			context.atlas[TEX_SPLAT].SetColorMod(uint8(255*enemy.r), uint8(255*enemy.g), uint8(255*enemy.b))
			if enemy.animstage > enemy.animmax/2 {
				context.renderer.CopyEx(context.atlas[TEX_SPLAT], nil, enemy.rect(), 0.0, nil, sdl.FLIP_HORIZONTAL)
			} else {
				context.renderer.CopyEx(context.atlas[TEX_SPLAT], nil, enemy.rect(), 0.0, nil, sdl.FLIP_NONE)
			}
			context.atlas[TEX_SPLAT].SetColorMod(255, 255, 255)

		}
	}
}

func drawEnemyToRect(i int, r *sdl.Rect) {
	enemy := context.enemies[i]
	if enemy.alive {
		context.atlas[TEX_DUDE].SetColorMod(uint8(255*enemy.r), uint8(255*enemy.g), uint8(255*enemy.b))
		if enemy.animstage > enemy.animmax/2 {
			context.renderer.CopyEx(context.atlas[TEX_DUDE], nil, r, 0.0, nil, sdl.FLIP_HORIZONTAL)
		} else {
			context.renderer.CopyEx(context.atlas[TEX_DUDE], nil, r, 0.0, nil, sdl.FLIP_NONE)
		}
		if enemy.hp < enemy.hpmax {
			drawHPBar(r.X, r.Y-5, r.W, int32(float64(r.W)*(enemy.hp/enemy.hpmax)))
		}
		context.atlas[TEX_DUDE].SetColorMod(255, 255, 255)
	} else /* if enemy.splatTime > 0 */ {
		context.atlas[TEX_SPLAT].SetColorMod(uint8(255*enemy.r), uint8(255*enemy.g), uint8(255*enemy.b))
		if enemy.animstage > enemy.animmax/2 {
			context.renderer.CopyEx(context.atlas[TEX_SPLAT], nil, r, 0.0, nil, sdl.FLIP_HORIZONTAL)
		} else {
			context.renderer.CopyEx(context.atlas[TEX_SPLAT], nil, r, 0.0, nil, sdl.FLIP_NONE)
		}
		context.atlas[TEX_SPLAT].SetColorMod(255, 255, 255)

	}
}

func drawHPBar(x, y, w, g int32) {
	var h int32 = 6
	context.renderer.CopyEx(context.atlas[TEX_BARGREEN], nil, &sdl.Rect{x, y, g, h}, 0.0, nil, sdl.FLIP_NONE)
	context.renderer.CopyEx(context.atlas[TEX_BARRED], nil, &sdl.Rect{x + g, y, w - g, h}, 0.0, nil, sdl.FLIP_NONE)
	context.renderer.CopyEx(context.atlas[TEX_BAREND], nil, &sdl.Rect{x - 1, y, 2, h}, 0.0, nil, sdl.FLIP_NONE)
	context.renderer.CopyEx(context.atlas[TEX_BAREND], nil, &sdl.Rect{x + w - 1, y, 2, h}, 0.0, nil, sdl.FLIP_HORIZONTAL)
}

const SPLAT_TIME = 0.25

func killEnemy(i int) {
	// what is dead may never die
	if context.enemies[i].alive {
		context.enemies[i].alive = false
		context.enemies[i].splatTime = SPLAT_TIME
		context.money += 1
		context.chunks[CHUNK_DIE].Play(-1, 0)
	}
}

func doTournamentSelection() Chromosome {
	if len(context.parentGeneration) == 0 || len(context.parentGeneration) == 1 {
		panic("Tried to do selection with insufficient parent generation")
	}
	var a, b int
	for {
		a = rand.Intn(len(context.parentGeneration))
		b = rand.Intn(len(context.parentGeneration))
		if a != b {
			break
		}
	}

	scoreA := fitness(a)
	scoreB := fitness(b)
	if scoreA > scoreB {
		return context.parentGeneration[a].chromosome
	} else {
		return context.parentGeneration[b].chromosome
	}
}

func (c Chromosome) mutate() Chromosome {
	mutationRate := 0.5
	c.health += c.health * mutationRate * 2 * (rand.Float64() - 0.5)
	c.speed += c.speed * mutationRate * 2 * (rand.Float64() - 0.5)
	c.resistance[0] += c.resistance[0] * mutationRate * 2 * (rand.Float64() - 0.5)
	c.resistance[1] += c.resistance[1] * mutationRate * 2 * (rand.Float64() - 0.5)
	c.resistance[2] += c.resistance[2] * mutationRate * 2 * (rand.Float64() - 0.5)
	return c.norm()
}

// tower enemy, careful, its bad api design lol
func damage(towerIdx int, enemyIdx int) {
	amount := towerProperties[context.grid[towerIdx].tower.towerType].damage
	damageType := towerProperties[context.grid[towerIdx].tower.towerType].damageType
	damageAfterRes := amount * (1 - context.enemies[enemyIdx].res[damageType])
	damageBlocked := amount * context.enemies[enemyIdx].res[damageType]
	context.enemies[enemyIdx].hp -= damageAfterRes
	context.enemies[enemyIdx].fitness.damageReduced += damageBlocked / context.enemies[enemyIdx].hpmax // normalize to hp max or not? i reckon so

	//fmt.Printf("Attack enemy with raw damage %.0f of type %s, enemy res to that is %.2f, so damage that gets through is %.1f\n", amount, damageNames[damageType], context.enemies[enemyIdx].res[damageType], damageAfterRes)

	if context.enemies[enemyIdx].hp <= 0 {
		killEnemy(enemyIdx)
		context.grid[towerIdx].tower.kills += 1
	}
}

// this is actually the good way to do it
func fitness(idx int) float64 {
	fit := 0.0
	if idx != 0 {
		fit += context.parentGeneration[idx].distance - context.parentGeneration[idx-1].distance
	}
	if idx != len(context.parentGeneration)-1 {
		fit += context.parentGeneration[idx].distance - context.parentGeneration[idx+1].distance
	}
	return fit
}
