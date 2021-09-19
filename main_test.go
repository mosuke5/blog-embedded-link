package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/jarcoal/httpmock"
)

func TestBuildEmbededLink(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder("GET", "https://test.com/article",
		httpmock.NewStringResponder(200, httpmock.File("testdata/test.html").String()))

	oldExit := osExit

	defer func() { osExit = oldExit }()

	var status int
	exit := func(code int) {
		status = code
	}
	osExit = exit

	url := "https://test.com/article"
	actual := buildEmbededLink(url)

	if exp := 0; status != exp {
		t.Errorf("Expected exit code: %d, status: %d", exp, status)
	}

	expected := `<div class="belg-link">
  <div class="belg-left">
    <img src="https://test/ogp-image.png" />
  </div>
  <div class="belg-right">
    <div class="belg-title">test og title</div>
    <div class="belg-description">test og description</div>
    <div class="belg-site-name">my super blog</div>
  </div>
</div>`

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}
func TestBuildEmbededLinkWithWrongUrl(t *testing.T) {
	oldExit := osExit

	defer func() { osExit = oldExit }()

	var status int
	exit := func(code int) {
		status = code
	}
	osExit = exit

	url := "https://blog.mosuke.tech/notfound"
	buildEmbededLink(url)

	if exp := 1; status != exp {
		t.Errorf("Expected exit code: %d, status: %d", exp, status)
	}
}

func TestBuildEmbededLinkWithRedirectUrl(t *testing.T) {
	oldExit := osExit

	defer func() { osExit = oldExit }()

	var status int
	exit := func(code int) {
		status = code
	}
	osExit = exit

	url := "http://blog.mosuke.tech/"
	buildEmbededLink(url)

	if exp := 0; status != exp {
		t.Errorf("Expected exit code: %d, status: %d", exp, status)
	}
}

func TestBuildResultHtml(t *testing.T) {
	resultData := ResultData{
		Title:       "test og title",
		Description: "test og description",
		Favicon:     "/favicon.ico",
		Type:        "article",
		Image:       "https://test/ogp-image.png",
		SiteName:    "my super blog",
		Url:         "https://test/"}

	actual := buildResultHtml(resultData)
	expected := `<div class="belg-link">
  <div class="belg-left">
    <img src="https://test/ogp-image.png" />
  </div>
  <div class="belg-right">
    <div class="belg-title">test og title</div>
    <div class="belg-description">test og description</div>
    <div class="belg-site-name">my super blog</div>
  </div>
</div>`

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}

func TestBuildResultData(t *testing.T) {
	sitedata := SiteData{
		Title:         "test normal title",
		Description:   "test normal description",
		Favicon:       "/favicon.ico",
		OgType:        "article",
		OgImage:       "https://test/ogp-image.png",
		OgTitle:       "test og title",
		OgDescription: "test og description",
		OgSiteName:    "my super blog",
		OgUrl:         "https://test/"}

	actual := buildResultData(sitedata)
	expected := ResultData{
		Title:       "test og title",
		Description: "test og description",
		Favicon:     "/favicon.ico",
		Type:        "article",
		Image:       "https://test/ogp-image.png",
		SiteName:    "my super blog",
		Url:         "https://test/"}

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}

// すべてのデータが揃った状態でデータがとれるか
func TestGetSiteData(t *testing.T) {
	file, _ := ioutil.ReadFile("testdata/test.html")
	stringReader := strings.NewReader(string(file))
	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		t.Error("Can not load test data")
	}

	actual := getSiteData(doc)
	expected := SiteData{
		Title:         "test normal title",
		Description:   "test normal description",
		Favicon:       "/favicon.ico",
		OgType:        "article",
		OgImage:       "https://test/ogp-image.png",
		OgTitle:       "test og title",
		OgDescription: "test og description",
		OgSiteName:    "my super blog",
		OgUrl:         "https://test/"}

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}

// OGタグが一部なかった場合のテスト
func TestGetWithoutOGSiteData(t *testing.T) {
	file, _ := ioutil.ReadFile("testdata/test_without_og.html")
	stringReader := strings.NewReader(string(file))
	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		t.Error("Can not load test data")
	}

	actual := getSiteData(doc)
	expected := SiteData{
		Title:       "test normal title",
		Description: "test normal description",
		Favicon:     "/favicon.ico"}

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}
