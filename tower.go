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

type DamageType int32

const (
	DAMAGE_FIRE DamageType = iota
	DAMAGE_CHEMICAL
	DAMAGE_LIGHTNING
)

type TowerProperties struct {
	name          string
	cooldown      float64
	texture       *sdl.Texture
	attackType    AttackType
	attackRange   float64
	attackTexture *sdl.Texture
	damage        float64
	damageType    DamageType
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
			DAMAGE_CHEMICAL,
		},
		TowerProperties{
			"Laser Tower",
			0.25,
			context.atlas.laserTower,
			ATTACK_BEAM,
			200,
			context.atlas.laserBeam,
			1,
			DAMAGE_FIRE,
		},
		TowerProperties{
			"Fire Tower",
			2.0,
			context.atlas.fireTower,
			ATTACK_PROJECTILE,
			200,
			context.atlas.fireProjectile,
			9,
			DAMAGE_FIRE,
		},
		TowerProperties{
			"Lightning Tower",
			3,
			context.atlas.lightningTower,
			ATTACK_BEAM,
			350,
			context.atlas.lightningBeam,
			8,
			DAMAGE_LIGHTNING,
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

func damage(enemyIdx int, amount float64, damageType DamageType) {
	damageAfterRes := amount * (1 - context.enemies[enemyIdx].res[damageType])
	context.enemies[enemyIdx].hp -= damageAfterRes
	if context.enemies[enemyIdx].hp <= 0 {
		killEnemy(enemyIdx)
	}
}
