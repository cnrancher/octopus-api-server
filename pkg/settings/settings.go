package settings

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	settings       = map[string]Setting{}
	provider       Provider
	InjectDefaults string

	ServerURL        = NewSetting("server-url", "")
	Location         = NewSetting("location", "")
	SystemNamespaces = NewSetting("system-namespaces", "kube-system,kube-public,octopus-system")
	UIIndex          = NewSetting("ui-index", "https://rancher-octopus.s3-accelerate.amazonaws.com/ui/latest/index.html")
	UIPath           = NewSetting("ui-path", "")
	APIUIVersion     = NewSetting("api-ui-version", "1.1.9") // Please update the OCTOPUS_API_UI_VERSION in package/Dockerfile when updating the version here.
)

func init() {
	if InjectDefaults == "" {
		return
	}
	defaults := map[string]string{}
	if err := json.Unmarshal([]byte(InjectDefaults), &defaults); err != nil {
		return
	}
	for name, defaultValue := range defaults {
		value, ok := settings[name]
		if !ok {
			continue
		}
		value.Default = defaultValue
		settings[name] = value
	}
}

type Provider interface {
	Get(name string) string
	Set(name, value string) error
	SetIfUnset(name, value string) error
	SetAll(settings map[string]Setting) error
}

type Setting struct {
	Name     string
	Default  string
	ReadOnly bool
}

func (s Setting) SetIfUnset(value string) error {
	if provider == nil {
		return s.Set(value)
	}
	return provider.SetIfUnset(s.Name, value)
}

func (s Setting) Set(value string) error {
	if provider == nil {
		s, ok := settings[s.Name]
		if ok {
			s.Default = value
			settings[s.Name] = s
		}
	} else {
		return provider.Set(s.Name, value)
	}
	return nil
}

func (s Setting) Get() string {
	if provider == nil {
		s := settings[s.Name]
		return s.Default
	}
	return provider.Get(s.Name)
}

func (s Setting) GetInt() int {
	v := s.Get()
	i, err := strconv.Atoi(v)
	if err == nil {
		return i
	}
	logrus.Errorf("failed to parse setting %s=%s as int: %v", s.Name, v, err)
	i, err = strconv.Atoi(s.Default)
	if err != nil {
		return 0
	}
	return i
}

func SetProvider(p Provider) error {
	if err := p.SetAll(settings); err != nil {
		return err
	}
	provider = p
	return nil
}

func NewSetting(name, def string) Setting {
	s := Setting{
		Name:    name,
		Default: def,
	}
	settings[s.Name] = s
	return s
}

func GetEnvKey(key string) string {
	return "OCTOPUS_" + strings.ToUpper(strings.Replace(key, "-", "_", -1))
}
