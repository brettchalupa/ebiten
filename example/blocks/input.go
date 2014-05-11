package blocks

import (
	"github.com/hajimehoshi/go-ebiten/ui"
)

type Input struct {
	states map[ui.Key]int
}

func NewInput() *Input {
	states := map[ui.Key]int{}
	for key := ui.Key(0); key < ui.KeyMax; key++ {
		states[key] = 0
	}
	return &Input{
		states: states,
	}
}

func (i *Input) StateForKey(key ui.Key) int {
	return i.states[key]
}

func (i *Input) Update(keys []ui.Key) {
	pressedKeys := map[ui.Key]struct{}{}
	for _, key := range keys {
		pressedKeys[key] = struct{}{}
	}

	for key, _ := range i.states {
		if _, ok := pressedKeys[key]; !ok {
			i.states[key] = 0
			continue
		}
		i.states[key] += 1
	}
}