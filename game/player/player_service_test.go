package	player

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestPlayerService_AddAndGetPlayer(t *testing.T) {
	logger := zap.NewNop()
	ctx, canc := context.WithCancel(context.Background())

	ps := NewPlayerService(ctx, logger, nil, nil, canc)
	go ps.Start()

	proxy := NewProxy(ps.channels)

	const playerCount = 10
	players := make([]Player, playerCount)

	// add some players to the PlayerService via the PlayerServiceProxy
	for i := 0; i < playerCount; i++ {
		p := New(ctx, logger, nil, nil, nil)
		p.Details.Username = fmt.Sprintf("test-%d", i)

		players[i] = p

		proxy.AddPlayer(p)
	}

	// attempt to fetch the players from the PlayerService via the PlayerServiceProxy
	for _, p := range players {
		cachedPlayer := proxy.GetPlayer(p)
		require.Equal(t, p.Details.Username, cachedPlayer.Details.Username)
	}
}

func TestPlayerService_RemovePlayer(t *testing.T) {
	logger := zap.NewNop()
	ctx, canc := context.WithCancel(context.Background())

	ps := NewPlayerService(ctx, logger, nil, nil, canc)
	go ps.Start()

	proxy := NewProxy(ps.channels)

	const playerCount = 10
	players := make([]Player, playerCount)

	// add some players to the PlayerService via the PlayerServiceProxy
	for i := 0; i < playerCount; i++ {
		p := New(ctx, logger, nil, nil, nil)
		p.Details.Username = fmt.Sprintf("test-%d", i)

		players[i] = p
		proxy.AddPlayer(p)
	}

	// attempt to fetch the players from the PlayerService via the PlayerServiceProxy
	for _, p := range players {
		cachedPlayer := proxy.GetPlayer(p)
		require.Equal(t, p.Details.Username, cachedPlayer.Details.Username)
	}

	// remove the players from the PlayerService via the PlayerServiceProxy
	for _, p := range players {
		proxy.RemovePlayer(p)
	}

	// verify the players are removed from the PlayerService via the PlayerServiceProxy
	for _, p := range players {
		cachedPlayer := proxy.GetPlayer(p)
		require.Equal(t, Player{}, cachedPlayer)
	}
}