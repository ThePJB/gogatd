package main

import "github.com/veandco/go-sdl2/sdl"

/*
In an effort to keep this also somewhat decoupled as beam is, it could just call a callback when its done
*/

type Projectile struct {
	start, end    vec2f
	lifetime      float64
	timeRemaining float64
	texture       *sdl.Texture
	callback      func()
	alive         bool
}

// actual position is not stored, its just analytically determined
// as is rotation
func makeProjectile(start, end vec2f, texture *sdl.Texture, speed float64, callback func()) {
	p := Projectile{
		start, end,
		dist(start, end) / speed,
		dist(start, end) / speed,
		texture,
		callback,
		true,
	}
	context.projectiles = append(context.projectiles, p)
}

func (p *Projectile) update(dt float64) {
	p.timeRemaining -= dt
	if p.timeRemaining <= 0 && p.alive {
		p.alive = false
		p.callback()
	}
}

func (p Projectile) draw() {
	_, _, width, height, err := p.texture.Query()
	if err != nil {
		panic(err)
	}
	angle := angle(p.start, p.end)
	toRect := &sdl.Rect{int32(p.start[0] + (p.end[0]-p.start[0])*(p.lifetime-p.timeRemaining)/p.lifetime),
		int32(p.start[1] + (p.end[1]-p.start[1])*(p.lifetime-p.timeRemaining)/p.lifetime),
		width,
		height}
	context.renderer.CopyEx(p.texture, nil, toRect, RAD_TO_DEG*angle, nil, sdl.FLIP_NONE)

}
