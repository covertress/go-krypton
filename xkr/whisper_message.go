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

// Contains the external API representation of a whisper message.

package xkr

import (
	"time"

	"github.com/krypton/go-krypton/common"
	"github.com/krypton/go-krypton/crypto"
	"github.com/krypton/go-krypton/whisper"
)

// WhisperMessage is the external API representation of a whisper.Message.
type WhisperMessage struct {
	ref *whisper.Message

	Payload string `json:"payload"`
	To      string `json:"to"`
	From    string `json:"from"`
	Sent    int64  `json:"sent"`
	TTL     int64  `json:"ttl"`
	Hash    string `json:"hash"`
}

// NewWhisperMessage converts an internal message into an API version.
func NewWhisperMessage(message *whisper.Message) WhisperMessage {
	return WhisperMessage{
		ref: message,

		Payload: common.ToHex(message.Payload),
		From:    common.ToHex(crypto.FromECDSAPub(message.Recover())),
		To:      common.ToHex(crypto.FromECDSAPub(message.To)),
		Sent:    message.Sent.Unix(),
		TTL:     int64(message.TTL / time.Second),
		Hash:    common.ToHex(message.Hash.Bytes()),
	}
}
