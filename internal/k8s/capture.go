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
	// captureBinaryPath is the path of the capture binary inside the capture
	// image.
	//
	// The sidecar executes the binary directly instead of running it through a
	// shell, so the capture does not depend on a shell being part of the
	// image.
	captureBinaryPath = "/usr/local/bin/rustnet"

	// captureMountPath is the mount point of the victim traffic directory
	// inside the capture sidecar.
	//
	// It must not live below /root: the capture image ships a mode 0700
	// /root owned by root, so a mount point below it is unreachable for the
	// capture process and the pcap export fails before the first packet is
	// captured.
	captureMountPath = "/mnt"

	// captureGeoIPPath is the mount point of the GeoLite2 city database. It
	// sits outside /root for the same reason.
	captureGeoIPPath = "/etc/GeoLite2-City.mmdb"

	// captureDropUID and captureDropGID are the user the capture process
	// continues to run as after the capture binary dropped its own root
	// privileges.
	//
	// The capture binary prepares its output files before it drops the
	// privileges, but the traffic directory has to stay writable for the
	// unprivileged user as well.
	captureDropUID int64 = 65534
	captureDropGID int64 = 65534

	// captureInvokingUIDEnv and captureInvokingGIDEnv declare the user the
	// capture binary treats as the one that invoked it. Since 1.7 the capture
	// binary only accepts an output directory that belongs to root or to that
	// user, so without the declaration it refuses to write the pcap.
	captureInvokingUIDEnv = "SUDO_UID"
	captureInvokingGIDEnv = "SUDO_GID"
)

// captureTrafficSubPath returns the victim traffic directory relative to the
// root of the shared volume.
func captureTrafficSubPath(victim model.Victim) string {
	return strings.TrimPrefix(strings.TrimPrefix(victim.TrafficBasePath(), config.Env.Path), "/")
}

// prepareCaptureTrafficDir creates the victim traffic directory on the shared
// volume and hands it over to the unprivileged capture user.
//
// The container runtime creates the directory as root when the sidecar mounts
// it, so the capture process cannot write into it once it dropped its
// privileges. The owner has to be the user the sidecar declares as the
// invoking user (see captureInvokingUIDEnv).
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
// so it captures every packet of the victim.
//
// The container runs as root on purpose. The capture needs CAP_NET_RAW for
// the raw socket and CAP_BPF and CAP_PERFMON for the eBPF process attribution
// of the captured connections. The container runtime drops the ambient
// capability set and the capture image only carries a CAP_NET_RAW file
// capability, so a capture running as the unprivileged image user ends up
// without any capability at all and loses the process attribution. Root keeps
// the capability set of the container across the exec of the capture binary,
// no matter whether the binary itself carries a file capability.
//
// The container is not privileged and does not allow privilege escalation, and
// it only gets the capabilities the capture needs. CAP_SETUID and CAP_SETGID
// let the capture binary drop its own root privileges after the capture is
// initialized, CAP_SETPCAP lets it remove the capabilities from its bounding
// set.
//
// The container also declares the unprivileged capture user as the invoking
// user, see captureInvokingUIDEnv.
func captureContainer(victim model.Victim, pcapName string) corev1.Container {
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
			// The capture binary validates the owner of the traffic directory
			// against the user it drops its privileges to.
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
				SubPath:   captureTrafficSubPath(victim),
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
