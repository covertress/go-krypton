// Copyright 2015 The go-krypton Authors
// This file is part of the go-krypton library.
//
// The go-krypton library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-krypton library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-krypton library. If not, see <http://www.gnu.org/licenses/>.

package trie

import (
	"bytes"
	"testing"

	"github.com/krypton/go-krypton/common"
	"github.com/krypton/go-krypton/crypto"
	"github.com/krypton/go-krypton/krdb"
)

func newEmptySecure() *SecureTrie {
	db, _ := krdb.NewMemDatabase()
	trie, _ := NewSecure(common.Hash{}, db)
	return trie
}

func TestSecureDelete(t *testing.T) {
	trie := newEmptySecure()
	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"krypton", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"krypton", ""},
		{"dog", "puppy"},
		{"shaman", ""},
	}
	for _, val := range vals {
		if val.v != "" {
			trie.Update([]byte(val.k), []byte(val.v))
		} else {
			trie.Delete([]byte(val.k))
		}
	}
	hash := trie.Hash()
	exp := common.HexToHash("29b235a58c3c25ab83010c327d5932bcf05324b7d6b1185e650798034783ca9d")
	if hash != exp {
		t.Errorf("expected %x got %x", exp, hash)
	}
}

func TestSecureGetKey(t *testing.T) {
	trie := newEmptySecure()
	trie.Update([]byte("foo"), []byte("bar"))

	key := []byte("foo")
	value := []byte("bar")
	seckey := crypto.Sha3(key)

	if !bytes.Equal(trie.Get(key), value) {
		t.Errorf("Get did not return bar")
	}
	if k := trie.GetKey(seckey); !bytes.Equal(k, key) {
		t.Errorf("GetKey returned %q, want %q", k, key)
	}
}
