// Copyright (c) 2021 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "testing"

func TestIndexCGIPermissions(t *testing.T) {
	f, err := static("index.cgi")()
	if err != nil {
		t.Fatal(err)
	}
	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode() != 0755 {
		t.Errorf("mode = %v; want 0755", fi.Mode())
	}
}
