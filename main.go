package main

import (
	v1handler "golang-backend-fundamental-1/internal/api/v1/handler"
	v2handler "golang-backend-fundamental-1/internal/api/v2/handler"
	"golang-backend-fundamental-1/middlewares"
	"golang-backend-fundamental-1/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)



func main() {
 // Create a Gin router with default middleware (logger and recovery)


//  r.Use(middlewares.SimPleMiddleware())

 if err := utils.RegisterValidator(); err != nil {
	panic(err)
 }

 err := godotenv.Load()
 if err != nil {
 log.Println("Error loading .env file")
 }
 r := gin.Default()

 go middlewares.CleanupOldClients()

 r.Use(middlewares.LoggerMiddleware(),middlewares.ApiKeyMiddleware(), middlewares.RateLimitingMiddleware() )


 v1 := r.Group("/api/v1")
 v2 := r.Group("/api/v2")

 usersV1 := v1.Group("/users")
 productsV1 := v1.Group("/products")
 usersV2 := v2.Group("/users")
 categoriesV1 := v1.Group("/categories").Use(middlewares.SimPleMiddleware())
 newsV1 := v1.Group("/news")

 usersHandlerV1 := v1handler.NewUsersHandler()
 productsHandlerV1 := v1handler.NewProductsHandler()
 usersHandlerV2 := v2handler.NewUsersHandler()
 categoriesHandlerV1 := v1handler.NewCategoriesHandler()
 newsHandlerV1 := v1handler.NewNewsHandler()

 categoriesV1.GET("/", categoriesHandlerV1.GetCategoriesV1)
 categoriesV1.GET("/:id", categoriesHandlerV1.GetCategoryByIdV1)
 categoriesV1.POST("/", categoriesHandlerV1.CreateCategoryV1)
 categoriesV1.PUT("/:id", categoriesHandlerV1.UpdateCategoryV1)
 categoriesV1.DELETE("/:id", categoriesHandlerV1.DeleteCategoryV1)
 categoriesV1.GET("/code/:category_code", categoriesHandlerV1.GetCategoryByCategoryCodeV1)


 newsV1.GET("/", newsHandlerV1.GetNewsV1)
 newsV1.GET("/:slug", newsHandlerV1.GetNewsV1)
 newsV1.POST("/", newsHandlerV1.CreateNewsV1)
 newsV1.POST("/upload-flie", newsHandlerV1.UploadFileNewV1)
 newsV1.POST("/upload-multiple-files", newsHandlerV1.UploadMultipleFilesNewV1)

 usersV1.GET("/", usersHandlerV1.GetUsersV1)
 usersV1.GET("/:id", middlewares.SimPleMiddleware(), usersHandlerV1.GetUserByIdV1)
 usersV1.GET("/admin/:uuid", usersHandlerV1.GetUserByUuidV1)
 usersV1.POST("/", usersHandlerV1.CreateUserV1)
 usersV1.PUT("/:id", usersHandlerV1.UpdateUserV1)
 usersV1.DELETE("/:id", usersHandlerV1.DeleteUserV1)

 productsV1.GET("/", productsHandlerV1.GetProductsV1)
 productsV1.GET("/:id", productsHandlerV1.GetProductByIdV1)
 productsV1.POST("/", productsHandlerV1.CreateProductV1)
 productsV1.PUT("/:id", productsHandlerV1.UpdateProductV1)
 productsV1.DELETE("/:id", productsHandlerV1.DeleteProductV1)
 productsV1.GET("/products/:slug", productsHandlerV1.GetProductBySlugV1)

 usersV2.GET("/", usersHandlerV2.GetUsersV2)
 usersV2.GET("/:id", usersHandlerV2.GetUserByIdV2)
 usersV2.POST("/", usersHandlerV2.CreateUserV2)
 usersV2.PUT("/:id", usersHandlerV2.UpdateUserV2)
 usersV2.DELETE("/:id", usersHandlerV2.DeleteUserV2)


 r.StaticFS("/uploads", gin.Dir("./uploads", true))

 // Start server on port 8080 (default)
 // Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
 r.Run()
}
