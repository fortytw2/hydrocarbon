package discollect

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func DownloadImages(textIn string, c *http.Client, fs FileStore) (string, error) {
	doc, err := html.Parse(strings.NewReader(textIn))
	if err != nil {
		return "", err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			var err error
			var src string
			src, err = findSrc(n.Attr)
			if err != nil {
				panic(err)
			}

			resp, err := c.Get(src)
			if err != nil {
				panic(err)
			}

			buf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				resp.Body.Close()
				panic(err)
			}

			resp.Body.Close()

			newSrc, err := fs.Put(src, buf)
			if err != nil {
				panic(err)
			}

			n.Attr = setSrc(n.Attr, newSrc)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	var b bytes.Buffer
	err = html.Render(&b, doc)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func findSrc(attrs []html.Attribute) (string, error) {
	for _, a := range attrs {
		if a.Key == "src" {
			return a.Val, nil
		}
	}
	return "", errors.New("no src found")
}

func setSrc(attrs []html.Attribute, src string) []html.Attribute {
	var newVals []html.Attribute
	for _, a := range attrs {
		if a.Key == "src" {
			newVals = append(newVals, html.Attribute{Key: "src", Val: src})
		} else {
			newVals = append(newVals, a)
		}
	}
	return newVals
}
