/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package resourceid

type ResourceID uint32

type ResourceNameLookup interface {
	LookupResourceId(resourceName string) ResourceID
	SortedResourceNames() []string
}

type ResourceNameLookupImpl struct {
	lookup map[string]ResourceID
	stored []string
}

func NewResourceNameLookupImpl() *ResourceNameLookupImpl {
	return &ResourceNameLookupImpl{
		lookup: make(map[string]ResourceID),
	}
}

func (r *ResourceNameLookupImpl) LookupResourceId(resourceName string) ResourceID {
	existingResourceID, hasID := r.lookup[resourceName]
	if hasID {
		return existingResourceID
	}

	newResourceID := ResourceID(len(r.lookup))
	r.lookup[resourceName] = newResourceID
	r.stored = append(r.stored, resourceName)

	return newResourceID
}

func (r *ResourceNameLookupImpl) SortedResourceNames() []string {
	return r.stored
}
