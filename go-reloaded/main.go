package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"go-reloaded/pkg"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run . <input_file> <output_file>")
		return
	}

	inputFile := os.Args[1]  
	outputFile := os.Args[2]

	
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	text:= string(data)


	text = pkg.Finally(text) 
	text=pkg.Final(text)
	text = pkg.ApplyCap(text)             
	text = pkg.ApplyUpper(text)         
	text = pkg.ApplyLow(text)             
	text = pkg.FormatPunctuation(text) 
	text = pkg.FixArticle(text)
	

	
	err = ioutil.WriteFile(outputFile, []byte(text), 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}

	fmt.Println("Processing complete. Check", outputFile)
}
