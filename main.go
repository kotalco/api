package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kotalco/api/api"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"os"
)

func main() {
	app := fiber.New()

	// register middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	api.MapUrl(app)

	// routing groups
	//api := app.Group("api")
	//v1 := api.Group("v1")
	//
	//

	//// register handlers
	//app.Get("/", func(c *fiber.Ctx) error {
	//	return c.SendString("Kotal API")
	//})
	//
	////chainlinkHandlers.NewNodeHandler().Register(chainlinkNodes)
	//polkadotHandlers.NewNodeHandler().Register(polkadotNodes)
	//ethereumHandlers.NewNodeHandler().Register(nodes)
	//handlers3.NewPeerHandler().Register(peers)
	//handlers3.NewClusterPeerHandler().Register(clusterpeers)
	//filecoinHandlers.NewNodeHandler().Register(filecoinNodes)
	//handlers.NewBeaconNodeHandler().Register(beaconnodes)
	//handlers.NewValidatorHandler().Register(validators)
	//handlers2.NewSecretHandler().Register(secrets)
	//handlers2.NewStorageClassHandler().Register(storageClasses)
	//nearHandlers.NewNodeHandler().Register(nearNodes)

	port := os.Getenv("KOTAL_API_SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		panic(err)
	}
}
