package main

import (
	_ "image/png"
	"log"
	"os"
	"time"

	asciigif "github.com/avbrocklehurst/go-ascii-gif"
)

func main() {
	ag, err := asciigif.New(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	//_ = ag
	ag.Start()
	time.Sleep(time.Second * 5)
	ag.Stop()
	time.Sleep(time.Second * 1)
}
