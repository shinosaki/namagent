package utils

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func ExecCommand(command []string, sc *SignalContext) {
	var proc *exec.Cmd

	commandString := strings.Join(command, " ")
	sc.AddTask(commandString, func() {
		if proc != nil {
			if err := proc.Process.Signal(os.Interrupt); err != nil {
				log.Println("Exec Command: send interrupt to process", err)
			}
			if err := proc.Wait(); err != nil {
				log.Println("Exec Command: terminated process", err)
			}
		}
	})
	log.Println("Exec Command:", commandString)

	go func() {
		defer sc.CancelTask(commandString)

		proc = exec.Command(command[0], command[1:]...)

		if err := proc.Start(); err != nil {
			log.Println("Exec Command: failed to start", err)
			return
		}

		go func() {
			<-sc.Context().Done()
			log.Println("Exec Command: receive interrupt, send interrupt to", command[0])
			// if err := proc.Process.Signal(os.Interrupt); err != nil {
			// 	log.Printf("Exec Command: %s's interrupt message, %v", command[0], err)
			// }
			sc.CancelTask(commandString)
		}()

		if err := proc.Wait(); err != nil {
			log.Println("Exec Command: running error", err)
		}
	}()
}
