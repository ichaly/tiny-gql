package core

import (
	"time"
)

func (my *Engine) initDBWatcher() error {
	ke := my.Load().(*kernel)

	// no schema polling in production
	if !ke.conf.Debug {
		return nil
	}

	ps := ke.conf.PollDuration

	switch {
	case ps < (1 * time.Second):
		return nil

	case ps < (5 * time.Second):
		ps = 10 * time.Second
	}

	go func() {
		my.startDBWatcher(ps)
	}()
	return nil
}

func (my *Engine) startDBWatcher(ps time.Duration) {
	ticker := time.NewTicker(ps)
	defer ticker.Stop()

	for range ticker.C {
		ke := my.Load().(*kernel)

		di, err := GetDBInfo(ke.db, ke.dialect, ke.conf.Blocklist)
		if err != nil {
			ke.log.Println(err)
			continue
		}

		if di.Hash() == ke.di.Hash() {
			continue
		}

		ke.log.Println("database change detected. reinitializing...")

		if err := my.reload(di); err != nil {
			ke.log.Println(err)
		}

		select {
		case <-my.done:
			return
		default:
		}
	}
}
