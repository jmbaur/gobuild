package builder

import (
	"bytes"
	"log"
	"os/exec"
	"time"
)

type Build struct {
	Output string    `json:"output"`
	Status Status    `json:"status"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}

func (b *Build) marshalMap() map[string]interface{} {
	return map[string]interface{}{
		"start":  b.Start.Format(time.RFC3339),
		"end":    b.End.Format(time.RFC3339),
		"output": b.Output,
		"status": b.Status.String(),
	}
}

func (b *Build) Do(cmd *exec.Cmd, done chan<- bool) {
	defer func() {
		b.End = time.Now()
		done <- true
		log.Println("finished build")
	}()
	log.Println("started build")
	b.Start = time.Now()

	buf := bytes.Buffer{}
	cmd.Stderr = &buf
	cmd.Stdout = &buf

	b.Status = StatusRunning
	err := cmd.Run()
	b.Output = buf.String()

	if err != nil {
		b.Status = StatusFailed
	} else {
		b.Status = StatusFinished
	}
}
