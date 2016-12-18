package web

import "net/http"

func renderFeed(w http.ResponseWriter, r *http.Request) {
	out, err := TMPLERRfeed("Hydrocarbon", false, 0)
	if err != nil {
		panic(err)
	}

	_, err = w.Write([]byte(out))
	if err != nil {
		panic(err)
	}
}

func renderPost(w http.ResponseWriter, r *http.Request) {
	out, err := TMPLERRpost("Hydrocarbon", false, 0)
	if err != nil {
		panic(err)
	}

	_, err = w.Write([]byte(out))
	if err != nil {
		panic(err)
	}
}

func reorderFeeds(w http.ResponseWriter, r *http.Request) {

}

func addFeed(w http.ResponseWriter, r *http.Request) {

}

func deleteFeed(w http.ResponseWriter, r *http.Request) {

}
