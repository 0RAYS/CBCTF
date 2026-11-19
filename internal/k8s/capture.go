package k8s

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"

	"CBCTF/internal/config"
	"CBCTF/internal/model"
)

const (
	// The sidecar executes the binary directly, so the capture image needs no
	// shell.
	captureBinaryPath = "/usr/local/bin/rustnet"

	// The mount point must stay outside /root: the capture image ships a mode
	// 0700 /root owned by root, so a mount below it is unreachable for the
	// capture process and the pcap export fails before the first packet.
	captureMountPath = "/mnt"

	// The GeoLite2 city database mount point, outside /root for the same reason.
	captureGeoIPPath = "/etc/GeoLite2-City.mmdb"

	// The user the capture process keeps running as after dropping root. The
	// traffic directory has to stay writable for it.
	captureDropUID int64 = 65534
	captureDropGID int64 = 65534

	// The capture binary only accepts an output directory owned by root or by
	// the user it treats as the invoker, so without this declaration it
	// refuses to write the pcap.
	captureInvokingUIDEnv = "SUDO_UID"
	captureInvokingGIDEnv = "SUDO_GID"
)

// prepareCaptureTrafficDir creates the victim traffic directory on the shared
// volume and hands it over to the unprivileged capture user, because the
// runtime creates the directory as root. Must run before the pods are created.
func prepareCaptureTrafficDir(victim model.Victim) error {
	path := victim.TrafficBasePath()
	if err := os.MkdirAll(path, 0o755); err != nil {
		return err
	}
	if err := os.Chown(path, int(captureDropUID), int(captureDropGID)); err != nil {
		return err
	}
	return nil
}

// captureContainer builds the traffic capture sidecar of a victim pod.
//
// The sidecar shares the pod network namespace with the challenge containers,
// so it captures every packet of the victim. It runs as root because the
// runtime drops the ambient capability set, and the capture needs CAP_BPF and
// CAP_PERFMON for the eBPF process attribution.
func captureContainer(victim model.Victim, pcapName string) corev1.Container {
	// The traffic directory is mounted relative to the root of the shared volume.
	trafficSubPath := strings.TrimPrefix(strings.TrimPrefix(victim.TrafficBasePath(), config.Env.Path), "/")
	args := []string{"-i", "any"}
	geoIPAvailable := false
	if _, err := os.Stat(filepath.Join(config.Env.Path, "GeoLite2-City.mmdb")); err == nil {
		args = append(args, "--geoip-city", captureGeoIPPath)
		geoIPAvailable = true
	}
	args = append(args, "--pcap-export", filepath.Join(captureMountPath, pcapName))
	capture := corev1.Container{
		Resources:       sidecarResources(),
		Name:            CaptureContainerName,
		Image:           config.Env.K8S.CaptureImage,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Command:         append([]string{captureBinaryPath}, args...),
		Env: []corev1.EnvVar{
			{
				Name:  captureInvokingUIDEnv,
				Value: strconv.FormatInt(captureDropUID, 10),
			},
			{
				Name:  captureInvokingGIDEnv,
				Value: strconv.FormatInt(captureDropGID, 10),
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      nfsVolumeName,
				MountPath: captureMountPath,
				SubPath:   trafficSubPath,
			},
		},
		SecurityContext: &corev1.SecurityContext{
			Privileged:               new(false),
			AllowPrivilegeEscalation: new(false),
			ReadOnlyRootFilesystem:   new(false),
			RunAsUser:                new(int64(0)),
			RunAsGroup:               new(int64(0)),
			Capabilities: &corev1.Capabilities{
				Drop: []corev1.Capability{"ALL"},
				Add: []corev1.Capability{
					// Raw socket for the packet capture.
					"NET_RAW",
					// eBPF program loading and the process attribution.
					"BPF", "PERFMON",
					// Dropping the root privileges of the capture binary.
					"SETUID", "SETGID",
					// Removing the capabilities from the bounding set.
					"SETPCAP",
				},
			},
			SeccompProfile: &corev1.SeccompProfile{
				Type: corev1.SeccompProfileTypeRuntimeDefault,
			},
		},
		Stdin: true,
		TTY:   true,
	}
	if geoIPAvailable {
		capture.VolumeMounts = append(capture.VolumeMounts, corev1.VolumeMount{
			Name:      nfsVolumeName,
			MountPath: captureGeoIPPath,
			SubPath:   "GeoLite2-City.mmdb",
			ReadOnly:  true,
		})
	}
	return capture
}
