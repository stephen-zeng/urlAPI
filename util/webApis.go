package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"golang.org/x/net/html"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"urlAPI/file"
)

// 返回给你一个二进制文件

func Bili(ABV string) ([]byte, error) {
	var url string
	if ABV[0] == 'a' {
		ABV = ABV[2:]
		url = "https://api.bilibili.com/x/web-interface/view?aid=" + ABV
	} else if ABV[0] == 'B' {
		url = "https://api.bilibili.com/x/web-interface/view?bvid=" + ABV
	} else {
		return nil, errors.New("Util Bili Invalid ABV")
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Join(errors.New("Util Bili"), err)
	}
	resp, err := GlobalHTTPClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Error getting video info")
		return nil, errors.Join(errors.New("Util Bili"), err, errors.New(resp.Status))
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(errors.New("Util Bili"), err)
	}
	var info BiliResp
	err = json.Unmarshal(jsonResp, &info)
	if err != nil {
		return nil, errors.Join(errors.New("Util Bili"), err)
	}
	picURL := info.Data.Pic
	name := info.Data.Title
	author := info.Data.Owner.Name
	description := info.Data.Desc
	view := biliGetStr(info.Data.Stat.View)
	favorite := biliGetStr(info.Data.Stat.Favorite)
	like := biliGetStr(info.Data.Stat.Like)
	coin := biliGetStr(info.Data.Stat.Coin)
	ret, err := DrawVideo(picURL, name, author, description, view, favorite, like, coin)
	if err != nil {
		log.Println("Error when drawing the img")
		return nil, errors.Join(errors.New("Util Bili"), err)
	}
	return ret, nil
}

func Ytb(ID, Token string) ([]byte, error) {
	url := "https://www.googleapis.com/youtube/v3/videos?part=snippet,statistics&id=" + ID + "&key=" + Token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Join(errors.New("Util Ytb"), err)
	}
	resp, err := GlobalHTTPClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.Join(errors.New("Util Bili"), err, errors.New(resp.Status))
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Error getting video info")
		return nil, errors.Join(errors.New("Util Bili"), err)
	}
	var info YtbResp
	err = json.Unmarshal(jsonResp, &info)
	name := info.Items[0].Snippet.Title
	author := info.Items[0].Snippet.ChannelTitle
	description := info.Items[0].Snippet.Description
	picURL := info.Items[0].Snippet.Thumbnails.Standard.URL
	view := info.Items[0].Statistics.ViewCount
	like := info.Items[0].Statistics.LikeCount
	favorite := "N/A"
	coin := "N/A"
	ret, err := DrawVideo(picURL, name, author, description, view, favorite, like, coin)
	if err != nil {
		log.Println("Error when drawing the img")
		return nil, errors.Join(errors.New("Util Bili"), err)
	}
	return ret, nil
}

func Arxiv(URL string) ([]byte, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, errors.Join(errors.New("Util Arxiv"), err)
	}
	resp, err := GlobalHTTPClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Error getting arxiv info")
		return nil, errors.Join(errors.New("Util Arxiv"), err, errors.New(resp.Status))
	}
	defer resp.Body.Close()
	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(errors.New("Util Arxiv"), err, errors.New(resp.Status))
	}
	doc, err := html.Parse(bytes.NewReader(rawResp))
	if err != nil {
		return nil, errors.Join(errors.New("Util Arxiv"), err)
	}
	id := URL[22:]
	title, author, description := traverseArxiv(doc, "", "", "")
	logoFile, err := file.Logos.Open("logo/arxiv_logo.png")
	logoImg, err := png.Decode(logoFile)
	ret, err := DrawArticle(logoImg, id, title, author, description, "")
	if err != nil {
		log.Println("Error drawing the img")
		return nil, errors.Join(errors.New("Util Arxiv"), err)
	}
	return ret, nil
}

func ITHome(URL, endpoint, token, model, context string) ([]byte, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, errors.Join(errors.New("Util ITHome"), err)
	}
	resp, err := GlobalHTTPClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Error getting article info")
		return nil, errors.Join(errors.New("Util Arxiv"), err, errors.New(resp.Status))
	}
	defer resp.Body.Close()
	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(errors.New("Util Arxiv"), err)
	}
	doc, err := html.Parse(bytes.NewReader(rawResp))
	if err != nil {
		return nil, errors.Join(errors.New("Util Arxiv"), err)
	}
	title, tim, content := traverseITHome(doc, "", "", "")
	description, err := Txt(endpoint, token, model, context, content)
	logoFile, err := file.Logos.Open("assets/logo/ithome_logo.png")
	logoImg, err := png.Decode(logoFile)
	ret, err := DrawArticle(logoImg, "", title, "", description, tim)
	if err != nil {
		log.Println("Error drawing the img")
		return nil, errors.Join(errors.New("Util ITHome"), err)
	}
	return ret, nil
}

func Repo(URL string, Token string) ([]byte, error) {
	URL = strings.ReplaceAll(URL, "https://github.com", "https://api.github.com/repos")
	URL = strings.ReplaceAll(URL, "https://gitee.com", "https://gitee.com/api/v5/repos")
	req, err := http.NewRequest("GET", URL, nil)
	if Token != "" && strings.HasPrefix(URL, "https://api.github.com/repos") {
		req.Header.Set("Authorization", "Bearer "+Token)
	}
	if err != nil {
		return nil, errors.Join(errors.New("Util Github"), err)
	}
	resp, err := GlobalHTTPClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Error getting github repo info")
		return nil, errors.Join(errors.New("Util Github"), err, errors.New(resp.Status))
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(errors.New("Util Github"), err)
	}
	var repo RepoResp
	err = json.Unmarshal(jsonResp, &repo)
	if err != nil {
		return nil, errors.Join(errors.New("Util Github"), err)
	}
	author := repo.Owner.Login
	name := repo.Name
	description := repo.Description
	forkCount := getRepoCount(repo.ForksCount)
	starCount := getRepoCount(repo.StargazersCount)
	bgFile, err := file.Logos.Open("assets/logo/github_logo.png")
	bgImg, err := png.Decode(bgFile)
	if err != nil {
		log.Println("Unable to decode github background image")
		return nil, errors.Join(errors.New("Util Github"), err)
	}
	ret, err := DrawRepo(bgImg, name, author, description, starCount, forkCount)
	if err != nil {
		log.Println("Error when drawing the img")
		return nil, errors.Join(errors.New("Util Github"), err)
	}
	return ret, nil
}

func biliGetStr(x float64) string {
	if x >= 10000 {
		return strconv.FormatFloat(x/10000.0, 'f', 1, 64) + "w"
	} else {
		return strconv.FormatFloat(x, 'f', -1, 64)
	}
}

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

func findItem(n *html.Node) string {
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			ret += strings.TrimSpace(c.Data)
		} else if c.Type == html.ElementNode {
			ret += findItem(c)
		}
	}
	return ret
}

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

func getRepoCount(x float64) string {
	if x >= 1000 {
		return strconv.FormatFloat(x/1000.0, 'f', 1, 64) + "k"
	} else {
		return strconv.FormatFloat(x, 'f', -1, 64)
	}
}
