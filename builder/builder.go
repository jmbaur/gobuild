package builder

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"text/template"

	"git.jmbaur.com/gobuild/queue"
)

var (
	//go:embed index.html
	getText string
	templ   = template.Must(template.New("gobuild").Parse(getText))
)

type Builder struct {
	Builds   []*Build `json:"builds"`
	Commands []string `json:"commands"`
	Name     string   `json:"name"`
	Url      string   `json:"-"`

	l         sync.Mutex
	busy      bool
	buildChan chan *Build
	queueChan chan bool
	queue     *queue.Queue[Build]
}

func NewBuilder(name string, url string, commands []string) *Builder {
	return &Builder{
		Name:      name,
		Url:       url,
		Commands:  commands,
		Builds:    []*Build{},
		buildChan: make(chan *Build),
		queueChan: make(chan bool),
		queue:     queue.New[Build](),
	}
}

func (b *Builder) Run() {
	defer func() {
		close(b.buildChan)
		close(b.queueChan)
		log.Printf("Ended builder '%s'", b.Name)
	}()

	log.Printf("Started builder '%s'", b.Name)

	for {
		select {
		case build := <-b.buildChan:
			if build == nil {
				continue
			}
			if b.busy {
				log.Println("builder is busy, enqueueing build")
				b.queue.Enqueue(build)
			} else {
				log.Println("builder is not busy, preparing build")
				b.prepareAndCall(build)
			}
		case <-b.queueChan:
			b.l.Lock()
			b.busy = false
			b.l.Unlock()
			if !b.queue.IsEmpty() {
				log.Println("queue is not empty, sending in next queued build")
				b.prepareAndCall(b.queue.Dequeue())
			} else {
				log.Println("queue is empty, waiting for next build")
			}
		}
	}
}

func (b *Builder) prepareAndCall(build *Build) {
	b.l.Lock()
	b.busy = true
	b.l.Unlock()
	cmd := exec.Command(b.Commands[0], b.Commands[1:]...)
	b.Builds = append(b.Builds, build)
	go build.Do(cmd, b.queueChan)
}

func (b *Builder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.Contains(r.UserAgent(), "Mozilla") {
			w.Header().Add("Content-Type", "text/html")
			templ.Execute(w, b.MarshalMap())
			return
		} else {
			w.Header().Add("Content-Type", "application/json")
			data, err := json.Marshal(b)
			if err != nil {
				fmt.Fprint(w, struct {
					Error string `json:"error"`
				}{Error: err.Error()})
				return
			}
			fmt.Fprint(w, string(data))
			return
		}
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		b.buildChan <- &Build{}
		data, _ := json.Marshal(struct {
			Message string `json:"message"`
		}{Message: "OK"})
		fmt.Fprint(w, string(data))
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (b *Builder) MarshalMap() map[string]interface{} {
	builds := []map[string]interface{}{}
	for _, build := range b.Builds {
		builds = append(builds, build.marshalMap())
	}
	return map[string]interface{}{
		"builds":   builds,
		"name":     b.Name,
		"commands": b.Commands,
	}
}
