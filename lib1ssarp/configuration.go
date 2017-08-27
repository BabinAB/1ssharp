package lib1ssarp

import (
	"fmt"
)

func init() {}


type Configuration struct {
	Language string `language`
	Version string
	Output string
}

func (c Configuration) String () string {
	return fmt.Sprintf("Configuration{Language: %s, Version: %s}", c.Language, c.Version)
}


type ConfigurationBuilder interface {
	Build(c Configuration)
}