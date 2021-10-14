package mutation

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// minLifespanTolerations is a container for mininum lifespan mutation
type minLifespanTolerations struct {
	Logger logrus.FieldLogger
}

// minLifespanTolerations implements the podMutator interface
var _ podMutator = (*minLifespanTolerations)(nil)

// Name returns the minLifespanTolerations short name
func (mpl minLifespanTolerations) Name() string {
	return "min_lifespan"
}

// Mutate returns a new mutated pod according to lifespan tolerations rules
func (mpl minLifespanTolerations) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	const (
		lifespanLabel = "acme.com/lifespan-requested"
		taintKey      = "acme.com/lifespan-remaining"
		taintMaxAge   = 14
	)

	mpl.Logger = mpl.Logger.WithField("mutation", mpl.Name())
	mpod := pod.DeepCopy()

	if pod.Labels == nil || pod.Labels[lifespanLabel] == "" {
		mpl.Logger.WithField("min_lifespan", 0).
			Printf("no lifespan label found, applying default lifespan toleration")

		tn := []corev1.Toleration{{
			Key:      taintKey,
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		}}

		mpod.Spec.Tolerations = appendTolerations(tn, mpod.Spec.Tolerations)
		return mpod, nil
	}

	ts := pod.Labels[lifespanLabel]
	minAge, err := strconv.Atoi(ts)
	if err != nil {
		return nil, fmt.Errorf("pod lifespan label %q is not an integer: %v", ts, err)
	}

	mpl.Logger.WithField("min_lifespan", ts).Printf("setting lifespan tolerations")

	t := []corev1.Toleration{}
	for i := taintMaxAge; i >= minAge; i-- {
		t = append(t, corev1.Toleration{
			Key:      taintKey,
			Operator: corev1.TolerationOpEqual,
			Effect:   corev1.TaintEffectNoSchedule,
			Value:    fmt.Sprint(i),
		})
	}

	mpod.Spec.Tolerations = appendTolerations(t, mpod.Spec.Tolerations)
	return mpod, nil
}

// appendTolerations appends existing to new without duplicating any tolerations
func appendTolerations(new, existing []corev1.Toleration) []corev1.Toleration {
	var toAppend []corev1.Toleration

	for _, n := range new {
		found := false
		for _, e := range existing {
			if reflect.DeepEqual(n, e) {
				found = true
			}
		}
		if !found {
			toAppend = append(toAppend, n)
		}
	}

	return append(existing, toAppend...)
}
