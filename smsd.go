package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"syscall"
)

var laddr string
var domain string
var mta string
var uid int
var gid int

func init() {
	flag.StringVar(&laddr, "l", "127.0.0.1:8080", "listen port")
	flag.StringVar(&mta, "m", "127.0.0.1:25", "mta")
	flag.StringVar(&domain, "d", "example.org", "domain")
	flag.IntVar(&uid, "u", 978, "user id")
	flag.IntVar(&gid, "g", 978, "group id")
}

type SmsHandler bool

func (s SmsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	msg := q.Get("message")

	if from == "" || to == "" || msg == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	to = to + "@" + domain
	from = from + "@" + domain
	bmsg := []byte("To: " + to + "\r\n" +
		"From:  " + from + "\r\n" +
		"Subject: SMS Message\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		msg)

	err := InsecureSendMail(mta, from, []string{to}, bmsg)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "OK")

	return
}

func main() {
	flag.Parse()

	err := syscall.Setgid(gid)
	if err != nil {
		log.Fatal(err)
	}

	err = syscall.Setuid(uid)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	s := SmsHandler(true)

	r.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, s))
	log.Fatal(http.ListenAndServe(laddr, r))
}
