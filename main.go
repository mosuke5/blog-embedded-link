package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type SiteData struct {
	OgType        string
	OgImage       string
	OgTitle       string
	OgDescription string
	OgSiteName    string
	OgUrl         string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("url arg is required.")
	}
	url := os.Args[1]

	log.Printf("url is %s", url)

	generateLink(url)
}

func generateLink(url string) {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("http status code is %d", resp.StatusCode)
		os.Exit(1)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	siteData := getSiteData(doc)
	fmt.Println(siteData)

	t, err := template.ParseFiles("template/result.html")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

	err = t.Execute(os.Stdout, siteData)
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

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
	return *siteData
}
