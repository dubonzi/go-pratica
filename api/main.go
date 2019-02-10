package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type dado struct {
	ID       int
	Mensagem string
}

var dados []dado

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/{id}", hello)
	r.Get("/", todos)
	c := make(chan dado)
	prepararDados(c)

	for d := range c {
		dados = append(dados, d)
	}

	http.ListenAndServe(":8000", r)
}

func todos(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(dados)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	retorno := dado{}
	for _, d := range dados {
		if id == d.ID {
			retorno = d
			break
		}
	}
	time.Sleep(time.Second * 3)
	json.NewEncoder(w).Encode(retorno)
}

func prepararDados(c chan dado) {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {

			defer wg.Done()
			d := dado{ID: id, Mensagem: uuid.New().String()}
			c <- d
		}(i)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
}
