package committee

// CommitteeResponse represents a committee in API responses
type CommitteeResponse struct {
	CommitteeID   string  `json:"committee_id" example:"CREATIVES"`
	CommitteeName string  `json:"committee_name" example:"Creatives Committee"`
	CommitteeHead *int32  `json:"committee_head,omitempty" example:"123"`
	DivisionID    *string `json:"division_id,omitempty" example:"INTERNALS"`
}

// GetAllCommitteesResponse is the response for the GET /committees endpoint
type GetAllCommitteesResponse struct {
	Committees []CommitteeResponse `json:"committees"`
}
