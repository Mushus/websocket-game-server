# websocket-game-server
game server over websocket

## Install

You can Download this binary from [Release Tab](https://github.com/Mushus/websocket-game-server/releases)

## Usage

### GET: `/playground`

Show Playground Page

#### Response Content-Type

text/html

### GET: `/rooms`

List all rooms.

#### Response Content-Type

application/json

#### Response body

|      key      |  type  | required | description |
|:-------------:|:------:|:--------:|:-----------:|
|    .rooms     | array  |   yes    | rooms data  |
|  .rooms[].id  | string |   yes    |   room id   |
| .rooms[].name | string |   yes    |  room name  |


### POST: `/rooms`

Create a room

#### Request Content-Type

application/json

#### Payload

| key  |  type  | required | description |
|:----:|:------:|:--------:|:-----------:|
| name | string |   yes    |  room name  |

#### Response Content-Type

application/json

#### Response body

|    key     |  type  | required | description |
|:----------:|:------:|:--------:|:-----------:|
|   .room    | object |   yes    |  room data  |
|  .room.id  | string |   yes    |   room id   |
| .room.name | string |   yes    |  room name  |

### WS: `/rooms/{id}`

Join `id` room

#### URL params

| key |  type  | required | description |
|:---:|:------:|:--------:|:-----------:|
| id  | string |   yes    |   room ID   |

#### GET params

| key  |  type  | required | description |
|:----:|:------:|:--------:|:-----------:|
| name | string |   yes    |  user name  |

#### Push messages format

|   key   |  type  | required |  description   |
|:-------:|:------:|:--------:|:--------------:|
| payload |  any   |   yes    | action payload |
|  type   | string |   yes    |  action type   |

#### Pull messages format

##### Generic

|   key   |  type  | required |  description   |
|:-------:|:------:|:--------:|:--------------:|
| userId  | string |   yes    |    user id     |
| payload |  any   |    no    | action payload |
|  type   | string |   yes    |  action type   |

##### Join room

|  key   |  type  | required | description |
|:------:|:------:|:--------:|:-----------:|
| userId | string |   yes    |   user id   |
|  type  | "join" |   yes    |             |

##### Leave room

|  key   |  type   | required | description |
|:------:|:-------:|:--------:|:-----------:|
| userId | string  |   yes    |   user id   |
|  type  | "leave" |   yes    |             |