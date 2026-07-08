package nvueschema

import (
	"encoding/json"
	"strings"
	"testing"
)

// A scalar union nested inside a single-branch anyOf wrapper (the OpenAPI
// idiom `{"anyOf":[{"$ref":...}],"nullable":true}` used to attach nullable to
// a $ref) must not be collapsed to its first branch. This mirrors the NVUE
// acl match/ip/protocol def, which is integer(0-255) | string enum
// (tcp,udp,...); collapsing it to integer alone wrongly rejects `protocol: tcp`.
func TestScalarUnion_NestedInSingleBranchWrapper_Preserved(t *testing.T) {
	min0, max255 := 0.0, 255.0
	protocol := &Config{
		Description: "IP protocol",
		Nullable:    true,
		AnyOf: []*Config{ // single-branch wrapper around the $ref target
			{
				AnyOf: []*Config{ // the referenced scalar union
					{Type: "integer", Minimum: &min0, Maximum: &max255, Nullable: true},
					{Type: "string", Enum: []any{"tcp", "udp", "icmp"}, Nullable: true},
				},
			},
		},
	}

	out := protocol.ToJSONSchema()

	// Must not have collapsed to a lone integer scalar.
	if _, collapsed := out["minimum"]; collapsed {
		t.Fatalf("union collapsed to a top-level integer scalar: %#v", out)
	}
	variants, ok := out["anyOf"].([]map[string]any)
	if !ok || len(variants) != 2 {
		t.Fatalf("expected anyOf with 2 variants, got %#v", out["anyOf"])
	}

	blob, _ := json.Marshal(out)
	s := string(blob)
	if !strings.Contains(s, "integer") {
		t.Error("integer branch lost")
	}
	if !strings.Contains(s, "\"tcp\"") {
		t.Error("string-enum branch (tcp) lost — this is the bug")
	}
	// Integer branch keeps its numeric bounds.
	if !strings.Contains(s, "255") {
		t.Error("integer branch dropped its maximum constraint")
	}
}
