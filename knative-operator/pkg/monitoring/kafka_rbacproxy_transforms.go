package monitoring

import (
	"context"
	"errors"

	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/util/sets"
	operatorv1alpha1 "knative.dev/operator/pkg/apis/operator/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	serverlessoperatorv1alpha1 "github.com/openshift-knative/serverless-operator/knative-operator/pkg/apis/operator/v1alpha1"
	"github.com/openshift-knative/serverless-operator/openshift-knative-operator/pkg/monitoring"
)

var (
	KafkaController = Component{
		Name:               "kafka-controller",
		ServiceAccountName: "kafka-controller",
	}
	KafkaWebhook = Component{
		Name:               "kafka-webhook-eventing",
		ServiceAccountName: "kafka-webhook-eventing",
	}

	KafkaBrokerReceiver = Component{
		Name:               "kafka-broker-receiver",
		ServiceAccountName: "knative-kafka-broker-data-plane",
	}
	KafkaBrokerDispatcher = Component{
		Name:               "kafka-broker-dispatcher",
		ServiceAccountName: "knative-kafka-broker-data-plane",
	}

	KafkaSinkReceiver = Component{
		Name:               "kafka-sink-receiver",
		ServiceAccountName: "knative-kafka-sink-data-plane",
	}

	KafkaChannelController = Component{
		Name:               "kafka-ch-controller",
		ServiceAccountName: "kafka-ch-controller",
	}
	KafkaChannelWebhook = Component{
		Name:               "kafka-webhook",
		ServiceAccountName: "kafka-webhook",
	}

	KafkaSourceController = Component{
		Name:               "kafka-controller-manager",
		ServiceAccountName: "kafka-controller-manager",
	}

	deployments = []string{
		KafkaController.Name,
		KafkaWebhook.Name,
		KafkaBrokerReceiver.Name,
		KafkaBrokerDispatcher.Name,
		KafkaSinkReceiver.Name,
		KafkaChannelController.Name,
		KafkaChannelWebhook.Name,
		KafkaSourceController.Name,
	}
)

// Component is a target for scraping metrics.
type Component struct {
	Name               string
	ServiceAccountName string
}

func AddRBACProxyToManifest(instance *serverlessoperatorv1alpha1.KnativeKafka, components ...Component) (*mf.Manifest, error) {
	proxyManifest := mf.Manifest{}
	// Only create the roles needed for the deployment service accounts as Prometheus has already
	// the rights needed due to eventing that is assumed to be installed.
	for _, c := range components {
		crbM, err := monitoring.CreateClusterRoleBindingManifest(c.ServiceAccountName, instance.GetNamespace())
		if err != nil {
			return nil, err
		}
		proxyManifest = proxyManifest.Append(*crbM)
		if err = monitoring.AppendManifestsForComponent(c.Name, instance.GetNamespace(), &proxyManifest); err != nil {
			return nil, err
		}
	}
	return &proxyManifest, nil
}

func GetRBACProxyInjectTransformer(apiClient client.Client) (mf.Transformer, error) {
	eventingList := &operatorv1alpha1.KnativeEventingList{}
	err := apiClient.List(context.Background(), eventingList)
	if err != nil {
		return nil, err
	}
	if len(eventingList.Items) == 0 {
		return nil, errors.New("eventing instance not found")
	}
	if monitoring.ShouldEnableMonitoring(eventingList.Items[0].GetSpec().GetConfig()) {
		return monitoring.InjectRbacProxyContainerToDeployments(sets.NewString(deployments...)), nil
	}
	return nil, nil
}
