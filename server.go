package main

import (
	"encoding/json"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Server struct {
	listenAddr string
	nameFolder string
	handler    http.Handler
	countries  []Country
}

func NewServer(la, nf string) *Server {
	a := Server{listenAddr: la, nameFolder: nf}
	m := http.NewServeMux()
	m.HandleFunc("/search", a.SearchHandler)
	m.HandleFunc("/ac", a.AutoCompleteHandler)

	m.HandleFunc("/countries", a.CountriesHandler)
	m.HandleFunc("/", a.RootHandler)

	c := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://indexbrain.org"},
		AllowedOrigins: []string{"*"},
	})
	h := c.Handler(m)
	a.handler = h

	a.GetCountries()

	return &a
}

func (s *Server) GetCountries() {
	dirEntries, err := os.ReadDir(s.nameFolder)
	if err != nil {
		panic(err)
	}
	d := []Country{}
	for _, de := range dirEntries {
		if de.IsDir() {
			d = append(d, NewCountry(de.Name()))
		}
	}
	s.countries = d
}

func (s *Server) LookupCountry(abbr string) (Country, bool) {
	for _, e := range s.countries {
		if e.Abbr == abbr {
			return e, true
		}
	}
	return Country{}, false
}

func (s *Server) Run() {
	log.Println("Names at", s.nameFolder)
	log.Println("Listening at", s.listenAddr)
	log.Println("Server listening at", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, s.handler))
}

// HANDLERS
func (s *Server) CountriesHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	e := json.NewEncoder(w)

	err := e.Encode(s.countries)
	if err != nil {
		panic(err)
	}
}

func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("IndexBrain Server is running. FIXME: add message about website"))
}

func (s *Server) SearchHandler(w http.ResponseWriter, r *http.Request) {
	s.SearchMiddle(w, r, false)
}

func (s *Server) AutoCompleteHandler(w http.ResponseWriter, r *http.Request) {
	s.SearchMiddle(w, r, true)
}

func (s *Server) SearchMiddle(w http.ResponseWriter, r *http.Request, ac bool) {
	query := r.URL.Query()
	q, hasQ := query["q"]
	// FIX HERE \/
	country, hasCountry := query["country"]
	entryType, hasEntryType := query["type"]

	// FIX HERE /\
	if !hasQ {
		w.WriteHeader(400)
		w.Write([]byte("No search term provided"))
		return
	}

	// actual search logic here
	c := ""
	if hasCountry {
		c = country[0]
	}
	t := ""
	if hasEntryType {
		t = entryType[0]
	}

	res, ok := s.SearchInternal(q[0], c, t, ac)

	if !ok {
		w.WriteHeader(500)
		w.Write([]byte("[]"))
	} else {
		w.WriteHeader(200)
		e := json.NewEncoder(w)
		err := e.Encode(res)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Server) SearchInternal(q string, country string, eType string, ac bool) ([]Entry, bool) {
	// You need to validate user input, like all of it...
	// FIXME: YOU NEED TO ESCAPE THIS, OR ELSE THERE WILL BE DIRE CONSEQUENCES

	// pass in country as an abbr

	// folder :=
	cfolder := ""
	if country != "" {
		c, ok := s.LookupCountry(country)
		if !ok {
			// Couldn't find the country
			log.Println("Couldn't find the country with abbreviation", country)
			return []Entry{}, false
		}
		cfolder = c.Folder()
	}

	// TODO: Down the line, we should probably filter the files we search by the type, but for now it's not worth it...

	folder := filepath.Join(s.nameFolder, cfolder)
	log.Println("Searching folder, ", folder)

	num := "10"

	if ac == false {
		num = "10000"
	}

	cmd := exec.Command("rg", "-m", num, q, folder)
	log.Println("COMMAND ARGS")
	log.Println(cmd.String())
	for _, i := range cmd.Args {
		log.Println(i)
	}

	out, err := cmd.Output()
	if err != nil {
		// Error with running rg
		if err.Error() == "exit status 1" {
			// No results
			return []Entry{}, true
		} else {
			log.Println(err)
			return []Entry{}, false
		}
	} else {
		entries := s.ParseEntries(string(out), cfolder, eType)
		return entries, true
	}
}

func (s *Server) ParseEntries(text string, c string, eType string) []Entry {
	ent := []Entry{}

	lines := strings.Split(text, "\n")
	for _, l := range lines {
		parts := strings.SplitN(l, ":", 2) // First part is path, second part actual name
		if len(parts) != 2 {
			log.Println("Parsing error: Returned length is not 2", parts, "original: ", l)
			continue
		}
		// Get country and type from path aIMHERE
		f, txt := filepath.Split(parts[0])
		f = filepath.Clean(f)
		_, cnf := filepath.Split(f)
		// cnf is country folder name, txt is file name (get type from that), parts[1] is the entry
		c := NewCountry(cnf)
		t, ok := verifyType(txt, c.Abbr)
		if !ok {
			log.Println("Rejecting line from ", txt)
			continue
		}

		ent = append(ent, Entry{
			Name:    parts[1],
			Type:    t,
			Country: c.Abbr,
		})
	}
	return ent
}

func verifyType(fn string, cab string) (string, bool) {
	prf := strings.Title(strings.ToLower(cab))
	res := strings.TrimPrefix(fn, prf)
	res = strings.TrimSuffix(res, ".txt")
	if len(res) != 1 {
		log.Println("INVALID FILE", fn)
		return "", false
	}
	switch res {
	case "N":
		return "name", true
	case "P":
		return "place", true
	case "M":
		return "misc", true
	default:
		log.Println("Unknown type", res)
		return "", false
	}
}
