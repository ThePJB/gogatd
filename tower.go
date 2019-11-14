package main

import "github.com/veandco/go-sdl2/sdl"

type TowerType int32

const (
	None TowerType = 0

	Skull TowerType = 1
	Fire  TowerType = 2
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
		TowerProperties{},
	}
}

type AttackType int32

const (
	ATTACK_BEAM AttackType = 0
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
	switch t.towerType {
	case None:
		break
	case Skull:
		context.renderer.CopyEx(context.atlas.skull, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
	case Fire:
		context.renderer.CopyEx(context.atlas.fire, nil, toRect, 0.0, nil, sdl.FLIP_NONE)
	}
}
