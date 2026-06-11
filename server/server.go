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

func (gs *gameServer) GetVisibleObjects(context.Context, *Empty) (*GameObjects, error) {
	gs.game.EnumerateVisibleObjects()
	wowObjects := gs.game.GetVisibleObjects()
	gameObjects := []*GameObject{}
	for _, wowObject := range wowObjects {
		gameObjects = append(gameObjects, wowObjectToGameObject(wowObject))
	}
	return &GameObjects{
		Objects: gameObjects,
	}, nil
}

func wowObjectToGameObject(wowObject game.WowObject) *GameObject {
	unitType := wowUnitTypeToUnitType(wowObject.Type)
	return &GameObject{
		Guid: &wowObject.Guid,
		Name: &wowObject.Name,
		Type: &unitType,
		Position: &Vec3{
			X: &wowObject.Position.X,
			Y: &wowObject.Position.Y,
			Z: &wowObject.Position.Z,
		},
		MaxHealth:     &wowObject.MaxHealth,
		CurrentHealth: &wowObject.CurrentHealth,
		MaxMana:       &wowObject.MaxMana,
		CurrentMana:   &wowObject.CurrentMana,
	}
}

func wowUnitTypeToUnitType(wowUnitType uint8) UnitType {
	switch wowUnitType {
	case game.None:
		return UnitType_None
	case game.Item:
		return UnitType_Item
	case game.Container:
		return UnitType_Container
	case game.Unit:
		return UnitType_Unit
	case game.Player:
		return UnitType_Player
	case game.GameObject:
		return UnitType_Object
	case game.DynamicObject:
		return UnitType_DynamicObject
	case game.Corpse:
		return UnitType_Corpse
	}
	return UnitType_None
}

func (gs *gameServer) MoveTo(ctx context.Context, pos *Vec3) (*Empty, error) {
	gs.game.MoveToPosition(game.Vec3{X: *pos.X, Y: *pos.Y, Z: *pos.Z})
	return &Empty{}, nil
}

func (gs *gameServer) Jump(context.Context, *Empty) (*Empty, error) {
	gs.game.Jump()
	return &Empty{}, nil
}

func (gs *gameServer) ToggleAttack(ctx context.Context, request *ToggleAttackRequest) (*Empty, error) {
	gs.game.ToggleAttack(*request.Attack)
	return &Empty{}, nil
}

func (gs *gameServer) CastSpellByName(ctx context.Context, request *CastSpellByNameRequest) (*Empty, error) {
	gs.game.CastSpellByName(*request.SpellName)
	return &Empty{}, nil
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
