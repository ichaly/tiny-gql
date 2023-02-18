package core

import (
	"database/sql"
	"github.com/spf13/afero"
	_log "log"
	"os"
	"path/filepath"
	"sync/atomic"
)

type kernel struct {
	dialect string
	conf    *Config
	done    chan bool
	db      *sql.DB
	di      *DBInfo
	fs      FS
	opts    []Option
	log     *_log.Logger
}

type Engine struct {
	atomic.Value
	done chan bool
}

type Option func(*kernel) error

func NewEngine(conf *Config, db *sql.DB, options ...Option) (e *Engine, err error) {
	fs, err := getFS(conf)
	if err != nil {
		return
	}

	e = &Engine{done: make(chan bool)}
	if err = e.newKernel(conf, db, nil, fs, options...); err != nil {
		return
	}

	if err = e.initDBWatcher(); err != nil {
		return
	}
	return
}

func (my *Engine) newKernel(
	conf *Config, db *sql.DB, di *DBInfo, fs FS, options ...Option,
) (err error) {
	if conf == nil {
		conf = &Config{Debug: true}
	}

	ke := &kernel{
		conf: conf,
		db:   db,
		di:   di,
		fs:   fs,
		done: my.done,
		log:  _log.New(os.Stdout, "", 0),
	}
	for _, op := range options {
		if err = op(ke); err != nil {
			return
		}
	}

	my.Store(ke)
	return
}

func (my *Engine) reload(di *DBInfo) (err error) {
	ke := my.Load().(*kernel)
	err = my.newKernel(ke.conf, ke.db, di, ke.fs, ke.opts...)
	return
}

func getFS(conf *Config) (fs FS, err error) {
	if v, ok := conf.FS.(FS); ok {
		fs = v
		return
	}

	v, err := os.Getwd()
	if err != nil {
		return
	}

	fs = newAferoFS(afero.NewOsFs(), filepath.Join(v, "config"))
	return
}
