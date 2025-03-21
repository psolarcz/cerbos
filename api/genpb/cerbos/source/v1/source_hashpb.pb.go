// Code generated by protoc-gen-go-hashpb. Do not edit.
// protoc-gen-go-hashpb v0.2.0
// Source: cerbos/source/v1/source.proto

package sourcev1

import (
	hash "hash"
)

// HashPB computes a hash of the message using the given hash function
// The ignore set must contain fully-qualified field names (pkg.msg.field) that should be ignored from the hash
func (m *Position) HashPB(hasher hash.Hash, ignore map[string]struct{}) {
	if m != nil {
		cerbos_source_v1_Position_hashpb_sum(m, hasher, ignore)
	}
}

// HashPB computes a hash of the message using the given hash function
// The ignore set must contain fully-qualified field names (pkg.msg.field) that should be ignored from the hash
func (m *Error) HashPB(hasher hash.Hash, ignore map[string]struct{}) {
	if m != nil {
		cerbos_source_v1_Error_hashpb_sum(m, hasher, ignore)
	}
}

// HashPB computes a hash of the message using the given hash function
// The ignore set must contain fully-qualified field names (pkg.msg.field) that should be ignored from the hash
func (m *SourceContext) HashPB(hasher hash.Hash, ignore map[string]struct{}) {
	if m != nil {
		cerbos_source_v1_SourceContext_hashpb_sum(m, hasher, ignore)
	}
}
