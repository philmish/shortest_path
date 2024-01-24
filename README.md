## Shortest Path

Trying to implement Dijkstra's algorithm in go.
This is a http server serving the shortest routes, starting from node "0", for a hardcoded graph

## Usage

To compile the binary the go toolkit must be installed on your system.

1. Build the server
```shell
make build
```
2. Run the server (default bind address ":8986")
```shell
bin/server
```

3. Curl the endpoint to receive json
```shell
curl http://localhost:8986/
```

4. Receive following json
```json
{
  "0": {
    "weight": 2,
    "route": [
      "1",
      "0"
    ]
  },
  "1": {
    "weight": 0,
    "route": [
      "1"
    ]
  },
  "2": {
    "weight": 5,
    "route": [
      "2"
    ]
  },
  "3": {
    "weight": 5,
    "route": [
      "1",
      "3"
    ]
  },
  "4": {
    "weight": 15,
    "route": [
      "1",
      "3",
      "4"
    ]
  },
  "5": {
    "weight": 20,
    "route": [
      "1",
      "3",
      "5"
    ]
  },
  "6": {
    "weight": 17,
    "route": [
      "1",
      "3",
      "4",
      "6"
    ]
  }
}

```
