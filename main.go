package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

type RoomParam struct {
	ID string `uri:"id" binding:"required"`
}
type PostRoomParam struct {
	Name string `json:"name"`
}

func main() {
	s := NewServer()
	r := gin.Default()
	r.GET("/playground", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>playground</title>
<style>
textarea {
	display: block;
	width: 100%;
	height: 5em;
	box-sizing: border-box;
}
</style>
</head>
<h1>playground</h1>
<section>
<h2>GET /rooms</h2>
<p>
Input:<br>
<button onclick="requestGetRooms()">req</button>
</p>
<p>
Result:<br>
<textarea id="roomsResult"></textarea>
</p>
</section>
<section>
<h2>POST /rooms</h2>
<p>
Input:<br>
<input id="postRoomsName">:name<br>
<button onclick="requestPostRooms()">req</button>
</p>
<p>
Result:<br>
<textarea id="postRoomsResult"></textarea>
</p>
</section>
<section>
<h2>WS /room/<span id="wsJoinRoomIdTmpl"></span></h2>
<p>
Input(join room):<br>
<input id="wsJoinRoomId" oninput="changeWsJoinRoomId(this.value)" >:id<br>
<input id="wsJoinRoomName">:name<br>
<button id="wsJoinRoomButton" onclick="requestWsJoinRoom()">req</button>
</p>
<p>
Input(send message):<br>
<input id="wsSendMessage">:message<br>
<button id="wsSendBtn" onclick="requestWsSendMessage()" disabled>req</button>
</p>
<p>
Input(leave room):<br>
<button id="wsLeaveBtn" onclick="requestWsLeaveRoom()" disabled>req</button>
</p>
<p>
Result:<br>
<textarea id="wsJoinRoomsResult"></textarea>
</p>
</section>
<script>
const header = {
	'Content-Type': 'application/json; charset=utf-8'
};

async function requestGetRooms() {
	const data = await fetch('/rooms', { method: 'GET' });
	document.querySelector('#roomsResult').value = await data.text();
}

async function requestPostRooms() {
	const name = document.querySelector('#postRoomsName').value;
	const data = await fetch('/rooms', { method: 'POST', header, body: JSON.stringify({ name }) });
	document.querySelector('#postRoomsResult').value = await data.text();
}

function changeWsJoinRoomId(value) {
	document.querySelector('#wsJoinRoomIdTmpl').innerText = value;
}

function toggleWsInputs(isOn) {
	const joinBtn = document.querySelector('#wsJoinRoomButton');
	joinBtn.disabled = isOn;
	wsSendBtn.disabled = !isOn;
	wsLeaveBtn.disabled = !isOn;
}

let con;
async function requestWsJoinRoom() {
	toggleWsInputs(true);
	const result = document.querySelector('#wsJoinRoomsResult');
	const id = document.querySelector('#wsJoinRoomId').value;
	const name = document.querySelector('#wsJoinRoomName').value;
	const params = new URLSearchParams({ name });
	result.value = "";
	con = new WebSocket(`+"`ws://${location.host}/rooms/${id}?${params.toString()}`"+`);
	con.onpoen = e => {
		result.value += "start\n";
		console.log(e);
	}
	con.onmessage = e => {
		result.value += "msg\n";
		console.log(e);
	}
	con.onerror = e => {
		result.value += "error\n";
		console.log(e);
	}
	con.onclose = () => {
		result.value += "close\n";
		toggleWsInputs(false);
	}
}

async function requestWsSendMessage() {
	if (!con) return;
	const value = document.querySelector('#wsSendMessage').value;
	con.send(value);
}

async function requestWsLeaveRoom() {
	if (!con) return;
	con.close()
}
</script>
</html>
`)
	})
	r.GET("/rooms/:id", func(c *gin.Context) {
		var prm RoomParam
		if err := c.ShouldBindUri(&prm); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}

		name := c.Query("userName")

		user, room, err := s.joinUser(prm.ID, name)
		if err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		websocket.Handler(func(ws *websocket.Conn) {
			go func() {
				for event := range user.EventCh() {
					b, _ := json.Marshal(event)
					ws.Write(b)
				}
			}()

			decoder := json.NewDecoder(ws)
			for {
				var in Action
				if err := decoder.Decode(&in); err != nil {
					fmt.Printf("err: %v\n", err)
					break
				}
				fmt.Printf("msg: %v\n", in)
				s.sendToRoom(room.ID, user.ID, in.Type, in.Payload)
			}
			defer s.leaveUser(room.ID, user.ID)
		}).ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/rooms", func(c *gin.Context) {
		rooms := s.listRooms()
		c.JSON(200, gin.H{
			"rooms": rooms,
		})
	})
	r.POST("/rooms", func(c *gin.Context) {
		var prm PostRoomParam
		if err := c.BindJSON(&prm); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		room := s.createRoom(prm.Name)
		c.JSON(200, gin.H{
			"room": room,
		})
	})
	r.Run(":8080")
}
