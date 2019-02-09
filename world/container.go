package world

import (
	"gopherlife/calc"
)

type TileContainer interface {
	Tile(x int, y int) (*Tile, bool)
}

type Basic2DContainer struct {
	grid   [][]*Tile
	x      int
	y      int
	width  int
	height int
}

func NewBasic2DContainer(x int, y int, width int, height int) Basic2DContainer {

	container := Basic2DContainer{
		x:      x,
		y:      y,
		width:  width,
		height: height}

	container.grid = make([][]*Tile, width)

	for i := 0; i < width; i++ {
		container.grid[i] = make([]*Tile, height)

		for j := 0; j < height; j++ {
			tile := Tile{nil, nil}
			container.grid[i][j] = &tile
		}
	}

	return container
}

func (container *Basic2DContainer) Tile(x int, y int) (*Tile, bool) {
	if x < container.x || x >= container.width+container.x || y < container.y || y >= container.height+container.y {
		return nil, false
	}

	return container.grid[x-container.x][y-container.y], true
}

//InsertGopher Inserts the given gopher into the tileMap at the specified co-ordinate
func (container *Basic2DContainer) InsertGopher(x int, y int, gopher *Gopher) bool {

	if tile, ok := container.Tile(x, y); ok {
		if !tile.HasGopher() {
			tile.SetGopher(gopher)
			return true
		}
	}

	return false

}

//InsertFood Inserts the given food into the tileMap at the specified co-ordinate
func (container *Basic2DContainer) InsertFood(x int, y int, food *Food) bool {

	if tile, ok := container.Tile(x, y); ok {
		if !tile.HasFood() {
			tile.SetFood(food)
			return true
		}
	}
	return false
}

//RemoveFoodFromWorld Removes food from the given coordinates. Returns the food value.
func (container *Basic2DContainer) RemoveGopher(x int, y int) bool {

	if mapPoint, ok := container.Tile(x, y); ok {
		if mapPoint.HasGopher() {
			mapPoint.ClearGopher()
			return true
		}
	}

	return false
}

//RemoveFoodFromWorld Removes food from the given coordinates. Returns the food value.
func (container *Basic2DContainer) RemoveFood(x int, y int) (*Food, bool) {

	if mapPoint, ok := container.Tile(x, y); ok {
		if mapPoint.HasFood() {
			var food = mapPoint.Food
			mapPoint.ClearFood()
			return food, true
		}
	}

	return nil, false
}

type TrackedTileContainer struct {
	x      int
	y      int
	width  int
	height int
	TileContainer
	gopherTileLocations map[int]*Tile
	foodTileLocations   map[int]*Tile
	Insertable
}

func NewTrackedTileContainer(x int, y int, width int, height int) TrackedTileContainer {
	b2dc := NewBasic2DContainer(x, y, width, height)

	return TrackedTileContainer{
		x:                   x,
		y:                   y,
		width:               width,
		height:              height,
		TileContainer:       &b2dc,
		gopherTileLocations: make(map[int]*Tile),
		foodTileLocations:   make(map[int]*Tile),
	}
}

func (container *TrackedTileContainer) Tile(x int, y int) (*Tile, bool) {
	return container.TileContainer.Tile(x, y)
}

func (container *TrackedTileContainer) ConvertToTrackedTileCoordinates(x int, y int) (gridX int, gridY int) {
	return (x - container.x), (y - container.y)
}

func (container *TrackedTileContainer) InsertGopher(x int, y int, gopher *Gopher) bool {
	if tile, ok := container.Tile(x, y); ok {
		if !tile.HasGopher() {
			tile.SetGopher(gopher)
			gopher.Position.X = x
			gopher.Position.Y = y
			x, y = container.ConvertToTrackedTileCoordinates(x, y)
			container.gopherTileLocations[calc.Hashcode(x, y)] = tile
			return true
		}
	}

	return false

}

func (container *TrackedTileContainer) InsertFood(x int, y int, food *Food) bool {
	if tile, ok := container.Tile(x, y); ok {
		if !tile.HasFood() {
			food.Position.X = x
			food.Position.Y = y
			tile.SetFood(food)
			x, y = container.ConvertToTrackedTileCoordinates(x, y)
			container.foodTileLocations[calc.Hashcode(x, y)] = tile
			return true
		}
	}

	return false

}

