package battleship

import (
	"strconv"
	"strings"
)

type Location string

func (l Location) Row() int {
	row := strings.Index("ABCDEFGHIJKLMNOPQRSTUVWXYZ", string(l[0:1]))
	return row
}

func (l Location) Column() int {
	column, _ := strconv.Atoi(string(l[1:]))
	return (column - 1)
}
