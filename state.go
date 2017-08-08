package bstools

import (
	"fmt"
	"math"
	"time"
)

type Coef struct {
	Gold  int
	Wood  int
	Stone int
}

var allBuilds = []string{
	Townhall,
	Storage,
	Houses,
	Farm,
	Sawmill,
	Mine,
	Barracks,
	Wall,
	Trebuchet,
}

var coefs = map[string]Coef{
	Townhall:  Coef{Gold: 500, Wood: 200, Stone: 200},
	Storage:   Coef{Gold: 200, Wood: 100, Stone: 100},
	Houses:    Coef{Gold: 200, Wood: 100, Stone: 100},
	Farm:      Coef{Gold: 100, Wood: 50, Stone: 50},
	Sawmill:   Coef{Gold: 100, Wood: 50, Stone: 50},
	Mine:      Coef{Gold: 100, Wood: 50, Stone: 50},
	Barracks:  Coef{Gold: 200, Wood: 100, Stone: 100},
	Wall:      Coef{Gold: 5000, Wood: 500, Stone: 1500},
	Trebuchet: Coef{Gold: 8000, Wood: 1000, Stone: 300},
}

type State map[string]int

func (s State) Apply(up *Upgrade) {
	s[up.Type]++
}

func (s State) BalancedUpgrade() *Upgrade {
	var minPayback time.Duration = math.MaxInt64
	var recommend string
	for build := range s {
		payback := s.calcPayback(build)
		if payback < minPayback {
			minPayback = payback
			recommend = build
		}
	}
	_, okHouse := s.storageFitUpgrade(Houses)
	_, okTownhall := s.storageFitUpgrade(Townhall)
	if !okHouse && !okTownhall {
		storageUp := s.calcUpgrade(Storage)
		housesUp := s.calcUpgrade(Houses)
		townhallUp := s.calcUpgrade(Townhall)
		price := storageUp.TotalCost() + housesUp.TotalCost() + townhallUp.TotalCost()
		delta := s.IncomeDelta(Townhall) + (s[Townhall]+1)*2 - 10
		payback := time.Duration(price/delta) * time.Minute
		if payback < minPayback {
			minPayback = payback
			recommend = Storage
		}
	}
	return s.calcUpgrade(recommend)
}

func (s State) RushUpgrade() *Upgrade {
	if s[Barracks] < s[Houses] {
		return s.calcUpgrade(Barracks)
	}
	if up, ok := s.storageFitUpgrade(Houses); ok {
		return up
	}
	return s.calcUpgrade(Storage)
}

