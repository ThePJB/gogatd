package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func drawSelectedEnemy() {
	var pad int32 = 10
	var textSize int32 = 14
	var selSize int32 = 200
	selRect := &sdl.Rect{pad, GAMEYRES + pad, selSize, selSize}
	drawEnemyToRect(context.selectedEnemy, selRect)
	e := context.enemies[context.selectedEnemy]
	drawText(pad+selSize+pad, GAMEYRES+pad, fmt.Sprintf("%.0f/%.0f HP", e.hp, e.hpmax), 2)
	drawText(pad+selSize+pad, GAMEYRES+pad+1*textSize+1*pad, fmt.Sprintf("%.0f%% Fire Resist", 100*e.res[DAMAGE_FIRE]), 2)
	drawText(pad+selSize+pad, GAMEYRES+pad+2*textSize+2*pad, fmt.Sprintf("%.0f%% Chem Resist", 100*e.res[DAMAGE_CHEMICAL]), 2)
	drawText(pad+selSize+pad, GAMEYRES+pad+3*textSize+3*pad, fmt.Sprintf("%.0f%% Lightning Resist", 100*e.res[DAMAGE_LIGHTNING]), 2)
	drawText(pad+selSize+pad, GAMEYRES+pad+4*textSize+4*pad, fmt.Sprintf("%.0f%% Speed", e.speedBase), 2)
}

// place mode shows grid and ghost tower

// draw qwer buttons for place towers
