/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type MerkleRootHash struct {
	Key
	Attributes MerkleRootHashAttributes `json:"attributes"`
}
type MerkleRootHashResponse struct {
	Data     MerkleRootHash `json:"data"`
	Included Included       `json:"included"`
}

type MerkleRootHashListResponse struct {
	Data     []MerkleRootHash `json:"data"`
	Included Included         `json:"included"`
	Links    *Links           `json:"links"`
}

// MustMerkleRootHash - returns MerkleRootHash from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustMerkleRootHash(key Key) *MerkleRootHash {
	var merkleRootHash MerkleRootHash
	if c.tryFindEntry(key, &merkleRootHash) {
		return &merkleRootHash
	}
	return nil
}
