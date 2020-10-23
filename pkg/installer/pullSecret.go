package installer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const redHatSSOURL = "https://sso.redhat.com/auth/realms/redhat-external/protocol/openid-connect/token"
const ocpApi = "https://api.openshift.com/api/accounts_mgmt/v1/access_token"

type token struct {
	AccessToken    string `json:"access_token"`
	Expires        int    `json:"expires_in"`
	RefreshExpires int    `json:"refresh_expires_in"`
	RefreshToken   string `json:"refresh_token"`
	TokenType      string `json:"token_type"`
	IdToken        string `json:"id_token"`
	NotBefore      int    `json:"not-before-policy"`
	SessionState   string `json:"session_state"`
	Scope          string `json:"scope"`
}

func GetPullSecret() string {

	var Token token

	fmt.Println(`Please provide a Red Hat credentials (the account used to login to the access.redhat.com portal) to get the pull-secret.
This pull-secret will be used during the installation to pull the container images from Red Hat registries.`)
	// Gather the username from a Red Hat account
	fmt.Print("Username: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')

	// Gather the password
	fmt.Print("Password: ")
	password, _ := terminal.ReadPassword(0)
	fmt.Println()

	client := &http.Client{}

	// Make a request to the Red Hat SSO to gather a token that will be used to the cloud openshift api
	r := strings.NewReader("scope=openid&username=" + username + "&password=" + string(password) + "&grant_type=password&client_id=admin-cli")
	request, err := http.NewRequest("POST", redHatSSOURL, r)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Fatalln(err)
	}

	// make the request
	resp, err := client.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Fatalln("ERROR Error trying to gather sso token: ", resp.Status)
	}

	// unmarshal the response to the token struct
	json.Unmarshal(body, &Token)
	accessToken := Token.AccessToken

	// make a request to gather the pull-secret from the cloud openshift api
	request, err = http.NewRequest("POST", ocpApi, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err = client.Do(request)
	body, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Fatalln("ERROR Error trying to gather the OpenShift pull-secret: ", resp.Status)
	}

	// create a file with the pull-secret at /tmp
	f, err := os.Create("/tmp/pull-secret.txt")
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	_, err = f.Write(body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("INFO Here is the pull-secret (there is also a copy of it on /tmp/pull-secret.txt): ")
	return string(body)
}
