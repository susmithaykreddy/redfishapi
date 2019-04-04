package redfishapi

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func queryData(c *IloClient, call string, link string) ([]byte, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(call, link, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(c.Username, c.Password))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		r, _ := regexp.Compile("dial tcp")
		if r.MatchString(err.Error()) == true {
			err := errors.New(StatusInternalServerError)
			return nil, err
		} else {
			return nil, err
		}
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 {
			err := errors.New(StatusUnauthorized)
			return nil, err
		}

	}
	defer resp.Body.Close()

	_body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return _body, nil
}
