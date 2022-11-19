package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	websitePattern := getWebsitePattern()

	r.Handle("/youtube/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		content, err := getSubtitlesForVideoId(id)
		if err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf(websitePattern, err.Error(), id)))
		}
		content = prepareContent(content)
		_, _ = w.Write([]byte(fmt.Sprintf(websitePattern, content, id)))
	}))

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
