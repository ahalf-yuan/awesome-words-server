package server

import (
	"net/http"
	"wordshub/services/conf"
	"wordshub/services/models"
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

	// router.Use(CORSMiddleware())

	// Create API route group
	api := router.Group("/api")
	api.Use(customErrors)
	{
		api.POST("/signup", gin.Bind(store.User{}), signUp)
		api.POST("/signin", gin.Bind(store.User{}), signIn)
		api.POST("/signout", signOut)
	}

	user := api.Group("/user")
	user.Use(authorization)
	{
		user.POST("/info", userInfo)
	}

	weapp := api.Group("/weapp")
	// weapp.Use(authorization)
	{
		weapp.GET("/test", func(c *gin.Context) {
			weappText(c, cfg)
		})
		weapp.GET("/refresh/login", func(c *gin.Context) {
			queryUserInfo(c, cfg)
		})
		weapp.GET("/wxcode", func(c *gin.Context) {
			getWxCode(c, cfg)
		})
	}

	words := api.Group("/")
	words.Use(authorization)
	{
		words.GET("/words", indexWords)
		words.POST("/words/create", gin.Bind(store.Words{}), createWord)
		words.DELETE("/words/:id", delWord)
		// selectedText,userId
	}

	trans := api.Group("/trans")
	trans.Use(authorization)
	{
		// trans.GET("/words", indexWords)
		trans.POST("/youdao", gin.Bind(models.YoudaoReq{}), youdao)
		trans.GET("/iciba", iciba)
		trans.GET("/iciba/sentence", icibaSentence)
		trans.GET("/iciba/sug", icibaSug)
		// trans.DELETE("/words/:id", delWord)
		// selectedText,userId
	}

	// Create API route group
	catalog := api.Group("/catalog")
	catalog.Use(authorization)
	{
		catalog.GET("/list", fetchCatalog)
		catalog.GET("/list/count", fetchCatalogAndCount)
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
