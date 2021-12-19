package server

import (
	"net/http"
	"wordshub/services/store"

	"github.com/gin-gonic/gin"
)

func createCatalog(ctx *gin.Context) {
	catalog := ctx.MustGet(gin.BindKey).(*store.Catalog)
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := store.AddCatalogNode(user, catalog); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Catalog created successfully.",
		"data": catalog,
	})
}
