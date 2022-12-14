package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/appclient"
	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"
)

//go:embed icon.png
var IconData []byte

var RootURL string = os.Getenv("MANIFEST_ROOT_URL")
var ServerPort string = os.Getenv("SERVER_PORT")

// Manifest declares the app's metadata. It must be provided for the app to be
// installable. In this example, the following permissions are requested:
//   - Create posts as a bot.
//   - Add icons to the channel header that will call back into your app when
//     clicked.
//   - Add a /-command with a callback.
var Manifest = apps.Manifest{
	// App ID must be unique across all Mattermost Apps.
	AppID: "welcome-bot",

	// App's release/version.
	Version: "v0.1.0",

	// A (long) display name for the app.
	DisplayName: "Welcome Bot",

	// The icon for the app's bot account, same icon is also used for bindings
	// and forms.
	Icon: "icon.png",

	// HomepageURL is required for an app to be installable.
	HomepageURL: "https://github.com/mattermost/mattermost-app-welcomebot",

	// Need ActAsBot to post back to the user.
	RequestedPermissions: []apps.Permission{
		apps.PermissionActAsBot,
		apps.PermissionActAsUser,
	},

	// Add UI elements: a /-command, and a channel header button.
	RequestedLocations: []apps.Location{
		apps.LocationChannelHeader,
		apps.LocationCommand,
	},

	// Running the app as an HTTP service is the only deployment option
	// supported.
	Deploy: apps.Deploy{
		HTTP: &apps.HTTP{
			RootURL: RootURL,
		},
	},
}

// The details for the App UI bindings
var Bindings = []apps.Binding{
	{
		Location: "/command",
		Bindings: []apps.Binding{
			{
				Icon:        "icon.png",
				Label:       "mybot",
				Description: "Welcome Bot app", // appears in autocomplete.
				Hint:        "[help|list|preview|set_channel_welcome|get_channel_welcome|delete_channel_welcome]",          // appears in autocomplete, usually indicates as to what comes after choosing the option.
				Bindings: []apps.Binding{
					{
						Label: "help", // displays usage information
						Submit: ShowHelp,
					},
					{
						Label: "list", // Lists the teams for which greetings were defined
						Submit:  ShowList,
					},
					{
						Label: "preview [team-name]", // Send ephemeral messages to the user
						Form:  &ShowPreviewForTeam,
					},
					{
						Label: "set_channel_welcome [team-name]", // Sets the given text as current's channel welcome message.
						Form:  &SetChannelWelcome,
					},
					{
						Label: "get_channel_welcome",  // Sets the current channel's welcome message
						Submit:  GetChannelWelcome,
					},
					{
						Label: "delete_channel_welcome",  // Deletes the current channel's welcome message.
						Submit:  DeleteChannelWelcome,
					},
				},
			},
		},
	},
}

var ShowPreviewForTeam = apps.Form{
	Title: "Welcome Bot",
	Icon:  "icon.png",
	Fields: []apps.Field{
		{
			Type: "text",
			Name: "Team Name",
		},
	},
	Submit: apps.NewCall("/preview").WithExpand(apps.Expand{ActingUserAccessToken: apps.ExpandAll}),
}

var SetChannelWelcome = apps.Form{
	Title: "Welcome Bot",
	Icon:  "icon.png",
	Fields: []apps.Field{
		{
			Type: "text",
			Name: "Team Name",
		},
	},
	Submit: apps.NewCall("/set_channel_welcome").WithExpand(apps.Expand{ActingUserAccessToken: apps.ExpandAll}),
}

var ShowHelp = apps.NewCall("/help").WithExpand(apps.Expand{ActingUserAccessToken: apps.ExpandAll})
var ShowList = apps.NewCall("/list")
var GetChannelWelcome = apps.NewCall("/get_channel_welcome")
var DeleteChannelWelcome = apps.NewCall("/delete_channel_welcome")

// main sets up the http server, with paths mapped for the static assets, the
// bindings callback, and the send function.
func main() {
	// Serve static assets: the manifest and the icon.
	http.HandleFunc("/manifest.json",
		httputils.DoHandleJSON(Manifest))
	http.HandleFunc("/static/icon.png",
		httputils.DoHandleData("image/png", IconData))

	// Bindinings callback.
	http.HandleFunc("/bindings",
		httputils.DoHandleJSON(apps.NewDataResponse(Bindings)))

	http.HandleFunc("/preview", PreviewCall)
	http.HandleFunc("/help", HelpCall)
	http.HandleFunc("/list", ListCall)
	http.HandleFunc("/set_channel_welcome", SetChannelWelcomeCall)
	http.HandleFunc("/get_channel_welcome", GetChannelWelcomeCall)
	http.HandleFunc("/delete_channel_welcome", DeleteChannelWelcomeCall)

	fmt.Printf("Use '/apps install http %s/manifest.json' to install the app\n", RootURL)
	log.Fatal(http.ListenAndServe(ServerPort, nil))
}

func HelpCall(w http.ResponseWriter, req *http.Request) {
	c := apps.CallRequest{}
	json.NewDecoder(req.Body).Decode(&c)

	message := "Help Call"
	appclient.AsBot(c.Context).DM(c.Context.ActingUser.Id, message)

	httputils.WriteJSON(w,
		apps.NewTextResponse("Shown Welcome Bot Help"))
}

func PreviewCall(w http.ResponseWriter, req *http.Request) {
	httputils.WriteJSON(w,
		apps.NewTextResponse("Shown Welcome Bot Preview"))
}

func ListCall(w http.ResponseWriter, req *http.Request) {
	httputils.WriteJSON(w,
		apps.NewTextResponse("Shown Welcome Bot List"))
}

func SetChannelWelcomeCall(w http.ResponseWriter, req *http.Request) {
	httputils.WriteJSON(w,
		apps.NewTextResponse("Shown Welcome Bot Set channel welcome"))
}

func GetChannelWelcomeCall(w http.ResponseWriter, req *http.Request) {
	httputils.WriteJSON(w,
		apps.NewTextResponse("Shown Welcome Bot Get Channel welcome"))
}

func DeleteChannelWelcomeCall(w http.ResponseWriter, req *http.Request) {
	httputils.WriteJSON(w,
		apps.NewTextResponse("Shown Welcome Bot Delete channel welcome"))
}