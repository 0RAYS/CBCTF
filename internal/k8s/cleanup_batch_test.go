package k8s

import (
	"context"
	"errors"
	"testing"

	"CBCTF/internal/model"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ktesting "k8s.io/client-go/testing"
)

func TestServiceCleanupAttemptsOtherObjectsAndReportsFailures(t *testing.T) {
	client := useFakePods(t)
	labels := VictimLabels(model.Victim{ID: 1})
	for _, name := range []string{"one", "two"} {
		if err := client.Tracker().Add(&corev1.Service{Name: name, Namespace: globalNamespace, Labels: labels}); err != nil {
			t.Fatal(err)
		}
	}
	attempts := 0
	client.PrependReactor("delete", "services", func(ktesting.Action) (bool, runtime.Object, error) {
		attempts++
		if attempts == 1 {
			return true, nil, errors.New("single object delete failure")
		}
		return false, nil, nil
	})
	ret := DeleteServiceCollection(context.Background(), labels)
	batch, ok := ret.Data.(*model.BatchResult)
	if ret.OK || !ok || attempts != 2 || batch.Failed != 1 || batch.Succeeded != 1 {
		t.Fatalf("cleanup abandoned siblings or hid failure: attempts=%d ret=%+v batch=%+v", attempts, ret, batch)
	}
}
