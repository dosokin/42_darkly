package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var SERVER_IP, _ = os.LookupEnv("SERVER_IP")

type safeCache struct {
	mu sync.Mutex
	alreadyParsed []string
}

var PAGES_FETCHED = 0

func checkReadMe(reader io.Reader) (string) {

	body, err := io.ReadAll(reader)
	if err != nil {
		return ""
	}

	re := regexp.MustCompile(`[a-zA-Z0-9]{64}`)
	match := re.Find(body)
	if match != nil {
		return string(match)
	}

	return ""
}

func crawl(url string, cache *safeCache) {

	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("http://%s/%s", SERVER_IP, url)
	}

	if !strings.Contains(url, SERVER_IP) {
		return
	}

	if slices.Contains(cache.alreadyParsed, url){
		return
	} else {
		cache.mu.Lock()
		cache.alreadyParsed = append(cache.alreadyParsed, url)
		cache.mu.Unlock()
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error get:", err)
		return
	} else if res.StatusCode != http.StatusOK {
		fmt.Println("wrong status:", res.StatusCode, "[", url, "]")
		return
	}
	defer res.Body.Close()

	if strings.Contains(url, "README") {
		flag := checkReadMe(res.Body)
		if flag != "" {
			fmt.Printf("Fetched %d pages", PAGES_FETCHED)
			fmt.Println("FLAG FOUND =>", flag)
			fmt.Println("URL =>", url)
			os.Exit(0)
		}
		return
	}

	PAGES_FETCHED += 1
	fmt.Printf("Fetched %d pages\r", PAGES_FETCHED)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("to doc: %w", err)
		return
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists || strings.Contains(link, ".."){
			return
		}
		crawl(url + link, cache)
	})
}

func main() {

	SERVER_IP, exists := os.LookupEnv("SERVER_IP")
	if !exists {
		fmt.Println("SERVER_IP not defined")
		return
	}

	crawl(fmt.Sprintf("http://%s/.hidden/", SERVER_IP), &safeCache{})
}