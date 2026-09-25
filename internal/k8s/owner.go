package k8s

import (
	"CBCTF/internal/model"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
)

type resourceOwnerKey struct{}

func withResourceOwner(ctx context.Context, owner metav1.OwnerReference) context.Context {
	owner.BlockOwnerDeletion = new(true)
	return context.WithValue(ctx, resourceOwnerKey{}, owner)
}

func resourceOwners(ctx context.Context) []metav1.OwnerReference {
	owner, ok := ctx.Value(resourceOwnerKey{}).(metav1.OwnerReference)
	if !ok {
		return nil
	}
	return []metav1.OwnerReference{owner}
}

func victimRootName(id uint, kind string) string { return fmt.Sprintf("victim-%d-%s", id, kind) }

func createVictimRoot(ctx context.Context, victim model.Victim, kind string) (context.Context, model.RetVal) {
	plan, _ := json.Marshal(victim.Spec.NetworkPlan)
	root, ret := CreateConfigMap(ctx, CreateConfigMapOptions{Name: victimRootName(victim.ID, kind), Labels: VictimLabels(victim, map[string]string{"cbctf.io/root": kind}), Data: map[string]string{"network_plan": string(plan)}})
	if !ret.OK {
		return ctx, ret
	}
	return withResourceOwner(ctx, metav1.OwnerReference{APIVersion: "v1", Kind: "ConfigMap", Name: root.Name, UID: root.UID}), ret
}

// Roots remain discoverable even if startup failed before creating a Pod.
func ListVictimRoots(ctx context.Context) ([]model.Victim, error) {
	c, err := cachedObjects(ctx, "configmaps")
	if err != nil {
		return nil, err
	}
	seen := make(map[uint]bool)
	victims := []model.Victim{}
	for _, obj := range c.informer.GetStore().List() {
		root := obj.(*corev1.ConfigMap)
		kind := root.Labels["cbctf.io/root"]
		if (kind != "workloads" && kind != "network") || root.Labels["cbctf.io/namespace"] != globalNamespace {
			continue
		}
		ids := make(map[string]uint)
		for _, key := range []string{"victim_id", "user_id", "team_id", "contest_id", "challenge_id", "contest_challenge_id"} {
			id, err := strconv.ParseUint(root.Labels[key], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid %s on root %s", key, root.Name)
			}
			ids[key] = uint(id)
		}
		if ids["victim_id"] == 0 || root.Name != victimRootName(ids["victim_id"], kind) {
			return nil, fmt.Errorf("invalid root identity %s", root.Name)
		}
		if seen[ids["victim_id"]] {
			continue
		}
		victim := model.Victim{BaseModel: model.BaseModel{ID: ids["victim_id"]}, UserID: ids["user_id"], ChallengeID: ids["challenge_id"], TeamID: sql.Null[uint]{V: ids["team_id"], Valid: ids["team_id"] > 0}, ContestID: sql.Null[uint]{V: ids["contest_id"], Valid: ids["contest_id"] > 0}, ContestChallengeID: sql.Null[uint]{V: ids["contest_challenge_id"], Valid: ids["contest_challenge_id"] > 0}}
		if err := json.Unmarshal([]byte(root.Data["network_plan"]), &victim.Spec.NetworkPlan); err != nil {
			return nil, fmt.Errorf("invalid network plan on %s: %w", root.Name, err)
		}
		seen[victim.ID] = true
		victims = append(victims, victim)
	}
	return victims, nil
}

// Foreground GC completes the workload tree before the independently owned
// network tree is deleted. Sibling resources cannot enforce this ordering.
func deleteVictimRoot(ctx context.Context, victim model.Victim, kind string) error {
	client := kubeClient.CoreV1().ConfigMaps(globalNamespace)
	name := victimRootName(victim.ID, kind)
	root, err := client.Get(ctx, name, metav1.GetOptions{})
	if apierror.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !labels.SelectorFromSet(VictimLabels(victim, map[string]string{"cbctf.io/root": kind})).Matches(labels.Set(root.Labels)) {
		return fmt.Errorf("root %s is not owned by victim %d", name, victim.ID)
	}
	uid := root.UID
	err = client.Delete(ctx, name, metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}, PropagationPolicy: new(metav1.DeletePropagationForeground)})
	if err != nil && !apierror.IsNotFound(err) {
		return err
	}
	c, err := cachedObjects(ctx, "configmaps")
	if err != nil {
		return err
	}
	return wait.PollUntilContextCancel(ctx, 100*time.Millisecond, true, func(ctx context.Context) (bool, error) {
		obj, exists, err := c.informer.GetStore().GetByKey(globalNamespace + "/" + name)
		if err != nil {
			return false, err
		}
		if exists {
			if obj.(metav1.Object).GetUID() != uid {
				return false, fmt.Errorf("root %s was replaced", name)
			}
			return false, nil
		}
		_, err = client.Get(ctx, name, metav1.GetOptions{})
		return apierror.IsNotFound(err), errUnlessNotFound(err)
	})
}

func deleteVictimVPC(ctx context.Context, victim model.Victim) error {
	name := victim.Spec.NetworkPlan.Name
	if name == "" {
		return nil
	}
	client := ovnClient.KubeovnV1().Vpcs()
	vpc, err := client.Get(ctx, name, metav1.GetOptions{})
	if apierror.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !labels.SelectorFromSet(VictimLabels(victim)).Matches(labels.Set(vpc.Labels)) {
		return fmt.Errorf("VPC %s is not owned by victim %d", name, victim.ID)
	}
	uid := vpc.UID
	if err = client.Delete(ctx, name, metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}, PropagationPolicy: new(metav1.DeletePropagationForeground)}); err != nil && !apierror.IsNotFound(err) {
		return err
	}
	c, err := cachedObjects(ctx, "vpcs")
	if err != nil {
		return err
	}
	return wait.PollUntilContextCancel(ctx, 100*time.Millisecond, true, func(ctx context.Context) (bool, error) {
		obj, exists, err := c.informer.GetStore().GetByKey(name)
		if err != nil {
			return false, err
		}
		if exists {
			if obj.(metav1.Object).GetUID() != uid {
				return false, fmt.Errorf("VPC %s was replaced", name)
			}
			return false, nil
		}
		_, err = client.Get(ctx, name, metav1.GetOptions{})
		return apierror.IsNotFound(err), errUnlessNotFound(err)
	})
}
