package handlers

import (
	"bufio"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

func InputSearch() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	term := scanner.Text()

	formattedTerm := FormatSearchTerm(term)

	return formattedTerm
}
