package server

import (
	"net/http"
	"strconv"
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

	ctx.JSON(http.StatusOK, errno.OK.WithData(words))
}

func delWord(ctx *gin.Context) {
	paramID := ctx.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errno.ErrWordID)
		return
	}
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errno.ErrUserContext)
		return
	}

	uWord, err := store.FetchWordById(id, user.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 用户不存在在 word_id
	if uWord == nil {
		ctx.AbortWithStatusJSON(http.StatusOK, errno.ErrWordExistWithUser)
		return
	}

	if err := store.DeleteWord(uWord.ID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, errno.OK.WithData(uWord))
}
