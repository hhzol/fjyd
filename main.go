package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

const (
	PORT = ":6688"
	BASE_URL = "http://zterr.hvs.fj.chinamobile.com:6060/PLTV/88888888/224/"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(
			w,
			r,
			"/tv",
			http.StatusFound,
		)
	})
	http.HandleFunc("/tv", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("interface.txt")
		if err != nil {
			http.Error(
				w,
				"interface.txt not found",
				http.StatusNotFound,
			)
			return
		}
		txt := strings.ReplaceAll(
			string(data),
			"${replace}",
			r.Host,
		)
		w.Header().Set(
			"Content-Type",
			"text/plain; charset=utf-8",
		)
		w.Write([]byte(txt))
	})

	http.HandleFunc("/fjmobile", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)
			json.NewEncoder(w).Encode(
				map[string]string{
					"error":"missing id",
				},
			)
			return
		}
		url := BASE_URL +
			id +
			"/index.m3u8"
		playseek := r.URL.Query().Get("playseek")
		if playseek != "" {
			url +=
				"?servicetype=3&playseek=" +
				playseek
		} else {
			url +=
				"?servicetype=1"
		}
		w.Header().Set(
			"Content-Type",
			"application/json",
		)
		json.NewEncoder(w).Encode(
			map[string]string{
				"url":url,
			},
		)
	})

	println("fjmobile_api running on", PORT)

	http.ListenAndServe(
		PORT,
		nil,
	)
}