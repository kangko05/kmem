package tests

import (
	"fmt"
	"kmem/internal/database"
	"kmem/internal/event"
	"kmem/internal/utils"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	err := godotenv.Load("../.env")
	assert.Nil(t, err)

	pg, err := database.Connect(t.Context())
	assert.Nil(t, err)
	defer pg.Close()

	t.Run("run event store", func(t *testing.T) {
		store := event.NewStore(t.Context(), pg)
		go store.Run()
		time.Sleep(200 * time.Millisecond)
	})

	t.Run("add test event", func(t *testing.T) {
		store := event.NewStore(t.Context(), pg)
		go store.Run()

		rchan := make(chan event.Result, 1)
		defer close(rchan)

		store.Register(event.TestEvent(0, event.WithResultChan(rchan)))

		result := <-rchan

		assert.Equal(t, result.Status(), utils.SUCCESS)
		assert.Equal(t, result.Message(), "test event 0")
	})

	t.Run("add multiple events", func(t *testing.T) {
		store := event.NewStore(t.Context(), pg)
		go store.Run()

		var wg sync.WaitGroup

		for i := range 50 {
			wg.Add(1)
			go func() {
				rchan := make(chan event.Result, 1)

				defer func() {
					wg.Done()
					close(rchan)
				}()

				store.Register(event.TestEvent(i, event.WithResultChan(rchan)))

				result := <-rchan

				assert.Equal(t, utils.SUCCESS, result.Status())
				assert.Equal(t, fmt.Sprintf("test event %d", i), result.Message())
			}()
		}

		wg.Wait()

		t.Context().Done()
	})

	t.Run("add multiple events without reschan", func(t *testing.T) {
		store := event.NewStore(t.Context(), pg)
		go store.Run()

		var wg sync.WaitGroup

		for i := range 50 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				store.Register(event.TestEvent(i))
			}()
		}

		wg.Wait()

		t.Context().Done()

		fmt.Println("done")
	})

	// t.Run("event cancel after timeout", func(t *testing.T) {
	// 	store := event.NewStore(t.Context(), pg)
	// 	go store.Run()
	//
	// 	rchan := make(chan event.Result, 1)
	// 	defer close(rchan)
	//
	// 	longEvent := event.TestEvent(0, event.WithResultChan(rchan), event.WithTimeout(time.Millisecond*100))
	//
	// 	store.Register(longEvent)
	// 	result := <-rchan
	//
	// 	assert.Equal(t, event.utils.FAIL, result.Status())
	// 	assert.Contains(t, result.Message(), "timed out")
	// })

	// t.Run("premature cancel", func(t *testing.T) {
	// 	store := event.NewStore(t.Context())
	// 	go store.Run()
	//
	// 	rchan := make(chan event.Result, 1)
	// 	defer close(rchan)
	//
	// 	store.Register(event.TestEvent(0, event.WithResultChan(rchan)))
	//
	// 	t.Context().Done()
	//
	// 	result := <-rchan
	//
	// 	fmt.Println(result)
	// })
}
