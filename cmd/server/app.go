package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/enicho/go-socket/util"
	"github.com/gomodule/redigo/redis"

	"github.com/gorilla/mux"
	grace "gopkg.in/paytm/grace.v1"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

var (
	matchName map[string]string
)

func main() {
	matchName = make(map[string]string)
	router := mux.NewRouter()
	router.HandleFunc("/v1/register", Register).Methods("GET")

	grace.Serve(":9002", router)
}

func Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "err: %v", err)
	}

	param1 := r.URL.Query().Get("name")
	if param1 != "" {
		result, err := redis.String(redis.DoWithTimeout(util.GetClient(), 10*time.Second, "GET", param1))
		if err == nil {
			if !strings.EqualFold(result, "") {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, result)
				return
			}
		}
		if val, ok := matchName[param1]; ok {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, val)
			return
		}
		matchID := RandStringBytes(20)
		matchName[param1] = matchID

		redis.DoWithTimeout(util.GetClient(), 10*time.Second, "SET", param1, matchID)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, matchName[param1])
		return

	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
	return
}
