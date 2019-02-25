package kube

import (
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubeContexts(path string) ([]string, error) {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: path}, &clientcmd.ConfigOverrides{}).RawConfig()
	if err != nil {
		return nil, err
	}
	items := []string{}
	for name, _ := range config.Contexts {
		items = append(items, name)
	}
	return items, nil
}
