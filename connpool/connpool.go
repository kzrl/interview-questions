package connpool

import (
	"fmt"
	"sync"
)

type Connection interface {
	Close() error
	Execute(query string) error
}

type ConnectionPool interface {
	GetConnection() (Connection, error)
}

// MyConnection implements Connection
type MyConnection struct {
	pool *MyConnectionPool
	num  int
}

// Close returns a connection to the pool, rather than closing the connection
func (m *MyConnection) Close() error {
	m.pool.mux.Lock()
	m.pool.used[m.num] = false //mark this connection as available
	m.pool.mux.Unlock()
	return nil
}

// Execute does not actually do anything in this example
func (m *MyConnection) Execute(query string) error {
	return nil
}

// MyConnectionPool implements ConnectionPool
type MyConnectionPool struct {
	connections []Connection
	used        []bool
	mux         sync.Mutex //Ensure only a single goroutine can modify the slices
}

// GetConnection returns a Connection if one is available in the pool.
// It is safe to call from multiple goroutines
func (p *MyConnectionPool) GetConnection() (Connection, error) {
	available := false
	var newConn MyConnection

	p.mux.Lock()
	for i, c := range p.connections {

		// Found an unused connection
		if p.used[i] == false {
			c.pool = p       // keep a pointer to this pool
			p.used[i] = true // mark this connection as used
			c.num = i
			newConn = c
			available = true
			break
		}
	}
	p.mux.Unlock()

	if available {
		return &newConn, nil
	}

	return nil, fmt.Errorf("no connections available")
}

// New creates a new MyConnectionPool with the specified number of connections
// Deviates from the python example, accepting an int rather than a slice
func New(conns []Connection) *MyConnectionPool {
	// TODO: what do we do with the connections to make this work?
	return &MyConnectionPool{
		connections: ?,
		used:        make([]bool, len(conns)),
	}
}
