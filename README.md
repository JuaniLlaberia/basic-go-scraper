﻿# Basic GO web scrapers
Implemented 3 types of scrapers in golang using "colly" and "chromedp" packages. All three scrapers scrape products from a demo page (https://www.scrapingcourse.com):
- #1: Scrapes server content.
- #2: Scrapes server content using go routines to run in parallel (more efficient).
- #3: Scrapes dynamic generated content, in this case the products are being loaded as we scroll.


