package connpool

import (
	"testing"
	"sync"
	"fmt"

)

// SomeBackendConnector implements connpool.Connection for some kind of backend.
// Copied from connectionpool.go so that we can test it here.
type SomeBackendConnector struct{}

func (c *SomeBackendConnector) Close() error {
	fmt.Println("I really closed the connection here")
	return nil
}

func (c *SomeBackendConnector) Execute(query string) error {
	fmt.Println("Executed query: ", query)
	return nil
}

func getConnections(num int) []Connection {
	conns := make([]Connection, num)
	for i := 0; i < num; i++ {
		conns[i] = &SomeBackendConnector{}
	}
	return conns
}


// Simple pool of 1 connection
func TestBasic(t *testing.T) {
	conns := getConnections(1)
	pool := New(conns)
	_, err := pool.GetConnection()
	if err != nil {
		t.Error("Unable to get connection")
	}
}

// Should get one connection successfully, then receive an error for the second
func TestErrorOnEmpty(t *testing.T) {
	conns := getConnections(1)
	pool := New(conns)
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
	conns := getConnections(5)
	pool := New(conns)
	
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

func TestUseAfterClose(t *testing.T) {
	conns := getConnections(5)
	pool := New(conns)
	c, err := pool.GetConnection()
	if err != nil {
		t.Error("Failed to get connection")
	}
	c.Close()
	err = c.Execute("This should fail")
	if err == nil {
		t.Error("Should not be allowed to Execute() after Closing the connection")
	}
}

func TestReuseConnection(t *testing.T) {
	conns := getConnections(1)
	pool := New(conns)
	c, err := pool.GetConnection()
	if err != nil {
		t.Error("Failed to get connection")
	}
	c.Close()
	err = c.Execute("This should fail")
	if err == nil {
		t.Error("Should not be allowed to Execute() after Closing the connection")
	}

	c2, err := pool.GetConnection()
	if err != nil {
		t.Error("Failed to get connection for reuse")
	}

	err = c2.Execute("This should work")
	if err != nil {
		t.Error("Should be able to reuse connection")
	}
}
