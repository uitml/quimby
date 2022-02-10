package user

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func basicAuth(username string, token string) string {
	auth := username + ":" + token
	return base64.StdEncoding.EncodeToString([]byte(auth))
	//return auth
}

// Reads default user config from a specified github repo. This can be private, so authentication is needed.
func GetDefaultConfig(repo string, file string, username string, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+repo+"/contents/"+file, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(username, token))
	req.Header.Add("Accept", "Accept: application/vnd.github.v3+json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Hack. We only need the download URL
	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}
	resp, err = http.Get(fmt.Sprint(m["download_url"]))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Finally, read and return
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
