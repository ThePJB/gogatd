package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type AttackType int32

const (
	ATTACK_BEAM AttackType = iota
	ATTACK_PROJECTILE_AOE
	ATTACK_PROJECTILE_ACCURATE
)

type Attack struct {
	attackType  AttackType
	attackSound ChunkID
	bp          BeamProperties
	pp          ProjectileProperties

	damage   float64
	dist     float64
	cooldown float64
}

type BeamProperties struct {
	texture  TextureID
	fadeTime float64
	width    int32
}

type ProjectileProperties struct {
	speed float64
	area  float64
	lead  bool

	texture       TextureID
	scale         float64
	flipInterval  float64
	rotationSpeed float64

	deathTexture      TextureID
	deathSound        ChunkID
	deathFade         float64
	deathFlipInterval float64
}

// make it attack the furthest forward, that would be less annoying
func tryAttack(i int) {
	props := towerProperties[context.grid[i].tower.towerType]
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
					context.chunks[props.attackSound].Play(-1, 0)
					damage(i, j)
				} else if props.attackType == ATTACK_PROJECTILE_AOE {
					target := context.enemies[j].position
					context.chunks[props.attackSound].Play(-1, 0)
					makeProjectileAoE(getTileCenter(int32(i)), target, props.attackTexture, 600, 80, i, props.attackLandSound)
				}
				context.grid[i].tower.cooldown = props.cooldown
				// you would play a sound or something too
				break // could have multishot
			}
		}
	}
}

var ProjectileSrcRect = sdl.Rect{0, 0, 100, 100}
var ProjectileDeathSrcRect = sdl.Rect{100, 0, 100, 100}

// right now they can all just look the same curve wise
func makeProjectileAoE(start, end vec2f, texture *sdl.Texture, speed float64, radius float64, fromTower int, sound ChunkID) {
	tt := dist(start, end) / speed
	// Travelling projectile
	context.tweens = append(context.tweens, Tween{
		context.simTime, context.simTime + tt, func(t float64) {
			flip := sdl.FLIP_NONE
			if int(t*5)%2 == 0 {
				flip = sdl.FLIP_VERTICAL
			}
			angle := angle(start, end)
			toRect := &sdl.Rect{int32(start[0] + (end[0]-start[0])*(t) - 50),
				int32(start[1] + (end[1]-start[1])*(t) - 50),
				100,
				100}
			context.renderer.CopyEx(texture, &ProjectileSrcRect, toRect, RAD_TO_DEG*angle, nil, flip)
		}})

	// After explodey projectile
	context.tweens = append(context.tweens, Tween{
		context.simTime + tt, context.simTime + tt + 0.4, func(t float64) {
			flip := sdl.FLIP_NONE
			if int(t*5)%2 == 0 {
				flip = sdl.FLIP_HORIZONTAL
			}
			angle := angle(start, end)
			toRect := &sdl.Rect{int32(end[0] - 50),
				int32(end[1] - 50),
				100,
				100}
			texture.SetAlphaMod(uint8(255 * slowStop2(1-t)))
			context.renderer.CopyEx(texture, &ProjectileDeathSrcRect, toRect, RAD_TO_DEG*angle, nil, flip)
			texture.SetAlphaMod(uint8(255))
		}})

	// Game impact
	context.events = append(context.events, Event{context.simTime + tt, func() {
		context.chunks[sound].Play(-1, 0)
		for k := range context.enemies {
			if !context.enemies[k].alive {
				continue
			}
			if dist(context.enemies[k].position, end) < radius {
				damage(fromTower, k) // dont think this is working, needs work anyway
			}
		}

	}, false})
}

/*
could expose fade as a parameter
make it look less shit one day
*/

type Beam struct {
	texture       *sdl.Texture
	start, end    vec2f
	lifetime      float64
	timeRemaining float64
}

func makeBeam(texture *sdl.Texture, start, end vec2f, lifetime float64) {
	context.beams = append(context.beams, Beam{texture, start, end, lifetime, lifetime})
}

func (b *Beam) update(dt float64) {
	// dec and remove
	b.timeRemaining -= dt
}

// width should really = height
func (b Beam) draw() {
	_, _, width, height, err := b.texture.Query()
	if err != nil {
		panic(err)
	}
	angle := angle(b.start, b.end)
	length := dist(b.start, b.end)
	nCopies := length / float64(width)
	var i int32 = 0
	alpha := b.timeRemaining / b.lifetime
	b.texture.SetAlphaMod(uint8(255 * alpha))
	for i = 0; i < int32(nCopies); i++ {
		// probably require offset
		toRect := &sdl.Rect{
			int32(b.start[0] + float64(i*width)*math.Cos(angle)),
			int32(b.start[1] + float64(i*width)*math.Sin(angle)),
			width,
			height,
		}
		context.renderer.CopyEx(b.texture, nil, toRect, RAD_TO_DEG*angle, nil, sdl.FLIP_NONE)
	}
	b.texture.SetAlphaMod(255)
	/*
		// this will probably be fucked but lets just see
		// expecting the last one to be cooked and stretched or anything
		remainingw := int32(nCopies - float64(int(nCopies))*float64(width))

		toRect := &sdl.Rect{
			int32(b.start[0] + float64(nCopies*float64(width))*math.Cos(angle)),
			int32(b.start[1] + float64(nCopies*float64(width))*math.Sin(angle)),
			remainingw,
			height,
		}
		context.renderer.CopyEx(b.texture, nil, toRect, DEG_TO_RAD*angle, nil, sdl.FLIP_NONE)
	*/
	// draw
}
