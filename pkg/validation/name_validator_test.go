package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNameValidatorValidate(t *testing.T) {
	t.Run("good name", func(t *testing.T) {
		pod := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name: "lifespan",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "lifespan",
					Image: "busybox",
				}},
			},
		}

		v, err := nameValidator{logger()}.Validate(pod)
		assert.Nil(t, err)
		assert.True(t, v.Valid)
	})

	t.Run("bad name", func(t *testing.T) {
		pod := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name: "lifespan-offensive",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "lifespan",
					Image: "busybox",
				}},
			},
		}

		v, err := nameValidator{logger()}.Validate(pod)
		assert.Nil(t, err)
		assert.False(t, v.Valid)
	})
}
