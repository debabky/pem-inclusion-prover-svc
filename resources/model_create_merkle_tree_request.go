/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateMerkleTreeRequest struct {
	Key
	Attributes CreateMerkleTreeRequestAttributes `json:"attributes"`
}
type CreateMerkleTreeRequestResponse struct {
	Data     CreateMerkleTreeRequest `json:"data"`
	Included Included                `json:"included"`
}

type CreateMerkleTreeRequestListResponse struct {
	Data     []CreateMerkleTreeRequest `json:"data"`
	Included Included                  `json:"included"`
	Links    *Links                    `json:"links"`
}

// MustCreateMerkleTreeRequest - returns CreateMerkleTreeRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateMerkleTreeRequest(key Key) *CreateMerkleTreeRequest {
	var createMerkleTreeRequest CreateMerkleTreeRequest
	if c.tryFindEntry(key, &createMerkleTreeRequest) {
		return &createMerkleTreeRequest
	}
	return nil
}
