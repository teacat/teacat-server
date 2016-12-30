package logging

import (
	"time"

	nsq "github.com/bitly/go-nsq"
)

// Uppercase 會紀錄 Uppercase 函式的相關資訊。
func (mw Middleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.Service.Uppercase(s)
	return
}

// Count 會紀錄 Count 函式的相關資訊。
func (mw Middleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

	n = mw.Service.Count(s)
	return
}

func (mw Middleware) Test(msg *nsq.Message) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "test",
			"input", msg.Body,
			"took", time.Since(begin),
		)
	}(time.Now())

	mw.Service.Test(msg)
	return
}
