package files

import (
	"fmt"
	"kmem/internal/database"
	"kmem/internal/database/query"
	"kmem/internal/models"
	"kmem/internal/utils"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// receive query params - sort, page, itemsPerPage
func GetItems(pg *database.Postgres, cache *database.Cache) func(*gin.Context) {
	return func(ctx *gin.Context) {
		// get user info
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

		sortBy := ctx.Query("sort")

		page, err := strconv.Atoi(ctx.Query("page"))
		if err != nil {
			page = 1 // default
		}

		itemsPerPage, err := strconv.Atoi(ctx.Query("itemsPerPage"))
		if err != nil {
			itemsPerPage = 8 // default
		}

		cacheVal, exists := cache.Get(username)
		if !exists {
			cacheVal, err = query.QueryUserFiles(pg, username)
			if err != nil {
				ctx.String(http.StatusInternalServerError, "failed to get files")
				return
			}

			cache.Add(username, cacheVal)
		}

		userfiles, ok := cacheVal.([]models.FileMetadata)
		if !ok {
			ctx.String(http.StatusInternalServerError, "failed to get files")
			return

		}

		sort.Slice(userfiles, func(i, j int) bool {
			if sortBy == "date" {
				return userfiles[i].UploadedAt.Unix() > userfiles[j].UploadedAt.Unix()
			}

			if sortBy == "name" {
				return userfiles[i].Filename < userfiles[j].Filename
			}

			return userfiles[i].UploadedAt.Unix() > userfiles[j].UploadedAt.Unix() // default by date
		})

		resp := buildItemsPage(userfiles, page, itemsPerPage)

		ctx.JSON(http.StatusOK, resp)
	}
}

// helpers ====================================================================

type itemResp struct {
	Totalpages int                   `json:"totalpages"`
	Totalitems int                   `json:"totalitems"`
	Items      []models.MetadataPart `json:"items"`
}

func buildItemsPage(userfiles []models.FileMetadata, page, itemsPerPage int) itemResp {

	var ir itemResp
	for _, uf := range userfiles {
		fp := fmt.Sprintf("/files/static/%s/%s", strings.TrimPrefix(uf.StoredPath, "/home/kang/Downloads/"), uf.Filename)

		ir.Items = append(ir.Items, models.MetadataPart{
			Filename:    fp,
			ContentType: uf.ContentType,
			UploadedAt:  uf.UploadedAt,
			Size:        uf.Size,
		})
	}

	// calc page items
	ir.Totalpages = len(ir.Items) / itemsPerPage

	start := (page - 1) * itemsPerPage
	end := min(start+itemsPerPage, len(ir.Items))

	resp := itemResp{
		Totalpages: ir.Totalpages,
		Totalitems: len(userfiles),
		Items:      ir.Items[start:end],
	}

	return resp
}

func min(a, b int) int {
	if a > b {
		return b
	}

	return a
}
