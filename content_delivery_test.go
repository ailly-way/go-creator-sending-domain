package main

import "testing"

func TestRenderUpdateKeepsProcessedAssetContent(t *testing.T) {
	asset := digitalAsset{Title: "Studio cut", Body: "  A short subscriber note.  "}
	want := "New from the studio: Studio cut\n\nA short subscriber note."
	if got := renderUpdate(asset); got != want {
		t.Fatalf("renderUpdate() = %q, want %q", got, want)
	}
}
