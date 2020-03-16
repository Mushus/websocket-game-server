package main

import (
	"time"
)

type room struct {
	in      chan interface{}
	name    string
	id      string
	users   map[string]*user
	created time.Time
}

func (r *room) run() {
	for in := range r.in {
		switch in := in.(type) {
		case roomToDataOrder:
			r.toDataProc(in)
		case joinUserOrder:
			r.joinUserProc(in)
		case leaveUserOrder:
			r.leaveUserProc(in)
		case sendOrder:
			r.sendProc(in)
		}
	}
}

type roomToDataOrder struct {
	out chan RoomData
}

func (r room) toData() RoomData {
	out := make(chan RoomData)
	r.in <- roomToDataOrder{
		out: out,
	}
	return <-out
}

func (r room) toDataProc(in roomToDataOrder) {
	in.out <- r.unsafeToData()
}

type joinUserOrder struct {
	out  chan RoomData
	user *user
}

func (r room) joinUser(user *user) RoomData {
	out := make(chan RoomData)
	r.in <- joinUserOrder{
		out:  out,
		user: user,
	}
	return <-out
}

func (r room) joinUserProc(in joinUserOrder) {
	r.users[in.user.id] = in.user
	in.out <- r.unsafeToData()
	ue := sendOrder{
		userID:  in.user.id,
		typ:     "join",
		payload: nil,
	}
	r.sendProc(ue)
}

type leaveUserOrder struct {
	userID string
}

func (r room) leaveUser(userID string) {
	r.in <- leaveUserOrder{
		userID: userID,
	}
}

func (r room) leaveUserProc(in leaveUserOrder) {
	delete(r.users, in.userID)
	ue := sendOrder{
		userID:  in.userID,
		typ:     "leave",
		payload: nil,
	}
	r.sendProc(ue)
}

func (r room) unsafeToData() RoomData {
	return RoomData{
		Name: r.name,
		ID:   r.id,
	}
}

type sendOrder struct {
	userID  string
	typ     string
	payload interface{}
}

func (r room) send(userID string, typ string, paylaod interface{}) {
	r.in <- sendOrder{
		userID:  userID,
		typ:     typ,
		payload: paylaod,
	}
}

type UserEvent struct {
	UserID  string      `json:"userId"`
	Payload interface{} `json:"payload,omitempty"`
	Type    string      `json:"type"`
}

func (r room) sendProc(in sendOrder) {
	ue := UserEvent{
		UserID:  in.userID,
		Payload: in.payload,
		Type:    in.typ,
	}

	for _, u := range r.users {
		u.send(ue)
	}
}
