package booking

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	redis "github.com/pangdfg/cinema/internal/adapters"
)

func TestConcurrentBooking_ExactlyOneWins(t *testing.T) {
	rdb, err := redis.NewClient("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}

	store := NewRedisStore(rdb)
	svc := NewService(store)

	const n = 1000

	var success atomic.Int64
	var wg sync.WaitGroup

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			_, err := svc.Book(Booking{
				MovieID: "m1",
				SeatID:  "A1",
				UserID:  uuid.New().String(),
			})

			if err == nil {
				success.Add(1)
			}
		}()
	}

	wg.Wait()

	if success.Load() != 1 {
		t.Fatalf("expected 1 success, got %d", success.Load())
	}
}