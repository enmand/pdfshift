package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/enmand/pdfshift"
)

func main() {
	ctx := context.Background()
	args := os.Args

	if len(args) != 2 {
		fmt.Printf("usage: html-file [api_key]\n")
	}

	url := "https://example.com"

	s := pdfshift.New(args[1])
	out, err := s.Convert(ctx, pdfshift.NewPDFBuilder().Sandbox(false).URL(url))
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("example.com.pdf", out, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
