package server

import (
	"errors"
	"fmt"
	"net/http"
	"wordshub/services/common/errno"
	"wordshub/services/store"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func authorization(ctx *gin.Context) {
	cookieJwt, err := ctx.Cookie("wordhub_jwt")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errno.ErrSignParam)
		return
	}
	userID, err := verifyJWT(cookieJwt)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errno.ErrSignParam)
		return
	}
	user, err := store.FetchUser(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errno.ErrUserIDNotExit)
		return
	}
	ctx.Set("user", user)
	ctx.Next()
}

func currentUser(ctx *gin.Context) (*store.User, error) {
	var err error
	_user, exists := ctx.Get("user")
	if !exists {
		err = errors.New("Current context user not set")
		log.Error().Err(err).Msg("")
		return nil, err
	}
	user, ok := _user.(*store.User)
	if !ok {
		err = errors.New("Context user is not valid type")
		log.Error().Err(err).Msg("")
		return nil, err
	}
	return user, nil
}

func customErrors(ctx *gin.Context) {
	ctx.Next()
	if len(ctx.Errors) > 0 {
		for _, err := range ctx.Errors {
			// Check error type
			switch err.Type {
			case gin.ErrorTypePublic:
				// Show public errors only if nothing has been written yet
				if !ctx.Writer.Written() {
					ctx.AbortWithStatusJSON(ctx.Writer.Status(), gin.H{"error": err.Error()})
				}
			case gin.ErrorTypeBind:
				errMap := make(map[string]string)
				if errs, ok := err.Err.(validator.ValidationErrors); ok {
					for _, fieldErr := range []validator.FieldError(errs) {
						errMap[fieldErr.Field()] = customValidationError(fieldErr)
					}
				}

				status := http.StatusBadRequest
				// Preserve current status
				if ctx.Writer.Status() != http.StatusOK {
					status = ctx.Writer.Status()
				}
				ctx.AbortWithStatusJSON(status, gin.H{"error": errMap})
			default:
				// Log other errors
				log.Error().Err(err.Err).Msg("Other error")
			}
		}

		// If there was no public or bind error, display default 500 message
		if !ctx.Writer.Written() {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		}
	}
}

func customValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required.", err.Field())
	case "min":
		return fmt.Sprintf("%s must be longer than or equal %s characters.", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s characters.", err.Field(), err.Param())
	default:
		return err.Error()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
