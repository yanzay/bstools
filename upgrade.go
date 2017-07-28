package bstools

import (
	"fmt"
	"time"
)

type Upgrade struct {
	Type     string
	Gold     int
	Stone    int
	Wood     int
	Duration time.Duration
}

func (u *Upgrade) String() string {
	return fmt.Sprintf("[%s] Type: %s, Gold: %d, Stone: %d, Wood: %d", u.Duration, u.Type, u.Gold, u.Stone, u.Wood)
}

func (u *Upgrade) TotalCost() int {
	return u.Gold + u.Stone*2 + u.Wood*2
}
