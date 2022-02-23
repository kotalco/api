package filecoin

import (
	"context"
	"fmt"
	restErrors "github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/logger"
	filecoinv1alpha1 "github.com/kotalco/kotal/apis/filecoin/v1alpha1"
	sharedAPIs "github.com/kotalco/kotal/apis/shared"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type filecoinService struct{}

type filecoinServiceInterface interface {
	Get(name string) (*filecoinv1alpha1.Node, *restErrors.RestErr)
	Create(dto *FilecoinDto) (*filecoinv1alpha1.Node, *restErrors.RestErr)
	Update(*FilecoinDto, *filecoinv1alpha1.Node) (*filecoinv1alpha1.Node, *restErrors.RestErr)
	List() (*filecoinv1alpha1.NodeList, *restErrors.RestErr)
	Delete(node *filecoinv1alpha1.Node) *restErrors.RestErr
	Count() (*int, *restErrors.RestErr)
}

var (
	FilecoinService filecoinServiceInterface
)

func init() { FilecoinService = &filecoinService{} }

// Get gets a single filecoin node by name
func (service filecoinService) Get(name string) (*filecoinv1alpha1.Node, *restErrors.RestErr) {
	node := &filecoinv1alpha1.Node{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(context.Background(), key, node); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, restErrors.NewNotFoundError(fmt.Sprintf("node by name %s doesn't exit", name))
		}
		go logger.Error("ERROR_IN_GET_FILECOIN", err)
		return nil, restErrors.NewInternalServerError(fmt.Sprintf("can't get node by name %s", name))
	}

	return node, nil
}

// Create creates filecoin node from spec
func (service filecoinService) Create(dto *FilecoinDto) (*filecoinv1alpha1.Node, *restErrors.RestErr) {
	node := &filecoinv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dto.Name,
			Namespace: "default",
		},
		Spec: filecoinv1alpha1.NodeSpec{
			Network: filecoinv1alpha1.FilecoinNetwork(dto.Network),
			Resources: sharedAPIs.Resources{
				StorageClass: dto.StorageClass,
			},
		},
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8s.Client().Create(context.Background(), node); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			return nil, restErrors.NewBadRequestError(fmt.Sprintf("node by name %s already exits", dto))
		}
		go logger.Error("ERROR_IN_CREATE_FILECOIN", err)
		return nil, restErrors.NewInternalServerError("failed to create node")
	}

	return node, nil
}

// Update updates filecoin node by name from spec
func (service filecoinService) Update(dto *FilecoinDto, node *filecoinv1alpha1.Node) (*filecoinv1alpha1.Node, *restErrors.RestErr) {
	if dto.API != nil {
		node.Spec.API = *dto.API
	}

	if dto.APIPort != 0 {
		node.Spec.APIPort = dto.APIPort
	}

	if dto.APIHost != "" {
		node.Spec.APIHost = dto.APIHost
	}

	if dto.APIRequestTimeout != 0 {
		node.Spec.APIRequestTimeout = dto.APIRequestTimeout
	}

	if dto.DisableMetadataLog != nil {
		node.Spec.DisableMetadataLog = *dto.DisableMetadataLog
	}

	if dto.P2PPort != 0 {
		node.Spec.P2PPort = dto.P2PPort
	}

	if dto.P2PHost != "" {
		node.Spec.P2PHost = dto.P2PHost
	}

	if dto.IPFSPeerEndpoint != "" {
		node.Spec.IPFSPeerEndpoint = dto.IPFSPeerEndpoint
	}

	if dto.IPFSOnlineMode != nil {
		node.Spec.IPFSOnlineMode = *dto.IPFSOnlineMode
	}

	if dto.IPFSForRetrieval != nil {
		node.Spec.IPFSForRetrieval = *dto.IPFSForRetrieval
	}

	if dto.CPU != "" {
		node.Spec.CPU = dto.CPU
	}
	if dto.CPULimit != "" {
		node.Spec.CPULimit = dto.CPULimit
	}
	if dto.Memory != "" {
		node.Spec.Memory = dto.Memory
	}
	if dto.MemoryLimit != "" {
		node.Spec.MemoryLimit = dto.MemoryLimit
	}
	if dto.Storage != "" {
		node.Spec.Storage = dto.Storage
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8s.Client().Update(context.Background(), node); err != nil {
		go logger.Error("ERROR_IN_UPDATE_FILECOIN", err)
		return nil, restErrors.NewInternalServerError(fmt.Sprintf("can't update node by name %s", node.Name))
	}

	return node, nil
}

// List returns all filecoin nodes
func (service filecoinService) List() (*filecoinv1alpha1.NodeList, *restErrors.RestErr) {
	nodes := &filecoinv1alpha1.NodeList{}
	if err := k8s.Client().List(context.Background(), nodes, client.InNamespace("default")); err != nil {
		go logger.Error("ERROR_IN_LIST_FILECOIN", err)
		return nil, restErrors.NewInternalServerError("failed to get all nodes")
	}
	return nodes, nil
}

// Count returns total number of filecoin nodes
func (service filecoinService) Count() (*int, *restErrors.RestErr) {

	nodes := &filecoinv1alpha1.NodeList{}
	if err := k8s.Client().List(context.Background(), nodes, client.InNamespace("default")); err != nil {
		go logger.Error("ERROR_IN_COUNT_FILECOIN", err)
		return nil, restErrors.NewInternalServerError("failed to count filecoin nodes")
	}

	length := len(nodes.Items)
	return &length, nil
}

// Delete deletes ethereum 2.0 filecoin node by name
func (service filecoinService) Delete(node *filecoinv1alpha1.Node) *restErrors.RestErr {
	if err := k8s.Client().Delete(context.Background(), node); err != nil {
		go logger.Error("ERROR_IN_DELETE_FILECOIN", err)
		return restErrors.NewInternalServerError(fmt.Sprintf("can't delte node by name %s", node.Name))
	}
	return nil
}
