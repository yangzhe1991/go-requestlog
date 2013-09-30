package main

import (
	requestlog ".."
	"log"
	"net/http"
)

func main() {
	loacl := requestlog.GetLocalRequestLogger("IMPublicAccount", log.Println)
	remote := requestlog.GetRemoteRequestLogger("IMPublicAccount", "")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		loacl.Log("test", r, nil, nil)
		remote.Log("test", r, nil, nil)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
