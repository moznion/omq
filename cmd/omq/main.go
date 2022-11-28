package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/moznion/omq"
)

func main() {
	var (
		encodeToJSON   bool
		decodeFromJSON bool
		query          string
	)

	flag.BoolVar(&encodeToJSON, "j", false, "encode OpenMetrics input to JSON")
	flag.BoolVar(&decodeFromJSON, "o", false, "decode JSON input to OpenMetrics")
	flag.StringVar(&query, "q", "", "query")

	flag.Parse()

	r := io.NopCloser(bufio.NewReader(os.Stdin))
	defer func() {
		_ = r.Close()
	}()

	if query != "" {
		queried, err := omq.Query(r, query)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", queried)
		return
	}

	if decodeFromJSON {
		openMetricsText, err := omq.ConvertJSONToOpenMetricsText(r)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", openMetricsText)
		return
	}

	if encodeToJSON {
		openMetricsJson, err := omq.ConvertOpenMetricsTextToJSON(r)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", openMetricsJson)
		return
	}

	log.Fatal("TODO")
}