func (s State) BattleUpgrade() *Upgrade {
	if s[Houses] < s[Barracks] {
		return s.calcUpgrade(Houses)
	}
	if s[Wall]%2 != 0 {
		if up, ok := s.storageFitUpgrade(Wall); ok {
			return up
		}
		return s.calcUpgrade(Storage)
	}
	if s[Trebuchet]%2 != 0 {
		if up, ok := s.storageFitUpgrade(Trebuchet); ok {
			return up
		}
		return s.calcUpgrade(Storage)
	}
	housesCost := s.calcUpgrade(Houses).TotalCost()
	barracksCost := s.calcUpgrade(Barracks).TotalCost()
	storageCost := s.calcUpgrade(Storage).TotalCost()
	wallCost := s.calcUpgrade(Wall).TotalCost()
	trebCost := s.calcUpgrade(Trebuchet).TotalCost()
	var storageForBarracks, storageForWall, storageForTreb bool

	// barracks
	var barracksPpw int
	if s[Houses] == s[Barracks] {
		barracksPpw = housesCost + barracksCost
		if _, ok := s.storageFitUpgrade(Houses); !ok {
			barracksPpw += storageCost
			storageForBarracks = true
		}
		barracksPpw = barracksPpw / 40
	} else {
		barracksPpw = barracksCost / 40
	}

	// wall
	state := copyState(s)
	if _, ok := state.storageFitUpgrade(Wall); !ok {
		wallCost += storageCost
		state[Storage]++
		storageForWall = true
	}
	state[Wall]++
	wallCost += state.calcUpgrade(Wall).TotalCost()
	if _, ok := state.storageFitUpgrade(Wall); !ok {
		wallCost += state.calcUpgrade(Storage).TotalCost()
	}
	wallPpw := wallCost / 100

	// treb
	state = copyState(s)
	if _, ok := state.storageFitUpgrade(Trebuchet); !ok {
		trebCost += storageCost
		state[Storage]++
		storageForTreb = true
	}
	state[Trebuchet]++
	trebCost += state.calcUpgrade(Trebuchet).TotalCost()
	if _, ok := state.storageFitUpgrade(Trebuchet); !ok {
		trebCost += state.calcUpgrade(Storage).TotalCost()
	}
	trebPpw := trebCost / 100

	// decision
	if barracksPpw < wallPpw && barracksPpw < trebPpw {
		if storageForBarracks {
			return s.calcUpgrade(Storage)
		}
		return s.calcUpgrade(Barracks)
	}
	if wallPpw < barracksPpw && wallPpw < trebPpw {
		if storageForWall {
			return s.calcUpgrade(Storage)
		}
		return s.calcUpgrade(Wall)
	}
	if storageForTreb {
		return s.calcUpgrade(Storage)
	}
	return s.calcUpgrade(Trebuchet)
}

func (s State) Merge(state State) {
	for _, build := range allBuilds {
		s[build] = max(s[build], state[build])
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func copyState(state State) State {
	return State{
		Townhall:  state[Townhall],
		Storage:   state[Storage],
		Houses:    state[Houses],
		Farm:      state[Farm],
		Sawmill:   state[Sawmill],
		Mine:      state[Mine],
		Barracks:  state[Barracks],
		Wall:      state[Wall],
		Trebuchet: state[Trebuchet],
	}
}

func (s State) calcPayback(upType string) time.Duration {
	up := s.calcUpgrade(upType)
	price := up.TotalCost()
	if _, ok := s.storageFitUpgrade(up.Type); !ok {
		price += s.calcUpgrade(Storage).TotalCost()
	}
	delta := s.IncomeDelta(upType)
	if delta != 0 {
		return time.Duration(price/delta) * time.Minute
	}
	return math.MaxInt64
}

func (s State) IncomeDelta(upType string) int {
	switch upType {
	case Townhall:
		return s[Houses] * 2
	case Houses:
		return s[Townhall]*2 - 10
	case Farm, Mine, Sawmill:
		return 20
	}
	return 0
}

func (s State) storageFitUpgrade(upType string) (*Upgrade, bool) {
	up := s.calcUpgrade(upType)
	storageLvl := s[Storage]
	storageCap := (storageLvl*50 + 1000) * storageLvl
	if up.Wood > storageCap || up.Stone > storageCap {
		return nil, false
	}
	return up, true
}

func (s State) calcUpgrade(upType string) *Upgrade {
	level := s[upType]
	k := (level + 1) * (level + 2) / 2
	up := &Upgrade{
		Type:  upType,
		Gold:  k * coefs[upType].Gold,
		Wood:  k * coefs[upType].Wood,
		Stone: k * coefs[upType].Stone,
	}
	up.Duration = time.Minute * time.Duration((up.TotalCost() / s.gpm()))
	return up
}

func (s State) gpm() int {
	gold := s[Houses] * (10 + s[Townhall]*2)
	food := (s[Farm] - s[Houses]) * 20
	wood := s[Sawmill] * 20
	stone := s[Mine] * 20
	return gold + food + wood + stone
}

func (s State) Valid() error {
	for _, build := range allBuilds {
		if s[build] == 0 {
			return fmt.Errorf("I need to know your %s", build)
		}
	}
	return nil
}
