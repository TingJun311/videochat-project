package handlers

import (
	"fmt"
	"os"
	"time"
	"videochat-project/pkg/webrtc"
	w "videochat-project/pkg/webrtc"

	"github.com/gofiber/fiber/v2"
	websocket "github.com/gofiber/websocket/v2"
	guuid "github.com/google/uuid"
)

type websocketMessage struct {
	Event string `json:"event"`
	data string `json:"data"`
}

func RoomCreate(c *fiber.Ctx) error {
	return c.Redirect(fmt.Sprintf("/room/%s", guuid.New().String()))
}

func RoomWebsocket(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}
	_, _, room := createOrGetRoom(uuid)
	w.RoomConn(c, room.Peers)
}

func Room(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		c.Status(400)
		return nil
	}
	ws := "ws"
	if os.Getenv("ENV") == "PROD" {
		ws = "wss"
	}
	uuid, suuid, _ := createOrGetRoom(uuid)
	return c.Render("peer", fiber.Map {
			"RoomWebSocketAddr": fmt.Sprintf("%s://%s/room/%s/websocket", ws, c.Hostname(), uuid),
			"RoomLink": fmt.Sprintf("%s://%s/room/%s", c.Protocol(), c.Hostname(), uuid),
			"ChatWebSocketAddr": fmt.Sprintf("%s://%s/room/%s/chat/websocket", ws, c.Hostname(), uuid),
			"ViewerWebSocketAddr": fmt.Sprintf("%s//%s/room/%s/viewer/websocket", ws, c.Hostname(), uuid),
			"StreamLink": fmt.Sprintf("%s://%s/stream/%s", c.Protocol(), c.Hostname(), suuid),
			"Type": "room",
		},
		"layouts/main",
	)
}

func createOrGetRoom(uuid string) (string, string, *w.Room) {
	w.RoomsLock.Lock()
	defer w.RoomsLock.Unlock()
	h := sha256.new()
	h.Write([]byte(uuid))
	suuid := fmt.Sprintf("%x", h.sum(nil))

	if room := w.Rooms[uuid]; room != nil {
		if _, ok := w.Streams[suuid]; !ok {
			w.Streams[suuid] = room
		}
		return uuid, suuid, room
	}
	hub := chat.NewHub()
	p := &w.Peers{}
	p.TrackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)
	room := &w.Room{
		Peers: p,
		Hub: hub,
	}
	w.Rooms[uuid] = room
	w.Streams[suuid] = room
	go hub.Run()
	return uuid, suuid, room
}

func RoomViewerWebsocket(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return 
	}

	w.RoomsLock.Lock()
	if peer, ok := w.Rooms[uuid]; ok {
		w.RoomsLock.Unlock()
		roomViewerConn(c. peer.Peers)
		return
	}
	w.RoomsLock.Unlock()
}

func roomViewerConn(c *websocket.Conn, p *w.Peers) {
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case <- ticker.C:
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(fmt.Sprintf("%d", len(p.Connections))))
		}
	}
}