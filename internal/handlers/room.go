package handlers

import (
	"github.com/gofiber/fiber/v2"
	websocket "github.com/gofiber/websocket/v2"
	guuid "github.com/google/uuid"
	"fmt"
	"os"
	"time"
)

func RoomCreate(c *fiber.Ctx) error {
	return c.Redirect(fmt.Sprintf("/room/%s", guuid.New().String()))
}

func RoomWebsocket(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}
	_, _, room := createOrGetRoom(uuid)
	
}

func Room(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		c.Status(400)
		return nil
	}

	uuid, suuid, _ := createPrGetRoom(uuid)
}

func createOrGetRoom(uuid string) (string, string, Room) {
	
}