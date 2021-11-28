package models

import polkadotv1alpha1 "github.com/kotalco/kotal/apis/polkadot/v1alpha1"

type Node struct {
	Name                     string `json:"name"`
	Network                  string `json:"network"`
	NodePrivateKeySecretName string `json:"nodePrivateKeySecretName"`
}

func FromPolkadotNode(node *polkadotv1alpha1.Node) *Node {
	return &Node{
		Name:                     node.Name,
		Network:                  node.Spec.Network,
		NodePrivateKeySecretName: node.Spec.NodePrivateKeySecretName,
	}
}