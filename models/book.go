package book

type Book struct{
	ID string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	ISBN string `json:"isbn"`
	Year int `json:"year"`
}

var Books = []Book{
	{ID: "1", Title: "Harry Potter", Author: "J. K. Rowling",Year: 2011},
	{ID: "2", Title: "The Lord of the Rings", Author: "J. R. R. Tolkien",Year: 1999},
	{ID: "3", Title: "The Wizard of Oz", Author: "L. Frank Baum",Year: 0001},
}