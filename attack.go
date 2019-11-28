package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type AttackType int32

const (
	ATTACK_BEAM AttackType = iota
	ATTACK_PROJECTILE
)

type Attack struct {
	attackType  AttackType
	attackSound ChunkID
	bp          BeamProperties
	pp          ProjectileProperties

	damage     float64
	damageType DamageType
	dist       float64
	cooldown   float64
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
	deathScale        float64
	deathSound        ChunkID
	deathFlipInterval float64
}

func tryAttack(i int) {
	props := towerProperties[context.grid[i].tower.towerType].attack
	if context.grid[i].tower.cooldown <= 0 {
		for j := range context.enemies {
			if !context.enemies[j].alive {
				continue
			}
			start := getTileCenter(int32(i))
			end := context.enemies[j].position
			if dist(start, end) < props.dist {
				context.chunks[props.attackSound].Play(-1, 0)
				if props.attackType == ATTACK_BEAM {
					//makeBeam(props.BeamProperties, getTileCenter(int32(i)), context.enemies[j].position, 0.4)
					context.chunks[props.attackSound].Play(-1, 0)
					damage(i, j)
				} else if props.attackType == ATTACK_PROJECTILE {
					tt := dist(start, end) / props.pp.speed
					if props.pp.lead {
						// calculate where enemy will be
						d := context.enemies[j].distance + context.enemies[j].speedBase*tt
						end = pathPos(d)
					}
					makeProjectile(start, end, props, i, context.enemies[j].uid)
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
func makeProjectile(start, end vec2f, attack Attack, fromTower int, toEnemyUID int) {
	//texture *sdl.Texture, speed float64, radius float64, fromTower int, sound ChunkID) {
	tt := dist(start, end) / attack.pp.speed
	// Travelling projectile
	context.tweens = append(context.tweens, Tween{
		context.simTime, context.simTime + tt, func(t float64) {
			flip := sdl.FLIP_NONE
			if math.Mod(t, 2*attack.pp.flipInterval) > attack.pp.flipInterval {
				flip = sdl.FLIP_VERTICAL
			}
			angle := angle(start, end)
			toRect := &sdl.Rect{int32(start[0] + (end[0]-start[0])*(t) - (50 * attack.pp.scale)),
				int32(start[1] + (end[1]-start[1])*(t) - (50 * attack.pp.scale)),
				int32(100.0 * attack.pp.scale),
				int32(100.0 * attack.pp.scale)}
			context.renderer.CopyEx(context.atlas[attack.pp.texture], &ProjectileSrcRect, toRect, RAD_TO_DEG*angle+t*attack.pp.rotationSpeed, nil, flip)
		}})

	// After explodey projectile
	context.tweens = append(context.tweens, Tween{
		context.simTime + tt, context.simTime + tt + 0.4, func(t float64) {
			flip := sdl.FLIP_NONE
			if math.Mod(t, 2*attack.pp.deathFlipInterval) > attack.pp.flipInterval {
				flip = sdl.FLIP_VERTICAL
			}
			angle := angle(start, end)
			toRect := &sdl.Rect{int32(end[0] - 50*attack.pp.deathScale),
				int32(end[1] - 50*attack.pp.deathScale),
				int32(100.0 * attack.pp.deathScale),
				int32(100.0 * attack.pp.deathScale)}
			tex := context.atlas[attack.pp.deathTexture]
			tex.SetAlphaMod(uint8(255 * slowStop2(1-t)))
			context.renderer.CopyEx(tex, &ProjectileDeathSrcRect, toRect, RAD_TO_DEG*angle, nil, flip)
			tex.SetAlphaMod(uint8(255))
		}})

	// Game impact
	context.events = append(context.events, Event{context.simTime + tt, func() {
		context.chunks[attack.pp.deathSound].Play(-1, 0)
		for k := range context.enemies {
			if !context.enemies[k].alive {
				continue
			}
			if attack.pp.area == 0 {
				if context.enemies[k].uid == toEnemyUID {
					damage(fromTower, k)
					break
				}
			} else {
			}
			if dist(context.enemies[k].position, end) < attack.pp.area {
				damage(fromTower, k)
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
