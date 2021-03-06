// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package chunks

import (
	"io"

	"github.com/attic-labs/noms/go/hash"
)

// ChunkStore is the core storage abstraction in noms. We can put data
// anyplace we have a ChunkStore implementation for.
type ChunkStore interface {
	// Get the Chunk for the value of the hash in the store. If the hash is
	// absent from the store EmptyChunk is returned.
	Get(h hash.Hash) Chunk

	// GetMany gets the Chunks with |hashes| from the store. On return,
	// |foundChunks| will have been fully sent all chunks which have been
	// found. Any non-present chunks will silently be ignored.
	GetMany(hashes hash.HashSet, foundChunks chan *Chunk)

	// Returns true iff the value at the address |h| is contained in the
	// store
	Has(h hash.Hash) bool

	// Returns a new HashSet containing any members of |hashes| that are
	// absent from the store.
	HasMany(hashes hash.HashSet) (absent hash.HashSet)

	// Put caches c in the ChunkSource. Upon return, c must be visible to
	// subsequent Get and Has calls, but must not be persistent until a call
	// to Flush(). Put may be called concurrently with other calls to Put(),
	// Get(), GetMany(), Has() and HasMany().
	Put(c Chunk)

	// Returns the NomsVersion with which this ChunkSource is compatible.
	Version() string

	// Rebase brings this ChunkStore into sync with the persistent storage's
	// current root.
	Rebase()

	// Root returns the root of the database as of the time the ChunkStore
	// was opened or the most recent call to Rebase.
	Root() hash.Hash

	// Commit atomically attempts to persist all novel Chunks and update the
	// persisted root hash from last to current. If last doesn't match the
	// root in persistent storage, returns false.
	// TODO: is last now redundant? Maybe this should just try to update from
	// the cached root to current?
	Commit(current, last hash.Hash) bool

	// Stats may return some kind of struct that reports statistics about the
	// ChunkStore instance. The type is implementation-dependent, and impls
	// may return nil
	Stats() interface{}

	io.Closer
}

// Factory allows the creation of namespaced ChunkStore instances. The details
// of how namespaces are separated is left up to the particular implementation
// of Factory and ChunkStore.
type Factory interface {
	CreateStore(ns string) ChunkStore

	// Shutter shuts down the factory. Subsequent calls to CreateStore() will fail.
	Shutter()
}
