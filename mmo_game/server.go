package main

import (
	"zinx/mmo_game/core"
	"zinx/ziface"
	"zinx/znet"
)

func main() {
	server := znet.NewServer()
	server.Server()
	server.SetOnConnStart(onConnectStart)
}

func onConnectStart(conn ziface.IConnection) {
	player := core.NewPlayer(conn)
	player.SyncPid()
	player.BroadCastStartPosition()
}
