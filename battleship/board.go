package battleship

import (
	"fmt"
	"strings"
)

const (
	MARK_EMPTY = " "
	MARK_MISS  = "-"
	MARK_HIT   = "x"
	MARK_SHIP  = "s"
)

type Board struct {
	size           int
	layout         [][]string
	ships          *Ships
	shipsRemaining int
}

func NewOpponentBoard(size int) *Board {
	return newBoard(size, true)
}

func NewMyBoard(size int) *Board {
	return newBoard(size, false)
}

func newBoard(size int, opponent bool) *Board {
	board := &Board{
		size:           size,
		layout:         make([][]string, size),
		ships:          NewShips(),
		shipsRemaining: 0,
	}
	for i := 0; i < size; i++ {
		board.layout[i] = make([]string, size)

		for j := range board.layout[i] {
			board.layout[i][j] = MARK_EMPTY
		}
	}
	return board
}

func (b *Board) MarkHit(location Location) {
	b.mark(location, MARK_HIT)
}

func (b *Board) MarkMiss(location Location) {
	b.mark(location, MARK_MISS)
}

func (b *Board) mark(location Location, mark string) {
	b.layout[location.Row()][location.Column()] = mark
}

func (b *Board) CheckForShip(location Location) *Ship {
	for _, ship := range b.ships.ships {
		if b.check(location, ship.Mark) {
			return ship
		}
	}
	return nil
}

func (b *Board) CheckLocationUnplayed(location Location) bool {
	for _, ship := range b.ships.ships {
		if b.check(location, ship.Mark) {
			return false
		}
	}
	return b.check(location, MARK_EMPTY)
}

func (b *Board) check(location Location, lookup string) bool {
	if b.layout[location.Row()][location.Column()] == lookup {
		return true
	}
	return false
}

func (b *Board) GetLocationsAround(location Location) []Location {
	locations := make([]Location, 0)
	row := location.Row()
	col := location.Column()
	row++
	col++
	if row > 1 {
		if b.CheckLocationUnplayed(cordinatesToLocation(row-1, col)) {
			locations = append(locations, cordinatesToLocation(row-1, col))
		}
	}
	if row <= b.size-1 {
		if b.CheckLocationUnplayed(cordinatesToLocation(row+1, col)) {
			locations = append(locations, cordinatesToLocation(row+1, col))
		}
	}
	if col > 1 {
		if b.CheckLocationUnplayed(cordinatesToLocation(row, col-1)) {
			locations = append(locations, cordinatesToLocation(row, col-1))
		}
	}
	if col <= b.size-1 {
		if b.CheckLocationUnplayed(cordinatesToLocation(row, col+1)) {
			locations = append(locations, cordinatesToLocation(row, col+1))
		}
	}
	return locations
}

func (b *Board) String() string {
	output := "  "
	for i := 1; i <= b.size; i++ {
		output += fmt.Sprint(i, " ")
	}
	output += fmt.Sprintln("")
	for i, row := range b.layout {
		output += fmt.Sprintf("%c ", 'A'+i)
		output += fmt.Sprintln(strings.Join(row, " "))
	}
	return output
}

func (b *Board) SetupShips() {
	for _, ship := range b.ships.ships {
		b.shipsRemaining++
		var startrow int
		var startcol int
		for {
			failed := false
			startrow = randomInt(9)
			startcol = randomInt(10 - ship.Length)
			for col := startcol; col < startcol+ship.Length; col++ {
				if b.layout[startrow][col] != MARK_EMPTY {
					failed = true
				}
			}
			if !failed {
				break
			}
		}

		for col := startcol; col < startcol+ship.Length; col++ {
			b.layout[startrow][col] = ship.Mark
		}
	}
}

func cordinatesToLocation(x, y int) Location {
	return Location(fmt.Sprintf("%c%v", 'A'+x-1, y))
}
