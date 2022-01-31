package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/api/handlers"
	sharedHandlers "github.com/kotalco/api/handlers/shared"
	"github.com/kotalco/api/k8s"
	models "github.com/kotalco/api/models/near"
	"github.com/kotalco/api/shared"
	nearv1alpha1 "github.com/kotalco/kotal/apis/near/v1alpha1"
	sharedAPIs "github.com/kotalco/kotal/apis/shared"
	"github.com/ybbus/jsonrpc/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NodeHandler is NEAR node handler
type NodeHandler struct{}

// NewNodeHandler creates a new NEAR node handler
func NewNodeHandler() handlers.Handler {
	return &NodeHandler{}
}

// Get gets a single NEAR node by name
func (n *NodeHandler) Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*nearv1alpha1.Node)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"node": models.FromNEARNode(node),
	})
}

// List returns all NEAR nodes
func (n *NodeHandler) List(c *fiber.Ctx) error {
	nodes := &nearv1alpha1.NodeList{}
	if err := k8s.Client().List(c.Context(), nodes, client.InNamespace("default")); err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get all nodes",
		})
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

	nodeModels := []models.Node{}

	page, _ := strconv.Atoi(c.Query("page"))

	start, end := shared.Page(uint(len(nodes.Items)), uint(page))
	sort.Slice(nodes.Items[:], func(i, j int) bool {
		return nodes.Items[j].CreationTimestamp.Before(&nodes.Items[i].CreationTimestamp)
	})

	for _, node := range nodes.Items[start:end] {
		nodeModels = append(nodeModels, *models.FromNEARNode(&node))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"nodes": nodeModels,
	})

}

// Create creates NEAR node from spec
func (n *NodeHandler) Create(c *fiber.Ctx) error {
	model := new(models.Node)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	node := &nearv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      model.Name,
			Namespace: "default",
		},
		Spec: nearv1alpha1.NodeSpec{
			Network: model.Network,
			Archive: model.Archive,
			Resources: sharedAPIs.Resources{
				StorageClass: model.StorageClass,
			},
		},
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8s.Client().Create(c.Context(), node); err != nil {
		log.Println(err)

		if errors.IsAlreadyExists(err) {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": fmt.Sprintf("node by name %s already exist", model.Name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create node",
		})

	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"node": models.FromNEARNode(node),
	})
}

// Delete deletes NEAR node by name
func (n *NodeHandler) Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*nearv1alpha1.Node)

	if err := k8s.Client().Delete(c.Context(), node); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't delete node by name %s", c.Params("name")),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates NEAR node by name from spec
func (n *NodeHandler) Update(c *fiber.Ctx) error {
	model := new(models.Node)
	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	name := c.Params("name")
	node := c.Locals("node").(*nearv1alpha1.Node)

	if model.NodePrivateKeySecretName != "" {
		node.Spec.NodePrivateKeySecretName = model.NodePrivateKeySecretName
	}

	if model.ValidatorSecretName != "" {
		node.Spec.ValidatorSecretName = model.ValidatorSecretName
	}

	if model.MinPeers != 0 {
		node.Spec.MinPeers = model.MinPeers
	}

	if model.P2PPort != 0 {
		node.Spec.P2PPort = model.P2PPort
	}

	if model.P2PHost != "" {
		node.Spec.P2PHost = model.P2PHost
	}

	if model.RPC != nil {
		node.Spec.RPC = *model.RPC
	}
	if node.Spec.RPC {
		if model.RPCPort != 0 {
			node.Spec.RPCPort = model.RPCPort
		}
		if model.RPCHost != "" {
			node.Spec.RPCHost = model.RPCHost
		}
	}

	if model.PrometheusPort != 0 {
		node.Spec.PrometheusPort = model.PrometheusPort
	}

	if model.PrometheusHost != "" {
		node.Spec.PrometheusHost = model.PrometheusHost
	}

	if model.TelemetryURL != "" {
		node.Spec.TelemetryURL = model.TelemetryURL
	}

	if bootnodes := model.Bootnodes; bootnodes != nil {
		node.Spec.Bootnodes = *bootnodes
	}

	if model.CPU != "" {
		node.Spec.CPU = model.CPU
	}
	if model.CPULimit != "" {
		node.Spec.CPULimit = model.CPULimit
	}
	if model.Memory != "" {
		node.Spec.Memory = model.Memory
	}
	if model.MemoryLimit != "" {
		node.Spec.MemoryLimit = model.MemoryLimit
	}
	if model.Storage != "" {
		node.Spec.Storage = model.Storage
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8s.Client().Update(c.Context(), node); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't update node by name %s", name),
		})
	}

	updatedModel := models.FromNEARNode(node)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"node": updatedModel,
	})
}

// Count returns total number of nodes
func (n *NodeHandler) Count(c *fiber.Ctx) error {
	nodes := &nearv1alpha1.NodeList{}
	if err := k8s.Client().List(c.Context(), nodes, client.InNamespace("default")); err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

	return c.SendStatus(http.StatusOK)
}

