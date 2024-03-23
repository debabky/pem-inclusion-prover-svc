/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type InclusionProof struct {
	Key
	Attributes InclusionProofAttributes `json:"attributes"`
}
type InclusionProofResponse struct {
	Data     InclusionProof `json:"data"`
	Included Included       `json:"included"`
}

type InclusionProofListResponse struct {
	Data     []InclusionProof `json:"data"`
	Included Included         `json:"included"`
	Links    *Links           `json:"links"`
}

// MustInclusionProof - returns InclusionProof from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustInclusionProof(key Key) *InclusionProof {
	var inclusionProof InclusionProof
	if c.tryFindEntry(key, &inclusionProof) {
		return &inclusionProof
	}
	return nil
}
