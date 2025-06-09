package router

import (
	"fmt"
	"kmem/internal/cache"
	"kmem/internal/db"
	"kmem/internal/models"
	"kmem/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUsage(pg *db.Postgres, cache *cache.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		v, ok := ctx.Get(utils.USERNAME_KEY)
		if !ok {
			models.ErrorResponse(
				http.StatusUnauthorized,
				models.ErrUnauthorized,
				"authentication required",
			).Send(ctx)

			return
		}

		username, ok := v.(string)
		if !ok {
			models.ErrorResponse(
				http.StatusUnauthorized,
				models.ErrUnauthorized,
				"authentication required",
			).Send(ctx)

			return
		}

		cacheKey := fmt.Sprintf("%s:stats:usage", username)

		val, ok := cache.Get(cacheKey)
		if ok {
			models.SuccessResponse(val).Send(ctx)
			return
		}

		totalCnt, totalSize, err := pg.GetUserFilesUsage(username)
		if err != nil {
			models.ErrorResponse(
				http.StatusInternalServerError,
				models.ErrDatabase,
				"failed to get total size and count",
			).Send(ctx)

			return
		}

		resp := map[string]any{
			"username":     username,
			"count":        totalCnt,
			"size":         totalSize,
			"readableSize": utils.GetReadableSize(totalSize),
		}

		models.SuccessResponse(resp).Send(ctx)

		cache.Set(cacheKey, resp)
	}
}
