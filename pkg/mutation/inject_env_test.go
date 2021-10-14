package mutation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInjectEnvMutate(t *testing.T) {
	want := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name: "test",
				Env: []corev1.EnvVar{
					{
						Name:  "KUBE",
						Value: "true",
					},
				},
			}},
			InitContainers: []corev1.Container{{
				Name: "inittest",
				Env: []corev1.EnvVar{
					{
						Name:  "KUBE",
						Value: "true",
					},
				},
			}},
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name: "test",
			}},
			InitContainers: []corev1.Container{{
				Name: "inittest",
			}},
		},
	}

	got, err := injectEnv{Logger: logger()}.Mutate(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}

func TestHasEnvVar(t *testing.T) {
	ey := corev1.EnvVar{
		Name:  "foo",
		Value: "sball",
	}

	en := corev1.EnvVar{
		Name:  "the_pope",
		Value: "of_nope",
	}

	c := corev1.Container{
		Name: "test",
		Env:  []corev1.EnvVar{ey},
	}

	assert.True(t, HasEnvVar(c, ey))
	assert.False(t, HasEnvVar(c, en))
}
