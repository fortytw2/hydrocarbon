package rss

import (
)

// func TestRSS(t *testing.T) {
// 	var cases = []struct {
// 		URL   string
// 		Title string
// 	}{
// 		{
// 			"http://feeds.arstechnica.com/arstechnica/features",
// 			"Features – Ars Technica",
// 		},
// 	}

// 	r := Reader{Client: http.DefaultClient}

// 	for _, c := range cases {
// 		t.Run(c.Title, func(t *testing.T) {
// 			title, url, err := r.Info(context.Background(), "http://feeds.arstechnica.com/arstechnica/features")
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			if title == "" || url == "" {
// 				t.Fatal("no title or URL found")
// 			}

// 			if title != c.Title {
// 				t.Fatal("titles do not match")
// 			}

// 			_, err = r.Fetch(context.Background(), url, time.Time{})
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 		})
// 	}
}
