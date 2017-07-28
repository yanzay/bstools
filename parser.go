package bstools

import (
	"fmt"
	"strings"

	"github.com/yanzay/log"
)

func ParseState(text string) (State, error) {
	lines := strings.Split(text, "\n")
	var townhall, storage, houses, farm, sawmill, mine, barracks int

	for _, line := range lines {
		log.Info(line)
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
	}
	state := State{
		Townhall: townhall,
		Storage:  storage,
		Houses:   houses,
		Farm:     farm,
		Sawmill:  sawmill,
		Mine:     mine,
		Barracks: barracks,
	}
	fmt.Println(state)
	for build, level := range state {
		if level == 0 {
			return nil, fmt.Errorf("Can't parse %s", build)
		}
	}

	return state, nil
}
