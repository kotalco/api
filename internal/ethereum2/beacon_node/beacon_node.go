package beacon_node

import (
	"github.com/kotalco/api/internal/models"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/shared"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
)

type BeaconNodeDto struct {
	models.Time
	k8s.MetaDataDto
	Network string `json:"network"`
	Client  string `json:"client"`
	// todo: required only for prysm and network is not mainnet
	Eth1Endpoints *[]string `json:"eth1Endpoints"`
	REST          *bool     `json:"rest"`
	RESTHost      string    `json:"restHost"`
	RESTPort      uint      `json:"restPort"`
	RPC           *bool     `json:"rpc"`
	RPCHost       string    `json:"rpcHost"`
	RPCPort       uint      `json:"rpcPort"`
	GRPC          *bool     `json:"grpc"`
	GRPCHost      string    `json:"grpcHost"`
	GRPCPort      uint      `json:"grpcPort"`
	CPU           string    `json:"cpu"`
	CPULimit      string    `json:"cpuLimit"`
	Memory        string    `json:"memory"`
	MemoryLimit   string    `json:"memoryLimit"`
	Storage       string    `json:"storage"`
	StorageClass  *string   `json:"storageClass"`
}
type BeaconNodeListDto []BeaconNodeDto

func (dto BeaconNodeDto) FromEthereum2BeaconNode(node *ethereum2v1alpha1.BeaconNode) *BeaconNodeDto {
	dto.Name = node.Name
	dto.Time = models.Time{CreatedAt: node.CreationTimestamp.UTC().Format(shared.JavascriptISOString)}
	dto.Network = node.Spec.Network
	dto.Client = string(node.Spec.Client)
	dto.Eth1Endpoints = &node.Spec.Eth1Endpoints
	dto.REST = &node.Spec.REST
	dto.RESTHost = node.Spec.RESTHost
	dto.RESTPort = node.Spec.RESTPort
	dto.RPC = &node.Spec.RPC
	dto.RPCHost = node.Spec.RPCHost
	dto.RPCPort = node.Spec.RPCPort
	dto.GRPC = &node.Spec.GRPC
	dto.GRPCHost = node.Spec.GRPCHost
	dto.GRPCPort = node.Spec.GRPCPort
	dto.CPU = node.Spec.CPU
	dto.CPULimit = node.Spec.CPULimit
	dto.Memory = node.Spec.Memory
	dto.MemoryLimit = node.Spec.MemoryLimit
	dto.Storage = node.Spec.Storage
	dto.StorageClass = node.Spec.StorageClass

	return &dto
}

func (nodes BeaconNodeListDto) FromEthereum2BeaconNode(beaconnodeList []ethereum2v1alpha1.BeaconNode) BeaconNodeListDto {
	result := make(BeaconNodeListDto, len(beaconnodeList))
	for index, v := range beaconnodeList {
		result[index] = *(BeaconNodeDto{}.FromEthereum2BeaconNode(&v))
	}
	return result
}
