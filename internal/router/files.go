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

	"github.com/gin-gonic/gin"
)

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
