package mutation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMinLifespanTolerationsNoLabel(t *testing.T) {
	want := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "lifespan",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "lifespan",
				Image: "busybox",
			}},
			Tolerations: []corev1.Toleration{
				{
					Key:      "acme.com/lifespan-remaining",
					Operator: corev1.TolerationOpExists,
					Effect:   corev1.TaintEffectNoSchedule,
				},
			},
		},
	}

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

	got, err := minLifespanTolerations{logger()}.Mutate(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}

func TestMinLifespanTolerationsLabel(t *testing.T) {
	want := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "lifespan",
			Labels: map[string]string{
				"acme.com/lifespan-requested": "7",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "lifespan",
				Image: "busybox",
			}},
			Tolerations: []corev1.Toleration{
				{
					Key:      "something-unrelated",
					Operator: corev1.TolerationOpExists,
					Effect:   corev1.TaintEffectNoSchedule,
				},
				{
					Key:      "acme.com/lifespan-remaining",
					Operator: corev1.TolerationOpEqual,
					Effect:   corev1.TaintEffectNoSchedule,
					Value:    "14",
				},
				{
					Key:      "acme.com/lifespan-remaining",
					Operator: corev1.TolerationOpEqual,
					Effect:   corev1.TaintEffectNoSchedule,
					Value:    "13",
				},
				{
					Key:      "acme.com/lifespan-remaining",
					Operator: corev1.TolerationOpEqual,
					Effect:   corev1.TaintEffectNoSchedule,
					Value:    "12",
				},
				{
					Key:      "acme.com/lifespan-remaining",
					Operator: corev1.TolerationOpEqual,
					Effect:   corev1.TaintEffectNoSchedule,
					Value:    "11",
				},
				{
					Key:      "acme.com/lifespan-remaining",
					Operator: corev1.TolerationOpEqual,
					Effect:   corev1.TaintEffectNoSchedule,
					Value:    "10",
				},
				{
					Key:      "acme.com/lifespan-remaining",
					Operator: corev1.TolerationOpEqual,
					Effect:   corev1.TaintEffectNoSchedule,
					Value:    "9",
				},
				{
					Key:      "acme.com/lifespan-remaining",
					Operator: corev1.TolerationOpEqual,
					Effect:   corev1.TaintEffectNoSchedule,
					Value:    "8",
				},
				{
					Key:      "acme.com/lifespan-remaining",
					Operator: corev1.TolerationOpEqual,
					Effect:   corev1.TaintEffectNoSchedule,
					Value:    "7",
				},
			},
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "lifespan",
			Labels: map[string]string{
				"acme.com/lifespan-requested": "7",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "lifespan",
				Image: "busybox",
			}},
			Tolerations: []corev1.Toleration{
				{
					Key:      "something-unrelated",
					Operator: corev1.TolerationOpExists,
					Effect:   corev1.TaintEffectNoSchedule,
				},
			},
		},
	}
	got, err := minLifespanTolerations{logger()}.Mutate(pod)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

func TestMinLifespanTolerationsIdempotence(t *testing.T) {
	want := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "lifespan",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "lifespan",
				Image: "busybox",
			}},
			Tolerations: []corev1.Toleration{
				{
					Key:      "acme.com/lifespan-remaining",
					Operator: corev1.TolerationOpExists,
					Effect:   corev1.TaintEffectNoSchedule,
				},
				{
					Key:      "something-unrelated",
					Operator: corev1.TolerationOpExists,
					Effect:   corev1.TaintEffectNoSchedule,
				},
			},
		},
	}

	got, err := minLifespanTolerations{logger()}.Mutate(want.DeepCopy())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}
