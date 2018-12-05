package pornhub

import (
	"testing"
)

const (
	testURL   = "https://www.pornhub.com/view_video.php?viewkey=ph5b978edce2df2"
	testTitle = "Roblox Hacker Hentai: Project Ambamby"
)

func testBasic(t *testing.T, url string) (hub Pornhub) {
	vidURL = url
	err := hub.GetPage()
	if err != nil {
		t.Fatalf("Fetch/Read/Url problem, url %s: '%s'", vidURL, err.Error())
	}
	err = hub.SetTitle()
	if err != nil {
		t.Fatalf("Title problem, url %s: '%s'", vidURL, err.Error())
	}
	return
}

func TestRandVideo(t *testing.T) {
	testBasic(t, "https://pornhub.com/random")
}

func TestVideo(t *testing.T) {
	hub := testBasic(t, testURL)
	if hub.URL != testURL {
		t.Errorf("Returned URL '%s' doesn't match test url '%s'", hub.URL, testURL)
	}
	if hub.Title != testTitle {
		t.Errorf("Returned title '%s' doesn't match test title '%s'", hub.Title, testTitle)
	}
}

func testIsInvalid(t *testing.T, url string) {
	var hub Pornhub

	vidURL = url
	err := hub.GetPage()
	if err == nil {
		t.Errorf("Problem with invalid URL checking. Test value: '%s' is considered valid", url)
	}
}

func TestInvalidURL(t *testing.T) {
	invalidURLs := []string{"", "https://", pornhubURL, pornhubVidFMT}

	for _, v := range invalidURLs {
		testIsInvalid(t, v)
	}
}
