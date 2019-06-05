package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("templates/index.html"))
}

type Calc struct {
	Start  int
	Result []int
}

//DoubleMethod doubles the last number in the result slice
func (c *Calc) DoubleMethod() []int {
	r := c.Result[len(c.Result)-1] * 2
	c.Result = append(c.Result, r)
	return c.Result
}
func (c *Calc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", c)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func main() {
	c := &Calc{
		1,
		[]int{1},
	}

	fmt.Println("main", c)
	http.Handle("/", index(c))
	http.Handle("/click", double(c, c))

	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		h.ServeHTTP(w, req)
	})
}

func double(h http.Handler, c *Calc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c.DoubleMethod()
		fmt.Println("double", c.Result)

		h.ServeHTTP(w, req)
	})
}
