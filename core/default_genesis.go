// Copyright 2014 The go-krypton Authors
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

package core

import (
	"compress/gzip"
	"encoding/base64"
	"io"
	"strings"
)

func NewDefaultGenesisReader() (io.Reader, error) {
	return gzip.NewReader(base64.NewDecoder(base64.StdEncoding, strings.NewReader(defaultGenesisBlock)))
}

const defaultGenesisBlock = "H4sIAAAJbogA/6yRzUrEMBSF3yXrWSRNk9zObn5EFyqCvsC9+XECaSttBipD391MuxBBFwPeRSDknO+cJBf23HfWsy3jE/8xsmIb9hZbP2ZsPxaB0oLLPUA5eMHBd/kBx9Mv1tunEO+mPOARMy5AIYg8OKq9rE1ZwTbSWGFBGu5rUt5JiY6kI91YcmRQghdecRIEVcACvMfxMbYxLzwtdtfaxxhCtOeUP9eU7/inOJ3+8TaHPnaE4/qwSKhBqEbxGirSFXBvKVClpSWtTBAN1CagKL5dSr1l28strqI+9O4aVQCvuR/wvWy6c0obtseE6/+Kv+vO8/wVAAD//zXruKEIAgAA";