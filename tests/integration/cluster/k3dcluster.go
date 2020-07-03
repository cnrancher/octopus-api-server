package cluster

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	ClusterName = "edge-integration-test"
	K3dVersion  = "v1.7.0"
)

func StepK3dCluster() error {
	if err := ValidateK3d(); err != nil {
		log.Println("ValidateK3d error", err)
		return err
	}

	if err := CreateK3dCluster(); err != nil {
		log.Println("CreateK3dCluster error", err)
		return err
	}

	if err := ValidateK3dCluster(); err != nil {
		log.Println("ValidateK3dCluster error", err)
		return err
	}
	return nil
}

func InstallK3d() error {
	log.Println("Begin Install k3d")
	osStr := runtime.GOOS
	arch := runtime.GOARCH
	url := fmt.Sprintf("https://github.com/rancher/k3d/releases/download/%s/k3d-%s-%s", K3dVersion, osStr, arch)
	cmd := exec.Command("curl", "-fL", url, "-o", "/tmp/k3d")
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	}
	if err = exec.Command("chmod", "+x", "/tmp/k3d").Run(); err != nil {
		log.Println(err)
		return err
	}
	if out, err := exec.Command("/bin/bash", "-c", "mv /tmp/k3d /usr/local/bin/k3d").CombinedOutput(); err != nil {
		log.Println(err, string(out))
		return err
	}
	return nil
}

func ValidateK3d() error {
	_, err := exec.LookPath("k3d")
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found in $PATH") {
			return InstallK3d()
		}
	}
	return err
}

func CreateK3dCluster() error {
	maxTry := 10
	reTry := 0
RANDOM:
	if reTry >= maxTry {
		return errors.New("CreateK3dCluster random port error")
	}
	rand.Seed(time.Now().UnixNano())
	randomPort := rand.Intn(10000) + 50000
	ncCmd := exec.Command("nc", "-z", "127.0.0.1", strconv.Itoa(randomPort))
	err := ncCmd.Run()
	if err == nil {
		reTry++
		goto RANDOM
	}

	cmd := exec.Command("k3d", "create", "--api-port", fmt.Sprintf("0.0.0.0:%d", randomPort), "--name", ClusterName, "--wait", "60")
	if errCreate := cmd.Run(); errCreate != nil {
		return errCreate
	}

	return nil
}

func ValidateK3dCluster() error {
	reTry := 0
RETRY:
	log.Println("wait k3d cluster ready...")
	if reTry >= 50 {
		return errors.New("ValidateK3dCluster error")
	}
	cmd := exec.Command("k3d", "get-kubeconfig", "--name", ClusterName)
	if _, err := cmd.Output(); err != nil {
		time.Sleep(3 * time.Second)
		reTry++
		goto RETRY
	}
	return nil
}

func CleanK3dCluster() error {
	cmd := exec.Command("k3d", "delete", "--name", ClusterName)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
