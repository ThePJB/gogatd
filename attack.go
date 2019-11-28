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
					//startX := int32(start[0] + 0.5*(end[0]-start[0]))
					//startY := int32(start[1] + 0.5*(end[1]-start[1]))
					toRect := &sdl.Rect{int32(start[0]), int32(start[1]), int32(dist(start, end)), props.bp.width}

					tex := context.atlas[props.bp.texture]
					beamAngle := angle(start, end)
					_, _, _, height, err := tex.Query()
					if err != nil {
						panic(err)
					}
					centerOfRotation := &sdl.Point{0, height / 2}

					context.tweens = append(context.tweens, Tween{
						context.simTime, context.simTime + props.bp.fadeTime, func(t float64) {
							tex.SetAlphaMod(uint8(255 * slowStop2(1-t)))
							context.renderer.CopyEx(tex, nil, toRect, RAD_TO_DEG*beamAngle, centerOfRotation, sdl.FLIP_NONE)
							tex.SetAlphaMod(uint8(255))
						},
					})
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
			context.renderer.CopyEx(context.atlas[attack.pp.texture], &ProjectileSrcRect, toRect, RAD_TO_DEG*(angle+t*attack.pp.rotationSpeed), nil, flip)
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
