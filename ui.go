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
	drawText(pad+selSize+pad, GAMEYRES+pad, fmt.Sprintf("%.0f/%.0f HP + %.2f HP/S", e.hp, e.hpmax, e.regen), 2)
	drawText(pad+selSize+pad, GAMEYRES+pad+1*textSize+1*pad, fmt.Sprintf("%.0f%% Fire Resist", 100*e.res[DAMAGE_FIRE]), 2)
	drawText(pad+selSize+pad, GAMEYRES+pad+2*textSize+2*pad, fmt.Sprintf("%.0f%% Chem Resist", 100*e.res[DAMAGE_CHEMICAL]), 2)
	drawText(pad+selSize+pad, GAMEYRES+pad+3*textSize+3*pad, fmt.Sprintf("%.0f%% Lightning Resist", 100*e.res[DAMAGE_LIGHTNING]), 2)
	drawText(pad+selSize+pad, GAMEYRES+pad+4*textSize+4*pad, fmt.Sprintf("%.0f%% Physical Resist", 100*e.res[DAMAGE_PHYSICAL]), 2)
	drawText(pad+selSize+pad, GAMEYRES+pad+5*textSize+5*pad, fmt.Sprintf("%.0f px/s Speed", e.speedBase), 2)
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
	context.renderer.CopyEx(context.atlas[TEX_CASH], nil, &sdl.Rect{cursorX, cursorY, 20, 20}, 0, nil, sdl.FLIP_NONE)
	drawText(cursorX+20+pad, cursorY, fmt.Sprintf("%d", props.cost), 2)
	cursorY += textSize + pad
	drawText(cursorX, cursorY, fmt.Sprintf("%.0f %s Damage", props.attack.damage, damageNames[props.attack.damageType]), 2)
	cursorY += textSize + pad
	drawText(cursorX, cursorY, fmt.Sprintf("%.0f range", props.attack.dist), 2)
	cursorY += textSize + pad
	drawText(cursorX, cursorY, fmt.Sprintf("%.2f second cooldown", props.attack.cooldown), 2)
	cursorY += textSize + pad
	drawText(cursorX, cursorY, fmt.Sprintf("%d kills", t.kills), 2)
	cursorY += textSize + pad
	indicateRange(context.selectedTower, props.attack.dist)
}

func drawTowerInfo(t TowerType) {
	pad := int32(10)
	var textSize int32 = 14
	var selSize int32 = 200
	selRect := &sdl.Rect{pad, GAMEYRES + pad, selSize, selSize}
	cursorX := pad + selSize + pad
	cursorY := GAMEYRES + pad

	props := towerProperties[t]

	context.renderer.CopyEx(context.atlas[towerProperties[t].texture], nil, selRect, 0.0, nil, sdl.FLIP_NONE)

	drawText(cursorX, cursorY, fmt.Sprintf("%s", props.name), 2)
	cursorY += textSize + pad

	drawText(cursorX, cursorY, fmt.Sprintf("%s", props.tooltip), 2)
	cursorY += textSize + pad

	context.renderer.CopyEx(context.atlas[TEX_CASH], nil, &sdl.Rect{cursorX, cursorY, 20, 20}, 0, nil, sdl.FLIP_NONE)
	drawText(cursorX+20+pad, cursorY, fmt.Sprintf("%d", props.cost), 2)
	cursorY += textSize + pad

	drawText(cursorX, cursorY, fmt.Sprintf("%.0f %s Damage", props.attack.damage, damageNames[props.attack.damageType]), 2)
	cursorY += textSize + pad

	drawText(cursorX, cursorY, fmt.Sprintf("%.0f range", props.attack.dist), 2)
	cursorY += textSize + pad

	drawText(cursorX, cursorY, fmt.Sprintf("%.2f second cooldown", props.attack.cooldown), 2)
	cursorY += textSize + pad
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

const TOWERBTN_S = int32(80)
const TOWERBTN_PAD = int32(10)

// btn frame, if pressed rotate even
func drawTowerBtn(hotkey string, t TowerType, x, y int32) {
	px := TOWERBTN_PAD + x*(TOWERBTN_S+TOWERBTN_PAD)
	py := TOWERBTN_PAD + y*(TOWERBTN_S+TOWERBTN_PAD) + GAMEYRES

	dstrect := &sdl.Rect{px, py, TOWERBTN_S, TOWERBTN_S}

	var col []uint8
	switch towerProperties[t].attack.damageType {
	case DAMAGE_CHEMICAL:
		col = []uint8{20, 200, 20, 255}
	case DAMAGE_LIGHTNING:
		col = []uint8{50, 50, 200, 255}
	case DAMAGE_FIRE:
		col = []uint8{200, 50, 0, 255}
	case DAMAGE_PHYSICAL:
		col = []uint8{200, 200, 200, 255}
	}

	context.renderer.SetDrawColorArray(col...)
	context.renderer.FillRect(dstrect)
	context.renderer.CopyEx(context.atlas[TextureID(t)+TEX_OFFSET_TOWERS-1], nil, dstrect, 0, nil, sdl.FLIP_NONE)
	context.renderer.CopyEx(context.atlas[TEX_BTN], nil, dstrect, 0, nil, sdl.FLIP_NONE)
	drawText(px+10, py+TOWERBTN_S-20, hotkey, 2)
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
