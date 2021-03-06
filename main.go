package main

import (
	"net/http"

	"github.com/alecthomas/template"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

// Counter is ...
type Counter struct {
	Name  string
	Count int
}

var handlerTemplate = template.Must(template.New("handler").Parse(handlerHTML))

const handlerHTML = `
{{range .}}
<p>{{.Name}}: {{.Count}}</p>
{{end}}
<p>Start a new counter:</p>
<form action="/" method="POST">
<input type="text" name="name">
<input type="submit" value="Add">
</form>
`

func init() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/worker", workHandler)
}

func main() {
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if name := r.FormValue("name"); name != "" {
		t := taskqueue.NewPOSTTask("/worker", map[string][]string{"name": {name}})
		if _, err := taskqueue.Add(ctx, t, ""); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	q := datastore.NewQuery("Counter")
	var counters []Counter
	if _, err := q.GetAll(ctx, &counters); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := handlerTemplate.Execute(w, counters); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func workHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	name := r.FormValue("name")
	key := datastore.NewKey(ctx, "Counter", name, 0, nil)
	var counter Counter
	if err := datastore.Get(ctx, key, &counter); err == datastore.ErrNoSuchEntity {
		counter.Name = name
	} else if err != nil {
		log.Errorf(ctx, "%v", err)
		return
	}
	counter.Count++
	if _, err := datastore.Put(ctx, key, &counter); err != nil {
		log.Errorf(ctx, "%v", err)
	}
}
