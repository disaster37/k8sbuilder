package k8sbuilder

import (
	"reflect"

	"github.com/imdario/mergo"
	"github.com/thoas/go-funk"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

type PodTemplateBuilder interface{
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
func(h *PodTemplateBuilderDefault) PodTemplate() *corev1.PodTemplateSpec {
	return h.podTemplate
}

// WithPodTemplateSpec permit to use existing podTemplateSpec
func(h *PodTemplateBuilderDefault) WithPodTemplateSpec(pts *corev1.PodTemplateSpec, opts ...WithOption) PodTemplateBuilder {
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
		if err := mergo.Merge(h.podTemplate, pts); err != nil {
			panic(err)
		}
		h.WithAffinity(*pts.Spec.Affinity, Merge).
		WithAnnotations(pts.Annotations, Merge).
		WithContainers(pts.Spec.Containers, Merge).
		WithImagePullSecrets(pts.Spec.ImagePullSecrets, Merge).
		WithInitContainers(pts.Spec.InitContainers, Merge).
		WithLabels(pts.Labels, Merge).
		WithNodeSelector(pts.Spec.NodeSelector, Merge).
		WithTerminationGracePeriodSeconds(*pts.Spec.TerminationGracePeriodSeconds, Merge).
		WithTolerations(pts.Spec.Tolerations, Merge).
		WithVolumes(pts.Spec.Volumes, Merge)
	}

	return h
}

// WithLabels permit to set labels
func (h * PodTemplateBuilderDefault) WithLabels(labels map[string]string, opts ...WithOption) PodTemplateBuilder {
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
		if err := mergo.Merge(h.podTemplate.Labels, labels); err != nil {
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
		if err := mergo.Merge(h.podTemplate.Annotations, annotations); err != nil {
			panic(err)
		}
	}
	
	return h
}

// WithImagePullSecrets permit to set ImagePullSecret
func (h *PodTemplateBuilderDefault) WithImagePullSecrets(ips []corev1.LocalObjectReference, opts ...WithOption) PodTemplateBuilder {
	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Spec.ImagePullSecrets == nil {
		h.podTemplate.Spec.ImagePullSecrets = ips
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Spec.ImagePullSecrets).Elem().IsZero() {
		h.podTemplate.Spec.ImagePullSecrets = ips
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, ref := range ips  {
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
	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Spec.Tolerations == nil {
		h.podTemplate.Spec.Tolerations = tolerations
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Spec.Tolerations).Elem().IsZero() {
		h.podTemplate.Spec.Tolerations = tolerations
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, toleration := range tolerations  {
			if !funk.Contains(h.podTemplate.Spec.Tolerations, toleration){
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
		if err := mergo.Merge(h.podTemplate.Spec.NodeSelector, nodeSelector); err != nil {
			panic(err)
		}
	}
	
	return h
}

// WithInitContainers permit to set init containers
func (h *PodTemplateBuilderDefault) WithInitContainers(containers []corev1.Container, opts ...WithOption) PodTemplateBuilder {
	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Spec.InitContainers == nil {
		h.podTemplate.Spec.InitContainers = containers
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Spec.InitContainers).Elem().IsZero() {
		h.podTemplate.Spec.InitContainers = containers
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, container := range containers  {
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
	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Spec.Containers == nil {
		h.podTemplate.Spec.Containers = containers
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Spec.Containers).Elem().IsZero() {
		h.podTemplate.Spec.Containers = containers
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, container := range containers  {
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
	// Overwrite
	if IsOverwrite(opts) || h.podTemplate.Spec.Volumes == nil {
		h.podTemplate.Spec.Volumes = volumes
		return h
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.podTemplate.Spec.Volumes).Elem().IsZero() {
		h.podTemplate.Spec.Volumes = volumes
		return h
	}

	// Merge
	if IsMerge(opts) {
		for _, volume := range volumes   {
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