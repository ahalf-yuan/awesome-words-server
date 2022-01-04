package server

import (
	"net/http"
	"wordshub/services/store"

	"github.com/gin-gonic/gin"
)

func createWord(ctx *gin.Context) {
	word := ctx.MustGet(gin.BindKey).(*store.Words)
	if word == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error"})
		return
	}

	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := store.AddUserWord(user, word); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Catalog created successfully.",
		"data": word,
	})
}
