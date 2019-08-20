### Some links

[Chord:A Scalable Peer-to-peer Lookup Protocol for Internet Applications](https://pdos.csail.mit.edu/papers/ton:chord/paper-ton.pdf)

[Chord: Building a DHT (Distributed Hash Table) in Golang](https://medium.com/techlog/chord-building-a-dht-distributed-hash-table-in-golang-67c3ce17417b) - [repo](https://github.com/arriqaaq/chord)

[another repo](https://github.com/r-medina/gmaj)

[Lecture on chord](https://www.youtube.com/watch?v=q29szpcnorA)

[Brown Assignment](http://cs.brown.edu/courses/cs138/s17/content/projects/chord.pdf)

### IDs

Using `[]byte` to represent ids for a few reasons:
  - can easily be truncated or padded, since `sha1` doesn't return a checksum of the same length for every input, apparently.  Chord uses a key of a fixed bit length (I think 32) which gives you a 2<sup>32</sup> possible keys, I think? In the video linked above they mention using a 160 bit sha hash but I don't think go does that and I'm not writing a custom hashing library for that shit.  As far as I know it returns either 32 or 64 bit int as a byte array.
  - can be converted to a `big.Int` `struct` or `interface` so the finger table math can be done to it easily since the integers involved will likely get huge.
  - Idk it's what everyone else did so seemed like a good idea.

Todo:
  - A between comparator (handling wrap around)
  - a greater than/less than comparator
  - padding ids if necessary (if `checksum` `[]byte` is shorter than `config.KeySize`)
  - Tests


### Fingertable

At the moment when a `Node` is constructed it's `FingerTable` will also be initialized, but each entry in the table will refer to itself until the `fixFingerTable()` runs (each `Node` will call this in a goroutine, and will run every few seconds to keep track of departures and arrivals).  Check `fingertable.go` for some psuedocode on what will go down.

Todo:
  - Actually write the `fixFingerTable()` func
  - Maybe a string() method for fun
  - Some more tests


### RPC functions
Here are some of the RPC's I can see `Node` needing:
  - `GetPredecessor()`: return the predecessor of the node
  - `GetSuccessor()`: return the successor of the node
  - `SetPredecessor()`: set the predecessor of the node (probably will be useful for exits)
  - `SetSuccessor()`: ditto
  - `FindSuccessor()` : Kind of the whole thing, given a key calculated by `fingerMath()`, find the successor to that key in the `Chord`
  - `GetKey()` : return value of the key at the node
  - `PutKey()` : store a `k,v` pair at the node

Todo:
  - Literally all of it


### All the server shit
Yeah idk where to even start with that since I don't really know the ins and outs of RPC.  But like
- How are these nodes going to find each other
- How are they going to pass their `k,v` pairs along when they fail
- How are nodes going to know they have to take `k,v` from a failed node
- what is the meaning of life

### Local dev

`make deps` - install dependencies
`make test` - run tests
`make serve` - start grpc server
