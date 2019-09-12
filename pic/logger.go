package pic

import (
	"log"
	"os"
)

type logger struct {
	f *os.File
}

func initLogger(o Options) {
	if o.LogFileName != "" {
		f, err := os.Create(o.LogFileName)
		if err == nil {
			log.SetOutput(&logger{f: f})
			if o.LogInfo() {
				log.Println("INFO: Initialised logger")
			}
		}
	}
}

func (l *logger) Write(p []byte) (n int, err error) {
	n, err = l.f.Write(p)
	if err == nil {
		err = l.f.Sync()
	}
	return
}
