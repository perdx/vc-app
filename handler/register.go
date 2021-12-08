package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// UserRegistration is the receiver for handling user onboarding requests.
type UserRegistration struct {
}

// TODO: Expand this to include anything needed for verifying and issuing a corresponding VC
type UserInfo struct {
	GivenName      string `json:"givenName"`
	FamilyName     string `json:"familyName"`
	BirthDate      string `json:"birthDate"`
	Sex            string `json:"sex"`
	Email          string `json:"email"`
	Mobile         string `json:"mobile"`
	City           string `json:"city"`
	Country        string `json:"country"`
	DriversLicense string `json:"drivers_license"`
	Password       string `json:"password"`
	Tenant         string `json:"tenant"`
}

func NewUserRegistration() *UserRegistration {
	return &UserRegistration{}
}

func (reg *UserRegistration) Create(w http.ResponseWriter, r *http.Request) {

	errRes := func(err error) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error %v", err)))
	}

	var u UserInfo
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		errRes(err)
		return
	}

	// TODO: Get Tenant from context

	// TODO: Proper validation
	if u.GivenName == "" {
		errRes(errors.New("given name is required"))
		return
	}
	if u.FamilyName == "" {
		errRes(errors.New("family name is required"))
		return
	}
	if u.BirthDate == "" {
		errRes(errors.New("birth date is required"))
		return
	}
	if u.Sex == "" {
		errRes(errors.New("sex is required"))
		return
	}
	if u.Email == "" {
		errRes(errors.New("email is required"))
		return
	}
	if u.Password == "" {
		errRes(errors.New("password is required"))
		return
	}
	if u.DriversLicense == "" {
		errRes(errors.New("drivers license is required"))
		return
	}
	// TODO: Store details and initiate validation process.

	type resBody struct {
		Id      string `json:"id"`
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	body := resBody{
		Id:      "da4e10b2-9209-427b-a64a-00b9e3fe2c3d",
		Message: fmt.Sprintf("User registration for %s %s received successfully", u.GivenName, u.FamilyName),
		Status:  "pending",
	}

	msg, err := json.Marshal(body)
	if err != nil {
		errRes(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

// Status returns the status of the registration process.
// TODO: Wire this up to a back end service.
func (reg *UserRegistration) Status(w http.ResponseWriter, r *http.Request) {

	// TODO: Implement proper authentication. Using basic auth implies checking
	// the password and id are a valid combination against a store of pending
	// onboard requests.
	auth := r.Header.Get("Authorization")
	if auth == "" {
		println("No authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokens := strings.Split(auth, " ")
	if tokens[0] != "Basic" {
		println("Invalid authorization header - no Basic keyword")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	pass, err := base64.StdEncoding.DecodeString(tokens[1])
	if err != nil {
		println("Invalid authorization header - invalid base64")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	println("pass:", string(pass))
	creds := strings.Split(string(pass), ":")
	if len(creds) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// TODO: creds[0] is the onboard request id created by the server
	//       creds[1] is the token created by the mobile app when the onboard request was made
	//       Match these against a store of pending onboard requests to ensure the requester of status is the same as
	//       the one who originally applied to onboard.

	// TODO: This hack is for demo purposes. Returns an "issued" status by default otherwise reflects query string.
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "issued"
	}

	type resBody struct {
		Id      string `json:"id"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	body := resBody{
		Id:      creds[0],
		Status:  status,
		Message: "",
	}

	if status == "declined" {
		body.Message = "The information you provided was unable to be used to validate your ID"
	}

	msg, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error %v", err)))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}
