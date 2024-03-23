/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Pem struct {
	Key
	Attributes PemAttributes `json:"attributes"`
}
type PemResponse struct {
	Data     Pem      `json:"data"`
	Included Included `json:"included"`
}

type PemListResponse struct {
	Data     []Pem    `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustPem - returns Pem from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustPem(key Key) *Pem {
	var pEM Pem
	if c.tryFindEntry(key, &pEM) {
		return &pEM
	}
	return nil
}
