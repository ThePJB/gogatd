package main

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Atlas struct {
	path      *sdl.Texture
	buildable *sdl.Texture
	portal    *sdl.Texture
	orb       *sdl.Texture
	wall      *sdl.Texture
	wallTop   *sdl.Texture

	fire      *sdl.Texture
	skull     *sdl.Texture
	lightning *sdl.Texture

	dude *sdl.Texture

	font *sdl.Texture
	cash *sdl.Texture

	barend   *sdl.Texture
	bargreen *sdl.Texture
	barred   *sdl.Texture

	whitebeam *sdl.Texture
}

func loadAssets() Atlas {
	atlas := Atlas{}
	atlas.path = loadTexture("assets/path.png")
	atlas.buildable = loadTexture("assets/buildable.png")
	atlas.skull = loadTexture("assets/spooky.png")
	atlas.fire = loadTexture("assets/fire.png")
	atlas.dude = loadTexture("assets/dude.png")
	atlas.wall = loadTexture("assets/brick.png")
	atlas.wallTop = loadTexture("assets/brickPerspective.png")
	atlas.orb = loadTexture("assets/magicOrb.png")
	atlas.portal = loadTexture("assets/portal.png")
	atlas.font = loadTexture("assets/custombold.png")
	atlas.barend = loadTexture("assets/hpend.png")
	atlas.bargreen = loadTexture("assets/hpgreen.png")
	atlas.barred = loadTexture("assets/hpred.png")
	atlas.whitebeam = loadTexture("assets/whitebeam.png")

	return atlas
}

func loadTexture(path string) *sdl.Texture {
	image, err := img.Load(path)
	if err != nil {
		panic(err)
	}
	defer image.Free()
	image.SetColorKey(true, 0xffff00ff)
	texture, err := context.renderer.CreateTextureFromSurface(image)
	if err != nil {
		panic(err)
	}
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	return texture
}

// also could do colour etc
// 8x32 and 8 px
func drawText(x, y int32, text string, scale int32) {
	w := int32(7)
	h := int32(8)
	fw := int32(32)
	for i := range text {
		destRect := sdl.Rect{x + int32(i)*w*scale, y, w * scale, h * scale}
		char := int32(text[i])
		sx := char % fw
		sy := char / fw
		srcRect := sdl.Rect{sx * w, sy * h, w, h}
		context.renderer.CopyEx(context.atlas.font, &srcRect, &destRect, 0.0, nil, sdl.FLIP_NONE)
	}

}
