package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Subscriber func(mux *http.ServeMux)

func NewSubscriber() Subscriber {
	videoSplit := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("ReadAll: %v", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		output := fmt.Sprintf("Completed task: payload(%s)", string(body))
		log.Println(output)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, output)
	}

	return func(mux *http.ServeMux) {
		mux.HandleFunc("/video_split", videoSplit)
	}
}
