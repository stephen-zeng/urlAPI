package web

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	id          string
	title       string
	author      string
	description string
)

func findItem(n *html.Node) string {
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			ret += strings.TrimSpace(c.Data)
		} else if c.Type == html.ElementNode && c.Data == "a" {
			ret += findItem(c)
		}
	}
	return ret
}

func traverse(n *html.Node) {
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
		traverse(c)
	}
}

func Arxiv(URL, From, UUID string) (WebResponse, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return WebResponse{}, err
	}
	client := &http.Client{
		Timeout: time.Second * 5,
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
	id = URL[22:]
	traverse(doc)
	fmt.Println(id, title, author, description)
	return WebResponse{}, nil
}