func (e *NodeHandler) Stats(c *websocket.Conn) {
	defer c.Close()

	type Result struct {
		Error string `json:"error,omitempty"`
		// network_info call
		ActivePeersCount       uint `json:"activePeersCount,omitempty"`
		MaxPeersCount          uint `json:"maxPeersCount,omitempty"`
		SentBytesPerSecond     uint `json:"sentBytesPerSecond,omitempty"`
		ReceivedBytesPerSecond uint `json:"receivedBytesPerSecond,omitempty"`
		// status call
		LatestBlockHeight   uint `json:"latestBlockHeight,omitempty"`
		EarliestBlockHeight uint `json:"earliestBlockHeight,omitempty"`
		Syncing             bool `json:"syncing,omitempty"`
	}

	// Mock serever
	if os.Getenv("MOCK") == "true" {
		var activePeersCount, sentBytesPerSecond, receivedBytesPerSecond, latestBlockHeight, earliestBlockHeight uint
		for {
			activePeersCount++
			sentBytesPerSecond += 100
			receivedBytesPerSecond += 100
			latestBlockHeight += 36
			earliestBlockHeight += 3

			r := &Result{
				ActivePeersCount:       activePeersCount,
				MaxPeersCount:          40,
				SentBytesPerSecond:     sentBytesPerSecond,
				ReceivedBytesPerSecond: receivedBytesPerSecond,
				LatestBlockHeight:      latestBlockHeight,
				EarliestBlockHeight:    earliestBlockHeight,
				Syncing:                true,
			}

			var msg []byte

			if activePeersCount > 40 {
				activePeersCount = 10
				r = &Result{
					Error: "rpc is not enabled",
				}
			}

			msg, _ = json.Marshal(r)
			c.WriteMessage(websocket.TextMessage, []byte(msg))
			time.Sleep(time.Second)
		}
	}

	name := c.Params("name")
	node := &nearv1alpha1.Node{}
	key := types.NamespacedName{
		Namespace: "default",
		Name:      name,
	}

	for {

		err := k8s.Client().Get(context.Background(), key, node)
		if errors.IsNotFound(err) {
			c.WriteJSON(fiber.Map{
				"error": fmt.Sprintf("node by name %s doesn't exist", name),
			})
			return
		}

		if !node.Spec.RPC {
			c.WriteJSON(fiber.Map{
				"error": "rpc is not enabled",
			})
			return
		}

		client := jsonrpc.NewClient(fmt.Sprintf("http://%s:%d", node.Name, node.Spec.RPCPort))

		type NodeStatus struct {
			SyncInfo struct {
				LatestBlockHeight   uint `json:"latest_block_height"`
				EarliestBlockHeight uint `json:"earliest_block_height"`
				Syncing             bool `json:"syncing"`
			} `json:"sync_info"`
		}

		// node status rpc call
		nodeStatus := &NodeStatus{}
		err = client.CallFor(nodeStatus, "status")
		if err != nil {
			fmt.Println(err)
		}

		type NetworkInfo struct {
			ActivePeersCount       uint `json:"num_active_peers"`
			MaxPeersCount          uint `json:"peer_max_count"`
			SentBytesPerSecond     uint `json:"sent_bytes_per_sec"`
			ReceivedBytesPerSecond uint `json:"received_bytes_per_sec"`
		}

		// network info rpc call
		networkInfo := &NetworkInfo{}
		err = client.CallFor(networkInfo, "network_info")
		if err != nil {
			fmt.Println(err)
		}

		c.WriteJSON(fiber.Map{
			"activePeersCount":       networkInfo.ActivePeersCount,
			"maxPeersCount":          networkInfo.MaxPeersCount,
			"sentBytesPerSecond":     networkInfo.SentBytesPerSecond,
			"receivedBytesPerSecond": networkInfo.ReceivedBytesPerSecond,
			"latestBlockHeight":      nodeStatus.SyncInfo.LatestBlockHeight,
			"earliestBlockHeight":    nodeStatus.SyncInfo.EarliestBlockHeight,
			"syncing":                nodeStatus.SyncInfo.Syncing,
		})

		time.Sleep(time.Second)
	}
}

// validateNodeExist validates NEAR node by name exist
func validateNodeExist(c *fiber.Ctx) error {
	name := c.Params("name")
	node := &nearv1alpha1.Node{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(c.Context(), key, node); err != nil {

		if errors.IsNotFound(err) {
			return c.Status(http.StatusNotFound).JSON(map[string]string{
				"error": fmt.Sprintf("node by name %s doesn't exist", name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't get node by name %s", name),
		})
	}

	c.Locals("node", node)

	return c.Next()
}

// Register registers all handlers on the given router
func (n *NodeHandler) Register(router fiber.Router) {
	router.Post("/", n.Create)
	router.Head("/", n.Count)
	router.Get("/", n.List)
	router.Get("/:name", validateNodeExist, n.Get)
	router.Get("/:name/logs", websocket.New(sharedHandlers.Logger))
	router.Get("/:name/status", websocket.New(sharedHandlers.Status))
	router.Get("/:name/stats", websocket.New(n.Stats))
	router.Put("/:name", validateNodeExist, n.Update)
	router.Delete("/:name", validateNodeExist, n.Delete)
}
