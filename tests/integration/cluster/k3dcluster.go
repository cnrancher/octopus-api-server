package cluster

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	ClusterName = "octopus-api-integration-test"
	K3dVersion  = "v1.7.0"
)

var proxy = []string{"http_proxy", "https_proxy", "HTTP_PROXY", "HTTPS_PROXY"}

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
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("curl download k3d error: %v, out: %s", err, string(out))
	}
	if out, err := exec.Command("chmod", "+x", "/tmp/k3d").CombinedOutput(); err != nil {
		return fmt.Errorf("chmod k3d error: %v, out: %s", err, string(out))
	}
	if out, err := exec.Command("/bin/bash", "-c", "mv /tmp/k3d /usr/local/bin/k3d").CombinedOutput(); err != nil {
		return fmt.Errorf("mv k3d error: %v, out: %s", err, string(out))
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
	for i := range proxy {
		os.Unsetenv(proxy[i])
	}
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
	if out, errCreate := cmd.CombinedOutput(); errCreate != nil {
		return fmt.Errorf("k3d create commond error: %v, out: %s", err, string(out))
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
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("k3d delete commond error: %v, out: %s", err, string(out))
	}
	return nil
}
