package lexicons_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bluesky-social/indigo/atproto/atdata"
	"github.com/bluesky-social/indigo/atproto/lexicon"
)

func TestLexicons(t *testing.T) {
	cat := lexicon.NewBaseCatalog()
	if err := cat.LoadDirectory("com"); err != nil {
		t.Fatalf("loading lognship lexicons: %v", err)
	}
	if err := cat.LoadDirectory("testdata/external"); err != nil {
		t.Fatalf("loading external (atproto-core) lexicons: %v", err)
	}

	cases := []struct {
		name  string
		file  string
		valid bool
	}{
		{
			name:  "log entry: minimal valid",
			file:  "com/lognship/log/entry/minimal-valid.json",
			valid: true,
		},
		{
			name:  "log entry: with mention and atLink facets",
			file:  "com/lognship/log/entry/with-facets-valid.json",
			valid: true,
		},
		{
			name:  "log entry: missing required text",
			file:  "com/lognship/log/entry/missing-text-invalid.json",
			valid: false,
		},
		{
			name:  "log entry: empty features array on a facet",
			file:  "com/lognship/log/entry/empty-features-invalid.json",
			valid: false,
		},
		{
			name:  "log entry: malformed at-uri in atLink feature",
			file:  "com/lognship/log/entry/malformed-atlink-invalid.json",
			valid: false,
		},
		{
			name:  "log label: minimal valid",
			file:  "com/lognship/log/label/minimal-valid.json",
			valid: true,
		},
		{
			name:  "log label: missing required subject",
			file:  "com/lognship/log/label/missing-subject-invalid.json",
			valid: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join("testdata", tc.file))
			if err != nil {
				t.Fatalf("reading fixture: %v", err)
			}

			data, err := atdata.UnmarshalJSON(raw)
			if err != nil {
				t.Fatalf("parsing fixture as atproto data: %v", err)
			}

			nsid, _ := data["$type"].(string)
			if nsid == "" {
				t.Fatal("fixture missing $type")
			}

			err = lexicon.ValidateRecord(cat, data, nsid, 0)
			switch {
			case tc.valid && err != nil:
				t.Errorf("expected valid, got error: %v", err)
			case !tc.valid && err == nil:
				t.Error("expected invalid, got no error")
			}
		})
	}
}
