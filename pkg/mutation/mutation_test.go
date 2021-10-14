package mutation

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMutatePodPatch(t *testing.T) {
	m := NewMutator(logger())
	got, err := m.MutatePodPatch(pod())
	if err != nil {
		t.Fatal(err)
	}

	p := patch()
	g := string(got)
	assert.Equal(t, p, g)
}

func BenchmarkMutatePodPatch(b *testing.B) {
	m := NewMutator(logger())
	pod := pod()

	for i := 0; i < b.N; i++ {
		_, err := m.MutatePodPatch(pod)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func pod() *corev1.Pod {
	return &corev1.Pod{
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
		},
	}
}

func patch() string {
	patch := `[
		{"op":"add","path":"/spec/containers/0/env","value":[
			{"name":"KUBE","value":"true"}
		]},
		{"op":"add","path":"/spec/tolerations","value":[
			{"effect":"NoSchedule","key":"acme.com/lifespan-remaining","operator":"Equal","value":"14"},
			{"effect":"NoSchedule","key":"acme.com/lifespan-remaining","operator":"Equal","value":"13"},
			{"effect":"NoSchedule","key":"acme.com/lifespan-remaining","operator":"Equal","value":"12"},
			{"effect":"NoSchedule","key":"acme.com/lifespan-remaining","operator":"Equal","value":"11"},
			{"effect":"NoSchedule","key":"acme.com/lifespan-remaining","operator":"Equal","value":"10"},
			{"effect":"NoSchedule","key":"acme.com/lifespan-remaining","operator":"Equal","value":"9"},
			{"effect":"NoSchedule","key":"acme.com/lifespan-remaining","operator":"Equal","value":"8"},
			{"effect":"NoSchedule","key":"acme.com/lifespan-remaining","operator":"Equal","value":"7"}
		]}
]`

	patch = strings.ReplaceAll(patch, "\n", "")
	patch = strings.ReplaceAll(patch, "\t", "")
	patch = strings.ReplaceAll(patch, " ", "")

	return patch
}

func logger() *logrus.Entry {
	mute := logrus.StandardLogger()
	mute.Out = ioutil.Discard
	return mute.WithField("logger", "test")
}
