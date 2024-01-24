package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Node struct {
    Value string
    Neighbours []*Neighbour
}

type Neighbour struct {
    Cost int
    Node *Node
}

type Edge struct {
    Weight int `json:"weight"`
    Route []string `json:"route"`
}

var (
    ZERO = &Node{
        Value: "0",
        Neighbours: make([]*Neighbour, 0),
    }
    ONE = &Node{
        Value: "1",
        Neighbours: make([]*Neighbour, 0),
    }
    TWO = &Node{
        Value: "2",
        Neighbours: make([]*Neighbour, 0),
    }
    THREE = &Node{
        Value: "3",
        Neighbours: make([]*Neighbour, 0),
    }
    FOUR = &Node{
        Value: "4",
        Neighbours: make([]*Neighbour, 0),
    }
    FIVE = &Node{
        Value: "5",
        Neighbours: make([]*Neighbour, 0),
    }
    SIX = &Node{
        Value: "6",
        Neighbours: make([]*Neighbour, 0),
    }
)

type Graph struct {
    nodes []*Node
}

func (g Graph)getShortestPaths(startAddr string) map[string]*Edge {
    inGraph := false
    for _, n := range g.nodes {
        if n.Value == startAddr {
            inGraph = true
            break
        }
    }
    if !inGraph {
        log.Fatalf("%s is not in the graph.", startAddr)
    }
    nodeMap := make(map[string]*Edge, len(g.nodes))

    for _, node := range g.nodes {
        if node.Value == startAddr {
            nodeMap[node.Value] = &Edge{Weight: 0, Route: []string{startAddr}}
        } else {
            nodeMap[node.Value] = &Edge{Weight: -1, Route: make([]string, 0)}
        }
    }
    for _, gnode := range g.nodes {
        currCost := nodeMap[gnode.Value].Weight
        currRoute := nodeMap[gnode.Value].Route
        for _, neighbour := range gnode.Neighbours {
            if nodeMap[neighbour.Node.Value].Weight == -1 {
                routeToNeighbour := make([]string, 0)
                routeToNeighbour = append(routeToNeighbour, currRoute...)
                routeToNeighbour = append(routeToNeighbour, neighbour.Node.Value)
                nodeMap[neighbour.Node.Value].Weight = currCost+neighbour.Cost
                nodeMap[neighbour.Node.Value].Route = routeToNeighbour
            } else {
                tmpCost := nodeMap[neighbour.Node.Value].Weight
                tmpRoute := nodeMap[neighbour.Node.Value].Route

                if nodeMap[gnode.Value].Weight > tmpCost+neighbour.Cost {
                    nodeMap[gnode.Value].Route = append(tmpRoute, gnode.Value)
                    nodeMap[gnode.Value].Weight = tmpCost+neighbour.Cost
                }

            }
        } 
    }
    return nodeMap
}

func generateGraph() Graph {
    var (
        ZERO = &Node{
            Value: "0",
            Neighbours: make([]*Neighbour, 0),
        }
        ONE = &Node{
            Value: "1",
            Neighbours: make([]*Neighbour, 0),
        }
        TWO = &Node{
            Value: "2",
            Neighbours: make([]*Neighbour, 0),
        }
        THREE = &Node{
            Value: "3",
            Neighbours: make([]*Neighbour, 0),
        }
        FOUR = &Node{
            Value: "4",
            Neighbours: make([]*Neighbour, 0),
        }
        FIVE = &Node{
            Value: "5",
            Neighbours: make([]*Neighbour, 0),
        }
        SIX = &Node{
            Value: "6",
            Neighbours: make([]*Neighbour, 0),
        }
    )
    ZERO.Neighbours = append(
        ZERO.Neighbours, 
        &Neighbour{Cost: 2, Node: ONE}, 
        &Neighbour{Cost: 6, Node: TWO},
    )
    ONE.Neighbours = append(
        ONE.Neighbours, 
        &Neighbour{Cost: 2, Node: ZERO}, 
        &Neighbour{Cost: 5, Node: THREE},
    )
    TWO.Neighbours = append(
        TWO.Neighbours, 
        &Neighbour{Cost: 6, Node: ZERO}, 
        &Neighbour{Cost: 8, Node: THREE},
    )
    THREE.Neighbours = append(
        THREE.Neighbours, 
        &Neighbour{Cost: 5, Node: ONE}, 
        &Neighbour{Cost: 8, Node: TWO},
        &Neighbour{Cost: 10, Node: FOUR},
        &Neighbour{Cost: 15, Node: FIVE},
    )
    FOUR.Neighbours = append(
        FOUR.Neighbours, 
        &Neighbour{Cost: 10, Node: THREE},
        &Neighbour{Cost: 6, Node: FIVE},
        &Neighbour{Cost: 2, Node: SIX},
    )
    FIVE.Neighbours = append(
        FIVE.Neighbours, 
        &Neighbour{Cost: 15, Node: THREE},
        &Neighbour{Cost: 6, Node: FOUR},
        &Neighbour{Cost: 6, Node: SIX},
    )
    SIX.Neighbours = append(
        FIVE.Neighbours, 
        &Neighbour{Cost: 2, Node: FOUR},
        &Neighbour{Cost: 6, Node: FIVE},
    )

    return Graph{
        nodes: []*Node{
            ZERO, ONE, TWO, THREE, FOUR, FIVE, SIX,
        },
    }
}


func main() {
    var listenAddr = flag.String("listen", ":8986", "Address to bind the server to")
    flag.Parse()
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        graph := generateGraph()
        nodeMap := graph.getShortestPaths("1")
        resp, err := json.Marshal(nodeMap)
        if err != nil {
            http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.Write(resp)
    })

    httpServer := http.Server{
        Addr: *listenAddr,
    }

    idleConnectionsClosed := make(chan struct{})
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, os.Interrupt, syscall.SIGKILL, syscall.SIGINT)
        <-sigChan
        if err := httpServer.Shutdown(context.Background()); err != nil {
            log.Printf("HTTP Server shutdown failed: %s", err.Error())
        }
        close(idleConnectionsClosed)
    }()

    if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("HTTP server listen and serve error: %s", err.Error())
    }
    <-idleConnectionsClosed

    log.Println("Server Shutdown")
}
