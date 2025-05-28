package router

import (
	"fmt"
	"io"
	"kmem/internal/config"
	"kmem/internal/models"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func upload(conf *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// v, ok := ctx.Get(utils.USERNAME_KEY)
		// if !ok {
		// 	models.APIResponse{
		// 		Status:  http.StatusUnauthorized,
		// 		Message: "failed to authorize",
		// 	}.Send(ctx)
		//
		// 	return
		// }
		//
		// username, ok := v.(string)
		// if !ok {
		// 	models.APIResponse{
		// 		Status:  http.StatusUnauthorized,
		// 		Message: "failed to authorize",
		// 	}.Send(ctx)
		//
		// 	return
		// }

		username := "testuser"

		encodedName := ctx.Query("filename")
		if len(encodedName) == 0 {
			models.APIResponse{
				Status:  http.StatusBadRequest,
				Message: "need filename",
			}.Send(ctx)

			return
		}

		// process filename
		// originalname, safename, mimetype, err := utils.ProcessFilename(encodedName)

		dst := fmt.Sprintf("%s/%s/%s", conf.UploadPath(), username, encodedName)

		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			models.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to create upload dir: %v", err),
			}.Send(ctx)

			return
		}

		file, err := os.Create(dst)
		if err != nil {
			models.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to create file: %v", err),
			}.Send(ctx)

			return
		}
		defer file.Close()

		_, err = io.Copy(file, ctx.Request.Body)
		if err != nil {
			models.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to stream request body: %v", err),
			}.Send(ctx)

			return
		}

		models.APIResponse{
			Status:  http.StatusOK,
			Message: "ok",
		}.Send(ctx)
	}
}
