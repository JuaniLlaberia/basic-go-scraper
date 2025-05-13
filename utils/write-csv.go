package utils

import (
	"encoding/csv"
	"os"
	"scraper/structs"
)

func WriteCsv(path string, products []structs.Product) error {
	file, err := os.Create("products.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"Url",
		"Image",
		"Name",
		"Price",
	}
	writer.Write(headers)

	for _, product := range products {
		record := []string{
			product.Url,
			product.Image,
			product.Name,
			product.Price,
		}

		writer.Write(record)
	}
	defer writer.Flush()

	return nil
}