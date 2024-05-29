package main

import (
	"log"
	"peekaboo_tools/gemini"
)

func main() {
	//sensitive.Do()
	log.SetFlags(log.Ldate | log.Lshortfile)
	//dynamo.NewDynamoDbV1Test().Do()
	//dynamo.NewDynamoDbTest().Do()
	//gemini.NewGeminiWorker().Do()
	gemini.ParseGeminiAccountJson()
}
