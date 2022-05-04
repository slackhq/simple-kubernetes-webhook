package mutation

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// Mutator is a container for mutation
type Mutator struct {
	Logger *logrus.Entry
}

// NewMutator returns an initialised instance of Mutator
func NewMutator(logger *logrus.Entry) *Mutator {
	return &Mutator{Logger: logger}
}

// nodeMutators is an interface used to group functions mutating nodes
type nodeMutator interface {
	Mutate(*corev1.Node) (*corev1.Node, error)
	Name() string
}
type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// MutateNodePatch returns a json patch containing all the mutations needed for
// a given node
func (m *Mutator) MutateNodePatch(node *corev1.Node) ([]byte, error) {

	newTaint := corev1.Taint{
		Key:    "foo",
		Value:  "bar",
		Effect: "NoSchedule",
	}

	node.Spec.Taints = append(node.Spec.Taints, newTaint)

	// generate json patch
	op := []PatchOperation{
		{
			Op:    "replace",
			Path:  "/spec/taints",
			Value: node.Spec.Taints,
		},
	}
	patchPayload, _ := json.Marshal(op)

	return patchPayload, nil
}
