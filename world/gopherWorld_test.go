package world

import (
	"gopherlife/geometry"
	"testing"
)

func TestPartitionGopherWorld_MoveGopher(t *testing.T) {

	settings := GopherWorldSettings{
		Dimensions:      Dimensions{10, 10},
		Population:      Population{0, 100},
		NumberOfFood:    20,
		GopherBirthRate: 7,
	}

	var tileMap = CreateGopherWorldGridPartition(settings)

	gopher := NewGopher("a", geometry.NewCoordinate(1, 2))
	tileMap.InsertGopher(1, 2, &gopher)

	tileMap.ActionQueuer.Add(func() {
		tileMap.MoveGopher(&gopher, 0, 1)
	})
	tileMap.Update()

	des, _ := tileMap.Tile(1, 3)

	if des.Gopher == nil {
		t.Errorf("Destination is empty")
	}

	prev, _ := tileMap.Tile(1, 2)

	if prev.Gopher != nil {
		t.Errorf("Previous Destination is not empty")
	}

}

func TestPartitionGopherWorld_RemoveGopher(t *testing.T) {

	settings := GopherWorldSettings{
		Dimensions:      Dimensions{10, 10},
		Population:      Population{0, 100},
		NumberOfFood:    20,
		GopherBirthRate: 7,
	}

	var tileMap = CreateGopherWorldGridPartition(settings)

	gopher := NewGopher("a", geometry.Coordinates{1, 2})
	tileMap.InsertGopher(1, 2, &gopher)

	_, bool := tileMap.RemoveGopher(1, 2)

	if !bool {
		t.Errorf("Gopher is not removed")
	}

	tile, _ := tileMap.Tile(1, 2)

	if tile.Gopher != nil {
		t.Errorf("Gopher is not removed")
	}
}
