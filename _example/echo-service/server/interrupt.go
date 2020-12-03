package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Interrupt server interrupt
func Interrupt(errc chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)

	// Place whatever shutdown handling you want here

	errc <- terminateError
}
