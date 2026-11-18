package k8s

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	corev1 "k8s.io/api/core/v1"

	"CBCTF/internal/config"
	"CBCTF/internal/model"
)

func TestCaptureContainerRunsOutsideRootAndKeepsCapabilities(t *testing.T) {
	previousConfig := config.Env
	path := t.TempDir()
	config.Env = &config.Config{Path: path}
	t.Cleanup(func() { config.Env = previousConfig })

	victim := model.Victim{BaseModel: model.BaseModel{ID: 7}}
	capture := captureContainer(victim, "pod-victim-7.pcap")

	if capture.Command[0] != captureBinaryPath {
		t.Fatalf("capture command must exec the capture binary directly, got %q", capture.Command[0])
	}
	if slices.Contains(capture.Command, "/bin/sh") {
		t.Fatalf("capture command must not run through a shell: %v", capture.Command)
	}
	wantPcap := filepath.Join(captureMountPath, "pod-victim-7.pcap")
	if !slices.Contains(capture.Command, wantPcap) {
		t.Fatalf("capture command must export %q, got %v", wantPcap, capture.Command)
	}
	for _, mount := range capture.VolumeMounts {
		if mount.MountPath == captureMountPath && mount.MountPath != "/root/mnt" {
			continue
		}
		if mount.MountPath == "/root/mnt" || mount.MountPath == "/root/GeoLite2-City.mmdb" {
			t.Fatalf("capture mounts must not live below the private /root: %v", capture.VolumeMounts)
		}
	}
	if len(capture.VolumeMounts) != 1 {
		t.Fatalf("capture mounts the traffic directory only, got %v", capture.VolumeMounts)
	}
	if capture.VolumeMounts[0].MountPath != captureMountPath {
		t.Fatalf("traffic directory must be mounted at %s, got %s", captureMountPath, capture.VolumeMounts[0].MountPath)
	}

	securityContext := capture.SecurityContext
	if securityContext == nil {
		t.Fatal("capture needs a security context")
	}
	if securityContext.Privileged == nil || *securityContext.Privileged {
		t.Fatal("capture must not be privileged")
	}
	if securityContext.AllowPrivilegeEscalation == nil || *securityContext.AllowPrivilegeEscalation {
		t.Fatal("capture must not allow privilege escalation")
	}
	// The capture needs CAP_BPF and CAP_PERFMON for the eBPF process
	// attribution. Only a root container keeps them, because the container
	// runtime drops the ambient capability set and the capture image only
	// carries a CAP_NET_RAW file capability.
	if securityContext.RunAsUser == nil || *securityContext.RunAsUser != 0 {
		t.Fatalf("capture must run as root to keep the eBPF capabilities, got %v", securityContext.RunAsUser)
	}
	if securityContext.RunAsGroup == nil || *securityContext.RunAsGroup != 0 {
		t.Fatalf("capture must run with the root group, got %v", securityContext.RunAsGroup)
	}
	if securityContext.Capabilities == nil {
		t.Fatal("capture needs an explicit capability set")
	}
	if !slices.Contains(securityContext.Capabilities.Drop, corev1.Capability("ALL")) {
		t.Fatalf("capture must drop all capabilities, got %v", securityContext.Capabilities.Drop)
	}
	for _, capability := range []corev1.Capability{"NET_RAW", "BPF", "PERFMON", "SETUID", "SETGID", "SETPCAP"} {
		if !slices.Contains(securityContext.Capabilities.Add, capability) {
			t.Fatalf("capture needs %s, got %v", capability, securityContext.Capabilities.Add)
		}
	}
}

func TestCaptureContainerMountsGeoIPDatabaseOutsideRoot(t *testing.T) {
	previousConfig := config.Env
	path := t.TempDir()
	config.Env = &config.Config{Path: path}
	t.Cleanup(func() { config.Env = previousConfig })

	if err := os.WriteFile(filepath.Join(path, "GeoLite2-City.mmdb"), []byte("db"), 0o644); err != nil {
		t.Fatal(err)
	}
	capture := captureContainer(model.Victim{}, "frpc.pcap")
	if !slices.Contains(capture.Command, "--geoip-city") {
		t.Fatalf("capture must use the GeoIP database, got %v", capture.Command)
	}
	if !slices.Contains(capture.Command, captureGeoIPPath) {
		t.Fatalf("capture must read the GeoIP database from %s, got %v", captureGeoIPPath, capture.Command)
	}
	found := false
	for _, mount := range capture.VolumeMounts {
		if mount.MountPath == captureGeoIPPath {
			found = true
			if !mount.ReadOnly {
				t.Fatal("the GeoIP database mount must be read only")
			}
		}
	}
	if !found {
		t.Fatalf("capture must mount the GeoIP database, got %v", capture.VolumeMounts)
	}
}

func TestPrepareCaptureTrafficDirHandsDirectoryToCaptureUser(t *testing.T) {
	previousConfig := config.Env
	path := t.TempDir()
	config.Env = &config.Config{Path: path}
	t.Cleanup(func() { config.Env = previousConfig })

	victim := model.Victim{BaseModel: model.BaseModel{ID: 3}}
	if err := prepareCaptureTrafficDir(victim); err != nil {
		if runtime.GOOS == "windows" {
			// Ownership changes are a linux only concept.
			t.Skip("ownership changes are not supported on windows")
		}
		t.Fatal(err)
	}
	info, err := os.Stat(victim.TrafficBasePath())
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o200 == 0 {
		t.Fatalf("the capture user must be able to write the traffic directory, got %v", info.Mode().Perm())
	}
	if runtime.GOOS == "windows" {
		return
	}
	uid, gid, ok := captureDirOwner(info)
	if !ok {
		t.Skip("ownership information is not available")
	}
	if uid != int(captureDropUID) || gid != int(captureDropGID) {
		t.Fatalf("the traffic directory must belong to the capture user %d:%d, got %d:%d",
			captureDropUID, captureDropGID, uid, gid)
	}
}
