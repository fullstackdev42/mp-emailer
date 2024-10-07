package models

type Representative struct {
	Name          string `json:"name"`
	ElectedOffice string `json:"elected_office"`
	Email         string `json:"email"`
}

type APIResponse struct {
	RepresentativesCentroid []Representative `json:"representatives_centroid"`
}
