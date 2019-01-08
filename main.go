package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type UsersResponse struct {
	Users []User `json:"users"`
}

type User struct {
	Logo string `json:"logo"`
}

func main() {
	httpClient := &http.Client{
		Timeout: time.Second,
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		login := r.URL.Query().Get("login")
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.twitch.tv/v5/users/?login=%s", url.QueryEscape(login)), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		req.Header.Set("Client-Id", "dkgtqxwkpbsxg6ffitog9l33mdtnif")
		response, err := httpClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Println(err)
			return
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		var users UsersResponse
		err = json.Unmarshal(body, &users)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if len(users.Users) != 1 {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("missing %s", login)
			return
		}
		w.Header().Set("location", users.Users[0].Logo)
		w.WriteHeader(http.StatusFound)
	})
	if err := http.ListenAndServe(":8040", nil); err != nil {
		panic(err)
	}
}
