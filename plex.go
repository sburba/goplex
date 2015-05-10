package goplex

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

const plexTVURL = "https://plex.tv"
const clientIdentifier = "plextrack"

type devicesResp struct {
	XMLName       xml.Name `xml:"MediaContainer"`
	PublicAddress HTTPURL  `xml:"publicAddress,attr"`
	Devices       []Device `xml:"Device"`
}

type sessionsResp struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Videos  []Video  `xml:"Video"`
}

// Hook to override for tests
var client = http.DefaultClient

func GetUser(username, password string) (User, error) {
	req, err := http.NewRequest("POST", plexTVURL+"/users/sign_in.xml", nil)
	if err != nil {
		return User{}, err
	}
	req.SetBasicAuth(username, password)
	req.Header.Add("X-Plex-Client-Identifier", clientIdentifier)

	resp, err := fetchContent(req, http.StatusCreated)
	if err != nil {
		return User{}, err
	}

	user := User{}
	err = xml.Unmarshal(resp, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (user User) GetDevices() ([]Device, error) {
	req, err := http.NewRequest("GET", plexTVURL+"/devices.xml", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Plex-Client-Identifier", clientIdentifier)
	req.Header.Add("X-Plex-Token", user.AuthToken)

	content, err := fetchContent(req, http.StatusOK)

	resp := &devicesResp{}

	err = xml.Unmarshal(content, resp)
	if err != nil {
		return nil, err
	}

	for i := range resp.Devices {
		resp.Devices[i].Owner = user
	}

	return resp.Devices, nil
}

func (user User) GetServers() ([]Server, error) {
	devices, err := user.GetDevices()
	if err != nil {
		return nil, err
	}

	var servers []Server
	for _, device := range devices {
		server, err := device.toServer()
		if err == nil {
			servers = append(servers, server)
		}
	}

	return servers, nil
}

func (server Server) GetActivity() ([]Video, error) {
	server.PublicAddress.Path = "/status/sessions"

	req, err := http.NewRequest("GET", server.PublicAddress.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Plex-Client-Identifier", clientIdentifier)
	req.Header.Add("X-Plex-Token", server.Owner.AuthToken)

	resp, err := fetchContent(req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	container := &sessionsResp{}
	if err := xml.Unmarshal(resp, container); err != nil {
		return nil, err
	}

	return container.Videos, nil
}

func fetchContent(req *http.Request, expectedStatusCode int) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatusCode {
		return nil, errors.New("Received status: " + strconv.Itoa(resp.StatusCode) +
			" expected status: " + strconv.Itoa(expectedStatusCode))
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
