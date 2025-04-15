package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {
	urls := []string{
		"https://www.bing.com", "https://duckduckgo.com", "https://search.brave.com",
		"https://www.yahoo.com", "https://www.startpage.com", "https://www.qwant.com",
		"https://www.yandex.com", "https://www.mojeek.com", "https://www.gibiru.com",
		"https://www.ecosia.org",
	}
	jobs := make(chan int, len(urls))
	result := map[string]int{}

	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go worker(jobs, result, &wg, urls, &mu)
	}
	for i := 0; i < len(urls); i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	for word, count := range result {
		fmt.Printf("%v,%d\n", word, count)
	}
}

func worker(jobs chan int, results map[string]int, waitGroup *sync.WaitGroup, urls []string, mutex *sync.Mutex) {
	defer waitGroup.Done()

	client := &http.Client{Timeout: 100 * time.Second}

	for index := range jobs {
		siteMap := make(map[string]int)
		url := urls[index]
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println("Failed to fetch from API")
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("Failed to read the body")
		}
		// convert string body to splitted string slice.
		splittedString := string(body)
		words := strings.Split(splittedString, " ")
		// fmt.Println(words)
		// count the word and add it to siteMap
		for _, word := range words {
			siteMap[word] = siteMap[word] + 1
		}

		mutex.Lock()
		// range the sitemap and add to the result
		for word, count := range siteMap {
			results[word] = results[word] + count
		}
		mutex.Unlock()
	}
}
