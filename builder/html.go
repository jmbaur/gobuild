package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"
)

const (
	ModeJson = iota
	ModeHtml
)

var (
	getTempl  = template.Must(template.New("gobuild").Parse(getText))
	postTempl = template.Must(template.New("gobuild").Parse(postText))
)

func (b *Builder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var mode int
	if strings.Contains(r.UserAgent(), "Mozilla") {
		mode = ModeHtml
		w.Header().Add("Content-Type", "text/html")
	} else {
		mode = ModeJson
		w.Header().Add("Content-Type", "application/json")
	}

	switch r.Method {
	case http.MethodGet:
		if mode == ModeJson {
			data, err := json.Marshal(b)
			if err != nil {
				fmt.Fprint(w, struct {
					Error string `json:"error"`
				}{Error: err.Error()})
				return
			}
			fmt.Fprint(w, string(data))
			return
		} else {
			getTempl.Execute(w, b)
			return
		}
	case http.MethodPost:
		if mode == ModeJson {
			b.buildChan <- &Build{
				Status: StatusRunning,
				Time:   time.Now(),
			}
			data, _ := json.Marshal(struct {
				Message string `json:"message"`
			}{Message: "OK"})
			fmt.Fprint(w, string(data))
			return
		} else {
			postTempl.Execute(w, nil)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

const getText = `<html>
<h1>{{.Name}}</h1>
<h3>{{.Commands}}</h3>
{{range $build := .Builds}}
<div>
<h5>Time: {{$build.Time}}</h5>
<h5>Status: {{$build.Status}}</h5>
<p>Output: {{$build.Output}}</p>
</div>
{{end}}
</html>`
const postText = `<html>OK</html>`
