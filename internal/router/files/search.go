package files

import (
	"fmt"
	"kmem/internal/database"
	"kmem/internal/database/query"
	"kmem/internal/event"
	"kmem/internal/models"
	"kmem/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Search(store *event.Store, pg *database.Postgres, cache *database.Cache) func(*gin.Context) {
	return func(ctx *gin.Context) {
		val, exists := ctx.Get(utils.USERNAME_KEY)
		if !exists {
			ctx.String(http.StatusUnauthorized, "user unauthorized")
			return
		}

		username, ok := val.(string)
		if !ok {
			ctx.String(http.StatusUnauthorized, "user unauthorized")
			return
		}

		// parse query
		searchQuery := ctx.Query("search")

		page, err := strconv.Atoi(ctx.Query("page"))
		if err != nil {
			page = 1 // default
		}

		itemsPerPage, err := strconv.Atoi(ctx.Query("itemsPerPage"))
		if err != nil {
			itemsPerPage = 8 // default
		}

		// get files
		searched, exists := cache.Get(fmt.Sprintf("%s:search:%s", username, searchQuery))
		if !exists {
			searched, err = query.QueryPatternMatching(pg, username, searchQuery)
			if err != nil {
				ctx.String(http.StatusBadRequest, "failed to search:", err)
				return
			}

			cache.Add(fmt.Sprintf("%s:search:%s", username, searchQuery), searched)
		}

		files, ok := searched.([]models.FileMetadata)
		if !ok {
			ctx.String(http.StatusInternalServerError, "something went wrong while searching data:", err)
			return
		}

		resp := buildItemsPage(files, page, itemsPerPage)

		ctx.JSON(http.StatusOK, resp)

	}
}
