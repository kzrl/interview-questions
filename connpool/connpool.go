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

// MyConnection embeds a Connection
type MyConnection struct {
	Connection
	pool   *MyConnectionPool
	num    int
	closed bool
}

// Close returns a connection to the pool, rather than closing the connection
func (m *MyConnection) Close() error {
	m.pool.mux.Lock()

	// Create a new MyConnection
	var newconn MyConnection
	newconn.pool = m.pool
	newconn.num = m.num
	newconn.closed = false
	newconn.Connection = m.Connection //copy the old embedded Connection

	m.pool.connections[m.num] = newconn // replace the connection

	//This old connection is now closed forever
	m.closed = true

	m.pool.used[m.num] = false //mark the new connection as available
	m.pool.mux.Unlock()
	return nil
}

// Check if this connection is closed
func (m *MyConnection) IsClosed() bool {
	return m.closed
}

// Execute method is actually implemented on the embedded Connection.
// Calling it here to check if the connection is already closed
func (m *MyConnection) Execute(query string) error {
	if m.IsClosed() {
		return fmt.Errorf("connection is closed")
	}
	return m.Connection.Execute(query)
}

// MyConnectionPool implements ConnectionPool
type MyConnectionPool struct {
	connections []MyConnection
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

// New creates a new MyConnectionPool with the given []Connection
func New(conns []Connection) *MyConnectionPool {
	// Make a slice of MyConnections
	myConnections := make([]MyConnection, len(conns))

	// Copy each passed Connection in conns to the new slice
	for i, val := range conns {
		myConnections[i].Connection = val
	}

	return &MyConnectionPool{
		connections: myConnections,
		used:        make([]bool, len(conns)),
	}
}
