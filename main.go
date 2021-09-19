package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type SiteData struct {
	Title         string
	Description   string
	Favicon       string
	OgType        string
	OgImage       string
	OgTitle       string
	OgDescription string
	OgSiteName    string
	OgUrl         string
}

type ResultData struct {
	Title       string
	Description string
	Favicon     string
	Type        string
	Image       string
	SiteName    string
	Url         string
}

var osExit = os.Exit

func main() {
	if len(os.Args) < 2 {
		log.Fatal("url arg is required.")
	}
	url := os.Args[1]
	fmt.Print(buildEmbededLink(url))
}

func buildEmbededLink(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("http status code is %d", resp.StatusCode)
		osExit(1)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	siteData := getSiteData(doc)
	resultData := buildResultData(siteData)
	return buildResultHtml(resultData)
}

func buildResultHtml(resultData ResultData) string {
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

// SiteDataをResultDataに変換する
func buildResultData(siteData SiteData) ResultData {
	resultData := ResultData{}

	if siteData.OgTitle != "" {
		resultData.Title = siteData.OgTitle
	}

	if siteData.OgDescription != "" {
		resultData.Description = siteData.OgDescription
	}

	resultData.Favicon = siteData.Favicon
	resultData.Type = siteData.OgType
	resultData.Image = siteData.OgImage
	resultData.SiteName = siteData.OgSiteName
	resultData.Url = siteData.OgUrl

	return resultData
}

func getSiteData(d *goquery.Document) SiteData {
	siteData := new(SiteData)
	metaSelection := d.Find("head meta")
	metaSelection.Each(func(index int, s *goquery.Selection) {
		attr, _ := s.Attr("property")
		switch attr {
		case "og:type":
			ogType, _ := s.Attr("content")
			siteData.OgType = ogType
		case "og:image":
			ogImage, _ := s.Attr("content")
			siteData.OgImage = ogImage
		case "og:title":
			ogTitle, _ := s.Attr("content")
			siteData.OgTitle = ogTitle
		case "og:description":
			ogDescription, _ := s.Attr("content")
			siteData.OgDescription = ogDescription
		case "og:site_name":
			ogSiteName, _ := s.Attr("content")
			siteData.OgSiteName = ogSiteName
		case "og:url":
			ogUrl, _ := s.Attr("content")
			siteData.OgUrl = ogUrl
		}
	})

	siteData.Title = getTitle(d)
	siteData.Description = getDescription(d)
	siteData.Favicon = getFavicon(d)
	return *siteData
}

func getTitle(d *goquery.Document) string {
	return d.Find("head title").Text()
}

func getDescription(d *goquery.Document) string {
	description := ""
	metaSelection := d.Find("meta")
	metaSelection.Each(func(index int, s *goquery.Selection) {
		attr, _ := s.Attr("name")
		switch attr {
		case "description":
			description, _ = s.Attr("content")
		}
	})
	return description
}

func getFavicon(d *goquery.Document) string {
	favicon := ""
	linkSelection := d.Find("head link")
	linkSelection.Each(func(index int, s *goquery.Selection) {
		attr, _ := s.Attr("rel")
		switch attr {
		case "icon":
			favicon, _ = s.Attr("href")
		}
	})
	return favicon
}
