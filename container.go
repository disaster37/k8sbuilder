package k8sbuilder

import (
	"reflect"

	"github.com/thoas/go-funk"
	corev1 "k8s.io/api/core/v1"
)

type ContainerBuilder interface {
	Container() *corev1.Container
	WithContainer(container *corev1.Container, opts ...WithOption) ContainerBuilder
	WithEnvFrom(envFroms []corev1.EnvFromSource, opts ...WithOption) ContainerBuilder
	WithEnv(envs []corev1.EnvVar, opts ...WithOption) ContainerBuilder
	WithImage(image string, opts ...WithOption) ContainerBuilder
	WithImagePullPolicy(pullPolicy corev1.PullPolicy, opts ...WithOption) ContainerBuilder
	WithPort(ports []corev1.ContainerPort, opts ...WithOption) ContainerBuilder
	WithResource(ressources *corev1.ResourceRequirements, opts ...WithOption) ContainerBuilder
	WithSecurityContext(sc *corev1.SecurityContext, opts ...WithOption) ContainerBuilder
	WithVolumeMount(volumeMounts []corev1.VolumeMount, opts ...WithOption) ContainerBuilder
	WithLivenessProbe(probe *corev1.Probe, opts ...WithOption) ContainerBuilder
	WithReadinessProbe(probe *corev1.Probe, opts ...WithOption) ContainerBuilder
	WithStartupProbe(probe *corev1.Probe, opts ...WithOption) ContainerBuilder
}

type ContainerBuilderDefault struct {
	container *corev1.Container
}

// NewContainerBuilder permit to get new container builder
func NewContainerBuilder() ContainerBuilder {
	return &ContainerBuilderDefault{
		container: &corev1.Container{},
	}
}

// Container permit to get current container
func (h *ContainerBuilderDefault) Container() *corev1.Container {
	return h.container
}

// WithContainer permit to set existing container
func (h *ContainerBuilderDefault) WithContainer(container *corev1.Container, opts ...WithOption) ContainerBuilder {

	if container == nil {
		return h
	}

	// Overwrite
	if IsOverwrite(opts) {
		h.container = container
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.container).Elem().IsZero() {
		h.container = container
		return h
	}

	// Merge
	if IsMerge(opts) {
		orgContainer := h.container.DeepCopy()
		if err := MergeK8s(h.container, h.container, container); err != nil {
			panic(err)
		}
		h.WithEnv(orgContainer.Env).
			WithEnv(container.Env, Merge).
			WithEnvFrom(orgContainer.EnvFrom).
			WithEnvFrom(container.EnvFrom, Merge).
			WithPort(orgContainer.Ports).
			WithPort(container.Ports, Merge).
			WithVolumeMount(orgContainer.VolumeMounts).
			WithVolumeMount(container.VolumeMounts, Merge)
	}

	return h
}

// WithEnvFrom permit to set envFrom
func (h *ContainerBuilderDefault) WithEnvFrom(envFroms []corev1.EnvFromSource, opts ...WithOption) ContainerBuilder {

	// Copy to avoid overwrite envFroms
	tmpEnvFrom := make([]corev1.EnvFromSource, len(envFroms))
	copy(tmpEnvFrom, envFroms)

	// Overwrite
	if IsOverwrite(opts) || h.container.EnvFrom == nil {
		h.container.EnvFrom = tmpEnvFrom
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.container.EnvFrom).Elem().IsZero() {
		h.container.EnvFrom = tmpEnvFrom
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, envFrom := range tmpEnvFrom {
			if !funk.Contains(h.container.EnvFrom, envFrom) {
				h.container.EnvFrom = append(h.container.EnvFrom, envFrom)
			}
		}
	}

	return h
}

// WithEnv permit to set env
func (h *ContainerBuilderDefault) WithEnv(envs []corev1.EnvVar, opts ...WithOption) ContainerBuilder {

	// Copy to avoid overwrite envFroms
	tmpEnvs := make([]corev1.EnvVar, len(envs))
	copy(tmpEnvs, envs)

	// Overwrite
	if IsOverwrite(opts) || h.container.Env == nil {
		h.container.Env = tmpEnvs
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.container.Env).Elem().IsZero() {
		h.container.Env = tmpEnvs
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, env := range tmpEnvs {
			if !funk.Contains(h.container.Env, env) {
				h.container.Env = append(h.container.Env, env)
			}
		}
	}

	return h
}

// WithImage permit to set image
func (h *ContainerBuilderDefault) WithImage(image string, opts ...WithOption) ContainerBuilder {
	// Overwrite
	if IsOverwrite(opts) || IsMerge(opts) || h.container.Image == "" {
		h.container.Image = image
		return h
	}

	return h
}

// WithImagePullPolicy permit to set image pull policy
func (h *ContainerBuilderDefault) WithImagePullPolicy(pullPolicy corev1.PullPolicy, opts ...WithOption) ContainerBuilder {
	// Overwrite
	if IsOverwrite(opts) || IsMerge(opts) || string(h.container.ImagePullPolicy) == "" {
		h.container.ImagePullPolicy = pullPolicy
		return h
	}

	return h
}

