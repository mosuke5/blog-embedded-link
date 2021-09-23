package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
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

	url := "https://test.com/article"
	resultData := ResultData{
		Title:       "test og title",
		Description: "test og description",
		Favicon:     "https://test.com/favicon.ico",
		Type:        "article",
		Image:       "https://test.com/ogp-image.png",
		SiteName:    "my super blog",
		Url:         "https://test.com/article",
		BaseUrl:     "https://test.com/"}

	actual := buildEmbededLink(url)
	expected := expectedHtml(resultData)

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}

func TestBuildEmbededLinkWithoutOg(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder("GET", "https://test.com/article",
		httpmock.NewStringResponder(200, httpmock.File("testdata/test_without_og.html").String()))

	url := "https://test.com/article"
	resultData := ResultData{
		Title:       "test normal title",
		Description: "test normal description",
		Favicon:     "https://test.com/favicon.ico",
		SiteName:    "test.com",
		Url:         "https://test.com/article",
		BaseUrl:     "https://test.com/"}

	actual := buildEmbededLink(url)
	expected := expectedHtml(resultData)

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}

func TestBuildEmbededLinkNoFavicon(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder("GET", "https://test.com/article",
		httpmock.NewStringResponder(200, httpmock.File("testdata/test_no_favicon.html").String()))

	url := "https://test.com/article"
	resultData := ResultData{
		Title:       "test og title",
		Description: "test og description",
		Type:        "article",
		Image:       "https://test.com/ogp-image.png",
		SiteName:    "my super blog",
		Url:         "https://test.com/article",
		BaseUrl:     "https://test.com/"}

	actual := buildEmbededLink(url)
	expected := expectedHtml(resultData)

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}

func TestBuildEmbededLinkWithNotFoundUrl(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder("GET", "https://test.com/notfound",
		httpmock.NewStringResponder(404, "not found"))

	oldExit := osExit
	defer func() { osExit = oldExit }()

	var status int
	exit := func(code int) {
		status = code
	}
	osExit = exit

	url := "https://test.com/notfound"
	buildEmbededLink(url)

	if exp := 1; status != exp {
		t.Errorf("Expected exit code: %d, status: %d", exp, status)
	}
}

// リダイレクトするサイトの場合、エラーを返さずリダイレクト先のページで情報を取得すること
func TestBuildEmbededLinkWithRedirectUrl(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://test.com/article",
		httpmock.NewStringResponder(200, httpmock.File("testdata/test.html").String()))
	httpmock.RegisterResponder("GET", "https://short-url.com/xxx",
		responderWithLocationHeader(301, "foo", "https://test.com/article"))

	url := "https://short-url.com/xxx"
	resultData := ResultData{
		Title:       "test og title",
		Description: "test og description",
		Favicon:     "https://test.com/favicon.ico",
		Type:        "article",
		Image:       "https://test.com/ogp-image.png",
		SiteName:    "my super blog",
		Url:         "https://test.com/article",
		BaseUrl:     "https://test.com/"}

	buildEmbededLink(url)
	actual := buildEmbededLink(url)
	expected := expectedHtml(resultData)

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}

// ResultDataから適切なHTMLを生成できること
//func TestBuildResultHtml(t *testing.T) {
//	resultData := ResultData{
//		Title:       "test og title",
//		Description: "test og description",
//		Favicon:     "https://test.com/favicon.ico",
//		Type:        "article",
//		Image:       "https://test.com/ogp-image.png",
//		SiteName:    "my super blog",
//		Url:         "https://test.com/article",
//		BaseUrl:     "https://test.com"}
//
//	actual := buildResultHtml(resultData)
//	expected := `<div class="belg-link">
//  <div class="belg-left">
//    <img src="https://test.com/ogp-image.png" />
//  </div>
//  <div class="belg-right">
//    <div class="belg-title">
//      <a href="https://test.com/article" target="_blank">test og title</a>
//    </div>
//    <div class="belg-description">test og description</div>
//    <div class="belg-site">
//      <img src="https://test.com/favicon.ico" class="belg-site-icon">
//      <span class="belg-site-name">my super blog</span>
//    </div>
//  </div>
//</div>`
//
//	if actual != expected {
//		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
//	}
//}

