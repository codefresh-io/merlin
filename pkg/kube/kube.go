package kube

import (
	"fmt"

	"github.com/codefresh-io/merlin/pkg/config"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeConfig "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type (
	Kubernetes interface {
		EnsureNamespaceNotExist(string) error
	}

	k struct {
		client *kubeConfig.Clientset
		log    *logrus.Entry
	}
)

func New(cnf *config.Config, log *logrus.Entry) (Kubernetes, error) {

	config := clientcmd.GetConfigFromFileOrDie(cnf.Kube.Path)

	override := &clientcmd.ConfigOverrides{
		ClusterInfo: api.Cluster{
			Server: "",
		},
	}

	c := clientcmd.NewNonInteractiveClientConfig(*config, cnf.Kube.Context, override, nil)
	clientCnf, err := c.ClientConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubeConfig.NewForConfig(clientCnf)
	if err != nil {
		return nil, err
	}
	log.Debug("Create kubernetes config")
	return &k{
		client: clientset,
		log:    log,
	}, err
}

func (k *k) EnsureNamespaceNotExist(name string) error {
	var err error
	k.log.WithField("namespace", name).Debug("Making sure the namespace is not exist yet")
	res, err := k.client.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil {
		if statusError, errIsStatusError := err.(*errors.StatusError); errIsStatusError {
			if statusError.ErrStatus.Reason == metav1.StatusReasonNotFound {
				return nil
			}
		}
		return err
	}
	if res != nil {
		return fmt.Errorf("Unknown error")
	}
	return nil
}
