package k8s

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api_core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	// Install special auth plugins like GCP Plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
)

// debug governs whether to print verbose logs or not
// It can be set by Environment Variable `CITF_VERBOSE_LOG``
var debug bool

func init() {
	debugEnv := os.Getenv("CITF_VERBOSE_LOG")

	if strings.ToLower(debugEnv) == "true" {
		debug = true
	} else {
		debug = false
	}
}

// K8S is a struct which will be the driver for all the methods related to kubernetes
type K8S struct {
	Config     *rest.Config
	Clientset  *kubernetes.Clientset
	RESTClient *rest.RESTClient
}

// NewK8S returns K8S struct
func NewK8S() K8S {
	config, err := GetClientConfig()
	Expect(err).NotTo(HaveOccurred())

	clientset, err := GetClientsetFromConfig(config)
	Expect(err).NotTo(HaveOccurred())

	restClient, err := rest.RESTClientFor(config)
	Expect(err).NotTo(HaveOccurred())

	return K8S{
		Config:     config,
		Clientset:  clientset,
		RESTClient: restClient,
	}
}

// Different phases of Namespace

// NsGoodPhases is an array of phases of the Namespace which are considered to be good
var NsGoodPhases = []api_core_v1.NamespacePhase{"Active"}

// Different phases of Pod

// PodWaitPhases is an array of phases of the Pod which are considered to be waiting
var PodWaitPhases = []string{"Pending"}

// PodGoodPhases is an array of phases of the Pod which are considered to be good
var PodGoodPhases = []string{"Running"}

// PodBadPhases is an array of phases of the Pod which are considered to be bad
var PodBadPhases = []string{"Error"}

// Different states of Pod

// PodWaitStates is an array of the states of the Pod which are considered to be waiting
var PodWaitStates = []string{"ContainerCreating", "Pending"}

// PodGoodStates is an array of the states of the Pod which are considered to be good
var PodGoodStates = []string{"Running"}

// PodBadStates is an array of the states of the Pod which are considered to be bad
var PodBadStates = []string{"CrashLoopBackOff", "ImagePullBackOff", "RunContainerError", "ContainerCannotRun"}

// GetClientConfig first tries to get a config object which uses the service account kubernetes gives to pods,
// if it is called from a process running in a kubernetes environment.
// Otherwise, it tries to build config from a default kubeconfig filepath if it fails, it fallback to the default config.
// Once it get the config, it returns the same.
func GetClientConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		if debug {
			fmt.Printf("Unable to create config. Error: %+v\n", err)
		}
		err1 := err
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			err = fmt.Errorf("InClusterConfig as well as BuildConfigFromFlags Failed. Error in InClusterConfig: %+v\nError in BuildConfigFromFlags: %+v", err1, err)
			return nil, err
		}
	}

	return config, nil
}

// GetClientsetFromConfig takes REST config and Create a clientset based on that and return that clientset
func GetClientsetFromConfig(config *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		err = fmt.Errorf("failed creating clientset. Error: %+v", err)
		return nil, err
	}

	return clientset, nil
}

// GetClientset first tries to get a config object which uses the service account kubernetes gives to pods,
// if it is called from a process running in a kubernetes environment.
// Otherwise, it tries to build config from a default kubeconfig filepath if it fails, it fallback to the default config.
// Once it get the config, it creates a new Clientset for the given config and returns the clientset.
func GetClientset() (*kubernetes.Clientset, error) {
	config, err := GetClientConfig()
	if err != nil {
		return nil, err
	}

	return GetClientsetFromConfig(config)
}

// GetRESTClient first tries to get a config object which uses the service account kubernetes gives to pods,
// if it is called from a process running in a kubernetes environment.
// Otherwise, it tries to build config from a default kubeconfig filepath if it fails, it fallback to the default config.
// Once it get the config, it
func GetRESTClient() (*rest.RESTClient, error) {
	config, err := GetClientConfig()
	if err != nil {
		return &rest.RESTClient{}, err
	}

	return rest.RESTClientFor(config)
}
