package server

import (
	"net/http"
	"wordshub/services/conf"
	"wordshub/services/store"

	"github.com/gin-gonic/gin"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"

	ginSwagger "github.com/swaggo/gin-swagger"

	// swagger embed files
	_ "wordshub/docs"
)

func setRouter(cfg conf.Config) *gin.Engine {
	// Creates default gin router with Logger and Recovery middleware already attached
	router := gin.Default()

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	router.RedirectTrailingSlash = true

	// Serve static files to frontend if server is started in production environment
	if cfg.Env == "prod" {
		// router.Use(static.Serve("/", static.LocalFile("./assets/build", true)))
	}

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "success",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create API route group
	api := router.Group("/api")
	api.Use(customErrors)
	{
		api.POST("/signup", gin.Bind(store.User{}), signUp)
		api.POST("/signin", gin.Bind(store.User{}), signIn)
		api.POST("/signout", signOut)
	}

	// api.Use(authorization)
	// {
	// 	api.POST("/userinfo", userInfo)
	// }

	words := api.Group("/words")
	words.Use(authorization)
	{
		words.POST("/create", gin.Bind(store.Words{}), createWord)
	}

	// Create API route group
	catalog := api.Group("/catalog")
	catalog.Use(authorization)
	{
		catalog.GET("/list", fetchCatalog)
		catalog.POST("/create", gin.Bind(store.Catalog{}), createCatalog)
		catalog.PUT("/update", gin.Bind(store.Catalog{}), updateCatalog)
		catalog.DELETE("/delete/:id", deleteCatalog)
	}

	authorized := api.Group("/")
	authorized.Use(authorization)
	{
		authorized.GET("/posts", indexPosts)
		authorized.POST("/posts", gin.Bind(store.Post{}), createPost)
		authorized.PUT("/posts", gin.Bind(store.Post{}), updatePost)
		authorized.DELETE("/posts/:id", deletePost)
	}

	router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })

	return router
}
