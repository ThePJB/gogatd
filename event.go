package main

// maybe it should be called the tweener or something
// esp for graphics
// I could make this just literally a graphics one, since its order matters and everything as well
type DoLater struct {
	from   float64
	to     float64
	update func(t float64) // analytical 0..1
}

// this can be composed into multiple stages e.g. by making 0 to 0.2 something and then expanding that to 1 and then you could slowStart it or whatever

// then for gameplay events it would just be a one shot event

// or a (dt) one if i actually needed it
