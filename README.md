# websocket-game-server
game server over websocket

## Install

You can Download this binary from [Release Tab](https://github.com/Mushus/websocket-game-server/releases)

## Usage

### GET: `/rooms`

List all rooms.

### POST: `/rooms`

Create a room

#### Content-Type

application/json

#### Payload

| key  |  type  | required | description |
|:----:|:------:|:--------:|:-----------:|
| name | string |   yes    |  room name  |

### WS: `/room/{id}`

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