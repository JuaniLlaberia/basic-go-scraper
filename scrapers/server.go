package scrapers

import (
	"fmt"
	"log"
	"scraper/structs"
	"scraper/utils"
	"sync"

	"github.com/gocolly/colly"
)

func ServerScraper() {
	var products []structs.Product

	var visitedUrls sync.Map

	c := colly.NewCollector(
		colly.AllowedDomains("www.scrapingcourse.com"),
	)

	// Faking an user agent
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

	// Called before a HTPP request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// When an error occures
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong", err)
	})

	// Fired when the server responds
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited", r.Request.URL)
	})

	// Fired when a CSS selector matches an element
	c.OnHTML(".product", func(e *colly.HTMLElement) {
		product := structs.Product{}

		product.Url = e.ChildAttr("a", "href")
		product.Image = e.ChildAttr("img", "src")
		product.Name = e.ChildText(".product-name")
		product.Url = e.ChildText(".price")

		products = append(products, product)
	})
	c.OnHTML("a.next", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")

		if _, found := visitedUrls.Load(nextPage); !found {
			visitedUrls.Store(nextPage, struct{}{})
			e.Request.Visit(nextPage)
		}
	})

	// When the scraping finishes
	c.OnScraped(func(r *colly.Response) {
		err := utils.WriteCsv("products.csv", products)
		if err != nil {
			log.Fatalln("Failed to create output CSV file", err.Error())
		}
	})

	c.Visit("https://www.scrapingcourse.com/ecommerce")
}