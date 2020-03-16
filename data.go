package main

type RoomsData []RoomData

// RoomData Room Data
type RoomData struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type UserData struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	event <-chan UserEvent
}

func (u UserData) EventCh() <-chan UserEvent {
	return u.event
}

type Action struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
