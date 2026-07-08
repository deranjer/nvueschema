package nvueschema

import (
	"cmp"
	"slices"
	"strings"
)

// propertyEntry holds a property name and its schema, used when iterating
// over sorted properties in output generators.
type propertyEntry struct {
	name   string
	schema *Config
}

// sortedProperties returns the properties of a schema as a sorted slice.
func sortedProperties(s *Config) []propertyEntry {
	if s == nil || s.Properties == nil {
		return nil
	}
	props := make([]propertyEntry, 0, len(s.Properties))
	for k, v := range s.Properties {
		props = append(props, propertyEntry{k, v})
	}
	slices.SortFunc(props, func(a, b propertyEntry) int {
		return cmp.Compare(a.name, b.name)
	})
	return props
}

// scalarUnionVariants returns the anyOf or oneOf variants for a scalar union,
// preferring anyOf. Returns nil if neither is set.
//
// A single-branch anyOf/oneOf is the OpenAPI idiom for attaching nullable or a
// description to a $ref (e.g. `{"anyOf":[{"$ref":...}],"nullable":true}`); it
// is a passthrough, not a union. When that single branch is itself a scalar
// union, its variants are surfaced so the union is preserved rather than
// collapsed to its first branch by FlattenComposite.
func scalarUnionVariants(s *Config) []*Config {
	variants := s.AnyOf
	if len(variants) == 0 {
		variants = s.OneOf
	}
	if len(variants) == 1 {
		if inner := scalarUnionVariants(variants[0]); len(inner) > 1 {
			return inner
		}
	}
	return variants
}

// splitIdentifier splits a string on '-', '_', and '.' delimiters.
func splitIdentifier(s string) []string {
	return strings.FieldsFunc(s, func(c rune) bool {
		return c == '-' || c == '_' || c == '.'
	})
}

// toPascal converts a kebab-case or snake_case name to PascalCase.
func toPascal(s string) string {
	parts := splitIdentifier(s)
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}
