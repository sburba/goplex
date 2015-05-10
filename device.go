package goplex

import (
	"fmt"
	"net/url"
	"strings"
)

type Device struct {
	Name          string              `xml:"name,attr"`
	PublicAddress HttpUrl             `xml:"publicAddress,attr"`
	Product       string              `xml:"product,attr"`
	Provides      CommaSeperatedSlice `xml:"provides,attr"`
	Connections   []Connection        `xml:"Connection"`
	Owner         User
}

type Connection struct {
	Address HttpUrl `xml:"uri,attr"`
}

type Server struct {
	Device
}

func (device *Device) ProvidesFeature(feature string) bool {
	for _, providedFeature := range device.Provides {
		if providedFeature == feature {
			return true
		}
	}
	return false
}

func (device Device) toServer() (Server, error) {
	if !device.ProvidesFeature("server") {
		return Server{}, fmt.Errorf("Device %s is not a server", device.Name)
	}

	server := Server{device}
	// The public address for a server is missing the port, but the connection that matches has it
	fullPublicAddr, err := findMatchingConnection(server.PublicAddress.URL, server.Connections)
	if err == nil {
		server.PublicAddress = fullPublicAddr
	}
	return server, nil
}

func findMatchingConnection(address url.URL, connections []Connection) (HttpUrl, error) {
	for _, connection := range connections {
		if strings.Contains(connection.Address.Host, address.Host) {
			return connection.Address, nil
		}
	}
	return HttpUrl{}, fmt.Errorf("No matching connection found for %s in %v", address, connections)
}
