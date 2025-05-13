package scrapers

import (
	"fmt"
	"log"
	"scraper/structs"
	"scraper/utils"

	"github.com/gocolly/colly"
)

func ServerParallelScraper() {
	pagesToScrape := []string{
		"https://www.scrapingcourse.com/ecommerce/page/1/",
        "https://www.scrapingcourse.com/ecommerce/page/2/",
        "https://www.scrapingcourse.com/ecommerce/page/3/",
        "https://www.scrapingcourse.com/ecommerce/page/4/",
        "https://www.scrapingcourse.com/ecommerce/page/5/",
        "https://www.scrapingcourse.com/ecommerce/page/6/",
        "https://www.scrapingcourse.com/ecommerce/page/7/",
        "https://www.scrapingcourse.com/ecommerce/page/8/",
        "https://www.scrapingcourse.com/ecommerce/page/9/",
        "https://www.scrapingcourse.com/ecommerce/page/10/",
        "https://www.scrapingcourse.com/ecommerce/page/11/",
        "https://www.scrapingcourse.com/ecommerce/page/12/",
    }
	
	var products []structs.Product

	// Set colly to run in parallel
	c := colly.NewCollector(
		colly.AllowedDomains("www.scrapingcourse.com"),
		colly.Async(true),
	)
	// Limit colly to only have 4 parralel channels at a time
	c.Limit(&colly.LimitRule{
		Parallelism: 4,
	})

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

	// Scrapes all product information
	c.OnHTML(".product", func(e *colly.HTMLElement) {
		product := structs.Product{}

		product.Url = e.ChildAttr("a", "href")
		product.Image = e.ChildAttr("img", "src")
		product.Name = e.ChildText(".product-name")
		product.Url = e.ChildText(".price")

		products = append(products, product)
	})

	// Register all pages to scrape
	for _, pageToScrape := range pagesToScrape {
        c.Visit(pageToScrape)

        c.OnScraped(func(r *colly.Response) {
            err := utils.WriteCsv("products.csv", products)
		if err != nil {
			log.Fatalln("Failed to create output CSV file", err.Error())
		}
        })
    }

    // Wait for Colly to visit all pages
    c.Wait()
}