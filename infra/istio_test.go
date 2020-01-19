package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"testing"
)

// TODO: parameterise this
// Cluster Count
// Zone
// Cluster Prefix

// Note: you need to be logged in to gcloud with access to the cluster credentials
func TestIstioInstall(t *testing.T) {

	// Install
	runPerCluster(t, func(t *testing.T) {
		cmd := exec.Command("istioctl", "manifest", "apply", "--set", "profile=demo", "--set", "values.global.mtls.enabled=true", "--set", "values.global.controlPlaneSecurityEnabled=true")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			t.Errorf("istio install failed: %v", err)
		}
	})

	// Verify install and external LB IP
	runPerCluster(t, func(t *testing.T) {
		c1 := exec.Command("istioctl", "manifest", "generate", "--set", "profile=demo", "--set", "values.global.mtls.enabled=true", "--set", "values.global.controlPlaneSecurityEnabled=true")
		c2 := exec.Command("istioctl", "verify-install", "-f", "-")

		c2.Stdin, _ = c1.StdoutPipe()
		c2.Stdout = os.Stdout
		_ = c2.Start()
		_ = c1.Run()
		if err := c2.Wait(); err != nil {
			t.Errorf("unable to verify istio install: %v", err)
		}

		cmd := exec.Command("kubectl", "get", "service", "-n", "istio-system", "istio-ingressgateway", "-o", "jsonpath={.status.loadBalancer.ingress[0].ip}")
		output, err := cmd.Output()
		if err != nil {
			t.Errorf("external loadbalancer request failed: %v", err)
		}
		if net.ParseIP(string(output)) == nil {
			t.Errorf("external loadbalancer check failed received: %s", output)
		}
		t.Logf("valid external LB IP detected: %s", output)
	})

	// Teardown
	runPerCluster(t, func(t *testing.T) {
		cmd := exec.Command("kubectl", "delete", "namespace", "istio-system", "--ignore-not-found=true")
		if err := cmd.Run(); err != nil {
			t.Errorf("istio delete failed: %v", err)
		}
	})
}

func runPerCluster(t *testing.T, f func(t *testing.T)) {
	zone := "us-central1-a"

	for i := 0; i < 65; i++ {
		cluster := fmt.Sprintf("nist-2020-%03d", i)
		credzCmd := exec.Command("gcloud", "container", "clusters", "get-credentials", cluster, "--zone", zone, "--project", cluster)
		credzCmd.Stdout = os.Stdout
		credzCmd.Stderr = os.Stderr
		if err := credzCmd.Run(); err != nil {
			t.Errorf("get credz failed: %v", err)
		}
		f(t)
	}
}
