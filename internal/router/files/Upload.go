package files

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"kmem/internal/event"
	"kmem/internal/utils"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func Upload(store *event.Store) func(*gin.Context) {
	return func(ctx *gin.Context) {
		username, exists := ctx.Get(utils.USERNAME_KEY)
		if !exists {
			ctx.String(http.StatusUnauthorized, "user unauthorized")
			return
		}

		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("failed to parse multi part form: %v", err))
			return
		}

		// create temp upload dir
		tmpDir := filepath.Join(os.TempDir(), utils.UPLOADDIR_TEMP, username.(string))
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("failed to create temp dir: %v", err))
			return
		}

		finalDir := filepath.Join(utils.UPLOADDIR_FINAL, username.(string), "tmp")

		// chunk tracker
		uploadStatus := make(map[string]bool)

		// handle chunks here
		for fieldName, fileHeader := range form.File {
			if !strings.HasPrefix(fieldName, "file-") {
				continue
			}

			// 1. parse field name
			fieldParts, err := parseFieldName(fieldName)
			if err != nil {
				// TODO: think about how to handle this error
				continue
			}

			// 2. save chunk(tmp)
			chunkPath := filepath.Join(tmpDir, fieldParts.filename, fieldParts.chunkIdxStr)

			if err := saveChunk(ctx, fileHeader, chunkPath); err != nil {
				uploadStatus[fieldParts.filename] = false
				continue
			}

			// 3. validate chunk (md5)
			if err := validateChunk(chunkPath, fieldParts.md5hash); err != nil {
				os.Remove(chunkPath)
				uploadStatus[fieldParts.filename] = false
				continue
			}

			// 4. move tmp to final path
			finalPath := filepath.Join(utils.UPLOADDIR_FINAL, username.(string), "tmp", fieldParts.filename, fieldParts.chunkIdxStr)

			if err := finalizeChunk(chunkPath, finalPath); err != nil {
				log.Println(err)
				uploadStatus[fieldParts.filename] = false
				continue
			}

			if err := os.Remove(chunkPath); err != nil {
				log.Println(err)
			}

			uploadStatus[fieldParts.filename] = true
		}

		// register to event store
		store.Register(event.FileUploaded(username.(string), finalDir))

		ctx.JSON(http.StatusOK, uploadStatus)
	}
}

// client will send form field as 'file-{filename}-{md5hash}-{chunk idx}'
type parts struct {
	filename    string
	md5hash     string
	chunkIdxStr string
}

func parseFieldName(fieldName string) (*parts, error) {
	split := strings.Split(fieldName, "-")
	if len(split) < 4 {
		return nil, fmt.Errorf("bad request: not enough parts")
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(split[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode filename: %v", err)
	}

	filename, err := url.QueryUnescape(string(decodedBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to unescape filename: %v", err)
	}

	return &parts{
		filename:    filename,
		md5hash:     split[2],
		chunkIdxStr: split[3],
	}, nil
}

func saveChunk(ctx *gin.Context, fileHeader []*multipart.FileHeader, chunkPath string) error {
	// prepare directories
	if err := os.MkdirAll(filepath.Dir(chunkPath), 0755); err != nil {
		return err
	}

	if err := ctx.SaveUploadedFile(fileHeader[0], chunkPath); err != nil {
		return err
	}

	return nil
}

func validateChunk(chunkPath, md5hash string) error {
	chunk, err := os.Open(chunkPath)
	if err != nil {
		return err
	}
	defer chunk.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, chunk); err != nil {
		return err
	}

	hashedStr := hex.EncodeToString(hash.Sum(nil))

	if md5hash != hashedStr {
		return fmt.Errorf("md5 hash doesn't match")
	}

	return nil
}

func finalizeChunk(chunkPath, finalPath string) error {
	if err := os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
		return err
	}

	chunk, err := os.Open(chunkPath)
	if err != nil {
		return fmt.Errorf("failed to open chunk file: %v", err)
	}

	final, err := os.Create(finalPath)
	if err != nil {
		chunk.Close()
		return fmt.Errorf("failed to create final file: %v", err)
	}

	if _, err := io.Copy(final, chunk); err != nil {
		chunk.Close()
		final.Close()
		return fmt.Errorf("failed to copy chunk to final path: %v", err)
	}

	chunk.Close()
	final.Close()

	return nil
}
