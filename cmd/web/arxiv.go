package web

import (
	"bytes"
	"embed"
	"errors"
	"golang.org/x/net/html"
	"image/png"
	"io"
	"log"
	"net/http"
	"time"
)

//go:embed arxiv_logo.png
var arxivFS embed.FS

func traverseArxiv(n *html.Node, title, author, description string) (string, string, string) {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if n.Data == "h1" && attr.Key == "class" && attr.Val == "title mathjax" {
				title = findItem(n)
			} else if n.Data == "div" && attr.Key == "class" && attr.Val == "authors" {
				author = findItem(n)
			} else if n.Data == "blockquote" && attr.Key == "class" && attr.Val == "abstract mathjax" {
				description = findItem(n)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title, author, description = traverseArxiv(c, title, author, description)
	}
	return title, author, description
}

func arxiv(URL, From, UUID string) (WebResponse, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return WebResponse{}, err
	}
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		return WebResponse{}, err
	}
	defer resp.Body.Close()
	rawResp, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Error getting arxiv info")
		return WebResponse{}, errors.Join(err, errors.New(resp.Status))
	}
	doc, err := html.Parse(bytes.NewReader(rawResp))
	if err != nil {
		return WebResponse{}, err
	}
	id := URL[22:]
	title, author, description := traverseArxiv(doc, "", "", "")
	logoFile, err := arxivFS.Open("arxiv_logo.png")
	logoImg, err := png.Decode(logoFile)
	err = DrawArticle(logoImg, id, title, author, description, UUID, "")
	return WebResponse{
		Target: URL,
		URL:    From + "/download?img=" + UUID,
	}, err
}
