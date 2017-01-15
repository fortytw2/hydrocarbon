package web

import (
	"net/http"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/internal/httputil"
)

func renderFeed(s *hydrocarbon.Store) httputil.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		feedID := r.URL.Query().Get("id")

		if feedID == "" {
			fs, err := s.Feeds.GetFeeds(&hydrocarbon.Pagination{
				Page:     0,
				PageSize: 20,
			})
			if err != nil {
				return httputil.Wrap(err, 404)
			}

			out := TMPLfeeds("Hydrocarbon", false, 0, fs)

			_, err = w.Write([]byte(out))
			return err
		}

		f, err := s.Feeds.GetFeed(feedID)
		if err != nil {
			return httputil.Wrap(err, 404)
		}

		posts, err := s.Posts.GetPosts(f.ID, &hydrocarbon.Pagination{
			Page:     0,
			PageSize: 10,
		})
		if err != nil {
			return err
		}

		out := TMPLfeed("Hydrocarbon", false, 0, f, posts)

		_, err = w.Write([]byte(out))
		return err
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
