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
	TOWER_ARROW
	TOWER_BLACKSMITH
	TOWER_TREBUCHET
	NUM_TOWERS
)

type DamageType int32

const (
	DAMAGE_FIRE DamageType = iota
	DAMAGE_CHEMICAL
	DAMAGE_LIGHTNING
	DAMAGE_PHYSICAL
)

var damageNames = [...]string{
	"fire",
	"chemical",
	"lightning",
	"physical",
}

type TowerProperties struct {
	name    string
	tooltip string
	texture TextureID
	attack  Attack
	cost    int
}

var towerProperties []TowerProperties

func initTowerProps() {
	towerProperties = []TowerProperties{
		TowerProperties{},
		TowerProperties{},
		TowerProperties{},
		/*
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
		*/
		TowerProperties{
			name:    "Fire Tower",
			tooltip: "Nasty AoE damage fire tower",
			texture: TEX_TOWER_FIRE,
			cost:    10,
			attack: Attack{
				attackType:  ATTACK_PROJECTILE,
				attackSound: CHUNK_FIRE_LAUNCH,
				damage:      6,
				damageType:  DAMAGE_FIRE,
				dist:        200,
				cooldown:    2,
				bp:          BeamProperties{},
				pp: ProjectileProperties{
					speed:             600,
					area:              100,
					lead:              false,
					texture:           TEX_PROJECTILE_FIRE,
					scale:             1.0,
					flipInterval:      0.2,
					rotationSpeed:     0,
					deathTexture:      TEX_PROJECTILE_FIRE,
					deathScale:        1.0,
					deathSound:        CHUNK_FIRE_EXPLODE,
					deathFlipInterval: 0.2,
				},
			},
		},
		TowerProperties{},
		TowerProperties{},
		TowerProperties{},
		TowerProperties{},
		/*
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
			TowerProperties{
				"Arrow Tower",
				1,
				context.atlas[TEX_TOWER_ARROW],
				ATTACK_PROJECTILE_ACCURATE,
				250,
				context.atlas[TEX_PROJECTILE_ARROW],
				CHUNK_LIGHTNING,
				CHUNK_LIGHTNING,
				8,
				DAMAGE_PHYSICAL,
				6,
			},
			TowerProperties{
				"Blacksmith",
				1.5,
				context.atlas[TEX_TOWER_BLACKSMITH],
				ATTACK_PROJECTILE_AOE,
				250,
				context.atlas[TEX_PROJECTILE_HAMMER],
				CHUNK_LIGHTNING,
				CHUNK_LIGHTNING,
				10,
				DAMAGE_PHYSICAL,
				12,
			},
			TowerProperties{
				"Trebuchet",
				3,
				context.atlas[TEX_TOWER_TREBUCHET],
				ATTACK_PROJECTILE_AOE,
				500,
				context.atlas[TEX_PROJECTILE_ROCK],
				CHUNK_FIRE_LAUNCH,
				CHUNK_FIRE_EXPLODE,
				12,
				DAMAGE_PHYSICAL,
				14,
			},*/
		TowerProperties{},
	}
}

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
	context.renderer.CopyEx(context.atlas[towerProperties[t.towerType].texture], nil, toRect, 0.0, nil, sdl.FLIP_NONE)
}
