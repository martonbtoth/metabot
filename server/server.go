package server

import (
	context "context"
	"fmt"
	"net"
	"superbot/game"
	"superbot/logger"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gameServer struct {
	UnimplementedGameServiceServer

	game game.Game
}

func (gs *gameServer) GetPlayerGuid(context.Context, *Empty) (*PlayerGuid, error) {
	playerGuid := gs.game.GetPlayerGuid()
	return &PlayerGuid{PlayerGuid: &playerGuid}, nil
}

func (gs *gameServer) GetSpellbook(context.Context, *Empty) (*Spellbook, error) {
	spellsString := gs.game.GetAvailableSpells()
	spells := []*Spell{}
	for _, spellString := range spellsString {
		spells = append(spells, &Spell{Name: &spellString})
	}
	return &Spellbook{Spells: spells}, nil
}

func newServer(game game.Game) *gameServer {
	return &gameServer{game: game}
}

func Listen(game game.Game) {
	listenOn := fmt.Sprintf("0.0.0.0:%d", 3333)
	lis, err := net.Listen("tcp", listenOn)
	if err != nil {
		logger.Log(fmt.Sprintf("failed to listen: %v", err))
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	RegisterGameServiceServer(grpcServer, newServer(game))
	reflection.Register(grpcServer)

	go func() {
		logger.Log("gRPC listening on " + listenOn)
		grpcServer.Serve(lis)
	}()
}
