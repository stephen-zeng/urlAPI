package web

import (
	"bytes"
	"errors"
	"golang.org/x/net/html"
	"image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
	"urlAPI/cmd/txt"
	"urlAPI/internal/file"
	"urlAPI/internal/server"
)

func traverseITHome(n *html.Node, title, tim, content string) (string, string, string) {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if n.Data == "img" && attr.Key == "title" {
				title = attr.Val
			} else if n.Data == "div" && attr.Key == "class" && attr.Val == "post_content" {
				content = findItem(n)
			} else if n.Data == "span" && attr.Key == "id" && attr.Val == "pubtime_baidu" {
				tim = findItem(n)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title, tim, content = traverseITHome(c, title, tim, content)
	}
	return title, tim, content
}

func ithome(URL, From, UUID, IP, Device string, Referer *url.URL) (WebResponse, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return WebResponse{}, err
	}
	resp, err := server.GlobalHTTPClient.Do(req)
	if err != nil {
		return WebResponse{}, err
	}
	defer resp.Body.Close()
	rawResp, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Error getting article info")
		return WebResponse{}, errors.Join(err, errors.New(resp.Status))
	}
	doc, err := html.Parse(bytes.NewReader(rawResp))
	if err != nil {
		return WebResponse{}, err
	}
	title, tim, content := traverseITHome(doc, "", "", "")
	sumRet, err := txt.SumRequest(IP, From, "", "", content, Device, Referer)
	if err != nil {
		return WebResponse{}, err
	}
	description := sumRet.Response
	logoFile, err := file.LogoFS.Open("ithome_logo.png")
	logoImg, err := png.Decode(logoFile)
	err = DrawArticle(logoImg, "", title, "", description, UUID, tim)
	return WebResponse{
		Target: URL,
		URL:    From + "/download?img=" + UUID,
	}, nil
}
