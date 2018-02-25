package main

import (
	"fmt"
	"github.com/kzrl/interview-questions/connpool"
	"sync"
)




func main() {

	pool := connpool.New(10)

	for i := 0; i< 20; i++ {
		_, err := pool.GetConnection()
		if err != nil {
			fmt.Printf("%d %s\n", i, err)
			continue
		}
	}
	var wg sync.WaitGroup
	for i := 0; i< 20; i++ {
		wg.Add(1)
		go func() {
			_, err := pool.GetConnection()
			fmt.Printf("%s\n", err)
			wg.Done()
		}()
	}
	wg.Wait()

	
	
}
