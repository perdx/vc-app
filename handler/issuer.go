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

type IssuerStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Payload string `json:"payload"`
}

type Issuer struct {
	IssuerAuthority string
	Tenant          string
	callbackStatus  map[string]IssuerStatus
}

func NewIssuer() *Issuer {
	return &Issuer{
		IssuerAuthority: os.Getenv("AUTHORITY"),
		Tenant:          os.Getenv("TENANT"),
		callbackStatus:  make(map[string]IssuerStatus),
	}
}
func (i *Issuer) Issue(w http.ResponseWriter, r *http.Request) {

	// TODO: retrieve from request
	const vcType = "ATIdentity"
	host := "https://" + r.Host

	// VC schema for vcType
	mfst := fmt.Sprintf("https://beta.did.msidentity.com/v1.0/%s/verifiableCredential/contracts/%s", i.Tenant, vcType)

	req := IssuanceRequest{
		RequestCore: RequestCore{
			IncludeQRCode: true,
			Callback: Callback{
				URL:   host + "/issuer/callback",
				State: uuid.New().String(),
			},
			Authority: i.IssuerAuthority,
			Registration: Registration{
				ClientName: "Verifiable Credential Sample",
			},
		},
		Issuance: Issuance{
			Type:     vcType,
			Manifest: mfst,
			PIN: PIN{
				Value:  "1234",
				Length: 4,
			},
			// TODO: retrieve data for authenticated user (from JWT?)
			Claims: Claims{
				GivenName:      "Bob",
				FamilyName:     "Ross",
				BirthDate:      "07/07/1907",
				Sex:            "Male",
				Email:          "bob.ross@perdx.io",
				Mobile:         "89632786496",
				City:           "Auckland",
				Country:        "New Zealand",
				DriversLicense: "DJ0095438",
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

	url := fmt.Sprintf("https://beta.did.msidentity.com/v1.0/%s/verifiablecredentials/request", i.Tenant)
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
	status := IssuerStatus{
		Status:  "awaiting_issuance",
		Message: "Scan the QR code to receive your login credential",
	}
	i.callbackStatus[req.Callback.State] = status

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

func (i *Issuer) Callback(w http.ResponseWriter, r *http.Request) {

	var data struct {
		State string
		Code  string
		Error struct {
			Code    string
			Message string
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

	status := IssuerStatus{
		Status: data.Code,
	}

	switch data.Code {
	case "request_retrieved":
		status.Message = "QR Code is scanned. Waiting for issuance..."
	case "issuance_successful":
		status.Message = "You can now login using the AT Identity credential stored in Microsoft Authenticator"
	case "issuance_error":
		status.Message = data.Error.Message
		status.Payload = data.Error.Code
	}

	i.callbackStatus[data.State] = status
}

func (i *Issuer) Status(w http.ResponseWriter, r *http.Request) {
	state := chi.URLParam(r, "state")
	status := i.callbackStatus[state]

	rspByt, err := json.Marshal(status)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %v", err)))
		return
	}

	w.Write(rspByt)
}
