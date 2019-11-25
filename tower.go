package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type TowerType int32

const (
	None TowerType = iota
	TOWER_SKULL
	TOWER_LASER
	TOWER_FIRE
	TOWER_LIGHTNING
	NUM_TOWERS
)

type DamageType int32

const (
	DAMAGE_FIRE DamageType = iota
	DAMAGE_CHEMICAL
	DAMAGE_LIGHTNING
)

var damageNames = [...]string{
	"fire",
	"chemical",
	"lightning",
}

type TowerProperties struct {
	name            string
	cooldown        float64
	texture         *sdl.Texture
	attackType      AttackType
	attackRange     float64
	attackTexture   *sdl.Texture
	attackSound     ChunkID
	attackLandSound ChunkID
	damage          float64
	damageType      DamageType
	cost            int
}

var towerProperties []TowerProperties

func initTowerProps() {
	towerProperties = []TowerProperties{
		TowerProperties{},
		TowerProperties{
			"Skull Tower",
			3.0,
			context.atlas[TEX_TOWER_SKULL],
			ATTACK_BEAM,
			200,
			context.atlas[TEX_BEAM_WHITE],
			CHUNK_FIRE_LAUNCH,
			CHUNK_FIRE_LAUNCH,
			4,
			DAMAGE_CHEMICAL,
			4,
		},
		TowerProperties{
			"Laser Tower",
			0.25,
			context.atlas[TEX_TOWER_LASER],
			ATTACK_BEAM,
			200,
			context.atlas[TEX_BEAM_LASER],
			CHUNK_LASER,
			CHUNK_LASER,
			1,
			DAMAGE_FIRE,
			5,
		},
		TowerProperties{
			"Fire Tower",
			2.0,
			context.atlas[TEX_TOWER_FIRE],
			ATTACK_PROJECTILE,
			200,
			context.atlas[TEX_PROJECTILE_FIRE],
			CHUNK_FIRE_LAUNCH,
			CHUNK_FIRE_EXPLODE,
			9,
			DAMAGE_FIRE,
			10,
		},
		TowerProperties{
			"Lightning Tower",
			3,
			context.atlas[TEX_TOWER_LIGHTNING],
			ATTACK_BEAM,
			350,
			context.atlas[TEX_BEAM_LIGHTNING],
			CHUNK_LIGHTNING,
			CHUNK_LIGHTNING,
			8,
			DAMAGE_LIGHTNING,
			6,
		},
		TowerProperties{},
	}
}

type AttackType int32

const (
	ATTACK_BEAM AttackType = iota
	ATTACK_PROJECTILE
)

type Tower struct {
	towerType TowerType
	cooldown  float64
	kills     int
}

func makeTower(tt TowerType) Tower {
	return Tower{
		towerType: tt,
		cooldown:  0,
	}
}

func drawTower(t Tower, toRect *sdl.Rect) {
	context.renderer.CopyEx(towerProperties[t.towerType].texture, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
}
