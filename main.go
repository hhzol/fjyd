package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	PORT     = ":6688"
	BASE_URL = "http://zterr.hvs.fj.chinamobile.com:6060/PLTV/88888888/224/"
)

func main() {
	// 根路径重定向到 /tv
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(
			w,
			r,
			"/tv",
			http.StatusFound,
		)
	})

	// 读取接口配置文件并替换域名
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

	// 解析播放请求并直接 302 跳转到实际 m3u8 流地址
	http.HandleFunc("/fjmobile", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(
				w,
				"missing id parameter",
				http.StatusBadRequest,
			)
			return
		}

		url := BASE_URL + id + "/index.m3u8"
		playseek := r.URL.Query().Get("playseek")
		if playseek != "" {
			url += "?servicetype=3&playseek=" + playseek
		} else {
			url += "?servicetype=1"
		}

		// 执行 302 重定向到实际直播流地址
		http.Redirect(
			w,
			r,
			url,
			http.StatusFound,
		)
	})

	fmt.Println("fjmobile_api running on", PORT)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println("Server start failed:", err)
	}
}
