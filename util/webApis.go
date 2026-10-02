package util

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"urlAPI/file"
)

// 返回给你一个二进制文件

// Upstream API endpoints; variables so tests can point them at a local server.
var (
	bilibiliViewAPI = "https://api.bilibili.com/x/web-interface/view"
	youTubeVideoAPI = "https://www.googleapis.com/youtube/v3/videos"
	arxivAbsBase    = "https://arxiv.org/abs/"
	githubReposAPI  = "https://api.github.com/repos/"
	giteeReposAPI   = "https://gitee.com/api/v5/repos/"
)

// Bili renders a card for a Bilibili video given its BV or av identifier.
func Bili(ABV string) ([]byte, error) {
	query := url.Values{}
	switch {
	case biliBVPattern.MatchString(ABV):
		query.Set("bvid", ABV)
	case biliAVPattern.MatchString(ABV):
		query.Set("aid", ABV[2:])
	default:
		return nil, errors.New("Util Bili Invalid ABV")
	}
	req, err := http.NewRequest("GET", bilibiliViewAPI+"?"+query.Encode(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp, err := GlobalHTTPClient.Do(req)
	switch {
	case err != nil:
		return nil, errors.WithStack(err)
	case resp.StatusCode != http.StatusOK:
		return nil, errors.WithStack(errors.New(resp.Status))
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var info BiliResp
	err = json.Unmarshal(jsonResp, &info)
	if err != nil {
		return nil, errors.WithStack(err)
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
		return nil, errors.WithStack(err)
	}
	return ret, nil
}

// Ytb renders a card for a YouTube video given its id.
func Ytb(ID, Token string) ([]byte, error) {
	if !youTubeIDPattern.MatchString(ID) {
		return nil, errors.New("Util Ytb Invalid video id")
	}
	query := url.Values{}
	query.Set("part", "snippet,statistics")
	query.Set("id", ID)
	query.Set("key", Token)
	req, err := http.NewRequest("GET", youTubeVideoAPI+"?"+query.Encode(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp, err := GlobalHTTPClient.Do(req)
	switch {
	case err != nil:
		return nil, errors.WithStack(err)
	case resp.StatusCode != http.StatusOK:
		return nil, errors.WithStack(errors.New(resp.Status))
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var info YtbResp
	if err = json.Unmarshal(jsonResp, &info); err != nil {
		return nil, errors.WithStack(err)
	}
	if len(info.Items) == 0 {
		return nil, errors.New("YouTube video not found")
	}
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
		return nil, errors.WithStack(err)
	}
	return ret, nil
}

// Arxiv renders a card for an arXiv paper given its identifier.
func Arxiv(id string) ([]byte, error) {
	if !arxivNewIDPattern.MatchString(id) && !arxivOldIDPattern.MatchString(id) {
		return nil, errors.New("Util Arxiv Invalid id")
	}
	req, err := http.NewRequest("GET", arxivAbsBase+id, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp, err := GlobalHTTPClient.Do(req)
	switch {
	case err != nil:
		return nil, errors.WithStack(err)
	case resp.StatusCode != http.StatusOK:
		return nil, errors.WithStack(errors.New(resp.Status))
	}
	defer resp.Body.Close()
	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, resp.Status)
	}
	doc, err := html.Parse(bytes.NewReader(rawResp))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	title, author, description := traverseArxiv(doc, "", "", "")
	logoImg, err := loadLogo("logo/arxiv_logo.png")
	if err != nil {
		return nil, err
	}
	ret, err := DrawArticle(logoImg, id, title, author, description, "")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ret, nil
}

// ITHome renders a card for an ITHome article, summarised by a text model.
func ITHome(URL, endpoint, token, model, context string) ([]byte, error) {
	target, err := ParseWebTarget(URL)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(target.Hostname(), "www.ithome.com") {
		return nil, errors.New("Util ITHome Invalid URL")
	}
	// Only the path of the client URL is used; scheme and host are fixed.
	page := url.URL{Scheme: "https", Host: "www.ithome.com", Path: target.Path}
	req, err := http.NewRequest("GET", page.String(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp, err := GlobalHTTPClient.Do(req)
	switch {
	case err != nil:
		return nil, errors.WithStack(err)
	case resp.StatusCode != http.StatusOK:
		return nil, errors.WithStack(errors.New(resp.Status))
	}
	defer resp.Body.Close()
	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	doc, err := html.Parse(bytes.NewReader(rawResp))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	title, tim, content := traverseITHome(doc, "", "", "")
	description, err := Txt(endpoint, token, model, context, content)
	if err != nil {
		return nil, errors.WithMessage(err, "summarize article")
	}
	logoImg, err := loadLogo("logo/ithome_logo.png")
	if err != nil {
		return nil, err
	}
	ret, err := DrawArticle(logoImg, "", title, "", description, tim)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ret, nil
}

// Repo renders a card for a GitHub or Gitee repository URL.
func Repo(URL string, Token string) ([]byte, error) {
	target, err := ParseWebTarget(URL)
	if err != nil {
		return nil, err
	}
	owner, repoName, err := RepoFromURL(URL)
	if err != nil {
		return nil, err
	}
	var apiURL, logoURL string
	isGitHub := false
	switch strings.ToLower(target.Hostname()) {
	case "github.com":
		apiURL = githubReposAPI + url.PathEscape(owner) + "/" + url.PathEscape(repoName)
		logoURL = "logo/github_logo.png"
		isGitHub = true
	case "gitee.com":
		apiURL = giteeReposAPI + url.PathEscape(owner) + "/" + url.PathEscape(repoName)
		logoURL = "logo/gitee_logo.png"
	default:
		return nil, errors.New("Util Repo unsupported host")
	}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if Token != "" && isGitHub {
		req.Header.Set("Authorization", "Bearer "+Token)
	}
	resp, err := GlobalHTTPClient.Do(req)
	switch {
	case err != nil:
		return nil, errors.WithStack(err)
	case resp.StatusCode != http.StatusOK:
		return nil, errors.WithStack(errors.New(resp.Status))
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var repo RepoResp
	if err = json.Unmarshal(jsonResp, &repo); err != nil {
		return nil, errors.WithStack(err)
	}
	author := repo.Owner.Login
	name := repo.Name
	description := repo.Description
	forkCount := getRepoCount(repo.ForksCount)
	starCount := getRepoCount(repo.StargazersCount)
	bgImg, err := loadLogo(logoURL)
	if err != nil {
		return nil, err
	}
	ret, err := DrawRepo(bgImg, name, author, description, starCount, forkCount)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ret, nil
}

// loadLogo decodes an embedded PNG logo.
func loadLogo(name string) (image.Image, error) {
	logoFile, err := file.Logos.Open(name)
	if err != nil {
		return nil, errors.Wrapf(err, "open logo %s", name)
	}
	defer logoFile.Close()
	img, err := png.Decode(logoFile)
	if err != nil {
		return nil, errors.Wrapf(err, "decode logo %s", name)
	}
	return img, nil
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
