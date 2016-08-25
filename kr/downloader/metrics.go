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

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/krypton/go-krypton/metrics"
)

var (
	hashInMeter      = metrics.NewMeter("kr/downloader/hashes/in")
	hashReqTimer     = metrics.NewTimer("kr/downloader/hashes/req")
	hashDropMeter    = metrics.NewMeter("kr/downloader/hashes/drop")
	hashTimeoutMeter = metrics.NewMeter("kr/downloader/hashes/timeout")

	blockInMeter      = metrics.NewMeter("kr/downloader/blocks/in")
	blockReqTimer     = metrics.NewTimer("kr/downloader/blocks/req")
	blockDropMeter    = metrics.NewMeter("kr/downloader/blocks/drop")
	blockTimeoutMeter = metrics.NewMeter("kr/downloader/blocks/timeout")

	headerInMeter      = metrics.NewMeter("kr/downloader/headers/in")
	headerReqTimer     = metrics.NewTimer("kr/downloader/headers/req")
	headerDropMeter    = metrics.NewMeter("kr/downloader/headers/drop")
	headerTimeoutMeter = metrics.NewMeter("kr/downloader/headers/timeout")

	bodyInMeter      = metrics.NewMeter("kr/downloader/bodies/in")
	bodyReqTimer     = metrics.NewTimer("kr/downloader/bodies/req")
	bodyDropMeter    = metrics.NewMeter("kr/downloader/bodies/drop")
	bodyTimeoutMeter = metrics.NewMeter("kr/downloader/bodies/timeout")

	receiptInMeter      = metrics.NewMeter("kr/downloader/receipts/in")
	receiptReqTimer     = metrics.NewTimer("kr/downloader/receipts/req")
	receiptDropMeter    = metrics.NewMeter("kr/downloader/receipts/drop")
	receiptTimeoutMeter = metrics.NewMeter("kr/downloader/receipts/timeout")

	stateInMeter      = metrics.NewMeter("kr/downloader/states/in")
	stateReqTimer     = metrics.NewTimer("kr/downloader/states/req")
	stateDropMeter    = metrics.NewMeter("kr/downloader/states/drop")
	stateTimeoutMeter = metrics.NewMeter("kr/downloader/states/timeout")
)
