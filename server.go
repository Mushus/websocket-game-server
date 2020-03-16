package main

import (
	"errors"
	"sync"
	"time"

	"github.com/rs/xid"
)

type server struct {
	rooms map[string]*room
	in    chan interface{}
}

func NewServer() *server {
	s := &server{
		in:    make(chan interface{}, 10),
		rooms: map[string]*room{},
	}

	go s.run()

	return s
}

func (s server) run() {
	for in := range s.in {
		switch in := in.(type) {
		case listRoomsOrder:
			s.listRoomsProc(in)
		case createRoomOrder:
			s.createRoomProc(in)
		case joinUserToRoomOrder:
			s.joinUserProc(in)
		case leaveUserfromRoomOrder:
			s.leaveUserProc(in)
		case sendToRoomOrder:
			s.sendToRoomProc(in)
		}
	}
}

type listRoomsOrder struct {
	out chan listRoomResult
}

type listRoomResult struct {
	rooms []RoomData
}

func (s *server) listRooms() []RoomData {
	out := make(chan listRoomResult)
	s.in <- listRoomsOrder{
		out: out,
	}
	result := <-out
	return result.rooms
}

func (s *server) listRoomsProc(in listRoomsOrder) {
	wg := &sync.WaitGroup{}
	rooms := make([]RoomData, len(s.rooms))

	i := 0
	for roomID := range s.rooms {
		wg.Add(1)
		go func(i int, roomID string) {
			rooms[i] = s.rooms[roomID].toData()
			wg.Done()
		}(i, roomID)
		i++
	}
	wg.Wait()

	in.out <- listRoomResult{
		rooms: rooms,
	}
}

type createRoomOrder struct {
	out  chan RoomData
	name string
}

func (s *server) createRoom(name string) RoomData {
	out := make(chan RoomData)
	s.in <- createRoomOrder{
		out:  out,
		name: name,
	}
	return <-out
}

func (s *server) createRoomProc(in createRoomOrder) {
	id := xid.New().String()
	created := time.Now()
	r := &room{
		in:      make(chan interface{}, 10),
		name:    in.name,
		id:      id,
		created: created,
		users:   map[string]*user{},
	}
	go r.run()
	s.rooms[r.id] = r
	in.out <- r.toData()
}

type joinUserToRoomOrder struct {
	roomID   string
	userName string
	out      chan roomUserData
}

type roomUserData struct {
	user UserData
	room RoomData
	err  error
}

func (s *server) joinUser(roomID string, userName string) (UserData, RoomData, error) {
	out := make(chan roomUserData)
	s.in <- joinUserToRoomOrder{
		roomID:   roomID,
		userName: userName,
		out:      out,
	}
	ru := <-out
	return ru.user, ru.room, ru.err
}

func (s *server) joinUserProc(in joinUserToRoomOrder) {
	r, ok := s.rooms[in.roomID]
	if !ok {
		in.out <- roomUserData{
			err: errors.New("room not found"),
		}
		return
	}

	uid := xid.New().String()
	u := &user{
		in:    make(chan interface{}),
		event: make(chan UserEvent),
		name:  in.userName,
		id:    uid,
	}
	go u.run()

	roomData := r.joinUser(u)
	userData := u.toData()
	in.out <- roomUserData{
		room: roomData,
		user: userData,
	}
}

type leaveUserfromRoomOrder struct {
	roomID string
	userID string
}

func (s *server) leaveUser(roomID string, userID string) {
	s.in <- leaveUserfromRoomOrder{
		roomID: roomID,
		userID: userID,
	}
}

func (s *server) leaveUserProc(in leaveUserfromRoomOrder) {
	r, ok := s.rooms[in.roomID]
	if !ok {
		return
	}
	r.leaveUser(in.userID)
	return
}

type sendToRoomOrder struct {
	roomID     string
	userID     string
	actionType string
	payload    interface{}
}

func (s *server) sendToRoom(roomID string, userID string, actionType string, payload interface{}) {
	s.in <- sendToRoomOrder{
		roomID:     roomID,
		userID:     userID,
		actionType: actionType,
		payload:    payload,
	}
}
func (s *server) sendToRoomProc(in sendToRoomOrder) {
	r, ok := s.rooms[in.roomID]
	if !ok {
		return
	}
	r.send(in.userID, in.actionType, in.payload)
}
