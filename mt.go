package main

import (
	"fmt"
	"log"
	"net/http"
)

type SmsMtHandler Config
/*
 * Handle incoming message, from the outside world, to be terminated
 * on a mobile handset. HTTP request should be GET with the fields
 *   - oa
 *   - da
 *   - ud
 * set. The numbers should be in international format with country
 * code and no dialling prefix. The message should be encoded in
 * unicode.
 *
 * This HTTP handler looks up the configuration for the destination
 * number and takes the appropriate action(s).
 */
func (cfg SmsMtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	m := Message{}
	m.Src = q.Get("oa")
	m.Dst = q.Get("da")
	m.Msg = q.Get("ud")

	if m.Src == "" || m.Dst == "" || m.Msg == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	rtinfo, ok := cfg.Routing[m.Dst]
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	for _, rt := range rtinfo.Mt {
		h, ok := MessageHandlers[rt.Type]
		if !ok {
			log.Printf("MT SMS could not find handler of type %s for %s", rt.Type, m.Dst)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err := h(Config(cfg), rt, m)
		if err != nil {
			log.Printf("MT SMS error sending message: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "OK")
	return
}
