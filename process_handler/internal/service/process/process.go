package process

import "context"

type ProcessHandler struct {
}

func New(prefix string) *ProcessHandler {
	return &ProcessHandler{}
}

func (lh *ProcessHandler) Handle(ctx context.Context, message string) (string, error) {
	runes := []rune(message)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}
