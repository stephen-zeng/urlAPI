package op

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"urlAPI/file"
	"urlAPI/internal/database"
)

// EmptyImageID is the special identifier serving the built-in empty image.
const EmptyImageID = "empty"

// ErrInvalidImageID is returned for identifiers that are neither a generated
// image UUID nor EmptyImageID.
var ErrInvalidImageID = errors.New("invalid image id")

var imageIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ValidImageID reports whether id may be served by /download. Only
// canonical UUIDs (as generated for tasks) and EmptyImageID are accepted,
// which rules out path separators and traversal.
func ValidImageID(id string) bool {
	return id == EmptyImageID || imageIDPattern.MatchString(id)
}

// downloadURL is the server-relative URL of a generated image.
func downloadURL(id string) string {
	return "/download?img=" + url.QueryEscape(id)
}

func downloadImage(target string) ([]byte, string, error) {
	if !ValidImageID(target) {
		return nil, "", ErrInvalidImageID
	}
	var img []byte
	var err error
	switch target {
	case EmptyImageID:
		img, err = file.EmptyPNG.ReadFile("empty.png")
	default:
		img, err = os.ReadFile(filepath.Join(ImgPath, target+".png"))
	}
	if err != nil {
		return nil, database.SettingsStore.Get().Web.FallbackImageURL, err
	}
	return img, "", nil
}
