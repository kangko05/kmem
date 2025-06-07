package tests

import (
	"kmem/internal/config"
	"kmem/internal/db"
	"kmem/internal/models"
	"kmem/internal/queue"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {

	q := queue.New(t.Context())

	conf, err := config.Load("../config.yml")
	assert.Nil(t, err)

	pg, err := db.Connect(conf)
	assert.Nil(t, err)
	defer pg.Close()

	photo := models.File{
		ID:         6,
		StoredName: "1749169414955370447_5572.gif",
		FilePath:   "/mnt/tank/kmem/testuser/1749169414955370447_5572.gif",
		MimeType:   "image/gif",
	}

	q.Add(queue.GenThumbnail(pg, conf, photo))

	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			q.Add(queue.TestItem(i))
		}()
	}

	wg.Wait()

	time.Sleep(time.Second * 3)
}
