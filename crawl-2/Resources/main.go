package main

import (
	"fmt"
	"net/http"
	"os"
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

func crawl(url string, wg *sync.WaitGroup, cache *safeCache) {

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

	fmt.Println("fetching:", url)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error get: %w", err)
		return
	} else if res.StatusCode != http.StatusOK {
		fmt.Println("wrong status: %d", res.StatusCode)
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("to doc: %w", err)
		return
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists {
			return
		}

		wg.Go(func() {
			crawl(link, wg, cache)
		})
	})
}

var wg sync.WaitGroup

func main() {

	_, exists := os.LookupEnv("SERVER_IP")
	if !exists {
		fmt.Println("SERVER_IP not defined")
		return
	}

	wg.Go(func () {
		crawl(fmt.Sprintf("http://%s", SERVER_IP), &wg, &safeCache{})
	})

	wg.Wait()
}