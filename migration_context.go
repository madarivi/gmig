package main

import (
	"fmt"
	"path/filepath"

	"github.com/urfave/cli"
)

type migrationContext struct {
	lastApplied   string
	stateProvider StateProvider
}

func getMigrationContext(c *cli.Context) (ctx migrationContext, err error) {
	// allow target as folder name
	target := filepath.Base(c.Args().First())
	if len(target) == 0 {
		err = fmt.Errorf("missing target name in command line")
		return
	}
	stateProvider, err := getStateProvider(c)
	if err != nil {
		return
	}
	err = gcloudConfigSetProject(stateProvider.Config())
	if err != nil {
		return
	}
	lastApplied, err := stateProvider.LoadState()
	if err != nil {
		return
	}
	ctx.stateProvider = stateProvider
	ctx.lastApplied = lastApplied
	if len(lastApplied) > 0 {
		e := checkExists(lastApplied)
		if e != nil {
			err = e
			return
		}
	}
	return
}

func (m migrationContext) config() Config {
	return m.stateProvider.Config()
}