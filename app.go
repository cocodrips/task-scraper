package main

import (
	"./scraper"
)

func main() {
	//scraper.GetMessage(scraper.EnglishCondition, scraper.EnglishOutput)
	scraper.GetMessage(scraper.TaskCondition, scraper.TaskOutput)

}
