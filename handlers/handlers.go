package handlers

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func InputSearch() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	term := scanner.Text()

	formattedTerm := FormatSearchTerm(term)

	return formattedTerm
}

func MakeRequest(url string) *goquery.Document {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		panic(err)
	}

	return document
}

func GetPages(document *goquery.Document, url string) []string {
	pages := 0
	var pageLinks []string

	document.Find(".pagination__number").Each(func(i int, s *goquery.Selection) {
		pages++
	})

	if pages > 0 {
		for i := 1; i <= pages; i++ {
			newLink := url + "&page=" + strconv.Itoa(i)
			pageLinks = append(pageLinks, newLink)
		}
	} else {
		pageLinks = append(pageLinks, url)
	}

	return pageLinks
}

func DownloadJersey(link string, title string, refererLink string) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("User-Agent", "Mozilla/5.0")
	request.Header.Set("Referer", refererLink)

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	path := "/home/lucas/jerseys/" + title

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Baixado com sucesso\n")
}

func FormatSearchTerm(term string) string {
	var newTerm string

	s := strings.ToLower(term)

	for i := 0; i < len(s); i++ {
		if string(s[i]) == " " {
			newTerm += "+"
		} else {
			newTerm += string(s[i])
		}
	}

	return newTerm
}
