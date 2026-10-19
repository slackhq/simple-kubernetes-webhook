package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	admissionv1 "k8s.io/api/admission/v1"
)

func TestServePodReviews(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		podName string
		allowed bool
		patch   bool
	}{
		{"valid pod", ServeValidatePods, "example", true, false},
		{"rejected pod", ServeValidatePods, "offensive-pod", false, false},
		{"mutated pod", ServeMutatePods, "example", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := `{
				"apiVersion":"admission.k8s.io/v1",
				"kind":"AdmissionReview",
				"request":{
					"uid":"test-request",
					"kind":{"group":"","version":"v1","kind":"Pod"},
					"operation":"CREATE",
					"object":{
						"metadata":{"name":"` + tt.podName + `"},
						"spec":{"containers":[{"name":"app","image":"busybox"}]}
					}
				}
			}`
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			tt.handler(rec, req)

			if rec.Code != http.StatusOK || rec.Header().Get("Content-Type") != "application/json" {
				t.Fatalf("unexpected HTTP response: %d %s", rec.Code, rec.Body)
			}
			var review admissionv1.AdmissionReview
			if err := json.Unmarshal(rec.Body.Bytes(), &review); err != nil {
				t.Fatal(err)
			}
			if review.Kind != "AdmissionReview" || review.APIVersion != "admission.k8s.io/v1" ||
				review.Response == nil || review.Response.UID != "test-request" ||
				review.Response.Allowed != tt.allowed {
				t.Fatalf("unexpected admission response: %s", rec.Body)
			}
			if tt.patch {
				if review.Response.PatchType == nil || *review.Response.PatchType != admissionv1.PatchTypeJSONPatch {
					t.Fatalf("missing JSONPatch type: %s", rec.Body)
				}
				var patch []json.RawMessage
				if err := json.Unmarshal(review.Response.Patch, &patch); err != nil || len(patch) == 0 {
					t.Fatalf("invalid mutation patch: %s", review.Response.Patch)
				}
			} else if !tt.allowed && (review.Response.Result == nil || review.Response.Result.Message != `pod name contains "offensive"`) {
				t.Fatalf("unexpected rejection: %s", rec.Body)
			}
		})
	}
}

func TestServePodReviewsRejectMalformedRequests(t *testing.T) {
	for name, handler := range map[string]http.HandlerFunc{
		"validate": ServeValidatePods,
		"mutate":   ServeMutatePods,
	} {
		t.Run(name, func(t *testing.T) {
			for _, body := range []string{"", "{", `{"request":null}`} {
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				handler(rec, req)
				if rec.Code != http.StatusBadRequest {
					t.Fatalf("body %q returned %d, want 400", body, rec.Code)
				}
			}
		})
	}
}
