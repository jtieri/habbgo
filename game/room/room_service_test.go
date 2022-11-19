package room

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRoomService_CreateAndStop(t *testing.T) {
	logger := zap.NewNop()
	ctx, cancel := context.WithCancel(context.Background())

	rs := NewRoomService(ctx, logger, nil, nil, cancel)

	t.Log("Starting the room service...")
	go rs.Start()
	time.Sleep(10 * time.Second)
	require.True(t, rs.running)

	t.Log("Stopping the room service...")
	cancel()
	time.Sleep(5 * time.Second)
	require.False(t, rs.running)
}

func TestRoomService_AddRoom(t *testing.T) {
	logger := zap.NewNop()
	ctx, cancel := context.WithCancel(context.Background())

	rs := NewRoomService(ctx, logger, nil, nil, cancel)
	go rs.Start()

	proxy := NewProxy(rs.channels)

	room1 := NewRoom()
	room1.Details.Id = 1
	room1.Details.Name = "Room-1"

	t.Log("Sending a room to the service to be added to the cache...")
	proxy.AddRoom(room1)

	t.Log("Requesting the size of the room cache from the service...")
	size := proxy.RoomsCachedCount()
	require.Equal(t, 1, size)

	room2 := NewRoom()
	room2.Details.Id = 2
	room2.Details.Name = "Room-2"

	t.Log("Sending another room to the service to be added to the cache...")
	proxy.AddRoom(room2)

	t.Log("Requesting the size of the room cache from the service...")
	size = proxy.RoomsCachedCount()
	require.Equal(t, 2, size)

	t.Log("Trying to add a room with the same room ID as the previously added room...")
	room3 := NewRoom()
	room3.Details.Id = 2
	room3.Details.Name = "Room-3"

	proxy.AddRoom(room3)

	t.Log("Requesting the size of the room cache from the service...")
	size = proxy.RoomsCachedCount()
	require.Equal(t, 2, size)

	t.Log("Requesting a copy of the room cache...")
	rooms := proxy.Rooms()

	require.Equal(t, 2, len(rooms))
	require.NotEqual(t, "Room-3", rooms[0].Details.Name)
	require.NotEqual(t, "Room-3", rooms[1].Details.Name)
}

func TestRoomService_RemoveRoom(t *testing.T) {
	logger := zap.NewNop()
	ctx, cancel := context.WithCancel(context.Background())

	rs := NewRoomService(ctx, logger, nil, nil, cancel)
	go rs.Start()

	proxy := NewProxy(rs.channels)

	room1 := NewRoom()
	room1.Details.Id = 1
	room1.Details.Name = "Room-1"

	t.Log("Sending a room to the service to be added to the cache...")
	proxy.AddRoom(room1)

	t.Log("Requesting the size of the room cache from the service...")
	size := proxy.RoomsCachedCount()
	require.Equal(t, 1, size)

	t.Log("Sending request to delete the room from the cache...")
	proxy.RemoveRoom(1)

	t.Log("Requesting the size of the room cache from the service...")
	size = proxy.RoomsCachedCount()
	require.Equal(t, 0, size)
}
