package bstools

import (
	"math"
	"time"
)

type Coef struct {
	Gold  int
	Wood  int
	Stone int
}

var coefs = map[string]Coef{
	Townhall: Coef{Gold: 500, Wood: 200, Stone: 200},
	Storage:  Coef{Gold: 200, Wood: 100, Stone: 100},
	Houses:   Coef{Gold: 200, Wood: 100, Stone: 100},
	Farm:     Coef{Gold: 100, Wood: 50, Stone: 50},
	Sawmill:  Coef{Gold: 100, Wood: 50, Stone: 50},
	Mine:     Coef{Gold: 100, Wood: 50, Stone: 50},
	Barracks: Coef{Gold: 200, Wood: 100, Stone: 100},
}

type State struct {
	Buildings map[string]int
}

func (s *State) Apply(up *Upgrade) {
	s.Buildings[up.Type]++
}

func (s *State) BalancedUpgrade() *Upgrade {
	var minPayback time.Duration = math.MaxInt64
	var recommend string
	for build := range s.Buildings {
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
		delta := s.incomeDelta(Townhall) + (s.Buildings[Townhall]+1)*2 - 10
		payback := time.Duration(price/delta) * time.Minute
		if payback < minPayback {
			minPayback = payback
			recommend = Storage
		}
	}
	return s.calcUpgrade(recommend)
}

func (s *State) RushUpgrade() *Upgrade {
	if s.Buildings[Barracks] < s.Buildings[Houses] {
		return s.calcUpgrade(Barracks)
	}
	if up, ok := s.storageFitUpgrade(Houses); ok {
		return up
	}
	return s.calcUpgrade(Storage)
}

func copyState(state *State) *State {
	return &State{
		Buildings: map[string]int{
			Townhall: state.Buildings[Townhall],
			Storage:  state.Buildings[Storage],
			Houses:   state.Buildings[Houses],
			Farm:     state.Buildings[Farm],
			Sawmill:  state.Buildings[Sawmill],
			Mine:     state.Buildings[Mine],
			Barracks: state.Buildings[Barracks],
		},
	}
}

func (s *State) calcPayback(upType string) time.Duration {
	up := s.calcUpgrade(upType)
	price := up.TotalCost()
	if _, ok := s.storageFitUpgrade(up.Type); !ok {
		price += s.calcUpgrade(Storage).TotalCost()
	}
	delta := s.incomeDelta(upType)
	if delta != 0 {
		return time.Duration(price/delta) * time.Minute
	}
	return math.MaxInt64
}

func (s *State) incomeDelta(upType string) int {
	switch upType {
	case Townhall:
		return s.Buildings[Houses] * 2
	case Houses:
		return s.Buildings[Townhall]*2 - 10
	case Farm, Mine, Sawmill:
		return 20
	}
	return 0
}

func (s *State) storageFitUpgrade(upType string) (*Upgrade, bool) {
	up := s.calcUpgrade(upType)
	storageLvl := s.Buildings[Storage]
	storageCap := (storageLvl*50 + 1000) * storageLvl
	if up.Wood > storageCap || up.Stone > storageCap {
		return nil, false
	}
	return up, true
}

func (s *State) calcUpgrade(upType string) *Upgrade {
	level := s.Buildings[upType]
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

func (s *State) gpm() int {
	gold := s.Buildings[Houses] * (10 + s.Buildings[Townhall]*2)
	food := (s.Buildings[Farm] - s.Buildings[Houses]) * 20
	wood := s.Buildings[Sawmill] * 20
	stone := s.Buildings[Mine] * 20
	return gold + food + wood + stone
}
