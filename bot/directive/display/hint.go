package display

import (
	"github.com/johnnyeven/bot-sdk-go/bot/data"
	"github.com/johnnyeven/bot-sdk-go/bot/util"
)

type Hint struct {
	Type  string        `json:"type"`
	Hints []data.Speech `json:"hints"`
}

func NewHint(hint ...string) *Hint {
	h := &Hint{
		Type:  "Hint",
		Hints: make([]data.Speech, 0),
	}

	h.SetHints(hint)

	return h
}

func (h *Hint) SetHints(hints []string) *Hint {
	for _, value := range hints {
		h.AddHint(value)
	}

	return h
}

func (h *Hint) AddHint(hint string) *Hint {
	h.Hints = append(h.Hints, util.FormatSpeech(hint))

	return h
}
