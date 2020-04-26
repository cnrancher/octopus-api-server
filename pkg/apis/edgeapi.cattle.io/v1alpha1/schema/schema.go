package schema

import (
	"fmt"
)

const (
	version = "v1alpha1"
	group   = "edgeapi.cattle.io"
)

func SetAndGetCRDName(name string) string {
	return fmt.Sprintf("%s.%s/%s", name, group, version)
}
