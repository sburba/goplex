plex
===============

(Small subset of the) Plex api for Go

## Install

```go get github.com/sburba/goplex```

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

	// Get videos being watched on the given server right now
	videos, err := servers[0].GetActivity()
	if err != nil {
		log.Fatal("GetActivity: ", err)
	}
	fmt.Printf("%v\n\n", videos)
