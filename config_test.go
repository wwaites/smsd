package main

import (
	"testing"
)

var example_config = `
server:
  listen: 127.0.0.1:8080
  mta:    127.0.0.1:25
  domain: example.org
routing:
  15145551212:
    mt:
      - type: log
      - type: smtp
        dst: 15145551212@example.com
    mo:
      - type: log
      - type: aamo
        src: 441315551212
`

func TestConfig(t *testing.T) {
	cfg, err := ParseConfig([]byte(example_config))
	if err != nil {
		t.Fatalf("error parsing config: %v", err)
	}

	ent, ok := cfg.Routing["15145551212"]
	if !ok {
		t.Fatalf("config missing phone number")
	}
	if len(ent.Mt) != 2 {
		t.Fatalf("wrong number of MT params: %d", len(ent.Mt))
	}
	if len(ent.Mo) != 2 {
		t.Fatalf("wrong number of MO params: %d", len(ent.Mo))
	}
	for _, rt := range ent.Mo {
		_, ok := MessageHandlers[rt.Type]
		if !ok {
			t.Fatalf("could not find handler of type %s", rt.Type)
		}
	}
	for _, rt := range ent.Mt {
		_, ok := MessageHandlers[rt.Type]
		if !ok {
			t.Fatalf("could not find handler of type %s", rt.Type)
		}
	}
}
