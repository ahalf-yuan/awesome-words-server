package server

import (
	"net/http"
	"wordshub/services/common/errno"
	"wordshub/services/store"

	"github.com/gin-gonic/gin"
)

func signUp(ctx *gin.Context) {
	user := ctx.MustGet(gin.BindKey).(*store.User)

	fetchUser, err := store.FetchUserByName(user.Username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, errno.ErrServer)
		return
	}
	if fetchUser != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, errno.ErrUserNameNotUnique)
		return
	}

	if err := store.AddUser(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, errno.ErrServer)
		return
	}

	jwt := generateJWT(user)
	ctx.SetCookie("wordhub_jwt", jwt, 60*60*24, "/", "*", true, true) // 有效期一天

	ctx.JSON(http.StatusOK, errno.OK.WithData(gin.H{
		"jwt": jwt,
	}))
}

func signIn(ctx *gin.Context) {
	user := ctx.MustGet(gin.BindKey).(*store.User)
	user, err := store.Authenticate(user.Username, user.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, errno.ErrUserNameOrPassword)
		return
	}

	jwt := generateJWT(user)
	ctx.SetCookie("wordhub_jwt", jwt, 60*60*24, "/", "*", true, true) // 有效期一天

	ctx.JSON(http.StatusOK, errno.OK.WithData(gin.H{
		"jwt":      jwt,
		"id":       user.ID,
		"username": user.Username,
	}))
}

func signOut(ctx *gin.Context) {
	ctx.SetCookie("wordhub_jwt", "", 60*60*24, "/", "*", true, true)
	ctx.JSON(http.StatusOK, errno.OK)
}

func userInfo(ctx *gin.Context) {
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errno.ErrServer)
		return
	}

	ctx.JSON(http.StatusOK, errno.OK.WithData(gin.H{
		"id":       user.ID,
		"username": user.Username,
	}))
}
