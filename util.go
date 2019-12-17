package main

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
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
func vecSub(a, b [2]float64) [2]float64 {
	return [2]float64{a[0] - b[0], a[1] - b[1]}
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
func cross2d(a, b vec2f) float64 {
	return a[0]*b[1] - b[0]*a[1]
}
func rot90(a vec2f) vec2f {
	return vec2f{-a[1], a[0]}
}
func interp(a vec2f, b vec2f, c float64) vec2f {
	return vec2f{a[0] + c*(b[0]-a[0]), a[1] + c*(b[1]-a[1])}
}

func (a vec2f) norm() float64 {
	return math.Sqrt(a[0]*a[0] + a[1]*a[1])
}
func (a vec2f) unit() vec2f {
	return vecMulScalar(a, 1/math.Sqrt(a[0]*a[0]+a[1]*a[1]))
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

func getTileRect(idx int32) *sdl.Rect {
	return &sdl.Rect{
		(int32(idx) % context.gridw) * context.cellw,
		(int32(idx) / context.gridw) * context.cellh,
		context.cellw, context.cellh,
	}
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
func min(a, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}
func max(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}
func logisticFunction(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}
func slowStart2(x float64) float64 {
	return x * x
}
func slowStart3(x float64) float64 {
	return x * x * x
}
func slowStart4(x float64) float64 {
	return x * x * x * x
}
func slowStop2(x float64) float64 {
	return 1 - (slowStart2(1 - x))
}
func slowStop3(x float64) float64 {
	return 1 - (slowStart3(1 - x))
}
func slowStop4(x float64) float64 {
	return 1 - (slowStart4(1 - x))
}
func slowStartStop2(x float64) float64 {
	return x*slowStart2(x) + (1-x)*slowStop2(x)
}
func slowStartStop3(x float64) float64 {
	return x*slowStart3(x) + (1-x)*slowStop3(x)
}
func slowStartStop4(x float64) float64 {
	return x*slowStart4(x) + (1-x)*slowStop4(x)
}

// whats the curve for alpha?
/*
func init() {
	for x := 0.0; x <= 1; x += 0.1 {
		fmt.Printf("%.2f ", x)
	}
	fmt.Println()

	fmt.Println("slowstart2")
	for x := 0.0; x <= 1; x += 0.1 {
		fmt.Printf("%.2f ", slowStart2(x))
	}
	fmt.Println()
	fmt.Println("slowstop2")
	for x := 0.0; x <= 1; x += 0.1 {
		fmt.Printf("%.2f ", slowStop2(x))
	}
	fmt.Println()
	fmt.Println("slowss2")
	for x := 0.0; x <= 1; x += 0.1 {
		fmt.Printf("%.2f ", slowStartStop2(x))
	}
	fmt.Println()

	panic("breakpoint")
}
*/
