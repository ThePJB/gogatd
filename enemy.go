package main

import "github.com/veandco/go-sdl2/sdl"

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

	alive bool
	hpmax float64
	hp    float64
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
		hpmax:     50,
		hp:        50,
	}

	context.enemies = append(context.enemies, newEnemy)
}

func (e Enemy) rect() *sdl.Rect {
	return &sdl.Rect{int32(e.position[0]) - e.w/2, int32(e.position[1]) - e.h/2, e.w, e.h}
}

func drawEnemies() {
	for _, enemy := range context.enemies {
		if enemy.alive {
			if enemy.animstage > enemy.animmax/2 {
				context.renderer.CopyEx(context.atlas.dude, nil, enemy.rect(), 0.0, nil, sdl.FLIP_HORIZONTAL)
			} else {
				context.renderer.CopyEx(context.atlas.dude, nil, enemy.rect(), 0.0, nil, sdl.FLIP_NONE)
			}
			if enemy.hp < enemy.hpmax {
				r := enemy.rect()
				drawHPBar(r.X, r.Y-5, r.W, int32(float64(r.W)*(enemy.hp/enemy.hpmax)))
			}
		}
	}
}

func drawHPBar(x, y, w, g int32) {
	var h int32 = 6
	context.renderer.CopyEx(context.atlas.bargreen, nil, &sdl.Rect{x, y, g, h}, 0.0, nil, sdl.FLIP_NONE)
	context.renderer.CopyEx(context.atlas.barred, nil, &sdl.Rect{x + g, y, w - g, h}, 0.0, nil, sdl.FLIP_NONE)
	context.renderer.CopyEx(context.atlas.barend, nil, &sdl.Rect{x - 1, y, 2, h}, 0.0, nil, sdl.FLIP_NONE)
	context.renderer.CopyEx(context.atlas.barend, nil, &sdl.Rect{x + w - 1, y, 2, h}, 0.0, nil, sdl.FLIP_HORIZONTAL)
}
