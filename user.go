package main

type user struct {
	in    chan interface{}
	event chan UserEvent
	name  string
	id    string
}

func (u *user) run() {
	for in := range u.in {
		switch in := in.(type) {
		case userToDataOrder:
			u.toDataProc(in)
		}
	}
}

type userToDataOrder struct {
	out chan UserData
}

func (u user) toData() UserData {
	out := make(chan UserData)
	u.in <- userToDataOrder{
		out: out,
	}
	return <-out
}

func (u user) toDataProc(in userToDataOrder) {
	in.out <- UserData{
		Name:  u.name,
		ID:    u.id,
		event: u.event,
	}
}

func (u user) send(in UserEvent) {
	u.event <- in
}
