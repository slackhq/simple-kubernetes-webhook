package admission

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

func TestPod(t *testing.T) {
	want := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "lifespan",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "lifespan",
				Image: "busybox",
			}},
		},
	}

	raw, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}

	admreq := &admissionv1.AdmissionRequest{
		UID:  types.UID("test"),
		Kind: metav1.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
		Object: runtime.RawExtension{
			Raw:    raw,
			Object: runtime.Object(nil),
		},
	}

	a := Admitter{Request: admreq}
	got, err := a.Pod()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}

func TestReviewResponse(t *testing.T) {
	uid := types.UID("test")
	reason := "fail!"

	want := &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     uid,
			Allowed: false,
			Result: &metav1.Status{
				Code:    418,
				Message: reason,
			},
		},
	}

	got := reviewResponse(uid, false, http.StatusTeapot, reason)
	assert.Equal(t, want, got)
}

func TestPatchReviewResponse(t *testing.T) {
	uid := types.UID("test")
	patchType := admissionv1.PatchTypeJSONPatch
	patch := []byte(`not quite a real patch`)

	want := &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:       uid,
			Allowed:   true,
			PatchType: &patchType,
			Patch:     patch,
		},
	}

	got, err := patchReviewResponse(uid, patch)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}
