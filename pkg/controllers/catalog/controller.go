package catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/cnrancher/octopus-api-server/pkg/apis/octopusapi.cattle.io/v1alpha1"
	catalogcontroller "github.com/cnrancher/octopus-api-server/pkg/generated/controllers/octopusapi.cattle.io/v1alpha1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Name               = "catalog-controller"
	Namespace          = "kube-system"
	OctopusCatalogName = "octopus-catalog"
	OctopusCatalogURL  = "http://charts.cnrancher.com/octopus-catalog"
)

type Controller struct {
	catalogController catalogcontroller.CatalogController
	catalogCache      catalogcontroller.CatalogCache
	apply             apply.Apply
}

func Register(ctx context.Context, apply apply.Apply, catalogs catalogcontroller.CatalogController) {
	controller := &Controller{
		catalogController: catalogs,
		catalogCache:      catalogs.Cache(),
		apply:             apply,
	}
	catalogs.OnChange(ctx, Name, controller.OnCatalogChanged)
	catalogs.OnRemove(ctx, Name, controller.OnCatalogRemoved)
	if err := addMQTTCatalog(catalogs, OctopusCatalogName, OctopusCatalogURL); err != nil {
		logrus.Error(err)
	}
}

func addMQTTCatalog(controller catalogcontroller.CatalogController, name, url string) error {
	_, err := controller.Get(Namespace, name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		obj := &v1alpha1.Catalog{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: Namespace,
			},
			Spec: v1alpha1.CatalogSpec{
				URL: url,
			},
		}
		if _, err := controller.Create(obj); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) OnCatalogChanged(_ string, index *v1alpha1.Catalog) (*v1alpha1.Catalog, error) {
	if index == nil {
		return nil, nil
	}
	catalog := &Helm{
		catalogName: index.Name,
		url:         index.Spec.URL,
		username:    index.Spec.Username,
		password:    index.Spec.Password,
	}
	indexCopy := index.DeepCopy()
	v1alpha1.CatalogConditionCreated.True(indexCopy)
	file, err := catalog.downloadIndex(index.Spec.URL)
	if err != nil {
		setCatalogErrorState(indexCopy, err)
		return c.catalogController.Update(indexCopy)
	}
	indexCopy.Spec.IndexFile = file
	setCatalogRefreshed(indexCopy)
	return c.catalogController.Update(indexCopy)
}

func (c *Controller) OnCatalogRemoved(key string, index *v1alpha1.Catalog) (*v1alpha1.Catalog, error) {
	if key == "" {
		return index, nil
	}
	return index, c.catalogController.Delete(index.Namespace, index.Name, &metav1.DeleteOptions{})
}

func setCatalogErrorState(indexCopy *v1alpha1.Catalog, err error) {
	v1alpha1.CatalogConditionProcessed.True(indexCopy)
	v1alpha1.CatalogConditionRefreshed.SetError(indexCopy, fmt.Sprintf("Error syncing catalog %v", indexCopy.Name), err)
}

func setCatalogRefreshed(index *v1alpha1.Catalog) {
	index.Status.LastRefreshTimestamp = time.Now().Format(time.RFC3339)
	v1alpha1.CatalogConditionCreated.True(index)
	v1alpha1.CatalogConditionProcessed.True(index)
	v1alpha1.CatalogConditionRefreshed.True(index)
	v1alpha1.CatalogConditionRefreshed.Message(index, "")
	v1alpha1.CatalogConditionRefreshed.Reason(index, "")
}
