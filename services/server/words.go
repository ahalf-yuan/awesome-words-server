package server

import (
	"net/http"
	"wordshub/services/common/errno"
	"wordshub/services/store"

	"github.com/gin-gonic/gin"
)

func createWord(ctx *gin.Context) {
	word := ctx.MustGet(gin.BindKey).(*store.Words)

	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errno.ErrUserContext)
		return
	}
	if err := store.AddUserWord(user, word); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errno.ErrDB.WithData(gin.H{"error": err.Error()}))
		return
	}
	ctx.JSON(http.StatusOK, errno.OK.WithData(gin.H{}))
}

func indexWords(ctx *gin.Context) {
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	words, err := store.FetchUserWords(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, errno.OK.WithData(gin.H{
		"msg":  "Catalogs fetched successfully.",
		"data": words,
	}))
}
