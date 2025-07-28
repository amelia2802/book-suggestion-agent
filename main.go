package main

import (
	"errors"
	"fmt"
	"strings"
	"encoding/json"

	"github.com/hypermodeinc/modus/sdk/go/pkg/http"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
	"github.com/hypermodeinc/modus/sdk/go/pkg/dgraph"

	
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

type BookNode struct {
	UID         string   `json:"uid,omitempty"`
	DgraphType  []string `json:"dgraph.type,omitempty"`
	Title       string   `json:"Book.title"`
	Description string   `json:"Book.description,omitempty"`
	Category    string   `json:"Book.category,omitempty"`
	Authors     []AuthorNode `json:"Book.authors,omitempty"`
}

type AuthorNode struct {
	UID        string   `json:"uid,omitempty"`
	DgraphType []string `json:"dgraph.type,omitempty"`
	Name       string   `json:"Author.name"`
}

type BookSearchResult struct {
	Books []BookNode `json:"books"`
}

const modelname = "text-generator"
const dgraphConnection = "dgraph"

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

	// Extract raw content from model
	raw := strings.TrimSpace(response.Choices[0].Message.Content)
	fmt.Println("AI Response:\n", raw)

	// Try to parse response assuming format like:
	// 1. Title: Description
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

	// Optional: merge all into one long string
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

func StoreBookInGraph(book Book) (*string, error) {
	// Create author nodes
	var authors []AuthorNode
	for _, authorName := range book.AuthorName {
		authors = append(authors, AuthorNode{
			DgraphType: []string{"Author"},
			Name:       authorName,
		})
	}

	// Create book node
	bookNode := BookNode{
		DgraphType:  []string{"Book"},
		Title:       book.Title,
		Description: book.Description,
		Authors:     authors,
	}

	// Convert to JSON for mutation
	bookJson, err := json.Marshal(bookNode)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal book: %w", err)
	}

	// Execute mutation
	mutation := dgraph.NewMutation().WithSetJson(string(bookJson))
	response, err := dgraph.ExecuteMutations(dgraphConnection, mutation)
	if err != nil {
		return nil, fmt.Errorf("failed to store book in graph: %w", err)
	}

	result := fmt.Sprintf("Book stored with UID: %v", response.Uids)
	return &result, nil
}

func SearchBooksFromGraph(query string) (*BookSearchResult, error) {
	dqlQuery := dgraph.NewQuery(`
		query searchBooks($query: string) {
			books(func: alloftext(Book.title, $query)) {
				uid
				Book.title
				Book.description
				Book.category
				Book.authors {
					Author.name
				}
			}
		}
	`).WithVariable("query", query)

	response, err := dgraph.ExecuteQuery(dgraphConnection, dqlQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to search books: %w", err)
	}

	var result BookSearchResult
	err = json.Unmarshal([]byte(response.Json), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	return &result, nil
}

func SearchAndStoreBooks(query string) (*BookCategory, error) {
	// First, search from external API
	books, err := SearchBooks(query)
	if err != nil {
		return nil, err
	}

	// Categorize with AI
	category, err := CategorizeBook(query)
	if err != nil {
		return nil, err
	}

	// Store each book in knowledge graph
	for _, book := range books.Docs {
		_, err := StoreBookInGraph(book)
		if err != nil {
			fmt.Printf("Warning: failed to store book '%s': %v\n", book.Title, err)
		}
	}

	return category, nil
}

func GetBookRecommendations(authorName string) (*BookSearchResult, error) {
	dqlQuery := dgraph.NewQuery(`
		query getRecommendations($author: string) {
			var(func: eq(Author.name, $author)) {
				~Book.authors {
					similar_books as Book.authors {
						recommended_books as ~Book.authors @filter(NOT eq(Author.name, $author))
					}
				}
			}
			
			books(func: uid(recommended_books)) {
				uid
				Book.title
				Book.description
				Book.authors {
					Author.name
				}
			}
		}
	`).WithVariable("author", authorName)

	response, err := dgraph.ExecuteQuery(dgraphConnection, dqlQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}

	var result BookSearchResult
	err = json.Unmarshal([]byte(response.Json), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recommendations: %w", err)
	}

	return &result, nil
}


func CheckBook(){
	query := "thriller"

	books, err := SearchBooks(query)
	if err != nil{
		fmt.Printf("Error searching books: %v\n", err)
		return
	}

	fmt.Printf("Found %d books:\n", books.NumFound)
	for _, book := range books.Docs {
		fmt.Printf("- %s by %v\n", book.Title, book.AuthorName)
	}
	
	// Categorize books
	category, err := CategorizeBook(query)
	if err != nil {
		fmt.Printf("Error categorizing books: %v\n", err)
		return
	}
	
	fmt.Printf("\nCategory: %s\n", category.Category)
	fmt.Printf("Description: %s\n", category.Description)
}