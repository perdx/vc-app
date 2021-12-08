package handler

//-----------------------------------------
// Core structs
//-----------------------------------------
type Headers struct {
	APIKey string `json:"api-key"` // optional
}

type Callback struct { // used by the VC Request service
	URL     string  `json:"url"` // https://<ngrok uri>/api/issuer/issuanceCallback,
	State   string  `json:"state"`
	Headers Headers `json:"headers"`
}

type Registration struct {
	ClientName string `json:"clientName"`
	Purpose    string `json:"purpose"`
}

type RequestCore struct {
	IncludeQRCode bool         `json:"includeQRCode"`
	Callback      Callback     `json:"callback"`
	Authority     string       `json:"authority"` // <issuer DID>
	Registration  Registration `json:"registration"`
}

type ResponseCore struct {
	RequestID string `json:"requestId"`
	URL       string `json:"url"`
	Expiry    int    `json:"expiry"`
	QRCode    string `json:"qrCode"`
	State     string `json:"state"`
}

//-----------------------------------------
// Issuance structs
//-----------------------------------------
type PIN struct {
	Value  string `json:"value"`
	Length int    `json:"length"`
}

type Claims struct {
	GivenName      string `json:"given_name"`
	FamilyName     string `json:"family_name"`
	BirthDate      string `json:"birth_date"`
	Sex            string `json:"sex"`
	Email          string `json:"email"`
	Mobile         string `json:"mobile"`
	City           string `json:"city"`
	Country        string `json:"country"`
	DriversLicense string `json:"drivers_license"`
}

type Issuance struct {
	Type     string `json:"type"` // defined in rules.json
	Manifest string `json:"manifest"`
	PIN      PIN    `json:"pin"`
	Claims   Claims `json:"claims"` // payload from rules - rules map ID Token (hint) fields to VC
}

type IssuanceRequest struct {
	RequestCore
	Issuance Issuance `json:"issuance"`
}

//-----------------------------------------
// Presentation structs
//-----------------------------------------
type RequestedCredentials struct {
	Type            string   `json:"type"`
	Purpose         string   `json:"purpose"`
	AcceptedIssuers []string `json:"acceptedIssuers"`
}

type Presentation struct {
	IncludeReceipt       bool                   `json:"includeReceipt"`
	RequestedCredentials []RequestedCredentials `json:"requestedCredentials"`
}

type PresentationRequest struct {
	RequestCore
	Presentation Presentation `json:"presentation"`
}
