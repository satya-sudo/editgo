package data

import (
	"log"
	"time"
)

type AutoSave struct {
	FM       *FileManager
	Interval time.Duration
	Quit     chan struct{}
}

func NewAutoSave(fm *FileManager, interval time.Duration) *AutoSave {
	return &AutoSave{
		FM:       fm,
		Interval: interval,
		Quit:     make(chan struct{}),
	}
}

func (a *AutoSave) Start() {
	if a.FM.FilePath == "" {
		log.Println("auto save file not exist")
		return
	}
	go func() {
		ticker := time.NewTicker(a.Interval)
		defer ticker.Stop()

		log.Println("auto save file started")
		for {
			select {
			case <-ticker.C:
				if a.FM.Buffer.IsDirty() {
					err := a.FM.Save()
					if err != nil {
						log.Println("auto save file err:", err)
					} else {
						log.Println("auto save file done")
					}
				}
			case <-a.Quit:
				log.Println("auto save file quit")
				return
			}
		}
	}()
}

func (a *AutoSave) Stop() {
	close(a.Quit)
}
