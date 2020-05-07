package schema

import (
	"fmt"
)

const (
	Version = "v1alpha1"
	Group   = "edgeapi.cattle.io"
)

func SetAndGetCRDName(name string) string {
	return fmt.Sprintf("%s.%s/%s", name, Group, Version)
}
