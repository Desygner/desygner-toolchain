package extensions

import "github.com/goravel/framework/contracts/log"

type SlackLogger struct {
	logger log.Log
}

func NewSlackLogger(logger log.Log) *SlackLogger {
	return &SlackLogger{logger}
}

func (lgr *SlackLogger) Output(level int, message string) error {
	lgr.logger.Debug(message)
	return nil
}
