package server

import (
	"net/http"
	"strconv"
	"time"
	"wordshub/services/store"

	"github.com/gin-gonic/gin"
)

func fetchCatalog(ctx *gin.Context) {
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	catalog, err := store.FetchUserCatalogs(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	catalogCount, err := store.FetchUserCatalogAndCount(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// merge catalog with count data
	var res []interface{}
	res = mergeSlice(&catalog, &catalogCount)

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Catalogs fetched successfully.",
		"data": res,
	})
}

func mergeSlice(catalog *[]store.Catalog, catalogCount *[]store.CatalogCount) []interface{} {
	var interfaceSlice []interface{} = make([]interface{}, len(*catalogCount))

	catalogMap := make(map[int]store.Catalog)
	for _, item := range *catalog {
		catalogMap[item.ID] = item
	}

	for i, item := range *catalogCount {
		currCatalog := catalogMap[item.CatalogId]
		currCatalog.Count = item.Count
		currCatalog.ID = item.CatalogId
		interfaceSlice[i] = currCatalog
	}

	return interfaceSlice
}

func fetchCatalogAndCount(ctx *gin.Context) {
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	catalog, err := store.FetchUserCatalogAndCount(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Catalogs fetched successfully.",
		"data": catalog,
	})
}

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

func updateCatalog(ctx *gin.Context) {
	jsonCatalog := ctx.MustGet(gin.BindKey).(*store.Catalog)
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		return
	}

	dbCatalog, err := store.FetchCatalog(jsonCatalog.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.ID != dbCatalog.UserID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Not authorized."})
		return
	}

	jsonCatalog.ModifiedAt = time.Now()
	if err := store.UpdateCatalog(jsonCatalog); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Post updated successfully.",
		"data": jsonCatalog,
	})
}

func deleteCatalog(ctx *gin.Context) {
	paramID := ctx.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Not valid ID."})
		return
	}
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		return
	}
	catalog, err := store.FetchCatalog(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.ID != catalog.UserID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Not authorized."})
		return
	}
	if err := store.DeleteCatalog(catalog); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Catalog deleted successfully.", "data": catalog})
}
