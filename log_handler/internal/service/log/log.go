package log

import "context"

type LogHandler struct {
	prefix string
}

func New(prefix string) *LogHandler {
	return &LogHandler{prefix: prefix}
}

func (lh *LogHandler) Handle(ctx context.Context, message string) (string, error) {
	return lh.prefix + ":" + message, nil
}
