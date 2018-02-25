package connpool

import (
	"fmt"
	"sync"
)

const numConnections = 10

type Connection interface {
	Close() error
	Execute() error
}

type ConnectionPool interface {
	GetConnection() (Connection, error)
}

type MyConnection struct {
	pool *MyConnectionPool
	num  int
}

func (m *MyConnection) Close() error {
	m.pool.mux.Lock()
	m.pool.used[m.num] = false //mark this connection as available
	m.pool.mux.Unlock()
	return nil
}

func (m *MyConnection) Execute() error {
	return nil
}

type MyConnectionPool struct {
	connections []MyConnection
	used        []bool
	mux         sync.Mutex //Ensure only a single goroutine can modify the slices
}

func (p *MyConnectionPool) GetConnection() (Connection, error) {
	available := false
	var newConn Connection

	p.mux.Lock()
	for i, c := range p.connections {

		if p.used[i] == false {
			c.pool = p
			p.used[i] = true
			c.num = i
			newConn = &c
			available = true
			break
		}
	}
	p.mux.Unlock()

	if available {
		return newConn, nil
	}

	return nil, fmt.Errorf("no connections available")
}

func New(numConnections int) ConnectionPool {
	var p MyConnectionPool
	p.connections = make([]MyConnection, numConnections)
	p.used = make([]bool, numConnections)
	return &p
}
