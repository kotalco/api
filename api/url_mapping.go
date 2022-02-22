package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/api/api/handlers/chainlink"
	"github.com/kotalco/api/api/handlers/core/secret"
	"github.com/kotalco/api/api/handlers/ethereum"
	"github.com/kotalco/api/api/handlers/ethereum2/beacon_node"
	"github.com/kotalco/api/api/handlers/shared"
)

// MapUrl abstracted function to map and register all the url for the application
// Used in the main.go
// helps to keep all the endpoints' definition in one place
// the only place to interact with handlers and middlewares
func MapUrl(app *fiber.App) {
	// routing groups
	api := app.Group("api")
	v1 := api.Group("v1")

	// chainlink group
	chainlinkGroup := v1.Group("chainlink")
	chainlinkNodes := chainlinkGroup.Group("nodes")
	chainlinkNodes.Post("/", chainlink.Create)
	chainlinkNodes.Head("/", chainlink.Count)
	chainlinkNodes.Get("/", chainlink.List)
	chainlinkNodes.Get("/:name", chainlink.ValidateNodeExist, chainlink.Get)
	chainlinkNodes.Get("/:name/logs", websocket.New(shared.Logger))
	chainlinkNodes.Get("/:name/status", websocket.New(shared.Status))
	chainlinkNodes.Put("/:name", chainlink.ValidateNodeExist, chainlink.Update)
	chainlinkNodes.Delete("/:name", chainlink.ValidateNodeExist, chainlink.Delete)

	//ethereum group
	ethereumGroup := v1.Group("ethereum")
	ethereumNodes := ethereumGroup.Group("nodes")
	ethereumNodes.Post("/", ethereum.Create)
	ethereumNodes.Head("/", ethereum.Count)
	ethereumNodes.Get("/", ethereum.List)
	ethereumNodes.Get("/:name", ethereum.ValidateNodeExist, ethereum.Get)
	ethereumNodes.Get("/:name/logs", websocket.New(shared.Logger))
	ethereumNodes.Get("/:name/status", websocket.New(shared.Status))
	ethereumNodes.Get("/:name/stats", websocket.New(ethereum.Stats))
	ethereumNodes.Put("/:name", ethereum.ValidateNodeExist, ethereum.Update)
	ethereumNodes.Delete("/:name", ethereum.ValidateNodeExist, ethereum.Delete)

	//core group
	coreGroup := v1.Group("core")
	//secret group
	secrets := coreGroup.Group("secrets")
	secrets.Post("/", secret.Create)
	secrets.Head("/", secret.Count)
	secrets.Get("/", secret.List)
	secrets.Get("/:name", secret.ValidateSecretExist, secret.Get)
	secrets.Put("/:name", secret.ValidateSecretExist, secret.Update)
	secrets.Delete("/:name", secret.ValidateSecretExist, secret.Delete)

	//ethereum2 group
	ethereum2 := v1.Group("ethereum2")
	beaconnodes := ethereum2.Group("beaconnodes")
	beaconnodes.Post("/", beacon_node.Create)
	beaconnodes.Head("/", beacon_node.Count)
	beaconnodes.Get("/", beacon_node.List)
	beaconnodes.Get("/:name", beacon_node.ValidateBeaconNodeExist, beacon_node.Get)
	beaconnodes.Get("/:name/logs", websocket.New(shared.Logger))
	beaconnodes.Get("/:name/status", websocket.New(shared.Status))
	beaconnodes.Put("/:name", beacon_node.ValidateBeaconNodeExist, beacon_node.Update)
	beaconnodes.Delete("/:name", beacon_node.ValidateBeaconNodeExist, beacon_node.Delete)

}
