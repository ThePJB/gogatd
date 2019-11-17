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

	fireTower      *sdl.Texture
	laserTower     *sdl.Texture
	skull          *sdl.Texture
	lightningTower *sdl.Texture

	dude  *sdl.Texture
	splat *sdl.Texture

	font *sdl.Texture
	cash *sdl.Texture

	barend   *sdl.Texture
	bargreen *sdl.Texture
	barred   *sdl.Texture

	whitebeam      *sdl.Texture
	laserBeam      *sdl.Texture
	lightningBeam  *sdl.Texture
	fireProjectile *sdl.Texture
}

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

	TEX_PROJECTILE_FIRE

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

	"flameProjectile",
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
