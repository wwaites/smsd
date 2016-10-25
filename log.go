package main

import (
	"log"
)

func Log(cfg Config, rt Route, m Message) (err error){
	log.Printf("%s -> %s: %s", m.Src, m.Dst, m.Msg)
	return
}
