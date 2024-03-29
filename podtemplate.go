package k8sbuilder

import (
	"reflect"

	"github.com/imdario/mergo"
	"github.com/thoas/go-funk"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

type PodTemplateBuilder interface {
	WithPodTemplateSpec(pts *corev1.PodTemplateSpec, opts ...WithOption) PodTemplateBuilder
	WithLabels(labels map[string]string, opts ...WithOption) PodTemplateBuilder
	WithAnnotations(annotations map[string]string, opts ...WithOption) PodTemplateBuilder
	WithImagePullSecrets(ips []corev1.LocalObjectReference, opts ...WithOption) PodTemplateBuilder
	WithTerminationGracePeriodSeconds(nb int64, opts ...WithOption) PodTemplateBuilder
	WithTolerations(tolerations []corev1.Toleration, opts ...WithOption) PodTemplateBuilder
	WithNodeSelector(nodeSelector map[string]string, opts ...WithOption) PodTemplateBuilder
	WithInitContainers(containers []corev1.Container, opts ...WithOption) PodTemplateBuilder
	WithContainers(containers []corev1.Container, opts ...WithOption) PodTemplateBuilder
	WithVolumes(volumes []corev1.Volume, opts ...WithOption) PodTemplateBuilder
	WithAffinity(affinity corev1.Affinity, opts ...WithOption) PodTemplateBuilder
	WithSecurityContext(sc *corev1.PodSecurityContext, opts ...WithOption) PodTemplateBuilder
	PodTemplate() *corev1.PodTemplateSpec
}

type PodTemplateBuilderDefault struct {
	podTemplate *corev1.PodTemplateSpec
}

// NewPodTemplateBuilder permit to init pod template builder
func NewPodTemplateBuilder() PodTemplateBuilder {
	return &PodTemplateBuilderDefault{
		podTemplate: &corev1.PodTemplateSpec{},
	}
}

// PodTemplate permit to get current pod template
func (h *PodTemplateBuilderDefault) PodTemplate() *corev1.PodTemplateSpec {
	return h.podTemplate
}

// WithPodTemplateSpec permit to use existing podTemplateSpec
func (h *PodTemplateBuilderDefault) WithPodTemplateSpec(pts *corev1.PodTemplateSpec, opts ...WithOption) PodTemplateBuilder {
	if pts == nil {
		return h
	}

	// Overwrite
	if IsOverwrite(opts) {
		h.podTemplate = pts
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate).Elem().IsZero() {
		h.podTemplate = pts
		return h
	}

	// Merge
	if IsMerge(opts) {
		orgPts := h.podTemplate.DeepCopy()

		if err := MergeK8s(h.podTemplate, h.podTemplate, pts); err != nil {
			panic(err)
		}

		h.WithContainers(orgPts.Spec.Containers).
			WithContainers(pts.Spec.Containers, Merge).
			WithImagePullSecrets(orgPts.Spec.ImagePullSecrets).
			WithImagePullSecrets(pts.Spec.ImagePullSecrets, Merge).
			WithInitContainers(orgPts.Spec.InitContainers).
			WithInitContainers(pts.Spec.InitContainers, Merge).
			WithTolerations(orgPts.Spec.Tolerations).
			WithTolerations(pts.Spec.Tolerations, Merge).
			WithVolumes(orgPts.Spec.Volumes).
			WithVolumes(pts.Spec.Volumes, Merge)
	}

	return h
}

// WithLabels permit to set labels
func (h *PodTemplateBuilderDefault) WithLabels(labels map[string]string, opts ...WithOption) PodTemplateBuilder {
	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Labels == nil {
		h.podTemplate.Labels = labels
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Labels).Elem().IsZero() {
		h.podTemplate.Labels = labels
		return h
	}

	// Merge
	if IsMerge(opts) && labels != nil {
		if err := mergo.Merge(&h.podTemplate.Labels, labels); err != nil {
			panic(err)
		}
	}

	return h
}

// WithAnnotations permit to set annotations
func (h *PodTemplateBuilderDefault) WithAnnotations(annotations map[string]string, opts ...WithOption) PodTemplateBuilder {
	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Annotations == nil {
		h.podTemplate.Annotations = annotations
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Annotations).Elem().IsZero() {
		h.podTemplate.Annotations = annotations
		return h
	}

	// Merge
	if IsMerge(opts) && annotations != nil {
		if err := mergo.Merge(&h.podTemplate.Annotations, annotations); err != nil {
			panic(err)
		}
	}

	return h
}

// WithImagePullSecrets permit to set ImagePullSecret
func (h *PodTemplateBuilderDefault) WithImagePullSecrets(ips []corev1.LocalObjectReference, opts ...WithOption) PodTemplateBuilder {

	var tmpIps []corev1.LocalObjectReference

	// Avoid overwrite ips
	if ips != nil {
		tmpIps := make([]corev1.LocalObjectReference, len(ips))
		copy(tmpIps, ips)
	}

	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Spec.ImagePullSecrets == nil {
		h.podTemplate.Spec.ImagePullSecrets = tmpIps
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Spec.ImagePullSecrets).Elem().IsZero() {
		h.podTemplate.Spec.ImagePullSecrets = tmpIps
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, ref := range tmpIps {
			if !funk.Contains(h.podTemplate.Spec.ImagePullSecrets, func(o corev1.LocalObjectReference) bool {
				return ref.Name == o.Name
			}) {
				h.podTemplate.Spec.ImagePullSecrets = append(h.podTemplate.Spec.ImagePullSecrets, ref)
			}
		}
	}

	return h
}

// WithTerminationGracePeriodSeconds permit to set TerminationGracePeriodSeconds
func (h *PodTemplateBuilderDefault) WithTerminationGracePeriodSeconds(nb int64, opts ...WithOption) PodTemplateBuilder {
	// Overwrite
	if IsOverwrite(opts) || IsMerge(opts) || h.podTemplate.Spec.TerminationGracePeriodSeconds == nil {
		h.podTemplate.Spec.TerminationGracePeriodSeconds = pointer.Int64(nb)
		return h
	}

	return h
}

// WithTolerations permit to set tolerations
func (h *PodTemplateBuilderDefault) WithTolerations(tolerations []corev1.Toleration, opts ...WithOption) PodTemplateBuilder {

	var tmpTolerations []corev1.Toleration

	// To avoid to overwrite tolerations
	if tolerations != nil {
		tmpTolerations = make([]corev1.Toleration, len(tolerations))
		copy(tmpTolerations, tolerations)
	}

	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Spec.Tolerations == nil {
		h.podTemplate.Spec.Tolerations = tmpTolerations
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Spec.Tolerations).Elem().IsZero() {
		h.podTemplate.Spec.Tolerations = tmpTolerations
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, toleration := range tmpTolerations {
			if !funk.Contains(h.podTemplate.Spec.Tolerations, toleration) {
				h.podTemplate.Spec.Tolerations = append(h.podTemplate.Spec.Tolerations, toleration)
			}
		}
	}

	return h
}

// WithNodeSelector permit to set nodeSelector
func (h *PodTemplateBuilderDefault) WithNodeSelector(nodeSelector map[string]string, opts ...WithOption) PodTemplateBuilder {
	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Spec.NodeSelector == nil {
		h.podTemplate.Spec.NodeSelector = nodeSelector
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Spec.NodeSelector).Elem().IsZero() {
		h.podTemplate.Spec.NodeSelector = nodeSelector
		return h
	}

	// Merge
	if IsMerge(opts) && nodeSelector != nil {
		if err := mergo.Merge(&h.podTemplate.Spec.NodeSelector, nodeSelector); err != nil {
			panic(err)
		}
	}

	return h
}

// WithInitContainers permit to set init containers
func (h *PodTemplateBuilderDefault) WithInitContainers(containers []corev1.Container, opts ...WithOption) PodTemplateBuilder {

	var tmpContainers []corev1.Container

	// To avoid overwrite
	if containers != nil {
		tmpContainers = make([]corev1.Container, len(containers))
		copy(tmpContainers, containers)
	}

	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Spec.InitContainers == nil {
		h.podTemplate.Spec.InitContainers = tmpContainers
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Spec.InitContainers).Elem().IsZero() {
		h.podTemplate.Spec.InitContainers = tmpContainers
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, container := range tmpContainers {
			index := funk.IndexOf(h.podTemplate.Spec.InitContainers, func(o corev1.Container) bool {
				return container.Name == o.Name
			})
			if index == -1 {
				h.podTemplate.Spec.InitContainers = append(h.podTemplate.Spec.InitContainers, container)
			} else {
				h.podTemplate.Spec.InitContainers[index] = *NewContainerBuilder().
					WithContainer(&h.podTemplate.Spec.InitContainers[index]).
					WithContainer(&container, Merge).
					Container()

			}
		}
	}

	return h
}

// WithContainers permit to set containers
func (h *PodTemplateBuilderDefault) WithContainers(containers []corev1.Container, opts ...WithOption) PodTemplateBuilder {

	var tmpContainers []corev1.Container

	// To avoid overwrite
	if containers != nil {
		tmpContainers = make([]corev1.Container, len(containers))
		copy(tmpContainers, containers)
	}

	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Spec.Containers == nil {
		h.podTemplate.Spec.Containers = tmpContainers
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Spec.Containers).Elem().IsZero() {
		h.podTemplate.Spec.Containers = tmpContainers
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, container := range tmpContainers {
			index := funk.IndexOf(h.podTemplate.Spec.InitContainers, func(o corev1.Container) bool {
				return container.Name == o.Name
			})
			if index == -1 {
				h.podTemplate.Spec.Containers = append(h.podTemplate.Spec.Containers, container)
			} else {
				h.podTemplate.Spec.Containers[index] = *NewContainerBuilder().
					WithContainer(&h.podTemplate.Spec.Containers[index]).
					WithContainer(&container, Merge).
					Container()
			}
		}
	}

	return h
}

// WithContainers permit to set containers
func (h *PodTemplateBuilderDefault) WithVolumes(volumes []corev1.Volume, opts ...WithOption) PodTemplateBuilder {

	var tmpVolumes []corev1.Volume

	// To avoid to overwrite volumes
	if volumes != nil {
		tmpVolumes = make([]corev1.Volume, len(volumes))
		copy(tmpVolumes, volumes)
	}

	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Spec.Volumes == nil {
		h.podTemplate.Spec.Volumes = tmpVolumes
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Spec.Volumes).Elem().IsZero() {
		h.podTemplate.Spec.Volumes = tmpVolumes
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, volume := range tmpVolumes {
			index := funk.IndexOf(h.podTemplate.Spec.Volumes, func(o corev1.Volume) bool {
				return volume.Name == o.Name
			})
			if index == -1 {
				h.podTemplate.Spec.Volumes = append(h.podTemplate.Spec.Volumes, volume)
			} else {
				if err := MergeK8s(&h.podTemplate.Spec.Volumes[index], h.podTemplate.Spec.Volumes[index], volume); err != nil {
					panic(err)
				}
			}
		}
	}

	return h
}

// WithAffinity permit to set affinity
func (h *PodTemplateBuilderDefault) WithAffinity(affinity corev1.Affinity, opts ...WithOption) PodTemplateBuilder {
	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Spec.Affinity == nil {
		h.podTemplate.Spec.Affinity = &affinity
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Spec.Affinity).Elem().IsZero() {
		h.podTemplate.Spec.Affinity = &affinity
		return h
	}

	// Merge
	if IsMerge(opts) {
		if err := MergeK8s(h.podTemplate.Spec.Affinity, h.podTemplate.Spec.Affinity, affinity); err != nil {
			panic(err)
		}
	}

	return h
}

// WithSecurityContext permit to set security context
func (h *PodTemplateBuilderDefault) WithSecurityContext(sc *corev1.PodSecurityContext, opts ...WithOption) PodTemplateBuilder {
	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Spec.SecurityContext == nil {
		h.podTemplate.Spec.SecurityContext = sc
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Spec.SecurityContext).Elem().IsZero() {
		h.podTemplate.Spec.SecurityContext = sc
		return h
	}

	// Merge
	if IsMerge(opts) {
		if err := MergeK8s(h.podTemplate.Spec.SecurityContext, h.podTemplate.Spec.SecurityContext, sc); err != nil {
			panic(err)
		}
	}

	return h
}
