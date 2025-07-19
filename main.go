package main

import (
	"errors"
	"fmt"
	"strings"


	"github.com/hypermodeinc/modus/sdk/go/pkg/http"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
	
)

type Book struct {
	Title string `json:"title"`
	AuthorName []string `json:"author_name"`
	Subject []string `json:"subject"`
}

type SearchResponse struct{
	NumFound int `json:"numFound"`
	Docs []Book `json:"docs"`
}

func SearchBooks(query string) (*SearchResponse, error){
	url := fmt.Sprintf("https://openlibrary.org/search.json?q=%s&limit=5",query)

	response, err := http.Fetch(url)

	if err!= nil {
		return nil, err
	}

	if !response.Ok() {
		return nil, fmt.Errorf("search failed: %d %s", response.Status, response.StatusText)
	}

	var searchResult SearchResponse
	err = response.JSON(&searchResult)
	if err != nil {
		return nil, err
	}

	return &searchResult, nil
}

func CategorizeBook(query string) (*BookCategory, error){
	books, err := SearchBooks(query)
	if err != nil {
		return nil, err
	}
}

model, err := models.GetModel[]