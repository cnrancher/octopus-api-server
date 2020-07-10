package authtoken

import (
	"context"
	"time"

	"github.com/cnrancher/octopus-api-server/pkg/auth"
	corev1 "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/ticker"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
)

const (
	name            = "auth-token-controller"
	refreshInterval = 3600
)

type Controller struct {
	secretController corev1.SecretClient
	secretCache      corev1.SecretCache
	apply            apply.Apply
}

func Register(ctx context.Context, apply apply.Apply, secrets corev1.SecretController) {
	controller := &Controller{
		secretController: secrets,
		secretCache:      secrets.Cache(),
		apply:            apply,
	}
	secrets.AddGenericHandler(ctx, name, controller.tokenSyncHandler)

	// call runRefreshSecretToken to remove expired secret token in every 60 minutes
	go runRefreshSecretToken(ctx, refreshInterval, secrets)
}

func (c *Controller) tokenSyncHandler(key string, obj runtime.Object) (runtime.Object, error) {
	if key == "" || obj == nil {
		return nil, nil
	}

	secret := toSecretObject(obj)
	if checkExpiredToken(secret) && hasSecretTokenLabel(secret.Labels) {
		logrus.Infof("remove expired secret token:%s", secret.Name)
		return obj, c.secretController.Delete(secret.Namespace, secret.Name, &metav1.DeleteOptions{})
	}
	return obj, nil
}

func runRefreshSecretToken(ctx context.Context, interval int, secretController corev1.SecretController) {
	r, _ := labels.NewRequirement(auth.OctopusAPILabel, selection.Equals, []string{"true"})
	labels := labels.NewSelector().Add(*r)
	for range ticker.Context(ctx, time.Duration(interval)*time.Second) {
		logrus.Info("Run refresh secretToken")
		secrets, err := secretController.Cache().List(auth.TokenNamespace, labels)
		if err != nil {
			logrus.Error(err)
		}

		for _, secret := range secrets {
			if checkExpiredToken(secret) {
				secretController.Enqueue(auth.TokenNamespace, secret.Name)
			}
		}
	}
}

func checkExpiredToken(secret *v1.Secret) bool {
	expiresAt := string(secret.Data["ExpiresAt"])

	if expiresAt == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		logrus.Errorf("failed to parse expiresAt date: %s", err.Error())
		return false
	}
	if time.Now().After(t) {
		return true
	}
	return false
}

func hasSecretTokenLabel(labels map[string]string) bool {
	for k := range labels {
		if k == auth.OctopusAPILabel {
			return true
		}
	}
	return false
}

func toSecretObject(obj runtime.Object) *v1.Secret {
	if obj != nil {
		if r, ok := obj.(*v1.Secret); ok {
			return r
		}
	}
	return nil
}