// OGの情報が全て揃っている場合。該当項目は上書きを想定
func TestBuildResultData(t *testing.T) {
	sitedata := SiteData{
		RequestUrl:     "https://test.com/article",
		RequestBaseUrl: "https://test.com",
		Title:          "test normal title",
		Description:    "test normal description",
		Favicon:        "/favicon.ico",
		OgType:         "article",
		OgImage:        "https://test.com/ogp-image.png",
		OgTitle:        "test og title",
		OgDescription:  "test og description",
		OgSiteName:     "my super blog",
		OgUrl:          "https://test.com/article"}

	actual := buildResultData(sitedata)
	expected := ResultData{
		Title:       "test og title",
		Description: "test og description",
		Favicon:     "https://test.com/favicon.ico",
		Type:        "article",
		Image:       "https://test.com/ogp-image.png",
		SiteName:    "my super blog",
		Url:         "https://test.com/article",
		BaseUrl:     "https://test.com"}

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}

// OGの情報が全て揃っていない場合。該当項目のみ上書きを想定
func TestBuildResultDataWithoutOg(t *testing.T) {
	sitedata := SiteData{
		RequestUrl:     "https://test.com/article",
		RequestBaseUrl: "https://test.com",
		Title:          "test normal title",
		Description:    "test normal description",
		Favicon:        "/favicon.ico"}

	actual := buildResultData(sitedata)
	expected := ResultData{
		Title:       "test normal title",
		Description: "test normal description",
		Favicon:     "https://test.com/favicon.ico",
		Image:       "",
		SiteName:    "test.com",
		Url:         "https://test.com/article",
		BaseUrl:     "https://test.com"}

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}

// Faviconがないとき
func TestBuildResultDataNoFavicon(t *testing.T) {
	sitedata := SiteData{
		RequestUrl:     "https://test.com/article",
		RequestBaseUrl: "https://test.com",
		Title:          "test normal title",
		Description:    "test normal description"}

	actual := buildResultData(sitedata)
	expected := ResultData{
		Title:       "test normal title",
		Description: "test normal description",
		Favicon:     "",
		Image:       "",
		SiteName:    "test.com",
		Url:         "https://test.com/article",
		BaseUrl:     "https://test.com"}

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}

func TestBuildResultDataWithFaviconFromExternal(t *testing.T) {
	sitedata := SiteData{
		RequestUrl:     "https://test.com/article",
		RequestBaseUrl: "https://test.com",
		Title:          "test normal title",
		Description:    "test normal description",
		Favicon:        "//some-cdn.com/favicon.ico",
		OgType:         "article",
		OgImage:        "https://test.com/ogp-image.png",
		OgTitle:        "test og title",
		OgDescription:  "test og description",
		OgSiteName:     "my super blog",
		OgUrl:          "https://test.com/article"}

	actual := buildResultData(sitedata)
	expected := ResultData{
		Title:       "test og title",
		Description: "test og description",
		Favicon:     "//some-cdn.com/favicon.ico",
		Type:        "article",
		Image:       "https://test.com/ogp-image.png",
		SiteName:    "my super blog",
		Url:         "https://test.com/article",
		BaseUrl:     "https://test.com"}

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}

// HTMLから必要なデータを取得しSiteDataを作成できるか
func TestGetSiteData(t *testing.T) {
	file, _ := ioutil.ReadFile("testdata/test.html")
	stringReader := strings.NewReader(string(file))
	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		t.Error("Can not load test data")
	}
	siteData := SiteData{}
	actual := getSiteData(siteData, doc)
	expected := SiteData{
		Title:         "test normal title",
		Description:   "test normal description",
		Favicon:       "/favicon.ico",
		OgType:        "article",
		OgImage:       "https://test.com/ogp-image.png",
		OgTitle:       "test og title",
		OgDescription: "test og description",
		OgSiteName:    "my super blog",
		OgUrl:         "https://test.com/article"}

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}

// HTMLから必要なデータを取得しSiteDataを作成できるか（OGタグが一部なかった場合）
func TestGetWithoutOGSiteData(t *testing.T) {
	file, _ := ioutil.ReadFile("testdata/test_without_og.html")
	stringReader := strings.NewReader(string(file))
	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		t.Error("Can not load test data")
	}

	siteData := SiteData{}
	actual := getSiteData(siteData, doc)
	expected := SiteData{
		Title:       "test normal title",
		Description: "test normal description",
		Favicon:     "/favicon.ico"}

	if actual != expected {
		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
	}
}

func expectedHtml(resultData ResultData) string {
	t, err := template.ParseFiles("template/result.html")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

	var result bytes.Buffer
	err = t.Execute(&result, resultData)
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	return result.String()
}

func responderWithLocationHeader(s int, c string, location string) httpmock.Responder {
	resp := httpmock.NewStringResponse(s, c)
	resp.Header.Set("Location", location)
	return httpmock.ResponderFromResponse(resp)
}
