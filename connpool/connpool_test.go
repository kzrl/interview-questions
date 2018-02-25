package connpool

import (
	"testing"
	"sync"

)


// Simple pool of 1 connection
func TestBasic(t *testing.T) {
	pool := New(1)
	_, err := pool.GetConnection()
	if err != nil {
		t.Error("Unable to get connection")
	}
}

// Should get one connection successfully, then receive an error for the second
func TestErrorOnEmpty(t *testing.T) {
	pool := New(1)
	_, err := pool.GetConnection()
	if err != nil {
		t.Error("Should have been able to get 1 connecttion")
	}

	_, err = pool.GetConnection()
	if err == nil {
		t.Error("Requested more connections than in pool. Expected an error")
	}
}

func TestConcurrentConnections(t *testing.T) {
	pool := New(5)
	
	errorCh := make(chan error)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			pool.GetConnection() //deliberately ignoring value
			wg.Done()
		}()
	}
	wg.Wait()
	close(errorCh)
	
	for e := range errorCh {
		t.Log(e)
		if e != nil {
			t.Error("Failed to get a connection in a goroutine")
		}
	}
	
}
