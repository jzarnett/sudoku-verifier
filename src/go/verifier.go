package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type puzzle struct {
	Content [][]int
}

// Checks the following things to see if the solution is valid
// 1. There's no digit that is < 1 or > 9
// 2. If there are any duplicate numbers in the same row
// 3. If there are any duplicate numbers in the same column
// 4. If there are any duplicate numbers in the same 3x3 grid
// Based off https://ide.geeksforgeeks.org/Gs0uFu
func checkPuzzle(p puzzle) bool {
	var rows [9][10]int
	var col [9][10]int
	var grid [3][3][10]int

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			var num = p.Content[i][j]
			if num <= 0 || num > 9 {
				return false
			}

			if rows[i][num] < 1 {
				rows[i][num] = rows[i][num] + 1
			} else {
				return false
			}
			if col[j][num] < 1 {
				col[j][num] = col[j][num] + 1
			} else {
				return false
			}

			if grid[i/3][j/3][num] < 1 {
				grid[i/3][j/3][num] = grid[i/3][j/3][num] + 1
			} else {
				return false
			}
		}
	}
	return true
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic condition detected, but going to ignore "+
				"so the server does not die:", r)
		}
	}()
	defer r.Body.Close()
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
	// 50 ms sleep time to simulate actual travel of the data over the internet
	// although it is likely that in practice clients & servers will run inside
	// the local network of the university. But well.
	time.Sleep(50 * time.Millisecond)

	if r.Method != "POST" || r.ContentLength == 0 {
		http.Error(w, "No content received.", 400)
		return
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, fmt.Sprintf("Non-JSON content detected: %s", contentType), 400)
		return
	}

	var p puzzle
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if len(p.Content) != 9 || len(p.Content[0]) != 9 {
		http.Error(w, "Must be an array of 9x9 integers", 400)
		return
	}

	if checkPuzzle(p) {
		fmt.Println("Solution valid.")
		w.WriteHeader(200)
		w.Write([]byte("1"))
		return
	}
	fmt.Println("Solution invalid.")
	w.WriteHeader(200)
	w.Write([]byte("0"))
	return
}

func main() {
	for {
		srv := &http.Server{
			Addr:         ":4590",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		srv.SetKeepAlivesEnabled(false)
		http.HandleFunc("/verify", handler)
		println("Preparing to listen on port 4590...")
		if err := http.ListenAndServe(":4590", nil); err != nil {
			// We'll log the error just for the sake of it, but because
			// students can and do kill servers our primary strategy
			// is just to ignore it so nobody has to get woken up in the
			// night to get the server back on its feet.
			log.Println("Server problem...")
			log.Println(err)
		}
	}
}
