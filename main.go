package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kotalco/api/handlers"
	ethereumHandlers "github.com/kotalco/api/handlers/ethereum"
	ethereum2Handlers "github.com/kotalco/api/handlers/ethereum2"
	ipfsHandlers "github.com/kotalco/api/handlers/ipfs"
)

func main() {
	app := fiber.New()

	// register middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// routing groups
	api := app.Group("api")
	v1 := api.Group("v1")

	ethereum := v1.Group("ethereum")
	nodes := ethereum.Group("nodes")

	ipfs := v1.Group("ipfs")
	peers := ipfs.Group("peers")

	ethereum2 := v1.Group("ethereum2")
	beaconnodes := ethereum2.Group("beaconnodes")

	// register handlers
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Kotal API")
	})

	var nodeHandler, peerHandler, beaconHandler handlers.Handler

	if os.Getenv("MOCK") == "true" {
		nodeHandler = ethereumHandlers.NewNodeMockHandler()
		peerHandler = ipfsHandlers.NewPeerMockHandler()
		beaconHandler = ethereum2Handlers.NewBeaconNodeMockHandler()
	} else {
		nodeHandler = ethereumHandlers.NewNodeHandler()
		peerHandler = ipfsHandlers.NewPeerHandler()
		beaconHandler = ethereum2Handlers.NewBeaconNodeHandler()
	}

	nodeHandler.Register(nodes)
	peerHandler.Register(peers)
	beaconHandler.Register(beaconnodes)

	app.Listen(":3000")
}
