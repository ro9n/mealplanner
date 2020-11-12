package pkg

import (
	"log"
)

type reader interface {
	List(directory string) (*[]string, error)
	Read(filepath string) *Activity
}

// Worker represents file processor
type Worker struct {
	reader
}

// NewWorker creates worker
func NewWorker(r reader) *Worker {
	return &Worker{r}
}

func (w *Worker) process(directory string) <-chan Activity {
	files, err := w.List(directory)

	if err != nil {
		log.Fatalln("Unable to process directory", directory)
	}

	activities := make(chan Activity)

	go func() {
		for _, file := range *files {
			log.Println("Processing file", file)
			activities <- *w.Read(file)
		}
		close(activities)
	}()

	return activities
}
