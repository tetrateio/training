package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"testing"
)

var (
	clusterCount = 25
	zone         = "us-central1-a"
	workshopName = "flexport-nov-2020"

	installSettings = []string{"--set", "profile=demo", "--set", "meshConfig.accessLogFile=/dev/stdout"}
)

func installHelper(prefix ...string) []string {
	return append(prefix, installSettings...)
}

// Note: you need to be logged in to gcloud with access to the cluster credentials
func TestIstioInstall(t *testing.T) {

	// Install
	runPerCluster(t, install)

	// Verify install and external LB IP
	runPerCluster(t, func(t *testing.T) {
		c1 := exec.Command("istioctl", installHelper("manifest", "generate")...)
		c2 := exec.Command("istioctl", "verify-install", "-f", "-")

		c2.Stdin, _ = c1.StdoutPipe()
		c2.Stdout = os.Stdout
		c2.Stderr = os.Stderr
		_ = c2.Start()
		_ = c1.Run()
		if err := c2.Wait(); err != nil {
			t.Errorf("unable to verify istio install: %v", err)
		}

		cmd := exec.Command("kubectl", "get", "service", "-n", "istio-system", "istio-ingressgateway", "-o", "jsonpath={.status.loadBalancer.ingress[0].ip}")
		cmd.Stderr = os.Stderr
		output, err := cmd.Output()
		if err != nil {
			t.Errorf("external loadbalancer request failed: %v", err)
		}
		if net.ParseIP(string(output)) == nil {
			t.Errorf("external loadbalancer check failed received: %s", output)
		}
		if !t.Failed() {
			fmt.Printf("valid external LB IP detected: %s\n", output)
		}
	})

	runPerCluster(t, teardown)
}

func install(t *testing.T) {
	cmd := exec.Command("istioctl", installHelper("install")...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Errorf("istio install failed: %v", err)
	}
}

func teardown(t *testing.T) {
	cmd := exec.Command("istioctl", "x", "uninstall", "--purge", "-y", "--force")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Errorf("istio uninstall failed: %v", err)
	}
}

func runPerCluster(t *testing.T, f func(t *testing.T)) {
	for i := 0; i < clusterCount; i++ {
		cluster := fmt.Sprintf("%v-%03d", workshopName, i)
		credzCmd := exec.Command("gcloud", "container", "clusters", "get-credentials", cluster, "--zone", zone, "--project", cluster)
		credzCmd.Stdout = os.Stdout
		credzCmd.Stderr = os.Stderr
		if err := credzCmd.Run(); err != nil {
			t.Errorf("get credz failed: %v", err)
		}
		f(t)
	}
}
