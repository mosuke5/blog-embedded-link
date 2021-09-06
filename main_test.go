package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestGenerateLink(t *testing.T) {

}

func TestGetSiteData(t *testing.T) {
	file, _ := ioutil.ReadFile("testdata/test.html")
	stringReader := strings.NewReader(string(file))
	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		t.Error("Can not load test data")
	}

	actual := getSiteData(doc)
	expected := SiteData{
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

//func TestBuildHtml(t *testing.T) {
//	siteData := SiteData{
//		ogType:        "article",
//		ogImage:       "https://test/ogp-image.png",
//		ogTitle:       "test og title",
//		ogDescription: "test og description",
//		ogSiteName:    "my super blog",
//		ogUrl:         "https://test/"}
//
//	expected := buildHtml(siteData)
//	actual := `<div class="heml-link">
//  <div class="heml-left"><img src="https://test/ogp-image.png"></div>
//  <div class="heml-right">
//    <div class="heml-title">test og title</div>
//    <div class="heml-description">test og description</div>
//    <div class="heml-site-name">my super blog</div>
//  </div>
//</div>`
//
//	if actual != expected {
//		t.Errorf("getTitle() = '%s', but extected value is '%s'", actual, expected)
//	}
//}
