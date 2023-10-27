/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type Blob struct {
	Key
	Attributes    BlobAttributes     `json:"attributes"`
	Relationships *BlobRelationships `json:"relationships,omitempty"`
}
type BlobResponse struct {
	Data     Blob     `json:"data"`
	Included Included `json:"included"`
}

func (b BlobResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Data     Blob     `json:"data"`
			Included Included `json:"included"`
		}{
			Data:     b.Data,
			Included: b.Included,
		})
}

type BlobListResponse struct {
	Data     []Blob   `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustBlob - returns Blob from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustBlob(key Key) *Blob {
	var blob Blob
	if c.tryFindEntry(key, &blob) {
		return &blob
	}
	return nil
}
