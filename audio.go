package main

import (
	"github.com/veandco/go-sdl2/mix"
)

type ChunkID int

const (
	CHUNK_LASER = iota
	CHUNK_FIRE_LAUNCH
	CHUNK_ARROW_LAUNCH
	CHUNK_ARROW_HIT
	CHUNK_FIRE_EXPLODE
	CHUNK_LIGHTNING
	CHUNK_DIE
	CHUNK_WAVE
	CHUNK_LEAK
	NUM_CHUNKS
)

var chunkNames = [...]string{
	"laserAttack",
	"fireAttack",
	"arrowLaunch",
	"arrowHit",
	"woosh",
	"zzt",
	"pop",
	"snare",
	"leak",
}

func loadChunks() {
	for i := 0; i < int(NUM_CHUNKS); i++ {
		context.chunks = append(context.chunks, loadWav("assets/"+chunkNames[i]+".wav"))
	}
	context.chunks[CHUNK_LASER].Volume(35)
	context.chunks[CHUNK_FIRE_LAUNCH].Volume(20)
	context.chunks[CHUNK_FIRE_EXPLODE].Volume(40)
	context.chunks[CHUNK_DIE].Volume(50)
	context.chunks[CHUNK_LIGHTNING].Volume(25)
	context.chunks[CHUNK_WAVE].Volume(90)
	context.chunks[CHUNK_ARROW_LAUNCH].Volume(25)
	context.chunks[CHUNK_ARROW_HIT].Volume(30)
}

func loadWav(path string) *mix.Chunk {
	chunk, err := mix.LoadWAV(path)
	if err != nil {
		panic(err)
	}
	return chunk
}
