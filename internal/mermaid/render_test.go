package mermaid

import (
	"bytes"
	"testing"

	"github.com/goatx/goat-cli/internal/test"
	"github.com/google/go-cmp/cmp"
)

func TestRenderSequenceDiagram(t *testing.T) {
	t.Parallel()

	t.Run("spec", func(t *testing.T) {
		pkg := loadSpecPackage(t)
		var buf bytes.Buffer
		if err := RenderSequenceDiagram(pkg, &buf); err != nil {
			t.Fatalf("RenderSequenceDiagram returned error: %v", err)
		}

		got := buf.String()
		want := test.ReadGolden(t, "sequence_diagram.golden")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("rendered diagram mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("protobuf", func(t *testing.T) {
		pkg := loadProtobufPackage(t)
		var buf bytes.Buffer
		if err := RenderSequenceDiagram(pkg, &buf); err != nil {
			t.Fatalf("RenderSequenceDiagram returned error: %v", err)
		}
		got := buf.String()
		want := test.ReadGolden(t, "protobuf/sequence_diagram.golden")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("rendered diagram (protobuf) mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("openapi", func(t *testing.T) {
		pkg := loadOpenapiPackage(t)
		var buf bytes.Buffer
		if err := RenderSequenceDiagram(pkg, &buf); err != nil {
			t.Fatalf("RenderSequenceDiagram returned error: %v", err)
		}
		got := buf.String()
		want := test.ReadGolden(t, "openapi/sequence_diagram.golden")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("rendered diagram (openapi) mismatch (-want +got):\n%s", diff)
		}
	})
}
