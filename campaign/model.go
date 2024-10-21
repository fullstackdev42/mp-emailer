package campaign

import "time"

type Campaign struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Template    string    `json:"template"`
	OwnerID     int       `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tokens      []string  `json:"tokens"`
}

type Representative struct {
	Name              string   `json:"name"`
	DistrictName      string   `json:"district_name"`
	ElectedOffice     string   `json:"elected_office"`
	FirstName         string   `json:"first_name"`
	LastName          string   `json:"last_name"`
	Party             string   `json:"party_name"`
	Email             string   `json:"email"`
	URL               string   `json:"url"`
	PersonalURL       string   `json:"personal_url"`
	PhotoURL          string   `json:"photo_url"`
	Gender            string   `json:"gender"`
	Offices           []Office `json:"offices"`
	Extra             Extra    `json:"extra"`
	RepresentativeSet string   `json:"representative_set_name"`
}

type Office struct {
	Fax    string `json:"fax"`
	Tel    string `json:"tel"`
	Type   string `json:"type"`
	Postal string `json:"postal"`
}

type Extra struct {
	Roles              []string `json:"roles"`
	PreferredLanguages []string `json:"preferred_languages"`
}

type APIResponse struct {
	RepresentativesCentroid []Representative `json:"representatives_centroid"`
}
