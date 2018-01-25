package main

import (
	"gopkg.in/yaml.v2"
)

type Route struct {
	Type  string // handler type 
	Match string // regular expression to match dst (MO)
	Src   string // src number to use (MO)
	Dst   string // dst address/uri to use (SMTP)
	User  string // username for 3rd-party service
	Pass  string // password for 3rd-party service
	Key   string // key for allowing access to this route (MO)
}

type Config struct {
	Server struct {
		Listen string
		Domain string
		Pushover string
		Mta string
		Uid int
		Gid int
	}
	Dialplans map[string]Dialplan
	Routing map[string]struct {
		Dialplan string
		Mt []Route
		Mo []Route
	}
}

func ParseConfig(data []byte) (m Config, err error) {
	m = Config{}
	err = yaml.Unmarshal(data, &m)
	return
}
