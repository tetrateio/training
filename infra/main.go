package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

var (
	startingCluster = 10
	finalCluster    = 14
)

func main() {
	// If we deleted projects locally then we will need to re-import them...
	for i := startingCluster; i <= finalCluster; i++ {
		address := fmt.Sprintf("google_project.training[%v]", i)
		id := fmt.Sprintf("projects/%v-%03d", workshopName, i)
		cmd := exec.Command("terraform", "import", address, id)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("failed: %v", err)
		}
	}
}
