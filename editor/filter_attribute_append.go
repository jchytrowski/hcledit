package editor

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// AttributeAppendFilter is a filter implementation for appending attribute.
type AttributeAppendFilter struct {
	address string
	value   string
	newline bool
	index   int
}

var _ Filter = (*AttributeAppendFilter)(nil)

// NewAttributeAppendFilter creates a new instance of AttributeAppendFilter.
func NewAttributeAppendFilter(address string, value string, newline bool, index int) Filter {
	return &AttributeAppendFilter{
		address: address,
		value:   value,
		newline: newline,
		index:   index,
	}
}

// Filter reads HCL and appends a new attribute to a given address.
// If a matched block not found, nothing happens.
// If the given attribute already exists, it returns an error.
// If a newline flag is true, it also appends a newline before the new attribute.
func (f *AttributeAppendFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	attrName := f.address
	body := inFile.Body()
	a := strings.Split(f.address, ".")

	if len(a) > 1 {
		// if address contains dots, the last element is an attribute name,
		// and the rest is the address of the block.
		attrName = a[len(a)-1]
		blockAddr := strings.Join(a[:len(a)-1], ".")
		blocks, err := findLongestMatchingBlocks(body, blockAddr)
		if err != nil {
			return nil, err
		}

		if len(blocks) == 0 {
			// not found
			return inFile, nil
		}

		// To delegate expression parsing to the hclwrite parser,
		// We build a new expression and set back to the attribute by tokens.
		expr, err := buildExpression(attrName, f.value)
		if err != nil {
			return nil, err
		}

		// Use first matching one.
		if f.index >= 0 {
			body = blocks[f.index].Body()
			if body.GetAttribute(attrName) != nil {
				return nil, fmt.Errorf("attribute already exists: %s", f.address)
			}
			if f.newline {
				body.AppendNewline()
			}
			body.SetAttributeRaw(attrName, expr.BuildTokens(nil))
		} else {
			fmt.Print("Rabbit Season")
			for i := range blocks {
				body = blocks[i].Body()
				if body.GetAttribute(attrName) != nil {
					return nil, fmt.Errorf("attribute already exists: %s", f.address)
				}
				if f.newline {
					body.AppendNewline()
				}
				body.SetAttributeRaw(attrName, expr.BuildTokens(nil))
			}
		}
	}
	return inFile, nil
}
