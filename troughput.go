package troughput

import (
	"fmt"
	"net/http"
	"time"
)

// Stat stores info about a request
type Stat struct {
	request string
}

// Troughput is a middleware handler that logs the number of transactions per second
// of a handler
type Troughput struct {
	stats chan *Stat
}

// NewTroughput return a new Troughput instance
func NewTroughput() *Troughput {
	stats := make(chan *Stat, 16384)
	
	tr := &Troughput{
		stats: stats,
	}
	
	go tr.Log()
	
	return tr
}

// Log prints stats periodically
func (t *Troughput) Log() {
	stats := make(map[string]int, 256)
	
	tnext := time.Now().Add(1 * time.Second).Round(time.Second)
	for {
		stat := <- t.stats
		
		stats[stat.request]++
		
		// Print info, clear data, and reset the timer
		if time.Now().After(tnext) {
			for k, v := range stats {
				fmt.Printf("%s   %d TPS \n", k, v)
				delete(stats, k)
			}
			
			tnext = time.Now().Add(1 * time.Second).Round(time.Second)
		}
		
	}
}

func (t *Troughput) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
	
	t.stats <- &Stat{r.Method + " " + r.RequestURI}
}