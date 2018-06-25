package minikube

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

// waitForDotKubeDirToBeCreated waits for `.kube` to be created
func (minikube Minikube) waitForDotKubeDirToBeCreated() {
	homeDir := os.Getenv("HOME")

	fmt.Println("Waiting for `.kube` to be created...")
	for {
		if _, err := os.Stat(path.Join(homeDir, ".kube")); err == nil {
			fmt.Println(path.Join(homeDir, ".kube") + " created.")
			break
		} else if _, err := os.Stat("/root/.kube"); err == nil {
			fmt.Println("/root/.kube created.")
			break
		}
		time.Sleep(minikube.WaitTimeUnit)
	}
}

// waitForDotMinikubeDirToBeCreated waits for `.minikube` to be created
func (minikube Minikube) waitForDotMinikubeDirToBeCreated() {
	homeDir := os.Getenv("HOME")

	fmt.Println("Waiting for `.minikube` to be created...")
	for {
		if _, err := os.Stat(path.Join(homeDir, ".minikube")); err == nil {
			fmt.Println(path.Join(homeDir, ".minikube") + " created.")
			break
		} else if _, err := os.Stat("/root/.minikube"); err == nil {
			fmt.Println("/root/.minikube created.")
			break
		}
		time.Sleep(minikube.WaitTimeUnit)
	}
}

// checkStatus checks minikube status and parse it to a map .
// :return:   map: minikube status parsed into dict.
//          error: if any error occurs, otherwise nil
// Note: error can come when machine is stopped too. But in this case status will be filled too
func (minikube Minikube) checkStatus() (map[string]string, error) {
	// Caller of this function should have proper rights to check minikube status
	command := "minikube status"
	statusStr, err := execCommand(command)

	status := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(statusStr), "\n") {
		keyval := strings.SplitN(line, ":", 2)
		if len(keyval) == 1 {
			status[strings.TrimSpace(keyval[0])] = ""
		} else {
			status[strings.TrimSpace(keyval[0])] = strings.TrimSpace(keyval[1])
		}
	}
	return status, err
}

// Status checks the status and in case where it does not find "minikube:" in status string,
// it retries until timeout. Then it returns the last status as well as the last error.
func (minikube Minikube) Status() (map[string]string, error) {
	startTime := time.Now()
	var status map[string]string
	var err error
	for time.Since(startTime) < minikube.Timeout {
		status, err = minikube.checkStatus()
		if _, ok := status["minikube"]; !ok {
			time.Sleep(minikube.WaitTimeUnit)
		} else {
			break
		}
	}
	return status, err
}
