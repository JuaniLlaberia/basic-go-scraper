package scrapers

import (
	"context"
	"log"
	"scraper/structs"
	"scraper/utils"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func DynamicScraper() {
	// We are using Chromedp
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	var products []structs.Product

	// Create channel to receive products
	productChan := make(chan structs.Product)
	done := make(chan bool)

	// Go routine to collect products
	go func() {
		for product := range productChan {
			products = append(products, product)
		}
		done <- true
	}()

	err := chromedp.Run(ctx, chromedp.Navigate("https://www.scrapingcourse.com/infinite-scrolling"), scrapeProducts(productChan))
	if err != nil {
		log.Fatal(err)
	}

	close(productChan)
	<-done

	err = utils.WriteCsv("products.csv", products)
		if err != nil {
			log.Fatalln("Failed to create output CSV file", err.Error())
	}
}

func scrapeProducts(productChan chan<- structs.Product) chromedp.ActionFunc {
	return func (ctx context.Context) error {
		var prevHeight int

		// Get all nodes
		for {
			var nodes []*cdp.Node
			if err := chromedp.Nodes(".product-item", &nodes).Do(ctx); err != nil {
				return nil
			}

			// Extract data from nodes for each product
			for _, node := range nodes {
				var product structs.Product

				if err := chromedp.Run(ctx,
					chromedp.Text(".product-name", &product.Name, chromedp.ByQuery, chromedp.FromNode(node)),
					chromedp.Text(".product-price", &product.Price, chromedp.ByQuery, chromedp.FromNode(node)),
					chromedp.AttributeValue("img", "src", &product.Image, nil, chromedp.ByQuery, chromedp.FromNode(node)),
					chromedp.AttributeValue("a", "href", &product.Url, nil, chromedp.ByQuery, chromedp.FromNode(node)),
				); err != nil {
					continue
				}

				product.Price = strings.TrimSpace(product.Price)

				// Send product to channel if not empty
				if product.Name != "" {
					productChan <- product
				}
			}

			// Scroll to bottom
			var height int
			if err := chromedp.Evaluate(`document.documentElement.scrollHeight`, &height).Do(ctx); err != nil {
				return err
			}

			// Break program if we reach end of page
			if height == prevHeight {
				break
			}
			prevHeight = height

			// Scroll and wait for content to laod so we can retrieve it
			if err := chromedp.Run(ctx,
				chromedp.Evaluate(`window.scrollTo(0, document.documentElement.scrollHeight)`, nil),
				chromedp.Sleep(3*time.Second), // Wait for new content to load
			); err != nil {
				return err
			}
		}
		return nil
	}
}