func (container *TrackedTileContainer) RemoveGopher(x int, y int, gopher *Gopher) bool {
	if tile, ok := container.Tile(x, y); ok {
		if tile.HasGopher() {
			tile.ClearGopher()
			x, y = container.ConvertToTrackedTileCoordinates(x, y)
			delete(container.gopherTileLocations, calc.Hashcode(x, y))
			return true
		}
	}
	return false
}

func (container *TrackedTileContainer) RemoveFood(x int, y int, food *Food) bool {
	if tile, ok := container.Tile(x, y); ok {
		if tile.HasFood() {
			tile.ClearFood()
			x, y = container.ConvertToTrackedTileCoordinates(x, y)
			delete(container.foodTileLocations, calc.Hashcode(x, y))
			return true
		}
	}
	return false
}

type GridContainer interface {
	TileContainer
	Grid(x int, y int) (TileContainer, bool)
}

type BasicGridContainer struct {
	containers [][]*TrackedTileContainer
	gridWidth  int
	gridHeight int
	width      int
	height     int
}

func NewBasicGridContainer(width int, height int, gridWidth int, gridHeight int) BasicGridContainer {

	numberOfGridsX := width / gridWidth

	if numberOfGridsX*gridWidth < width {
		numberOfGridsX++
	}

	numberOfGridsY := height / gridHeight

	if numberOfGridsY*gridHeight < height {
		numberOfGridsY++
	}

	containers := make([][]*TrackedTileContainer, numberOfGridsX)

	for i := 0; i < numberOfGridsX; i++ {
		containers[i] = make([]*TrackedTileContainer, numberOfGridsY)

		for j := 0; j < numberOfGridsY; j++ {
			ttc := NewTrackedTileContainer(i*gridWidth,
				j*gridHeight,
				gridWidth,
				gridHeight)
			containers[i][j] = &ttc
		}
	}

	return BasicGridContainer{
		containers: containers,
		gridWidth:  gridWidth,
		gridHeight: gridHeight,
		width:      width,
		height:     height,
	}
}

func (container *BasicGridContainer) Tile(x int, y int) (*Tile, bool) {

	if grid, ok := container.Grid(x, y); ok {
		if tile, ok := grid.Tile(x, y); ok {
			return tile, ok
		}
	}
	return nil, false
}

//Takes an X and Y Position, and finds which grid it should be in
func (container *BasicGridContainer) Grid(x int, y int) (*TrackedTileContainer, bool) {

	if x < 0 || x >= container.width || y < 0 || y >= container.height {
		return nil, false
	}

	x, y = container.convertToGridCoordinates(x, y)

	val := container.containers[x][y]

	if val != nil {
		return val, true
	}

	return nil, false
}

func (container *BasicGridContainer) convertToGridCoordinates(x int, y int) (int, int) {
	gridX, gridY := x/container.gridWidth, y/container.gridHeight
	return gridX, gridY
}

type Insertable interface {
	InsertGopher(x int, y int, gopher *Gopher) bool
	InsertFood(x int, y int, food *Food) bool
	RemoveGopher(x int, y int) bool
	RemoveFood(x int, y int) (*Food, bool)
}

type GridInsertable struct {
	Grid [][]GridContainer
}

func (container *BasicGridContainer) InsertGopher(x int, y int, gopher *Gopher) bool {
	if grid, ok := container.Grid(x, y); ok {
		return grid.InsertGopher(x, y, gopher)
	}
	return false
}

func (container *BasicGridContainer) InsertFood(x int, y int, food *Food) bool {
	if grid, ok := container.Grid(x, y); ok {
		return grid.InsertFood(x, y, food)
	}
	return false
}

func (container *BasicGridContainer) RemoveGopher(x int, y int, gopher *Gopher) bool {
	if grid, ok := container.Grid(x, y); ok {
		return grid.RemoveGopher(x, y, gopher)
	}
	return false
}

func (container *BasicGridContainer) RemoveFood(x int, y int, food *Food) bool {
	if grid, ok := container.Grid(x, y); ok {
		return grid.RemoveFood(x, y, food)
	}
	return false
}
