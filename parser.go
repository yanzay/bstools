package bstools

import (
	"fmt"
	"strings"
)

func ParseState(text string) (State, error) {
	lines := strings.Split(text, "\n")
	var townhall, storage, houses, farm, sawmill, mine, barracks, wall, trebuchet int

	for _, line := range lines {
		if townhall == 0 {
			fmt.Sscanf(line, "🏤   %d", &townhall)
		}
		if storage == 0 {
			fmt.Sscanf(line, "🏚   %d", &storage)
		}
		if houses == 0 {
			fmt.Sscanf(line, "🏘   %d", &houses)
		}
		if farm == 0 {
			fmt.Sscanf(line, "🌻   %d", &farm)
		}
		if sawmill == 0 {
			fmt.Sscanf(line, "🌲   %d", &sawmill)
		}
		if mine == 0 {
			fmt.Sscanf(line, "⛏   %d", &mine)
		}
		if barracks == 0 {
			fmt.Sscanf(line, "🛡   %d", &barracks)
		}
		if wall == 0 {
			fmt.Sscanf(line, WallEmoji+"   %d", &wall)
		}
		if trebuchet == 0 {
			fmt.Sscanf(line, TrebuchetEmoji+"Требушет %d", &trebuchet)
		}
		if trebuchet == 0 {
			fmt.Sscanf(line, TrebuchetEmoji+"Trebuchet %d", &trebuchet)
		}
	}
	state := State{
		Townhall:  townhall,
		Storage:   storage,
		Houses:    houses,
		Farm:      farm,
		Sawmill:   sawmill,
		Mine:      mine,
		Barracks:  barracks,
		Wall:      wall,
		Trebuchet: trebuchet,
	}

	return state, nil
}
