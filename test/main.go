package main

import (
	"fmt"
	requestlog "github.com/yangzhe1991/go-requestlog"
	//requestlog ".."
	"log"
	"net/http"
)

func main() {
	//loacl := requestlog.GetLocalRequestLogger("IMPublicAccount", log.Println)
	//fmt.Println(time.Now().UnixNano() / 1000)
	//time.Now().UnixNano() / 1000 / 1000
	remote := requestlog.GetRemoteRequestLogger("IMPublicAccount", "")
	fmt.Println("start")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//loacl.Log("impublicaccount", r, nil, nil)
		remote.Log("impublicaccount", r, nil, nil)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
