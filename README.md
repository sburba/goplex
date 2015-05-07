goplex
===============

(Small subset of the) Plex api for Go

## Install

```go get github.com/sburba/goplex```

## Usage

// Get user information
user, err := goplex.GetUser(USERNAME, PASSWORD)
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

// Get sessions for given server
sessions, err := servers[0].GetSessions()
if err != nil {
	log.Fatal("GetSessions: ", err)
}
fmt.Printf("%v\n\n", sessions)

## Install
go get github.com/sburba/goplex