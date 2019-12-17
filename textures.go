package main

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type TextureID int

const (
	TEX_PATH TextureID = iota
	TEX_BUILDABLE
	TEX_PORTAL
	TEX_ORB
	TEX_WALL
	TEX_WALLTOP

	TEX_DUDE
	TEX_SPLAT

	TEX_FONT
	TEX_CASH
	TEX_BTN

	TEX_BAREND
	TEX_BARGREEN
	TEX_BARRED

	TEX_BEAM_WHITE
	TEX_BEAM_LASER
	TEX_BEAM_LIGHTNING

	TEX_TOWER_SKULL
	TEX_TOWER_LASER
	TEX_TOWER_FIRE
	TEX_TOWER_LIGHTNING
	TEX_TOWER_ARROW
	TEX_TOWER_BLACKSMITH
	TEX_TOWER_TREBUCHET

	TEX_PROJECTILE_FIRE
	TEX_PROJECTILE_ARROW
	TEX_PROJECTILE_HAMMER
	TEX_PROJECTILE_ROCK

	NUM_TEXTURES
)

const (
	TEX_OFFSET_TOWERS = TEX_TOWER_SKULL
	TEX_OFFSET_TILES  = TEX_PATH
)

var TextureNames = [...]string{
	"path",
	"buildable",
	"portal",
	"magicOrb",
	"brick",
	"brickPerspective",

	"dude",
	"splat",

	"custombold",
	"dolla",
	"btnFrame",

	"hpend",
	"hpgreen",
	"hpred",

	"whitebeam",
	"laserbeam",
	"lightningbeam",

	"spooky",
	"laserTower",
	"flametower",
	"lightningTower",
	"arrowTower",
	"blacksmith",
	"trebuchetTower",

	"flameProjectile",
	"arrowProjectile",
	"hammerProjectile",
	"rockProjectile",
}

func loadTextures() {
	for i := 0; i < int(NUM_TEXTURES); i++ {
		context.atlas = append(context.atlas, loadTexture("assets/"+TextureNames[i]+".png"))
	}
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
