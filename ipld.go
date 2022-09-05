package main

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/storage/memstore"
	"github.com/ipld/go-ipld-prime/traversal"
)

// Just import to register encoder for codec
var dagcborEncode = dagcbor.Encode

var frameHashIPLDSchema = []byte(`
	type FrameHash struct {
		Next optional &FrameHash
		Hash Bytes
	}
`)

type FrameHash struct {
	Next datamodel.Link
	Hash []byte
}

type SimpleDAG struct {
	lsys        linking.LinkSystem
	frameHashes [][]byte
}

func NewSimpleDAG(frameHashes [][]byte) *SimpleDAG {
	lsys := cidlink.DefaultLinkSystem()
	store := &memstore.Store{}
	lsys.SetWriteStorage(store)
	lsys.SetReadStorage(store)

	return &SimpleDAG{
		lsys:        lsys,
		frameHashes: frameHashes,
	}
}

func (dag *SimpleDAG) SaveFrameHash(fh *FrameHash) (datamodel.Link, error) {
	ts, err := ipld.LoadSchemaBytes(frameHashIPLDSchema)
	if err != nil {
		return nil, err
	}

	schemaType := ts.TypeByName("FrameHash")

	lp := cidlink.LinkPrototype{
		Prefix: cid.Prefix{
			Version:  1,
			Codec:    0x71,
			MhType:   0x13,
			MhLength: 64,
		},
	}

	n := bindnode.Wrap(fh, schemaType)
	nr := n.Representation()

	return dag.lsys.Store(linking.LinkContext{}, lp, nr)
}

func (dag *SimpleDAG) Save() (datamodel.Link, error) {
	var next datamodel.Link
	for i := len(dag.frameHashes) - 1; i >= 0; i-- {
		fh := &FrameHash{
			Next: next,
			Hash: dag.frameHashes[i],
		}

		link, err := dag.SaveFrameHash(fh)
		if err != nil {
			return nil, err
		}

		next = link
	}

	return next, nil
}

func (dag *SimpleDAG) PrintLinks(link datamodel.Link) error {
	currLink := link

	for currLink != nil {
		fmt.Println(currLink)

		np := basicnode.Prototype.Any
		n, err := dag.lsys.Load(linking.LinkContext{}, currLink, np)
		if err != nil {
			return err
		}

		links, err := traversal.SelectLinks(n)
		if err != nil {
			return err
		}

		if len(links) > 0 {
			// Assume only a single link
			currLink = links[0]
		} else {
			currLink = nil
		}
	}

	return nil
}
