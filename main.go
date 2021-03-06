package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

type database struct {
	data       map[string][]int
	lock       sync.Mutex
	lastValues map[string]string
	t          *template.Template
}
type tadata struct {
	U    string
	Vals map[string]string
}

func remove(s []int, i int) []int {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (db *database) pop(user string) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	a, ok := db.data[user]
	if !ok {
		return fmt.Errorf("ahh shit couldnt get that")
	}
	l := len(a)
	if l == 0 {
		db.lastValues[user] = "done"
		return fmt.Errorf("ahh shit we dont with this")
	}
	i := rand.Intn(l)

	db.lastValues[user] = fmt.Sprint(a[i])
	db.data[user] = remove(a, i)
	return nil
}

var users = []string{"matt", "noah", "mike"}

func createDatabase() (*database, error) {
	m := make(map[string][]int)
	l := make(map[string]string)
	array := []int{1, 2, 3, 4, 5, 6}
	for _, v := range users {
		m[v] = array
		l[v] = "not started"
	}
	tt, err := template.ParseFiles("web/index.html")
	if err != nil {
		return nil, err
	}
	return &database{m, sync.Mutex{}, l, tt}, nil
}
func (db *database) writeOut(w io.Writer, u string) {
	l := tadata{u, db.lastValues}

	if err := db.t.Execute(w, l); err != nil {
		w.Write([]byte(err.Error()))
	}
}
func createHomeHandler(db *database) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fmt.Println(vars)
		if r.Method == http.MethodPost {
			db.pop(vars["user"])
		}
		db.writeOut(w, vars["user"])
	}
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome to bag of beer"))

}
func main() {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	image_dir := os.Getenv("IMAGE_DIR")
	addr := fmt.Sprintf("%s:%s", host, port)
	r := mux.NewRouter()
	image_prefix := "/images/"
	image_server := http.FileServer(http.Dir(image_dir))
	r.PathPrefix(image_prefix).Handler(http.StripPrefix(image_prefix, image_server))
	mdb, err := createDatabase()
	if err != nil {
		log.Fatal(err)
	}
	r.HandleFunc("/{user}", createHomeHandler(mdb))
	r.HandleFunc("/", indexHandler)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Fatal(srv.ListenAndServe())
}
