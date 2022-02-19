package builder

import (
	"bytes"
	"log"
	"os/exec"
	"time"
)

const (
	StatusQueued = iota
	StatusRunning
	StatusFailed
	StatusFinished
)

type Build struct {
	Output string    `json:"output"`
	Status int       `json:"status"`
	Time   time.Time `json:"time"`
}

type Builder struct {
	Builds   []*Build `json:"builds"`
	Commands []string `json:"commands"`
	Name     string   `json:"name"`
	Url      string   `json:"-"`

	buildChan chan *Build
}

func NewBuilder(name string, url string, commands []string) *Builder {
	return &Builder{
		Name:      name,
		Url:       url,
		Commands:  commands,
		Builds:    []*Build{},
		buildChan: make(chan *Build),
	}
}

func (b *Builder) Run() {
	log.Printf("Started builder '%s'", b.Name)
	defer close(b.buildChan)
	for build := range b.buildChan {
		b.Builds = append(b.Builds, build)
		cmd := exec.Command(b.Commands[0], b.Commands[1:]...)
		buf := bytes.Buffer{}
		cmd.Stderr = &buf
		cmd.Stdout = &buf
		build.Status = StatusRunning
		err := cmd.Run()
		build.Output = buf.String()
		if err != nil {
			build.Status = StatusFailed
		} else {
			build.Status = StatusFinished
		}
	}
}
