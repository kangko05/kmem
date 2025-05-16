package files

import (
	"fmt"
	"io"
	"kmem/internal/event"
	"kmem/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Upload(store *event.Store) func(*gin.Context) {
	return func(ctx *gin.Context) {

		reader, err := ctx.Request.MultipartReader()
		if err != nil {
			ctx.String(http.StatusBadRequest, "failed to parse multipart: %v", err)
			return
		}

		claim, exists := ctx.Get(utils.USERNAME_KEY)
		if !exists {
			ctx.String(http.StatusUnauthorized, "user unauthorized")
			return
		}

		username, ok := claim.(string)
		if !ok {
			ctx.String(http.StatusBadRequest, "invalid token claim")
			return
		}

		successFiles := []string{}
		failedFiles := []string{}

		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				failedFiles = append(failedFiles, fmt.Sprintf("failed to read file part: %v", err))
				continue
			}

		}
	}
}
