package router

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"kmem/internal/config"
	"kmem/internal/db"
	"kmem/internal/models"
	"kmem/internal/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getLimitPageQuery(limitStr, pageStr string) (int, int) {
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = utils.DEAFULT_LIMIT
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}

	return limit, page
}

// get limit & offset & page through query
func servFiles(pg *db.Postgres, conf *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, ok := ctx.Get(utils.USERNAME_KEY)
		if !ok {
			models.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "failed to get username",
			}.Send(ctx)
			return
		}

		username := val.(string)

		limit, page := getLimitPageQuery(ctx.Query("limit"), ctx.Query("page"))

		sort := ctx.Query("sort")
		if len(sort) == 0 {
			sort = "date"
		}

		typeStr := ctx.Query("type")
		if len(typeStr) == 0 {
			typeStr = "all"
		}

		if len(username) == 0 {
			models.APIResponse{
				Status:  http.StatusBadRequest,
				Message: "failed to get username",
			}.Send(ctx)
			return
		}

		dbfiles, err := pg.GetFilesPage(username, page, limit, sort, typeStr)
		if err != nil {
			models.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}.Send(ctx)
			return
		}

		totalFiles, err := pg.GetFilesCount(username, typeStr)
		if err != nil {
			models.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}.Send(ctx)
			return
		}

		var files []models.File
		for _, f := range dbfiles {
			after, ok := strings.CutPrefix(f.FilePath, conf.UploadPath())
			if ok {
				files = append(files, models.File{
					OriginalName: f.OriginalName,
					FilePath:     after,
					MimeType:     f.MimeType,
				})
			}
		}

		type Page struct {
			Files    []models.File `json:"files"`
			HasNext  bool          `json:"hasNext"`
			NextPage int           `json:"nextPage"`
		}

		pageResponse := Page{
			Files:    files,
			HasNext:  (page+1)*limit < totalFiles,
			NextPage: page + 1,
		}

		ctx.JSON(http.StatusOK, pageResponse)
	}
}

func upload(pg *db.Postgres, conf *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		v, ok := ctx.Get(utils.USERNAME_KEY)
		if !ok {
			models.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "failed to authorize",
			}.Send(ctx)

			return
		}

		username, ok := v.(string)
		if !ok {
			models.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "failed to authorize",
			}.Send(ctx)

			return
		}

		encodedName := ctx.Query("filename")
		if len(encodedName) == 0 {
			models.APIResponse{
				Status:  http.StatusBadRequest,
				Message: "need filename",
			}.Send(ctx)

			return
		}

		// process filename
		originalName, safename, mimeType, err := utils.ProcessFilename(encodedName)

		dst := fmt.Sprintf("%s/%s/%s", conf.UploadPath(), username, safename)

		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			models.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to create upload dir: %v", err),
			}.Send(ctx)

			return
		}

		hasher := sha256.New()

		//
		file, err := os.Create(dst)
		if err != nil {
			models.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to create file: %v", err),
			}.Send(ctx)

			return
		}
		defer file.Close()

		size, err := io.Copy(io.MultiWriter(file, hasher), ctx.Request.Body)
		if err != nil {
			models.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to stream request body: %v", err),
			}.Send(ctx)

			return
		}

		hash := hex.EncodeToString(hasher.Sum(nil))

		// into db
		filemeta := models.File{
			Hash:         hash,
			Username:     username,
			OriginalName: originalName,
			StoredName:   safename,
			FilePath:     dst,
			FileSize:     size,
			MimeType:     mimeType,
		}

		err = pg.InsertFile(filemeta)
		if err != nil {
			os.Remove(dst)

			models.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to insert file metadata into db: %v", err),
			}.Send(ctx)

			return
		}

		models.APIResponse{
			Status:  http.StatusOK,
			Message: "ok",
		}.Send(ctx)
	}
}
