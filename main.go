package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	httpAddr   = flag.String("http", "localhost:8080", "Listen Address")
	pollPeriod = flag.Duration("poll", 10*time.Second, "Poll period")
	queryUrl   = flag.String("query", "http://globalchaosgaming.net/stuff/cwcounter.php?statsQuery", "Query Url")
	mode       = flag.String("mode", "default", "Game Mode")
	serverIp   = flag.String("serverip", "cw.gcg.io", "CW Public Server Address")
	serverName = flag.String("name", "GC Gaming cw.gcg.io connects to best of 12 servers", "Name of Server")
	location   = flag.String("location", "US", "Location of Server")
)

func main() {
	flag.Parse()
	go poll(*pollPeriod)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

type Server struct {
	Ip      string
	Max     int
	Current int
}

var data struct {
	sync.RWMutex
	current int
	max     int
}

func poll(period time.Duration) {
	for {
		log.Print("Polling")
		response, err := http.Get(*queryUrl)
		if err != nil {
			log.Print("Get Failed")
			log.Print(err)
			break
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Print("Read All failed")
			log.Print(err)
			break
		}

		response.Body.Close()

		var s struct {
			Servers []Server `json:"servers"`
		}

		err = json.Unmarshal(body, &s)
		if err != nil {
			log.Print(err)
			time.Sleep(period)
			continue
		}

		current := 0
		max := 0
		for _, server := range s.Servers {
			current += server.Current
			max += server.Max
		}

		data.Lock()
		data.current = current
		data.max = max
		data.Unlock()

		time.Sleep(period)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	data.RLock()
	response := struct {
		Current  int    `json:"players"`
		Max      int    `json:"max"`
		Name     string `json:"name"`
		Mode     string `json:"mode"`
		Ip       string `json:"ip"`
		Location string `json:"location"`
	}{
		Current:  data.current,
		Max:      data.max,
		Name:     *serverName,
		Mode:     *mode,
		Ip:       *serverIp,
		Location: *location,
	}
	data.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(response)
}
