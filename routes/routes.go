package routes

import(
	"github.com/gin-gonic/gin"
	funcs "github.com/AdityaSubrahmanyaBhat/golang/gin-go/functions"
)

func StartService(){
	router := gin.Default()
	api := router.Group("/api")
		
	api.GET("/books",funcs.GetAllBooks)
	api.POST("/books",funcs.CreateBook)
	api.GET("/book/:id",funcs.GetBook)
	api.PUT("/book/:id",funcs.UpdateBook)
	api.DELETE("/book/:id",funcs.DeleteBook)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.AbortWithStatus(404)
	})

	router.Run(":8080")
}