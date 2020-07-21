package v1

import "github.com/gin-gonic/gin"

// RegisteBook ...
func RegisteBook(g *gin.Engine) *gin.Engine {

	bookV1 := g.Group("/api/v1")

	bookV1.GET("/books", ListBokks)
	bookV1.GET("/books/:id", GetBook)
	bookV1.POST("/books", CreateBook)
	bookV1.PUT("/books/:id", UpdateBook)
	bookV1.DELETE("/books/:id", DeleteBook)

	return g
}

// GetBook Get a single article
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [get]
func GetBook(c *gin.Context) {}

func ListBokks(c *gin.Context) {}

func CreateBook(c *gin.Context) {}

func DeleteBook(c *gin.Context) {}

func UpdateBook(c *gin.Context) {}
