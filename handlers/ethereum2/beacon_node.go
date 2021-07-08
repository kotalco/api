package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	"github.com/kotalco/api/k8s"
	models "github.com/kotalco/api/models/ethereum2"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// BeaconNodeHandler is ethereum 2.0 beacon node handler
type BeaconNodeHandler struct{}

// NewBeaconNodeHandler creates a new ethereum 2.0 beacon node handler
func NewBeaconNodeHandler() handlers.Handler {
	return &BeaconNodeHandler{}
}

// Get gets a single ethereum 2.0 beacon node by name
func (p *BeaconNodeHandler) Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*ethereum2v1alpha1.BeaconNode)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"beaconnode": models.FromEthereum2BeaconNode(node),
	})
}

// List returns all ethereum 2.0 beacon nodes
func (p *BeaconNodeHandler) List(c *fiber.Ctx) error {
	nodes := &ethereum2v1alpha1.BeaconNodeList{}
	if err := k8s.Client().List(c.Context(), nodes); err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get all beacon nodes",
		})
	}

	nodeModels := []models.BeaconNode{}
	for _, node := range nodes.Items {
		nodeModels = append(nodeModels, *models.FromEthereum2BeaconNode(&node))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"beaconnodes": nodeModels,
	})
}

// Create creates ethereum 2.0 beacon node from spec
func (p *BeaconNodeHandler) Create(c *fiber.Ctx) error {
	model := new(models.BeaconNode)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	if model.Eth1Endpoints == nil {
		model.Eth1Endpoints = []string{}
	}

	client := ethereum2v1alpha1.Ethereum2Client(model.Client)

	beaconnode := &ethereum2v1alpha1.BeaconNode{
		ObjectMeta: metav1.ObjectMeta{
			Name:      model.Name,
			Namespace: "default",
		},
		Spec: ethereum2v1alpha1.BeaconNodeSpec{
			Join:          model.Network,
			Client:        client,
			Eth1Endpoints: model.Eth1Endpoints,
			RPC:           client == ethereum2v1alpha1.PrysmClient,
		},
	}

	if os.Getenv("MOCK") == "true" {
		beaconnode.Default()
	}

	if err := k8s.Client().Create(c.Context(), beaconnode); err != nil {
		log.Println(err)
		if errors.IsAlreadyExists(err) {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": fmt.Sprintf("beacon node by name %s already exist", model.Name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create beacon node",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"beaconnode": models.FromEthereum2BeaconNode(beaconnode),
	})
}

// Delete deletes ethereum 2.0 beacon node by name
func (p *BeaconNodeHandler) Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*ethereum2v1alpha1.BeaconNode)

	if err := k8s.Client().Delete(c.Context(), node); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't delete beacon node by name %s", c.Params("name")),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates ethereum 2.0 beacon node by name from spec
func (p *BeaconNodeHandler) Update(c *fiber.Ctx) error {
	model := new(models.BeaconNode)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	beaconnode := c.Locals("node").(*ethereum2v1alpha1.BeaconNode)

	if len(model.Eth1Endpoints) != 0 {
		beaconnode.Spec.Eth1Endpoints = model.Eth1Endpoints
	}

	if model.REST != nil {
		rest := *model.REST
		if rest {
			if model.RESTHost != "" {
				beaconnode.Spec.RESTHost = model.RESTHost
			}
			if model.RESTPort != 0 {
				beaconnode.Spec.RESTPort = model.RESTPort
			}
		}
		beaconnode.Spec.REST = rest
	}

	if model.RPC != nil {
		rpc := *model.RPC
		if rpc {
			if model.RPCHost != "" {
				beaconnode.Spec.RPCHost = model.RPCHost
			}
			if model.RPCPort != 0 {
				beaconnode.Spec.RPCPort = model.RPCPort
			}
		}
		beaconnode.Spec.RPC = rpc
	}

	if model.GRPC != nil {
		grpc := *model.GRPC
		if grpc {
			if model.GRPCHost != "" {
				beaconnode.Spec.GRPCHost = model.GRPCHost
			}
			if model.GRPCPort != 0 {
				beaconnode.Spec.GRPCPort = model.GRPCPort
			}
		}
		beaconnode.Spec.GRPC = grpc
	}

	if model.CPU != "" {
		beaconnode.Spec.CPU = model.CPU
	}
	if model.CPULimit != "" {
		beaconnode.Spec.CPULimit = model.CPULimit
	}
	if model.Memory != "" {
		beaconnode.Spec.Memory = model.Memory
	}
	if model.MemoryLimit != "" {
		beaconnode.Spec.MemoryLimit = model.MemoryLimit
	}
	if model.Storage != "" {
		beaconnode.Spec.Storage = model.Storage
	}

	if os.Getenv("MOCK") == "true" {
		beaconnode.Default()
	}

	if err := k8s.Client().Update(c.Context(), beaconnode); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't update beacon node by name %s", c.Params("name")),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"beaconnode": models.FromEthereum2BeaconNode(beaconnode),
	})
}

// validateBeaconNodeExist validate node by name exist
func validateBeaconNodeExist(c *fiber.Ctx) error {
	name := c.Params("name")
	node := &ethereum2v1alpha1.BeaconNode{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(c.Context(), key, node); err != nil {

		if errors.IsNotFound(err) {
			return c.Status(http.StatusNotFound).JSON(map[string]string{
				"error": fmt.Sprintf("beacon node by name %s doesn't exist", name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't get beacon node by name %s", name),
		})
	}

	c.Locals("node", node)

	return c.Next()

}

// Register registers all handlers on the given router
func (p *BeaconNodeHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", validateBeaconNodeExist, p.Get)
	router.Put("/:name", validateBeaconNodeExist, p.Update)
	router.Delete("/:name", validateBeaconNodeExist, p.Delete)
}
