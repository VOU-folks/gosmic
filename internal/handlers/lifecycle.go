package handlers

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Lifecycleable interface {
	Start() error
	Stop() error
}

type StarterFunc func()
type PreShutdownFunc func(os.Signal)
type ShutdownHandlerFunc func() error
type PostShutdownFunc func(os.Signal, error)

func waitForShutdownSignal(
	preShutdownFunc PreShutdownFunc,
	shutdownHandler ShutdownHandlerFunc,
	postShutdownFunc PostShutdownFunc,
) {
	log.Println("Press Control-C to stop")

	c := make(chan os.Signal, 1)
	signal.Notify(
		c,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	sig := <-c

	preShutdownFunc(sig)
	err := shutdownHandler()
	postShutdownFunc(sig, err)

	os.Exit(0)
}

func startRoutine(lifecycle Lifecycleable) StarterFunc {
	return func() {
		if err := lifecycle.Start(); err != nil {
			log.Fatalln(err.Error())
		}
	}
}

func stopRoutine(lifecycle Lifecycleable) ShutdownHandlerFunc {
	return lifecycle.Stop
}

func StartLifecycle(lifecycleable Lifecycleable) {
	go startRoutine(lifecycleable)()

	waitForShutdownSignal(
		func(sig os.Signal) {
			log.Println("Got", sig)
			log.Println("Initiating shutdown sequence...")
		},

		stopRoutine(lifecycleable),

		func(sig os.Signal, err error) {
			if err != nil {
				log.Println("Failed to gracefully end application lifecycle")
				log.Fatalln("Exiting with error:", err.Error())
			}
			log.Println("Application lifecycle gracefully ended")
			log.Println("Exiting...")
		},
	)
}
