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

// Contains the metrics collected by the fetcher.

package fetcher

import (
	"github.com/krypton/go-krypton/metrics"
)

var (
	propAnnounceInMeter   = metrics.NewMeter("kr/fetcher/prop/announces/in")
	propAnnounceOutTimer  = metrics.NewTimer("kr/fetcher/prop/announces/out")
	propAnnounceDropMeter = metrics.NewMeter("kr/fetcher/prop/announces/drop")
	propAnnounceDOSMeter  = metrics.NewMeter("kr/fetcher/prop/announces/dos")

	propBroadcastInMeter   = metrics.NewMeter("kr/fetcher/prop/broadcasts/in")
	propBroadcastOutTimer  = metrics.NewTimer("kr/fetcher/prop/broadcasts/out")
	propBroadcastDropMeter = metrics.NewMeter("kr/fetcher/prop/broadcasts/drop")
	propBroadcastDOSMeter  = metrics.NewMeter("kr/fetcher/prop/broadcasts/dos")

	blockFetchMeter  = metrics.NewMeter("kr/fetcher/fetch/blocks")
	headerFetchMeter = metrics.NewMeter("kr/fetcher/fetch/headers")
	bodyFetchMeter   = metrics.NewMeter("kr/fetcher/fetch/bodies")

	blockFilterInMeter   = metrics.NewMeter("kr/fetcher/filter/blocks/in")
	blockFilterOutMeter  = metrics.NewMeter("kr/fetcher/filter/blocks/out")
	headerFilterInMeter  = metrics.NewMeter("kr/fetcher/filter/headers/in")
	headerFilterOutMeter = metrics.NewMeter("kr/fetcher/filter/headers/out")
	bodyFilterInMeter    = metrics.NewMeter("kr/fetcher/filter/bodies/in")
	bodyFilterOutMeter   = metrics.NewMeter("kr/fetcher/filter/bodies/out")
)
