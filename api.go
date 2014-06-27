package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/skurtzemann/go-openvpn-api/vpn"
	"io/ioutil"
)

const (
	ApiName    = "go-openvpn-api"
	ApiVersion = "1.0.0"
)

// List the "openvpn client config dir" and return a slice of vpnUser
// we considers that a VpnUser is only a file not a directory
func listConfigDir(directory string) (users []string, err error) {

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			users = append(users, file.Name())
		}
	}
	return users, nil
}

func main() {
	// openvpn client config dir (ccd)
	ccdDir := "./ccd"

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		IndentJSON: true,
	}))

	// The default page return the name and the version of the api
	m.Get("/", func() string {
		return ApiName + " (" + ApiVersion + ")"
	})

	// Health of the API : for the moment return "true"
	m.Get("/_ping", func() string {
		return "true"
	})

	// Get all users
	m.Get("/users", func(r render.Render) {
		users, err := listConfigDir(ccdDir)

		if err != nil {
			r.JSON(404, map[string]string{
				"error": "OpenVPN client configuration directory not found",
			})
		} else {
			r.JSON(200, users)
		}
	})

	// Get the configuration of the given user
	m.Get("/users/:user", func(r render.Render, params martini.Params) {
		userConfigFile := ccdDir + "/" + params["user"]

		user := vpn.VpnUser{params["user"], true, "", ""}
		err := user.ParseConfigFile(userConfigFile)
		if err != nil {
			r.JSON(404, map[string]string{
				"error": "User retrieve error",
			})
		} else {
			r.JSON(200, user)
		}
	})

	m.Run()
}
