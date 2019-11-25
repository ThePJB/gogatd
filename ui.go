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
	drawText(pad+selSize+pad, GAMEYRES+pad+4*textSize+4*pad, fmt.Sprintf("%.0f px/s Speed", e.speedBase), 2)
}

func drawSelectedTower() {
	t := context.grid[context.selectedTower].tower
	pad := int32(10)
	var textSize int32 = 14
	var selSize int32 = 200
	selRect := &sdl.Rect{pad, GAMEYRES + pad, selSize, selSize}
	cursorX := pad + selSize + pad
	cursorY := GAMEYRES + pad

	props := towerProperties[t.towerType]
	drawTower(t, selRect)
	drawText(cursorX, cursorY, fmt.Sprintf("%.0f %s Damage", props.damage, damageNames[props.damageType]), 2)
	cursorY += textSize + pad
	drawText(cursorX, cursorY, fmt.Sprintf("%.0f range", props.attackRange), 2)
	cursorY += textSize + pad
	drawText(cursorX, cursorY, fmt.Sprintf("%.2f second cooldown", props.cooldown), 2)
	cursorY += textSize + pad
	drawText(cursorX, cursorY, fmt.Sprintf("%d kills", t.kills), 2)
	cursorY += textSize + pad
	indicateRange(context.selectedTower, towerProperties[t.towerType].attackRange)
}

// place mode shows grid and ghost tower

// draw qwer buttons for place towers

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

func waveAnnounce(n int, t float64) {
	dt := 6.0
	context.tweens = append(context.tweens, Tween{
		from: t,
		to:   t + dt,
		update: func(t float64) {
			var y int32 = 100
			if t < 0.2 {
				tn := t / 0.2
				tn = slowStop4(tn)
				y = int32(0.5 + tn*100)
			}
			if t > 0.6 {
				tn := (t - 0.6) / 0.6
				tn = slowStart4(tn)
				y = int32(0.5 + (1-tn)*100)
			}
			var alpha uint8 = 255
			if t > 0.8 {
				tn := (t - 0.8) / 0.8
				alpha = uint8(0.5 + 255.0*(1-tn))
				//fmt.Println(t, tn, alpha)
			}
			context.atlas[TEX_FONT].SetAlphaMod(alpha)
			w := int32(7)
			scale := int32(8)
			s := fmt.Sprintf("Wave %d", n)
			est_w := int32(len(s)) * w * scale
			drawText(GAMEXRES/2-est_w/2, y, s, scale)
			context.atlas[TEX_FONT].SetAlphaMod(255)
		},
	})
}

// colour tiles and use opacity to indicate the cutoff point
func indicateRange(aboutIdx int32, d float64) {
	context.renderer.SetDrawColor(255, 255, 255, 128)
	for i := range context.grid {
		if dist(getTileCenter(aboutIdx), getTileCenter(int32(i))) <= d {
			context.renderer.FillRect(getTileRect(int32(i)))
		}
	}
}
