package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type SmsMoHandler Config
/*
 * Handle incoming message, from the outside world, to be terminated
 * on a mobile handset. HTTP request should be GET with the fields
 *   - from
 *   - to
 *   - msg
 *   - key
 * set. The numbers should be in international format with country
 * code and no dialling prefix. The message should be encoded in
 * unicode.
 *
 * This HTTP handler looks up the configuration for the destination
 * number and takes the appropriate action(s).
 */
func (cfg SmsMoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("MO SMS could not parse form")
		http.Error(w, "Internal ServerError", http.StatusInternalServerError)
	}

	log.Printf("FORM: %v", r.Form)
	m := Message{}

	m.Src = r.Form.Get("oa")
	m.Dst = r.Form.Get("da")
	m.Msg = r.Form.Get("ud")

	if m.Src == "" || m.Dst == "" || m.Msg == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	rtinfo, ok := cfg.Routing[m.Src]
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if rtinfo.Dialplan != "" {
		dialplan, ok := cfg.Dialplans[rtinfo.Dialplan]
		if !ok {
			log.Printf("MO SMS could not find dialplan %s for %s", rtinfo.Dialplan, m.Src)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		m.Dst = dialplan.Normalise(m.Dst)
	}

	for _, rt := range rtinfo.Mo {
		h, ok := MessageHandlers[rt.Type]
		if !ok {
			log.Printf("MO SMS could not find handler of type %s for %s", rt.Type, m.Src)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if rt.Match != "" {
			matched, err := regexp.MatchString(rt.Match, m.Dst)
			if err != nil {
				log.Printf("MO SMS error matching destination for MO SMS: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if !matched {
				continue
			}
		}

		err := h(Config(cfg), rt, m)
		if err != nil {
			log.Printf("MO SMS error sending message: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "OK")
	return
}
