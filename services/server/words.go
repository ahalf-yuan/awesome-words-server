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
	ctx.JSON(http.StatusOK, errno.OK)
}
