package handler

import (
	"context"
	"fmt"
	"os"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)


func accessToken() (string, error) {
	const scope = "bbb94529-53a3-4be5-a069-7eaf2712b826/.default"
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tenant := os.Getenv("TENANT")

	cred, err := confidential.NewCredFromSecret(clientSecret)
	if err != nil {
		return "", err
	}

	auth := fmt.Sprintf("https://login.microsoftonline.com/%s", tenant)
	cli, err := confidential.New(clientID, cred, confidential.WithAuthority(auth))
	if err != nil {
		return "", err
	}

	res, err := cli.AcquireTokenByCredential(context.Background(), []string{scope})
	if err != nil {
		return "", err
	}

	return res.AccessToken, nil
}
