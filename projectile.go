package main

import "github.com/veandco/go-sdl2/sdl"

/*
No longer any need for the projectile struct

*/

type Projectile struct {
	start, end    vec2f
	lifetime      float64
	timeRemaining float64
	texture       *sdl.Texture
	callback      func()
	alive         bool
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
			if dist(context.enemies[k].position, end) < radius {
				damage(fromTower, k) // dont think this is working, needs work anyway
			}
		}

	}, false})
}

/*
// actual position is not stored, its just analytically determined
// as is rotation
func makeProjectile(start, end vec2f, texture *sdl.Texture, speed float64, callback func()) {
	p := Projectile{
		start, end,
		,
		dist(start, end) / speed,
		texture,
		callback,
		true,
	}
	context.projectiles = append(context.projectiles, p)
}
*/
func (p *Projectile) update(dt float64) {
	p.timeRemaining -= dt
	if p.timeRemaining <= 0 && p.alive {
		p.alive = false
		p.callback()
	}
}

func (p Projectile) draw() {
	//_, _, width, height, err := p.texture.Query()
	angle := angle(p.start, p.end)
	toRect := &sdl.Rect{int32(p.start[0] + (p.end[0]-p.start[0])*(p.lifetime-p.timeRemaining)/p.lifetime - 50),
		int32(p.start[1] + (p.end[1]-p.start[1])*(p.lifetime-p.timeRemaining)/p.lifetime - 50),
		100,
		100}
	context.renderer.CopyEx(p.texture, &ProjectileSrcRect, toRect, RAD_TO_DEG*angle, nil, sdl.FLIP_NONE)

}
