package main

import (
	"fmt"
	"github.com/kzrl/interview-questions/connpool"
	"os"
	"sync"
)

// SomeBackendConnector implements connpool.Connection for some kind of backend.
type SomeBackendConnector struct{}

func (c *SomeBackendConnector) Close() error {
	fmt.Println("I really closed the connection here")
	return nil
}

func (c *SomeBackendConnector) Execute(query string) error {
	fmt.Println("Executed query: ", query)
	return nil
}

func main() {
	// confirm SomeBackendConnector satisfies Connection interface
	var _ connpool.Connection = &SomeBackendConnector{}

	fmt.Println("Creating a pool of 10 connections")
	conns := make([]connpool.Connection, 10)
	for i := 0; i < 10; i++ {
		conns[i] = &SomeBackendConnector{}
	}
	pool := connpool.New(conns)

	fmt.Println("Try to get 20 connections")
	for i := 0; i < 20; i++ {
		c, err := pool.GetConnection()
		if err != nil {
			fmt.Printf("%d %s\n", i, err)
			continue
		}

		err = c.Execute("SELECT * FROM users;")
		if err != nil {
			fmt.Println(err)
		}

	}

	fmt.Println("Create a new pool of 15 connections")
	pool = connpool.New(conns)

	fmt.Println("Try to get 20 connections in goroutines")
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(j int) {
			_, err := pool.GetConnection()
			if err != nil {
				fmt.Printf("%d %s\n", j, err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("Done with goroutines")

	fmt.Println("Get a connection, close it and try to execute")
	pool = connpool.New(conns)
	c, err := pool.GetConnection()
	if err != nil {
		fmt.Println("Unable to get connection")
		os.Exit(1)
	}
	c.Close() //ignoring error
	err = c.Execute("This query should fail")
	if err != nil {
		fmt.Printf("Query failed: %s\n", err)
	}

}

// A pointless function to demonstrate my implementation satisfies the interface
func isConnectionPool(p connpool.ConnectionPool) {
	c, err := p.GetConnection()
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Close()
	return
}

// A pointless function to demonstrate my implementation satisfies the interfac
func isConnection(c connpool.Connection) {
	c.Execute("A query")
	c.Close()
}
