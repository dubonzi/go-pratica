package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func main() {
	c := make(chan string, 30)
	wg := sync.WaitGroup{}
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go chamarAPI(i, c, &wg)
	}
	wg.Wait()
	close(c)
	for v := range c {
		fmt.Println(v)
	}
}

func chamarAPI(id int, c chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	r, err := http.Get("http://localhost:8000/" + fmt.Sprint(id))
	if err != nil {
		log.Fatal(err)
	}
	dado, _ := ioutil.ReadAll(r.Body)
	c <- string(dado)
}
