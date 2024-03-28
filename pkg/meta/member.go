package meta

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/goccy/go-yaml"
)

// MemberTo is part of the meta model.
type MemberTo struct {
	Collections []string `yaml:"collections"`
	Field       string   `yaml:"field"`
}

// Member is part of the meta model.
type Member struct {
	Type                  string              `yaml:"type"`
	Description           string              `yaml:"description"`
	To                    *MemberTo           `yaml:"to"`
	Fields                *Fields             `yaml:"fields"`
	ReplacementCollection string              `yaml:"replacement_collection"`
	ReplacementEnum       []string            `yaml:"replacement_enum"`
	RestrictionMode       string              `yaml:"restriction_mode"`
	Required              bool                `yaml:"required"`
	Searchable            bool                `yaml:"-"`
	Analyzer              *string             `yaml:"-"`
	Relation              *CollectionRelation `yaml:"-"`
	Order                 int32               `yaml:"-"`
}

// UnmarshalYAML Parses yaml to MemberTo
func (mt *MemberTo) UnmarshalYAML(node []byte) error {
	// 1. string
	var s string
	if err := yaml.Unmarshal(node, &s); err == nil {
		mt.Field = s
		return nil
	}

	// 2. List of strings
	var collections []string
	if err := yaml.Unmarshal(node, &collections); err == nil {
		mt.Collections = collections
		return nil
	}

	// 3. struct
	var memberTo struct {
		Collections []string `yaml:"collections"`
		Field       string   `yaml:"field"`
	}
	if err := yaml.Unmarshal(node, &memberTo); err != nil {
		return fmt.Errorf("memberTo object without field: %w", err)
	}
	mt.Field = memberTo.Field
	mt.Collections = memberTo.Collections
	return nil
}

// UnmarshalYAML Parses yaml to Member
func (m *Member) UnmarshalYAML(node []byte) error {
	if node[0] == byte('&') {
		return nil
	}

	m.Order = fieldNum.Add(1)
	var s string
	if err := yaml.Unmarshal(node, &s); err == nil {
		m.Type = s
		return nil
	}
	var member struct {
		Type                  string    `yaml:"type"`
		Description           string    `yaml:"description"`
		To                    *MemberTo `yaml:"to"`
		Fields                *Fields   `yaml:"fields"`
		ReplacementCollection string    `yaml:"replacement_collection"`
		ReplacementEnum       []string  `yaml:"replacement_enum"`
		RestrictionMode       string    `yaml:"restriction_mode"`
		Required              bool      `yaml:"required"`
	}
	if err := yaml.Unmarshal(node, &member); err != nil {
		return fmt.Errorf("member object without type: %w", err)
	}
	m.Type = member.Type
	m.Description = member.Description
	m.To = member.To
	m.Fields = member.Fields
	m.ReplacementCollection = member.ReplacementCollection
	m.ReplacementEnum = member.ReplacementEnum
	m.RestrictionMode = member.RestrictionMode
	m.Required = member.Required
	return nil
}

// Clone returns a deep copy.
func (mt *MemberTo) Clone() *MemberTo {
	if mt == nil {
		return nil
	}
	return &MemberTo{
		Collections: copyStrings(mt.Collections),
		Field:       mt.Field,
	}
}

// Clone returns a deep copy.
func (m *Member) Clone() *Member {
	return &Member{
		Type:                  m.Type,
		Description:           m.Description,
		To:                    m.To.Clone(),
		Fields:                m.Fields.Clone(),
		ReplacementCollection: m.ReplacementCollection,
		ReplacementEnum:       copyStrings(m.ReplacementEnum),
		RestrictionMode:       m.RestrictionMode,
		Required:              m.Required,
		Order:                 m.Order,
	}
}

// RetainStrings returns a function which keeps string type fields in [Retain].
func RetainStrings() func(string, string, *Member) bool {
	return func(k, fk string, f *Member) bool {
		switch f.Type {
		case "string", "HTMLStrict", "text", "HTMLPermissive":
			return true
		default:
			log.Tracef("removing non-string %s.%s\n", k, fk)
			return false
		}
	}
}
