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

	if len(args) != 3 {
		fmt.Printf("usage: html-file [api_key]\n")
	}

	s := pdfshift.New(args[1], false)
	out, err := s.Convert(ctx, pdfshift.NewPDFBuilder().URL(args[2]))
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("golang.org.pdf", out, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
