[![Build Status](https://travis-ci.com/msmedes/scale.svg?branch=master)](https://travis-ci.com/msmedes/scale)

# Scale

Go implementation of Chord - DHT protocol

### Local Development

#### Setup

Have Go installed, all that stuff

BloomRPC is helpful to test RPC calls

```bash
brew cask install bloomrpc
```

#### Workflow

- `make` - run linting and tests
- `make serve` - start grpc server
- `make trace` - start the trace server (defaults to port 5000)
- `PORT=<port number> JOIN=<node to join ip> WEB=<port number> make serve`
- `make trace.codegen` - generate protobuf code for the trace server
- `make scale.codegen` - generate protobuf code for the scale server

Internode communication is on port 3000 by default. GraphQL API is on port 8000 by default.
Useful GraphQL queries:

```graphql
query {
  get(key: "hello") {
    value
    trace {
      addr
      functionCall
      duration
    }
  }
  metadata {
    node {
      id
      addr
      fingerTable
      predecessor {
        id
        addr
      }
      successor {
        id
        addr
      }
    }
  }
}

mutation {
  set(key: "hello", value: "world") {
    count
    trace {
      addr
      functionCall
      duration
    }
  }
}
```

### Resources

- [A different, more thorough version of the other chord paper](https://people.csail.mit.edu/karger/Papers/chord.pdf)
- [Chord:A Scalable Peer-to-peer Lookup Protocol for Internet Applications](https://pdos.csail.mit.edu/papers/ton:chord/paper-ton.pdf)
- [Chord: Building a DHT (Distributed Hash Table) in Golang](https://medium.com/techlog/chord-building-a-dht-distributed-hash-table-in-golang-67c3ce17417b)
- [another repo](https://github.com/r-medina/gmaj)
- [Lecture on chord](https://www.youtube.com/watch?v=q29szpcnorA)
- [Brown Assignment](http://cs.brown.edu/courses/cs138/s17/content/projects/chord.pdf)