func (h *ContainerBuilderDefault) WithPort(ports []corev1.ContainerPort, opts ...WithOption) ContainerBuilder {

	// Copy to avoid overwrite ports
	tmpPorts := make([]corev1.ContainerPort, len(ports))
	copy(tmpPorts, ports)

	// Overwrite
	if IsOverwrite(opts) || h.container.Ports == nil {
		h.container.Ports = tmpPorts
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.container.Ports).Elem().IsZero() {
		h.container.Ports = tmpPorts
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, port := range tmpPorts {
			index := funk.IndexOf(h.container.Ports, func(o corev1.ContainerPort) bool {
				return port.ContainerPort == o.ContainerPort
			})

			if index == -1 {
				h.container.Ports = append(h.container.Ports, port)
			} else {
				h.container.Ports[index] = port
			}
		}
	}

	return h
}

// WithResource permit to set resources
func (h *ContainerBuilderDefault) WithResource(resources *corev1.ResourceRequirements, opts ...WithOption) ContainerBuilder {
	if resources == nil {
		return h
	}

	// Overwrite
	if IsOverwrite(opts) {
		h.container.Resources = *resources
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.container.Resources).IsZero() {
		h.container.Resources = *resources
		return h
	}

	// Merge
	if IsMerge(opts) {
		if err := MergeK8s(&h.container.Resources, h.container.Resources, resources); err != nil {
			panic(err)
		}
	}

	return h
}

// WithSecurityContext permit to set security context
func (h *ContainerBuilderDefault) WithSecurityContext(sc *corev1.SecurityContext, opts ...WithOption) ContainerBuilder {
	// Overwrite
	if IsOverwrite(opts) || h.container.SecurityContext == nil {
		h.container.SecurityContext = sc
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.container.SecurityContext).Elem().IsZero() {
		h.container.SecurityContext = sc
		return h
	}

	// Merge
	if IsMerge(opts) {
		if err := MergeK8s(h.container.SecurityContext, h.container.SecurityContext, sc); err != nil {
			panic(err)
		}
	}

	return h
}

// WithVolumeMount permit to set volume mounts
func (h *ContainerBuilderDefault) WithVolumeMount(volumeMounts []corev1.VolumeMount, opts ...WithOption) ContainerBuilder {

	// Copy to avoid overwrite volumeMounts
	tmpVolumeMount := make([]corev1.VolumeMount, len(volumeMounts))
	copy(tmpVolumeMount, volumeMounts)

	// Overwrite
	if IsOverwrite(opts) || h.container.VolumeMounts == nil {
		h.container.VolumeMounts = tmpVolumeMount
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.container.VolumeMounts).Elem().IsZero() {
		h.container.VolumeMounts = tmpVolumeMount
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, volumeMount := range tmpVolumeMount {
			index := funk.IndexOf(h.container.VolumeMounts, func(o corev1.VolumeMount) bool {
				return volumeMount.MountPath == o.MountPath && volumeMount.SubPath == o.SubPath
			})

			if index == -1 {
				h.container.VolumeMounts = append(h.container.VolumeMounts, volumeMount)
			} else {
				h.container.VolumeMounts[index] = volumeMount
			}
		}
	}

	return h
}

func (h *ContainerBuilderDefault) WithLivenessProbe(probe *corev1.Probe, opts ...WithOption) ContainerBuilder {
	// Overwrite
	if IsOverwrite(opts) || h.container.LivenessProbe == nil {
		h.container.LivenessProbe = probe
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.container.LivenessProbe).Elem().IsZero() {
		h.container.LivenessProbe = probe
		return h
	}

	// Merge
	if IsMerge(opts) {
		if err := MergeK8s(&h.container.LivenessProbe, h.container.LivenessProbe, probe); err != nil {
			panic(err)
		}
	}

	return h
}

func (h *ContainerBuilderDefault) WithReadinessProbe(probe *corev1.Probe, opts ...WithOption) ContainerBuilder {
	// Overwrite
	if IsOverwrite(opts) || h.container.ReadinessProbe == nil {
		h.container.ReadinessProbe = probe
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.container.ReadinessProbe).Elem().IsZero() {
		h.container.ReadinessProbe = probe
		return h
	}

	// Merge
	if IsMerge(opts) {
		if err := MergeK8s(&h.container.ReadinessProbe, h.container.ReadinessProbe, probe); err != nil {
			panic(err)
		}
	}

	return h
}

func (h *ContainerBuilderDefault) WithStartupProbe(probe *corev1.Probe, opts ...WithOption) ContainerBuilder {
	// Overwrite
	if IsOverwrite(opts) || h.container.StartupProbe == nil {
		h.container.StartupProbe = probe
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.container.StartupProbe).Elem().IsZero() {
		h.container.StartupProbe = probe
		return h
	}

	// Merge
	if IsMerge(opts) {
		if err := MergeK8s(&h.container.StartupProbe, h.container.StartupProbe, probe); err != nil {
			panic(err)
		}
	}

	return h
}
