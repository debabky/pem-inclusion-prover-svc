/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type InclusionProofRequest struct {
	Key
	Attributes InclusionProofRequestAttributes `json:"attributes"`
}
type InclusionProofRequestResponse struct {
	Data     InclusionProofRequest `json:"data"`
	Included Included              `json:"included"`
}

type InclusionProofRequestListResponse struct {
	Data     []InclusionProofRequest `json:"data"`
	Included Included                `json:"included"`
	Links    *Links                  `json:"links"`
}

// MustInclusionProofRequest - returns InclusionProofRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustInclusionProofRequest(key Key) *InclusionProofRequest {
	var inclusionProofRequest InclusionProofRequest
	if c.tryFindEntry(key, &inclusionProofRequest) {
		return &inclusionProofRequest
	}
	return nil
}
