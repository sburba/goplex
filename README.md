plex
===============

(Small subset of the) Plex api for Go

## Install

```go get github.com/sburba/plex```

## Usage

// Get user information
user, err := plex.GetUser(USERNAME, PASSWORD)
if err != nil {
	log.Fatal("GetUser: ", err)
}
fmt.Printf("User: %v\n", user)

// Get servers for given user
servers, err := user.GetServers()
if err != nil {
	log.Fatal("GetServers: ", err)
}

if len(servers) == 0 {
	log.Fatal("Didn't find any servers!")
}

// Get currently watched videos for given server
videos, err := servers[0].GetActivity()
if err != nil {
	log.Fatal("GetActivity: ", err)
}
fmt.Printf("%v\n\n", videos)

## Install
go get github.com/sburba/plex