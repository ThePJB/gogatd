package main

import (
	"math"

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
		TowerProperties{
			name:    "Laser Tower",
			tooltip: "Rapid fire fire tower",
			texture: TEX_TOWER_LASER,
			cost:    8,
			attack: Attack{
				attackType:  ATTACK_BEAM,
				attackSound: CHUNK_LASER,
				damage:      1,
				damageType:  DAMAGE_FIRE,
				dist:        200,
				cooldown:    0.25,
				bp: BeamProperties{
					texture:  TEX_BEAM_LASER,
					fadeTime: 0.5,
					width:    6,
				},
				pp: ProjectileProperties{},
			},
		},
		TowerProperties{
			name:    "Fire Tower",
			tooltip: "Nasty AoE damage fire tower",
			texture: TEX_TOWER_FIRE,
			cost:    13,
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
					scale:             0.8,
					flipInterval:      0.2,
					rotationSpeed:     0,
					deathTexture:      TEX_PROJECTILE_FIRE,
					deathScale:        0.9,
					deathSound:        CHUNK_FIRE_EXPLODE,
					deathFlipInterval: 0.2,
					deathTime:         0.4,
				},
			},
		},
		TowerProperties{
			name:    "Lightning Tower",
			tooltip: "Long range beam tower",
			texture: TEX_TOWER_LIGHTNING,
			cost:    14,
			attack: Attack{
				attackType:  ATTACK_BEAM,
				attackSound: CHUNK_LIGHTNING,
				damage:      10,
				damageType:  DAMAGE_LIGHTNING,
				dist:        350,
				cooldown:    2,
				bp: BeamProperties{
					texture:  TEX_BEAM_LIGHTNING,
					fadeTime: 0.5,
					width:    8,
				},
				pp: ProjectileProperties{},
			},
		},
		TowerProperties{
			name:    "Arrow Tower",
			tooltip: "Medium range projectile tower",
			texture: TEX_TOWER_ARROW,
			cost:    9,
			attack: Attack{
				attackType:  ATTACK_PROJECTILE,
				attackSound: CHUNK_ARROW_LAUNCH,
				damage:      6,
				damageType:  DAMAGE_PHYSICAL,
				dist:        225,
				cooldown:    1.5,
				bp:          BeamProperties{},
				pp: ProjectileProperties{
					speed:             600,
					area:              0,
					lead:              true,
					texture:           TEX_PROJECTILE_ARROW,
					scale:             0.5,
					flipInterval:      0.25, // 3 = no flip
					rotationSpeed:     0,
					deathTexture:      TEX_PROJECTILE_ARROW,
					deathScale:        1.0,
					deathSound:        CHUNK_ARROW_HIT,
					deathFlipInterval: 0.05,
					deathTime:         0.2,
				},
			},
		},
		TowerProperties{
			name:    "Blacksmith",
			tooltip: "Buffs 4-adjacent physical towers (not implemented yet tho lol)",
			texture: TEX_TOWER_BLACKSMITH,
			cost:    20,
			attack: Attack{
				attackType:  ATTACK_PROJECTILE,
				attackSound: CHUNK_FIRE_LAUNCH,
				damage:      12,
				damageType:  DAMAGE_PHYSICAL,
				dist:        200,
				cooldown:    2,
				bp:          BeamProperties{},
				pp: ProjectileProperties{
					speed:             200,
					area:              60,
					lead:              true,
					texture:           TEX_PROJECTILE_HAMMER,
					scale:             0.5,
					flipInterval:      3, // 3 = no flip
					rotationSpeed:     0,
					deathTexture:      TEX_PROJECTILE_HAMMER,
					deathScale:        1.0,
					deathSound:        CHUNK_FIRE_EXPLODE,
					deathFlipInterval: 0.2,
					deathTime:         0.4,
				},
			},
		},
		TowerProperties{
			name:    "Trebuchet",
			tooltip: "Long range physical AoE tower",
			texture: TEX_TOWER_TREBUCHET,
			cost:    30,
			attack: Attack{
				attackType:  ATTACK_PROJECTILE,
				attackSound: CHUNK_FIRE_LAUNCH,
				damage:      16,
				damageType:  DAMAGE_PHYSICAL,
				dist:        350,
				cooldown:    3,
				bp:          BeamProperties{},
				pp: ProjectileProperties{
					speed:             150,
					area:              50,
					lead:              true,
					texture:           TEX_PROJECTILE_ROCK,
					scale:             0.5,
					flipInterval:      3, // 3 = no flip
					rotationSpeed:     6 * math.Pi,
					deathTexture:      TEX_PROJECTILE_ROCK,
					deathScale:        1.0,
					deathSound:        CHUNK_FIRE_EXPLODE,
					deathFlipInterval: 0.2,
					deathTime:         0.4,
				},
			},
		},
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
