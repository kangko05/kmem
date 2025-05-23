package event

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"kmem/internal/database"
	"kmem/internal/database/command"
	"kmem/internal/database/query"
	"kmem/internal/models"
	"kmem/internal/utils"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

/*
	handle(context.Context) Result

	setResultChannel(chan Result)
	getResultChannel() chan Result

	setTimeout(time.Duration)
*/

type fileUploaded struct {
	username  string
	chunksDir string // utils.uploadir_final/username/tmp -> ./filename/chunks...
	resultCh  chan Result
	timeout   time.Duration
}

func FileUploaded(username, chunksDir string, options ...eventOption) *fileUploaded {
	fileUploaded := &fileUploaded{
		username:  username,
		chunksDir: chunksDir,
		resultCh:  nil,
		timeout:   time.Duration(defaultTimeout),
	}

	for _, opt := range options {
		opt(fileUploaded)
	}

	return fileUploaded
}

func (f *fileUploaded) setResultChannel(resultChan chan Result) {
	f.resultCh = resultChan
}

func (f *fileUploaded) getResultChannel() chan Result {
	return f.resultCh
}

func (f *fileUploaded) setTimeout(dur time.Duration) {
	f.timeout = dur
}

func (f *fileUploaded) handle(ctx context.Context, pg *database.Postgres, cache *database.Cache) Result {
	// 1. aggregate chunks & store into db
	uploadedDir, err := f.aggregateChunks(pg, f.extractMetadata, f.storeMetadata)
	if err != nil {
		return newEventResult(utils.FAIL, fmt.Sprintf("failed to aggregate chunks: %v", err), nil)
	}

	// 2. check uploaded dir
	dinfo, err := os.Stat(uploadedDir)
	if err != nil {
		return newEventResult(utils.FAIL, fmt.Sprintf("failed to read aggregated chunks: %v", err), nil)
	}

	if dinfo.Size() == 0 {
		os.Remove(uploadedDir)
	}

	// 3. update cache
	userfiles, err := query.QueryUserFiles(pg, f.username)
	if err != nil {
		log.Printf("failed to update cache: %v\n", err)
	}

	cache.Delete(f.username)
	cache.Add(f.username, userfiles)

	// below will be handled in seperate module later
	// 2. compress uploaded files -> store in backup path
	// err = f.compressDir(uploadedDir, uploadedDir+".zip") // TODO: zip file path temp for now
	// if err != nil {
	// 	return newEventResult(utils.FAIL, fmt.Sprintf("failed to compress dir: %v", err), nil)
	// }

	return newEventResult(utils.SUCCESS, "upload success", nil)
}

// helpers ====================================================================
// returns (final dest path, error)
func (f *fileUploaded) aggregateChunks(
	pg *database.Postgres,
	extractMeta func(filename, timestampDir, md5hash, archivePath string, now time.Time, totalSize int64) models.FileMetadata,
	storeMeta func(pg *database.Postgres, metadta models.FileMetadata) error,
) (string, error) {
	entries, err := os.ReadDir(f.chunksDir)
	if err != nil {
		log.Printf("failed to read chunks dir: %v\n", err)
		return "", fmt.Errorf("failed to read chunks dir: %v", err)
	}

	// 1. get filenames
	filenames := make([]string, len(entries))
	for i, entry := range entries {
		filenames[i] = entry.Name()
	}

	// 2. prepare upload dir
	now := time.Now()
	timestamp := now.Unix()
	timestampDir := filepath.Join(filepath.Dir(f.chunksDir), fmt.Sprint(timestamp))
	if err := os.MkdirAll(timestampDir, 0755); err != nil {
		log.Printf("failed to create final dir: %v\n", err)
		return "", fmt.Errorf("failed to create upload dir: %v", err)
	}

	// 3. read dir (chunks)
	for _, filename := range filenames {
		fileDir := filepath.Join(f.chunksDir, filename)

		chunks, err := os.ReadDir(fileDir)
		if err != nil {
			log.Printf("failed to read file dir: %v\n", err)
			continue
		}

		isErr := false
		sort.Slice(chunks, func(i, j int) bool {
			ichunk, ierr := strconv.Atoi(chunks[i].Name())
			jchunk, jerr := strconv.Atoi(chunks[j].Name())

			if (ierr != nil) || (jerr != nil) {
				isErr = true
				return false
			}

			return ichunk < jchunk
		})

		if isErr {
			log.Printf("[%s]: failed to parse chunk name\n", filename)
			continue
		}

		// prepare md5
		hash := md5.New()
		finalPath := filepath.Join(timestampDir, filename)

		dest, err := os.Create(finalPath)
		if err != nil {
			log.Printf("[%s]: failed to create dest file: %v", filename, err)
			continue
		}

		var totalSize int64

		for _, chunk := range chunks {
			rb, err := os.ReadFile(filepath.Join(fileDir, chunk.Name()))
			if err != nil {
				log.Println(err)
				dest.Close()
				break
			}

			n, err := io.Copy(dest, bytes.NewReader(rb))
			if err != nil {
				log.Println(err)
				dest.Close()
				break
			}

			totalSize += n
			hash.Write(rb)
		}

		dest.Close()

		// extract & store -> if any error, remove dest
		md5hash := hex.EncodeToString(hash.Sum(nil))
		metadata := extractMeta(filename, timestampDir, md5hash, timestampDir+".zip", now, totalSize)
		err = storeMeta(pg, metadata)
		if err != nil {
			log.Printf("failed to store metadata into db: %v", err)
			os.Remove(finalPath)
		}

	}

	if err := os.RemoveAll(filepath.Join(f.chunksDir)); err != nil {
		log.Printf("failed to remove tmp folder")
	}

	return timestampDir, nil
}

func (f *fileUploaded) compressDir(srcDir, zipFilePath string) error {
	archive, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer archive.Close()

	zw := zip.NewWriter(archive)
	defer zw.Close()

	err = filepath.Walk(srcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relPath
		header.Method = zip.Deflate

		writer, err := zw.CreateHeader(header)

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func (f *fileUploaded) extractMetadata(filename, timestampDir, md5hash, archivePath string, now time.Time, totalSize int64) models.FileMetadata {
	var metadata models.FileMetadata

	metadata.Filename = filename
	metadata.ContentType = string(utils.GetContentType(filename))
	metadata.StoredPath = timestampDir
	metadata.UploadedBy = f.username
	metadata.ArchivePath = archivePath
	metadata.UploadedAt = now
	metadata.Size = totalSize
	metadata.Hash = md5hash

	return metadata
}

func (f *fileUploaded) storeMetadata(pg *database.Postgres, metadata models.FileMetadata) error {
	return command.InsertFileMetadata(pg, metadata)
}
