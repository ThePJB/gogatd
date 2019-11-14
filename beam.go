package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

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
