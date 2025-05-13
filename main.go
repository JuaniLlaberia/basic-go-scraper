package main

import (
	"fmt"
	"scraper/scrapers"
)

func main() {
	fmt.Println("What type of scraper you want to use:")
	fmt.Println("[1]: Server content")
	fmt.Println("[2]: Server content (Parallelly)")
	fmt.Println("[3]: Dynamic content")

	var option int64
	fmt.Scan(&option)

	if option == 1 {
		fmt.Println("Srapping server content sequentially...")
		scrapers.ServerScraper()
	} else if option == 2 {
		fmt.Println("Srapping server content in parallel...")
		scrapers.ServerParallelScraper()
	} else if option == 3 {
		fmt.Println("Scrapping dynamic content...")
		scrapers.DynamicScraper()
	} else {
		fmt.Println("No option with selected number...")
	}
}