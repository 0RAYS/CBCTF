package k8s

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"maps"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
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

func victimRootName(id uint, kind string) string {
	return fmt.Sprintf("victim-%d-%s", id, kind)
}

func createVictimRoot(ctx context.Context, victim model.Victim, kind string) (context.Context, model.RetVal) {
	plan, _ := json.Marshal(victim.Spec.NetworkPlan)
	root, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
		Name:   victimRootName(victim.ID, kind),
		Labels: VictimLabels(victim, map[string]string{"cbctf.io/root": kind}),
		Data:   map[string]string{"network_plan": string(plan)},
	})
	if !ret.OK {
		return ctx, ret
	}
	if root == nil {
		return ctx, model.FailRetVal(i18n.K8S.CreateError, "ConfigMap", victimRootName(victim.ID, kind), errMissingConfigMap)
	}
	return withResourceOwner(ctx, metav1.OwnerReference{
		APIVersion: "v1",
		Kind:       "ConfigMap",
		Name:       root.Name,
		UID:        root.UID,
	}), ret
}

// withVictimRootOwner 把 ctx 绑到实例的资源树上。
// 租约期创建的资源(如玩家网关)挂在实例的根对象下,实例销毁时一并回收。
func withVictimRootOwner(ctx context.Context, victim model.Victim, kind string) (context.Context, model.RetVal) {
	name := victimRootName(victim.ID, kind)
	root, ret := GetConfigMap(ctx, name)
	if !ret.OK {
		return ctx, ret
	}
	if root == nil {
		return ctx, model.FailRetVal(i18n.K8S.GetError, "ConfigMap", name, errMissingConfigMap)
	}
	return withResourceOwner(ctx, metav1.OwnerReference{
		APIVersion: "v1",
		Kind:       "ConfigMap",
		Name:       root.Name,
		UID:        root.UID,
	}), model.SuccessRetVal()
}

// VictimIdentityLabels 返回实例资源的稳定身份标签。
// 租约会改写队伍与用户,归属校验只比对不会变的部分。
func VictimIdentityLabels(victim model.Victim, tags ...map[string]string) map[string]string {
	labels := map[string]string{
		"cbctf.io/namespace": globalNamespace,
		"victim_id":          strconv.Itoa(int(victim.ID)),
		"challenge_id":       strconv.Itoa(int(victim.ChallengeID)),
	}
	if len(tags) > 0 {
		maps.Copy(labels, tags[0])
	}
	return labels
}

// hasVictimIdentity 判断资源的标签是否属于该实例。
func hasVictimIdentity(resourceLabels map[string]string, victim model.Victim, tags ...map[string]string) bool {
	for key, value := range VictimIdentityLabels(victim, tags...) {
		if resourceLabels[key] != value {
			return false
		}
	}
	return true
}

// ListVictimRoots finds the victims from their root ConfigMaps, so they stay
// discoverable even if startup failed before creating a Pod.
func ListVictimRoots(ctx context.Context) ([]model.Victim, error) {
	c, err := cachedObjects(ctx, "configmaps")
	if err != nil {
		return nil, err
	}
	seen := make(map[uint]bool)
	var victims []model.Victim
	for _, obj := range c.informer.GetStore().List() {
		root := obj.(*corev1.ConfigMap)
		kind := root.Labels["cbctf.io/root"]
		if (kind != "workloads" && kind != "network") || root.Labels["cbctf.io/namespace"] != globalNamespace {
			continue
		}
		ids := make(map[string]uint)
		for _, key := range []string{"victim_id", "user_id", "team_id", "contest_id", "challenge_id"} {
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
		victim := model.Victim{
			ID: ids["victim_id"],
			UserID: sql.Null[uint]{
				V:     ids["user_id"],
				Valid: ids["user_id"] > 0,
			},
			TeamID: sql.Null[uint]{
				V:     ids["team_id"],
				Valid: ids["team_id"] > 0,
			},
			ContestID: sql.Null[uint]{
				V:     ids["contest_id"],
				Valid: ids["contest_id"] > 0,
			},
		}
		if err := json.Unmarshal([]byte(root.Data["network_plan"]), &victim.Spec.NetworkPlan); err != nil {
			return nil, fmt.Errorf("invalid network plan on %s: %w", root.Name, err)
		}
		seen[victim.ID] = true
		victims = append(victims, victim)
	}
	return victims, nil
}

// deleteVictimRoot deletes a victim root with foreground GC, so its whole tree
// completes before the next one. Sibling resources cannot enforce that ordering.
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
	if !hasVictimIdentity(root.Labels, victim, map[string]string{"cbctf.io/root": kind}) {
		return fmt.Errorf("root %s is not owned by victim %d", name, victim.ID)
	}
	uid := root.UID
	err = client.Delete(ctx, name, metav1.DeleteOptions{
		Preconditions:     &metav1.Preconditions{UID: &uid},
		PropagationPolicy: new(metav1.DeletePropagationForeground),
	})
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
	if !hasVictimIdentity(vpc.Labels, victim) {
		return fmt.Errorf("VPC %s is not owned by victim %d", name, victim.ID)
	}
	uid := vpc.UID
	if err = client.Delete(ctx, name, metav1.DeleteOptions{
		Preconditions:     &metav1.Preconditions{UID: &uid},
		PropagationPolicy: new(metav1.DeletePropagationForeground),
	}); err != nil && !apierror.IsNotFound(err) {
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
