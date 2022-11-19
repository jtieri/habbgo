package service

/*
Proxies will serve as an intermediate API which Players utilize to make requests to the
various game services that are running. The Proxy abstraction allows us to handle the communication with
the game services behind a high level API so that Players don't have to manage writing/reading to channels.

Ex: The room service is started at startup and manages the global state for initialized/initializing rooms.
	It listens for requests for data on channels as well as requests to load rooms or update the cache of
	already loaded rooms. The room service proxy will provide high level requests,
	e.g. GetRoomForRoomID, GetPublicRooms, UpdatePlayerRoom, etc., that the Players use to fetch and store
	room data without having to know how this data is being handled.
*/
