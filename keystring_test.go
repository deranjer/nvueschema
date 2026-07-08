package nvueschema

import "testing"

// A key-string is an SSH public key, not a secret. It must not inherit
// secret-string's 64-char cap: an ed25519 public key is 68 base64 chars, and
// rsa keys are far longer. It maps to its own unbounded string def.
func TestKeyStringFormat_NoMaxLength(t *testing.T) {
	got := (&Config{Type: "string", Format: "key-string"}).ToJSONSchema()
	if ref, _ := got["$ref"].(string); ref != "#/$defs/key-string" {
		t.Fatalf("key-string $ref = %q, want #/$defs/key-string", ref)
	}

	defs := formatDefs()
	kdef, ok := defs["key-string"].(map[string]any)
	if !ok {
		t.Fatal("formatDefs missing key-string")
	}
	if _, capped := kdef["maxLength"]; capped {
		t.Errorf("key-string def must not set maxLength, got %#v", kdef)
	}
	if kdef["type"] != "string" {
		t.Errorf("key-string def type = %v, want string", kdef["type"])
	}

	// secret-string keeps its cap (unchanged).
	if sdef, _ := defs["secret-string"].(map[string]any); sdef["maxLength"] != 64 {
		t.Errorf("secret-string maxLength = %v, want 64", sdef["maxLength"])
	}
}
