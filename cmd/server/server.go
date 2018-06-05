package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/thijzert/unihand"
)

func main() {
	var zipPath string
	var listen string

	flag.StringVar(&zipPath, "zip", "Unihan.zip", "Path to Unihan.zip")
	flag.StringVar(&listen, "listen", "localhost:8978", "Listen address")
	flag.Parse()

	log.Printf("Initialising database...")
	if err := unihand.Initialise(zipPath); err != nil {
		log.Fatal(err)
	}
	log.Printf("Finished initialisation")
	log.Printf("%d characters loaded", unihand.LoadedCharacters())

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("Memory usage: %.2fM", float64(m.Alloc+m.StackSys)/1048576.0)

	mux := http.NewServeMux()
	mux.HandleFunc("/char/", characterInfo)
	mux.HandleFunc("/status", printStatus)

	log.Printf("Starting HTTP server on %s...", listen)
	log.Fatal(http.ListenAndServe(listen, mux))
}

func characterInfo(w http.ResponseWriter, r *http.Request) {
	var code uint32
	_, err := fmt.Sscanf(r.URL.Path, "/char/%x", &code)

	if err != nil {
		badRequest(w, err)
		return
	}

	char, err := unihand.Lookup(code)
	if err != nil {
		badRequest(w, err)
		return
	}

	writeJSON(w, char)
}

func printStatus(w http.ResponseWriter, r *http.Request) {
	var status struct {
		CharactersLoaded int
		MemoryUsed       float64
	}

	status.CharactersLoaded = unihand.LoadedCharacters()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	status.MemoryUsed = float64(m.Alloc+m.StackSys) / 1048576.0

	// Truncate to two decimals
	status.MemoryUsed = float64(int64(100.0*status.MemoryUsed)) * 0.01

	writeJSON(w, status)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	var enc = json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(v)
}

func badRequest(w http.ResponseWriter, e error) {
	w.WriteHeader(400)
	if e != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(e.Error()))
	}
}
func notFound(w http.ResponseWriter, e error) {
	w.WriteHeader(404)
	if e != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(e.Error()))
	} else {
		w.Write([]byte("The requested document could not be found"))
	}
}
