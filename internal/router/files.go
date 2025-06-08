package router

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"kmem/internal/config"
	"kmem/internal/db"
	"kmem/internal/models"
	"kmem/internal/queue"
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
func servFiles(pg *db.Postgres) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, ok := ctx.Get(utils.USERNAME_KEY)
		if !ok {
			models.ErrorResponse(
				http.StatusUnauthorized,
				models.ErrUnauthorized,
				"authentication required",
			).Send(ctx)

			return
		}

		username := val.(string)

		if len(username) == 0 {
			models.ErrorResponse(
				http.StatusUnauthorized,
				models.ErrUnauthorized,
				"authentication required",
			).Send(ctx)

			return
		}

		limit, page := getLimitPageQuery(ctx.Query("limit"), ctx.Query("page"))

		sort := ctx.Query("sort")
		if len(sort) == 0 {
			sort = "date"
		}

		typeStr := ctx.Query("type")
		if len(typeStr) == 0 {
			typeStr = "all"
		}

		dbfiles, err := pg.GetFilesPage(username, page, limit, sort, typeStr)
		if err != nil {

			fmt.Println(err)

			models.ErrorResponse(
				http.StatusInternalServerError,
				models.ErrDatabase,
				"failed to get uesr files",
			).Send(ctx)

			return
		}

		totalFiles, err := pg.GetFilesCount(username, typeStr)
		if err != nil {
			models.ErrorResponse(
				http.StatusInternalServerError,
				models.ErrDatabase,
				"failed to get user files",
			).Send(ctx)

			return
		}

		type Page struct {
			Files    []models.FileResponse `json:"files"`
			HasNext  bool                  `json:"hasNext"`
			NextPage int                   `json:"nextPage"`
		}

		pageResponse := Page{
			Files:    dbfiles,
			HasNext:  (page+1)*limit < totalFiles,
			NextPage: page + 1,
		}

		models.SuccessResponse(pageResponse).Send(ctx)
	}
}

func upload(pg *db.Postgres, conf *config.Config, q *queue.Queue) gin.HandlerFunc {
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

		encodedName := ctx.Query("filename")
		if len(encodedName) == 0 {
			models.ErrorResponse(
				http.StatusBadRequest,
				models.ErrInvalidInput,
				"filename required",
			).Send(ctx)

			return
		}

		// process filename
		originalName, safename, mimeType, err := utils.ProcessFilename(encodedName)
		if err != nil {
			models.ErrorResponse(
				http.StatusBadRequest,
				models.ErrInvalidInput,
				"invalid filename",
			).Send(ctx)

			return
		}

		dst := fmt.Sprintf("%s/%s/%s", conf.UploadPath(), username, safename)

		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			models.ErrorResponse(
				http.StatusInternalServerError,
				models.ErrDatabase,
				"failed to upload file",
			).Send(ctx)

			return
		}

		//
		hasher := sha256.New()

		file, err := os.Create(dst)
		if err != nil {
			models.ErrorResponse(
				http.StatusInternalServerError,
				models.ErrDatabase,
				"failed to upload file",
			).Send(ctx)

			fmt.Println(err)

			return
		}
		defer file.Close()

		size, err := io.Copy(io.MultiWriter(file, hasher), ctx.Request.Body)
		if err != nil {
			models.ErrorResponse(
				http.StatusInternalServerError,
				models.ErrDatabase,
				"failed to upload file",
			).Send(ctx)

			fmt.Println(err, 1)

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
			RelativePath: "/static" + strings.TrimPrefix(dst, conf.UploadPath()),
			FileSize:     size,
			MimeType:     mimeType,
		}

		fileId, err := pg.InsertFile(filemeta)
		if err != nil {
			os.Remove(dst)

			models.ErrorResponse(
				http.StatusInternalServerError,
				models.ErrDatabase,
				"failed to upload file",
			).Send(ctx)

			fmt.Println(err, 2)

			return
		}

		models.SuccessResponse(nil).Send(ctx)

		filemeta.ID = fileId

		q.Add(queue.GenThumbnail(pg, conf, filemeta))
	}
}
