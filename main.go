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
	Title	string	`json:"title"`
	AuthorName	[]string	`json:"author_name"`
	Description string   `json:"description,omitempty"`
}

type SearchResponse struct{
	NumFound	int	`json:"numFound"`
	Docs	[]Book	`json:"docs"`
}

type BookCategory struct {
	Books	[]Book	`json:"books"`
	Category	string	`json:"category"`
	Description	string	`json:"description"`
}



const modelname = "text-generator"


func SearchBooks(query string) (*SearchResponse, error){

	if strings.TrimSpace(query) == ""{
		return  nil, errors.New("query cannot be empty")
	}

	url := fmt.Sprintf("https://openlibrary.org/search.json?q=%s&limit=5",query)
	
	model, err := models.GetModel[openai.ChatModel](modelname)
	if err!= nil {
		return nil, fmt.Errorf("failed to fetch model: %w", err)
	}
	_ = model
	
	response, err := http.Fetch(url)

	if err!= nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	if !response.Ok() {
		return nil, fmt.Errorf("search failed: %d %s", response.Status, response.StatusText)
	}

	var searchResult SearchResponse
	err = response.JSON(&searchResult)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON data: %w", err)
	}

	return &searchResult, nil
}

func CategorizeBook(query string) (*BookCategory, error) {
	books, err := SearchBooks(query)
	if err != nil {
		return nil, err
	}

	model, err := models.GetModel[openai.ChatModel](modelname)
	if err != nil {
		return nil, fmt.Errorf("failed to get model: %w", err)
	}

	// Prepare book list for the prompt
	bookList := ""
	for i, book := range books.Docs {
		bookList += fmt.Sprintf("%d. %s by %s\n", i+1, book.Title, strings.Join(book.AuthorName, ", "))
	}

	systemPrompt := `You are a helpful book assistant. For each of the following books, write a one-sentence spoiler-free description. Do not include any extra commentary. Just list each book followed by its description.`
	userPrompt := fmt.Sprintf(`Books:\n%s`, bookList)

	input, err := model.CreateInput(
		openai.NewSystemMessage(systemPrompt),
		openai.NewUserMessage(userPrompt),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create model input: %w", err)
	}

	response, err := model.Invoke(input)
	if err != nil || len(response.Choices) == 0 {
		return &BookCategory{
			Books:       books.Docs,
			Category:    query,
			Description: "Could not generate descriptions due to model error.",
		}, nil
	}

	raw := strings.TrimSpace(response.Choices[0].Message.Content)
	fmt.Println("AI Response:\n", raw)

	lines := strings.Split(raw, "\n")
	for i, line := range lines {
		if i < len(books.Docs) {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				desc := strings.TrimSpace(parts[1])
				books.Docs[i].Description = desc
			}
		}
	}

	fullDescription := ""
	for _, book := range books.Docs {
		fullDescription += fmt.Sprintf("%s: %s\n", book.Title, book.Description)
	}

	return &BookCategory{
		Books:       books.Docs,
		Category:    query,
		Description: fullDescription,
	}, nil
}



