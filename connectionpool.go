package main

import (
	"fmt"
	"github.com/kzrl/interview-questions/connpool"
	"sync"
)

func main() {

	fmt.Println("Creating a pool of 10 connections")
	pool := connpool.New(10)

	fmt.Println("Try to get 20 connections")
	for i := 0; i < 20; i++ {
		_, err := pool.GetConnection()
		if err != nil {
			fmt.Printf("%d %s\n", i, err)
			continue
		}
	}

	fmt.Println("Create a new pool of 15 connections")
	pool = connpool.New(10)

	fmt.Println("Try to get 20 connections in goroutines")
	var wg sync.WaitGroup
	for i := 0; i< 20; i++ {
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
	fmt.Println("DONE")
	isConnectionPool(&pool)
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
	c.Execute()
	c.Close()
}
