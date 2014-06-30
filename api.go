package main

import (
	"flag"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	"github.com/skurtzemann/go-openvpn-api/vpn"
	"io/ioutil"
	"os"
)

const (
	ApiName    = "go-openvpn-api"
	ApiVersion = "1.0.0"
)

// General function which walk into the config directory and process the given callback
func EachConfig(directory string, callback func(os.FileInfo) bool) (err error) {
	if files, err := ioutil.ReadDir(directory); err == nil {
		for _, file := range files {
			if !file.IsDir() {
				if !callback(file) {
					break
				}
			}
		}
	}
	return err
}

// Returns a list of files into the OpenVPN client config dir
func ListConfigNames(directory string) (users []string, err error) {
	err = EachConfig(directory, func(file os.FileInfo) bool {
		users = append(users, file.Name())
		return true
	})
	return
}

// Returns a list of the all users informations
func ListConfigs(directory string) (users []vpn.VpnUser, err error) {
	err = EachConfig(directory, func(file os.FileInfo) bool {
		user := vpn.VpnUser{file.Name(), true, "", ""}
		if err = user.ParseConfigFile(directory + "/" + file.Name()); nil == err {
			users = append(users, user)
		}
		return nil == err
	})
	return
}

func main() {
	// params
	ccdDir := flag.String("ccddir", "./ccd", "openvpn's client config dir")
	flag.Parse()

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		IndentJSON: true,
	}))

	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
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
	m.Get("/v1/users", func(r render.Render) {
		users, err := ListConfigNames(*ccdDir)

		if err != nil {
			r.JSON(404, map[string]string{
				"error": "OpenVPN client configuration directory not found",
			})
		} else {
			r.JSON(200, users)
		}
	})

	// Get all users with the full details of them
	m.Get("/v1/users/_full", func(r render.Render) {
		users, err := ListConfigs(*ccdDir)

		if err != nil {
			r.JSON(404, map[string]string{
				"error": "OpenVPN client configuration directory not found",
			})
		} else {
			r.JSON(200, users)
		}
	})

	// Get the configuration of the given user
	m.Get("/v1/users/:user", func(r render.Render, params martini.Params) {
		userConfigFile := *ccdDir + "/" + params["user"]

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
