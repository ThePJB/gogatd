package main

import (
	"fmt"
	"math"
)

const RAD_TO_DEG = 180 / math.Pi
const DEG_TO_RAD = math.Pi / 180

type vec2f [2]float64
type vec2i [2]int32

func vecMul(a, b [2]float64) [2]float64 {
	return [2]float64{a[0] * b[0], a[1] * b[1]}
}
func vecMulScalar(a vec2f, b float64) vec2f {
	return vec2f{a[0] * b, a[1] * b}
}
func vecAdd(a, b [2]float64) [2]float64 {
	return [2]float64{a[0] + b[0], a[1] + b[1]}
}
func asF64(a [2]int32) [2]float64 {
	return [2]float64{float64(a[0]), float64(a[1])}
}
func asI32(a vec2f) vec2i {
	return vec2i{int32(a[0]), int32(a[1])}
}
func dist(a, b vec2f) float64 {
	return math.Sqrt((a[0]-b[0])*(a[0]-b[0]) + (a[1]-b[1])*(a[1]-b[1]))
}

// in radians
func angle(a, b vec2f) float64 {
	return math.Atan2(b[1]-a[1], b[0]-a[0])
}

func panics(a ...interface{}) {
	panic(fmt.Sprint(a...))
}

func getTileCenter(idx int32) vec2f {
	return vec2f{(float64(idx%GRIDW) + 0.5) * float64(context.cellw), (float64(idx/GRIDH) + 0.5) * float64(context.cellh)}
}

func getTileFromPos(pos vec2f) int32 {
	ivec := asI32(pos)
	gx := ivec[0] / context.cellw
	gy := ivec[1] / context.cellh

	if gx < 0 {
		panics("Tile gotten from", pos, " is out of bounds x < 0", gx)
	} else if gx > context.gridw {
		panics("Tile gotten from", pos, " is out of bounds x > gridw", gx, context.gridw)
	} else if gy > context.gridh {
		panics("Tile gotten from", pos, " is out of bounds y > gridh", gy, context.gridh)
	} else if gy < 0 {
		panics("Tile gotten from", pos, " is out of bounds y < 0", gy)
	}

	return gy*context.gridw + gx
}

func logisticFunction(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}
