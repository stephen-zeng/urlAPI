package util

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	biliBVPattern     = regexp.MustCompile(`^BV[0-9A-Za-z]{10}$`)
	biliAVPattern     = regexp.MustCompile(`^[aA][vV][0-9]{1,20}$`)
	youTubeIDPattern  = regexp.MustCompile(`^[0-9A-Za-z_-]{11}$`)
	arxivNewIDPattern = regexp.MustCompile(`^[0-9]{4}\.[0-9]{4,5}(v[0-9]+)?$`)
	arxivOldIDPattern = regexp.MustCompile(`^[a-z][a-z.-]*(\.[A-Z]{2})?/[0-9]{7}(v[0-9]+)?$`)
	repoSegment       = regexp.MustCompile(`^[0-9A-Za-z_.-]{1,100}$`)
)

// ParseWebTarget parses an absolute http(s) URL supplied by a client.
func ParseWebTarget(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errors.New("empty URL")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, errors.Wrap(err, "invalid URL")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.Errorf("unsupported URL scheme %q", u.Scheme)
	}
	if u.Hostname() == "" {
		return nil, errors.New("URL has no host")
	}
	return u, nil
}

// WebTargetHost returns the lower-cased host name (without port) of a client
// supplied URL, or "" when the URL is not a valid absolute http(s) URL.
func WebTargetHost(raw string) string {
	u, err := ParseWebTarget(raw)
	if err != nil {
		return ""
	}
	return strings.ToLower(u.Hostname())
}

// pathSegments splits a URL path into its non-empty segments.
func pathSegments(u *url.URL) []string {
	var ret []string
	for _, segment := range strings.Split(u.Path, "/") {
		if segment != "" {
			ret = append(ret, segment)
		}
	}
	return ret
}

// BilibiliVideoID extracts the BV or av identifier from a video URL such as
// https://www.bilibili.com/video/BV1xx411c7mD/?p=1.
func BilibiliVideoID(raw string) (string, error) {
	u, err := ParseWebTarget(raw)
	if err != nil {
		return "", err
	}
	segments := pathSegments(u)
	if len(segments) < 2 || segments[0] != "video" {
		return "", errors.New("not a Bilibili video URL")
	}
	id := segments[1]
	if !biliBVPattern.MatchString(id) && !biliAVPattern.MatchString(id) {
		return "", errors.Errorf("invalid Bilibili video id %q", id)
	}
	return id, nil
}

// YouTubeVideoID extracts the video id from https://www.youtube.com/watch?v=ID.
func YouTubeVideoID(raw string) (string, error) {
	u, err := ParseWebTarget(raw)
	if err != nil {
		return "", err
	}
	if strings.TrimSuffix(u.Path, "/") != "/watch" {
		return "", errors.New("not a YouTube watch URL")
	}
	id := u.Query().Get("v")
	if id == "" {
		return "", errors.New("YouTube URL has no video id")
	}
	if !youTubeIDPattern.MatchString(id) {
		return "", errors.Errorf("invalid YouTube video id %q", id)
	}
	return id, nil
}

// ArxivID extracts the paper id from https://arxiv.org/abs/ID, supporting
// both new-style (2101.00001v2) and old-style (hep-th/9901001) identifiers.
func ArxivID(raw string) (string, error) {
	u, err := ParseWebTarget(raw)
	if err != nil {
		return "", err
	}
	id, ok := strings.CutPrefix(u.Path, "/abs/")
	if !ok {
		return "", errors.New("not an arXiv abstract URL")
	}
	id = strings.TrimSuffix(id, "/")
	if id == "" {
		return "", errors.New("arXiv URL has no paper id")
	}
	if !arxivNewIDPattern.MatchString(id) && !arxivOldIDPattern.MatchString(id) {
		return "", errors.Errorf("invalid arXiv id %q", id)
	}
	return id, nil
}

// ValidRepoSegment reports whether s is a plausible GitHub/Gitee owner or
// repository name. It rejects path separators and dot-only names.
func ValidRepoSegment(s string) bool {
	return repoSegment.MatchString(s) && strings.Trim(s, ".") != ""
}

// SplitRepoInfo validates an "owner/repo" string.
func SplitRepoInfo(info string) (owner, repo string, err error) {
	owner, repo, ok := strings.Cut(info, "/")
	if !ok || !ValidRepoSegment(owner) || !ValidRepoSegment(repo) {
		return "", "", errors.Errorf("invalid repository %q, expected owner/repo", info)
	}
	return owner, repo, nil
}

// RepoFromURL extracts the owner and repository name from a repository URL
// such as https://github.com/owner/repo, https://github.com/owner/repo.git or
// https://gitee.com/owner/repo/tree/master.
func RepoFromURL(raw string) (owner, repo string, err error) {
	u, err := ParseWebTarget(raw)
	if err != nil {
		return "", "", err
	}
	segments := pathSegments(u)
	if len(segments) < 2 {
		return "", "", errors.New("repository URL must contain owner and repository")
	}
	owner, repo = segments[0], strings.TrimSuffix(segments[1], ".git")
	if !ValidRepoSegment(owner) || !ValidRepoSegment(repo) {
		return "", "", errors.Errorf("invalid repository path %q", u.Path)
	}
	return owner, repo, nil
}
