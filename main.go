package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	urls := flag.String("urls", "", "urls to parse")

	flag.Parse()
	if *urls == "" {
		fmt.Println("empty -urls")
		os.Exit(1)
	}
	fmt.Println("urls: ", *urls)

	urlsList := strings.Split(*urls, ",")

	wg := new(sync.WaitGroup)

	for _, url := range urlsList {
		wg.Add(1)
		go func() {
			err := sendRequest(url, wg)
			if err != nil {
				fmt.Println("send request to", url, err)
				os.Exit(1)
			}
		}()
		wg.Wait()
	}

}

func sendRequest(url string, wg *sync.WaitGroup) error {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(url, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected code %d", resp.StatusCode)
	}
	err = parseHTML(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func parseHTML(r io.Reader) error {
	utfBody, err := iconv.NewReader(r, "windows-1252", "utf-8")
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		return err
	}

	var links []Links
	doc.Find("h1").Each(func(i int, p *goquery.Selection) {
		link := Links{}
		link.header = p.Find("h1").Text()
		link.link = p.Find("a").Text()
		links = append(links, link)
	})

	doc.Find("h2").Each(func(i int, p *goquery.Selection) {
		link := Links{}
		link.header = p.Find("h2").Text()
		link.link = p.Find("a").Text()
		links = append(links, link)
	})

	doc.Find("h3").Each(func(i int, p *goquery.Selection) {
		link := Links{}
		link.header = p.Find("h3").Text()
		link.link = p.Find("a").Text()
		links = append(links, link)
	})

	doc.Find("h4").Each(func(i int, p *goquery.Selection) {
		link := Links{}
		link.header = p.Find("h4").Text()
		link.link = p.Find("a").Text()
		links = append(links, link)
	})

	doc.Find("h5").Each(func(i int, p *goquery.Selection) {
		link := Links{}
		link.header = p.Find("h5").Text()
		link.link = p.Find("a").Text()
		links = append(links, link)
	})
	fmt.Println(links)
	return nil
}

type Links struct {
	header, link string
}
