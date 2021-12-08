package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type VerifierStatus struct {
	Status         string
	Message        string
	Payload        string
	Subject        string
	GivenName      string
	FamilyName     string
	BirthDate      string
	Sex            string
	Email          string
	Mobile         string
	City           string
	Country        string
	DriversLicense string
}
type Verifier struct {
	IssuerAuthority   string
	VerifierAuthority string
	Tenant            string
	callbackStatus    map[string]VerifierStatus
}

func NewVerifier() *Verifier {
	return &Verifier{
		IssuerAuthority:   os.Getenv("AUTHORITY"),
		VerifierAuthority: os.Getenv("AUTHORITY"),
		Tenant:            os.Getenv("TENANT"),
		callbackStatus:    make(map[string]VerifierStatus),
	}
}

func (v *Verifier) Present(w http.ResponseWriter, r *http.Request) {
	const vcType = "ATIdentity"
	//host := r.Header.Get("X-Forwarded-Proto") + "://" + r.Header.Get("X-Forwarded-Host")
	host := "https://" + r.Host

	req := PresentationRequest{
		RequestCore: RequestCore{
			IncludeQRCode: true,
			Callback: Callback{
				URL:   host + "/verifier/callback",
				State: uuid.New().String(),
			},
			Authority: v.VerifierAuthority,
			Registration: Registration{
				ClientName: "Verifiable Credential Expert Verifier",
				Purpose:    "So we can see that you a verifiable credentials expert",
			},
		},
		Presentation: Presentation{
			IncludeReceipt: true,
			RequestedCredentials: []RequestedCredentials{
				{
					Type:            vcType,
					Purpose:         "So we can see that you a verifiable credentials expert",
					AcceptedIssuers: []string{v.IssuerAuthority},
				},
			},
		},
	}

	byt, err := json.Marshal(req)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %v", err)))
		return
	}

	tkn, err := accessToken()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %v", err)))
		return
	}

	url := fmt.Sprintf("https://beta.did.msidentity.com/v1.0/%s/verifiablecredentials/request", v.Tenant)
	hreq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(byt))
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %v", err)))
		return
	}
	hreq.Header.Add("content-type", "application/json")
	hreq.Header.Add("authorization", "Bearer "+tkn)

	res, err := http.DefaultClient.Do(hreq)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %v", err)))
		return
	}

	defer res.Body.Close()
	byt, err = ioutil.ReadAll(res.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %v", err)))
		return
	}

	if res.StatusCode != http.StatusCreated {
		w.Write([]byte(fmt.Sprintf("error %v", err)))
		return
	}

	// unpack response and add state id for later querying of callback status
	var rsp ResponseCore

	json.Unmarshal(byt, &rsp)
	rsp.State = req.Callback.State

	rspByt, err := json.Marshal(rsp)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %v", err)))
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(rspByt)
}

func (v *Verifier) Callback(w http.ResponseWriter, r *http.Request) {

	var data struct {
		State   string
		Code    string
		Subject string
		Issuers []struct {
			Claims struct {
				GivenName      string
				FamilyName     string
				BirthDate      string
				Sex            string
				Email          string
				Mobile         string
				City           string
				Country        string
				DriversLicense string
			}
		}
	}

	defer r.Body.Close()
	byt, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %v", err)))
		return
	}

	if err := json.Unmarshal(byt, &data); err != nil {
		w.Write([]byte(fmt.Sprintf("error %v", err)))
		return
	}
	// fmt.Print(string(byt))

	status := VerifierStatus{
		Status: data.Code,
	}

	switch data.Code {
	case "request_retrieved":
		status.Message = "QR Code is scanned. Waiting for validation..."
	case "presentation_verified":
		status.Message = "Presentation verified"
		// status.Payload =  data.Issuers
		status.Subject = data.Subject
		status.GivenName = data.Issuers[0].Claims.GivenName
		status.FamilyName = data.Issuers[0].Claims.FamilyName
		status.BirthDate = data.Issuers[0].Claims.BirthDate
		status.Sex = data.Issuers[0].Claims.Sex
		status.Email = data.Issuers[0].Claims.Email
		status.Mobile = data.Issuers[0].Claims.Mobile
		status.City = data.Issuers[0].Claims.City
		status.Country = data.Issuers[0].Claims.Country
		status.DriversLicense = data.Issuers[0].Claims.DriversLicense

	}

	v.callbackStatus[data.State] = status
}

func (v *Verifier) Status(w http.ResponseWriter, r *http.Request) {
	state := chi.URLParam(r, "state")
	status := v.callbackStatus[state]

	rspByt, err := json.Marshal(status)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %v", err)))
		return
	}

	w.Write(rspByt)
}
