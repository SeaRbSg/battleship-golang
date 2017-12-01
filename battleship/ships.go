package battleship

const (
	MARK_BATTLESHIP = "b"
	MARK_CRUISER    = "c"
	MARK_SUBMARINE  = "s"
	MARK_FRIGATE    = "f"
	MARK_DESTROYER  = "d"

	BOAT_BATTLESHIP = "battleship"
	BOAT_CRUISER    = "cruiser"
	BOAT_SUBMARINE  = "submarine"
	BOAT_FRIGATE    = "frigate"
	BOAT_DESTROYER  = "destroyer"
)

type Ship struct {
	Name   string
	Mark   string
	Length int
	Health int
}

func NewShips() *Ships {
	ships := map[string]*Ship{
		BOAT_BATTLESHIP: &Ship{BOAT_BATTLESHIP, MARK_BATTLESHIP, 5, 5},
		BOAT_CRUISER:    &Ship{BOAT_CRUISER, MARK_CRUISER, 4, 4},
		BOAT_SUBMARINE:  &Ship{BOAT_SUBMARINE, MARK_SUBMARINE, 3, 3},
		BOAT_FRIGATE:    &Ship{BOAT_FRIGATE, MARK_FRIGATE, 3, 3},
		BOAT_DESTROYER:  &Ship{BOAT_DESTROYER, MARK_DESTROYER, 2, 2},
	}
	return &Ships{ships: ships}
}

type Ships struct {
	ships map[string]*Ship
}

func (s *Ships) ShipSize(ship string) int {
	return s.ships[ship].Length
}

func (s *Ships) RemoveShip(ship string) {
	delete(s.ships, ship)
}

func (s *Ships) BiggestRemainingShip() int {
	biggest := 0
	for _, v := range s.ships {
		if v.Length > biggest {
			biggest = v.Length
		}
	}
	return biggest
}
