package main

import (
	"fmt"
	"github.com/kzrl/interview-questions/connpool"
)




func main() {
	fmt.Println("Yo")



	pool := connpool.New()

	for i := 0; i< 20; i++ {
		_, err := pool.GetConnection()
		if err != nil {
			fmt.Printf("%d %s\n", i, err)
		}
	}

	
	
}
