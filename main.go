package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	baseURL := "https://www.thesaurus.com.br/"
	c := colly.NewCollector()
	catCollector := colly.NewCollector()
	productCollector := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("nav.categories", func(e *colly.HTMLElement) {
		e.ForEach("a", func (i int, e *colly.HTMLElement) {
			title := e.Attr("title")
			link := e.Attr("href")
			fmt.Printf("%s %s\n", title, link)
			catCollector.Visit(e.Request.AbsoluteURL(link))
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	catCollector.OnHTML("ul.listItems", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, e *colly.HTMLElement) {
			link := e.ChildAttr("figure a", "href")
			productCollector.Visit(e.Request.AbsoluteURL(link))
		})
	})

	productCollector.OnHTML("div.detailProduct", func(e *colly.HTMLElement) {
		title := e.ChildText("h1")
		fmt.Println("Scrapping product " + title)
	})

	catCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Category", r.URL)
	})

	productCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Product", r.URL)
	})

	c.Visit(baseURL)
}
