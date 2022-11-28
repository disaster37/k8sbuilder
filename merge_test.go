package k8sbuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestMergeSliceOrDie(t *testing.T) {
	dst := make([]any, 0)
	src := []any {
		corev1.EnvFromSource {
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "test1",
				},
			},
		},
		corev1.EnvFromSource {
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "test2",
				},
			},
		},
		corev1.EnvFromSource {
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "test1",
				},
			},
		},
	}

	src2 := []any {
		corev1.EnvFromSource {
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "test1",
				},
			},
		},
		corev1.EnvFromSource {
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "test3",
				},
			},
		},
	}

	expected := []any {
		corev1.EnvFromSource {
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "test1",
				},
			},
		},
		corev1.EnvFromSource {
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "test2",
				},
			},
		},
		corev1.EnvFromSource {
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "test3",
				},
			},
		},
	}

	MergeSliceOrDie(&dst, "ConfigMapRef.LocalObjectReference.Name", src, src2)

	assert.Equal(t, expected, dst)


}