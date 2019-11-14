package main

import "github.com/veandco/go-sdl2/sdl"

type TowerType int32

const (
	None TowerType = iota
	Skull
	Laser
	Fire
	Lightning
	NUM_TOWERS
)

type TowerProperties struct {
	name          string
	cooldown      float64
	texture       *sdl.Texture
	attackType    AttackType
	attackRange   float64
	attackTexture *sdl.Texture
	damage        float64
}

var towerProperties []TowerProperties

func initTowerProps() {
	towerProperties = []TowerProperties{
		TowerProperties{},
		TowerProperties{
			"Skull Tower",
			3.0,
			context.atlas.skull,
			ATTACK_BEAM,
			200,
			context.atlas.whitebeam,
			4,
		},
		TowerProperties{
			"Laser Tower",
			0.25,
			context.atlas.laserTower,
			ATTACK_BEAM,
			200,
			context.atlas.laserBeam,
			1,
		},
		TowerProperties{
			"Fire Tower",
			2.0,
			context.atlas.fireTower,
			ATTACK_PROJECTILE,
			200,
			context.atlas.fireProjectile,
			12,
		},
		TowerProperties{
			"Lightning Tower",
			3,
			context.atlas.lightningTower,
			ATTACK_BEAM,
			350,
			context.atlas.lightningBeam,
			10,
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
