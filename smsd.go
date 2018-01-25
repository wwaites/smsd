package main

import (
	"flag"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"
)

var config_file string
type MessageHandler func(Config, Route, Message) error
var MessageHandlers map[string]MessageHandler

func init() {
	flag.StringVar(&config_file, "c", "", "Configuration File")
	MessageHandlers = make(map[string]MessageHandler)
	MessageHandlers["log"]  = Log
	MessageHandlers["smtp"] = SendMail
	MessageHandlers["aamt"] = AndrewsArnoldMt
	MessageHandlers["pushover"] = Pushover
	MessageHandlers["aamo"] = AndrewsArnoldMo
}

type Message struct {
	Src string
	Dst string
	Msg string
}

func main() {
	flag.Parse()

	if len(config_file) == 0 {
		flag.Usage()
		log.Fatal("config file is a required argument")
	}

	config_data, err := ioutil.ReadFile(config_file)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := ParseConfig(config_data)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Server.Gid > 0 {
		err := syscall.Setgid(cfg.Server.Gid)
		if err != nil {
			log.Fatal(err)
		}
	}

	if cfg.Server.Uid > 0 {
		err = syscall.Setuid(cfg.Server.Uid)
		if err != nil {
			log.Fatal(err)
		}
	}

	r := mux.NewRouter()
	mt := SmsMtHandler(cfg)
	mo := SmsMoHandler(cfg)

	r.Handle("/mt", handlers.CombinedLoggingHandler(os.Stderr, mt))
	r.Handle("/mo", handlers.CombinedLoggingHandler(os.Stderr, mo))
	log.Fatal(http.ListenAndServe(cfg.Server.Listen, r))
}
