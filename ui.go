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

func drawWaveText(n int) {
	w := int32(7)
	scale := int32(8)
	s := fmt.Sprintf("Wave %d", n)
	est_w := int32(len(s)) * w * scale
	drawText(GAMEXRES/2-est_w/2, 100, s, scale)
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
		context.renderer.CopyEx(context.atlas[TEX_FONT], &srcRect, &destRect, 0.0, nil, sdl.FLIP_NONE)
	}
}
