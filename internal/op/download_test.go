package op

import (
	"errors"
	"strings"
	"testing"
)

func TestValidImageID(t *testing.T) {
	valid := []string{"empty", "0f8fad5b-d9cb-469f-a165-70867728950e", "0F8FAD5B-D9CB-469F-A165-70867728950E"}
	invalid := []string{"", "EMPTY", "../empty", "0f8fad5b-d9cb-469f-a165-70867728950e/..", "0f8fad5bd9cb469fa16570867728950e",
		"{0f8fad5b-d9cb-469f-a165-70867728950e}", "../../etc/passwd", "..\\..\\x", "0f8fad5b-d9cb-469f-a165-70867728950e\x00"}
	for _, id := range valid {
		if !ValidImageID(id) {
			t.Errorf("ValidImageID(%q) = false", id)
		}
	}
	for _, id := range invalid {
		if ValidImageID(id) {
			t.Errorf("ValidImageID(%q) = true", id)
		}
		if _, _, err := downloadImage(id); !errors.Is(err, ErrInvalidImageID) {
			t.Errorf("downloadImage(%q) err = %v", id, err)
		}
	}
}

func TestDownloadURLIsRelative(t *testing.T) {
	got := downloadURL("0f8fad5b-d9cb-469f-a165-70867728950e")
	if !strings.HasPrefix(got, "/download?img=") {
		t.Fatalf("downloadURL = %q; generated URLs must be server-relative", got)
	}
}
