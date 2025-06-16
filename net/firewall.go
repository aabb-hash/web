package net

import (
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"
)

var heat map[string]uint8 = make(map[string]uint8)
var asyncBlocked []string

var blocked []string

var firewallLock sync.RWMutex
var analyzeLock sync.Mutex

func CheckRequest(next http.Handler) http.Handler {
	firewallLock.RLock()
	defer firewallLock.RUnlock()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		address := r.RemoteAddr

		fmt.Println(blocked)
		fmt.Println(address)
		fmt.Println(slices.Contains(blocked, address))

		if slices.Contains(blocked, address) {
			fmt.Println("hj")
			hijack(w)
			return
		}

		go analyze(address)
		next.ServeHTTP(w, r)
	})
}

func analyze(address string) {
	analyzeLock.Lock()

	heatVal := heat[address]
	heat[address] = heatVal + 1

	if heatVal > 5 {
		asyncBlocked = append(asyncBlocked, address)
	}

	analyzeLock.Unlock()
}

func hijack(w http.ResponseWriter) {
	hijacker, _ := w.(http.Hijacker)

	conn, _, err := hijacker.Hijack()
	if err != nil {
		close(w)
		return
	}

	conn.Close()
}

func close(w http.ResponseWriter) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}

func InitFirewall() {
	go decreaseHeat()
	go syncBlockedAddresses()
}

func decreaseHeat() {
	ticker := time.NewTicker(10 * time.Second)

	for {
		<-ticker.C
		go func() {
			analyzeLock.Lock()
			for address, value := range heat {
				if value > 0 {
					heat[address] = value - 1
				}
			}
			analyzeLock.Unlock()
		}()
	}
}
func syncBlockedAddresses() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		<-ticker.C
		go func() {
			firewallLock.Lock()
			blocked = asyncBlocked
			firewallLock.Unlock()
		}()
	}
}
