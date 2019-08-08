package scale

/*
(n + 2^(i-1)) mod (2^m) is the formula to remember for finger tables/chord
This will be used to calculate the offset for each entry of the finger table
from the root node.
n = number of nodes in the chord
i = the ith entry in the finger table
m = number of bits in the hash (I believe most implementations use 32 bit
sha1)
*/

type FingerTable []Finger

type Finger struct {
	ID   []byte
	Node *Node
}

func (f *Finger) String() string {
	var buf bytes.Buffer 

	buf.WriteString(fmt.Sprintf("id: %v node: %v", f.ID, Node.Id)

}

func NewFingerTable(n *Node, int m) {
	// TODO: A few finger table should take a node and an integer m
	// and construct a new finger table of fingers length m.
}

func NewFinger(id []byte, n *Node) *Finger {
	return &Finger{
		ID:   id,
		Node: n,
	}
}


