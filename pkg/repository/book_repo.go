package repository

import "apigateway/pkg/model"

// BookRepository ...
type BookRepository interface {
	// Get book by option
	GetBook() (model.Book, error)
	// // List books by option
	// ListBooks() ([]model.Book, error)
	// // Create a Book
	// CreateBook(book model.Book) (model.Book, error)
	// // Delete the Book
	// DeleteBook() error
	// // Update the Book
	// UpdateBook() (model.Book, error)
}

func (repo *repository) GetBook() (model.Book, error) {
	var book model.Book
	return book, nil
}
