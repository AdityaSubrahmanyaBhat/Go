package functions

import (
	"net/http"
	book "github.com/AdityaSubrahmanyaBhat/golang/gin-go/models"
	"github.com/gin-gonic/gin"
)

func GetAllBooks(c *gin.Context){
	c.JSON(200,book.Books)
}

func GetBook(c *gin.Context){
	id := c.Param("id")
	for _, item := range book.Books{
		if id == item.ID{
			c.JSON(200,item)
			break
		}
	}
}

func CreateBook(c *gin.Context){
	var newBook book.Book
	if err := c.ShouldBindJSON(&newBook); err != nil{
		c.AbortWithStatus(http.StatusBadRequest)
	}
	book.Books = append(book.Books,newBook)
	c.JSON(http.StatusCreated,newBook)
}

func UpdateBook(c *gin.Context){
	id := c.Param("id")
	found := false
	var updatedBook book.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil{
		c.AbortWithStatus(http.StatusBadRequest)
	}

	for index, item := range book.Books{
		if id == item.ID{
			book.Books = append(book.Books[:index], book.Books[index+1:]...)
			found = true
			break
		}
	}
	if !found {
		c.AbortWithStatusJSON(http.StatusBadRequest,gin.H{
			"status":http.StatusBadRequest,
			"message":"The object you want to update does not exist",
		})
	}else{
		book.Books = append(book.Books, updatedBook)
		c.JSON(http.StatusOK,book.Books)
	}
}

func DeleteBook(c *gin.Context){
	id := c.Param("id")
	found := false
	for index, item := range book.Books{
		if id == item.ID{
			book.Books = append(book.Books[:index], book.Books[index+1:]...)
			found = true
			break
		}
	}
	if !found {
		c.AbortWithStatusJSON(http.StatusBadRequest,gin.H{
			"status":http.StatusBadRequest,
			"message":"The object you want to delete does not exist",
		})
	}else{
		c.JSON(http.StatusOK,book.Books)
	}
}