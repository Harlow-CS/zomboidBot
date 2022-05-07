package bot

import (
	"os"
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
	"github.com/Harlow-CS/zomboidBot/zomboid"
)

var (
	integerOptionMinValue = 1.0
	dmPermission = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer
	serverName = os.Getenv("server_name")

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "server",
			Description: "Ping or start server",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "operation",
					Description: "(ping, start)",
					Required:    true,
				},
			},
		},
		{
			Name:        "server-settings",
			Description: "Edit or list current server settings",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "operation",
					Description: "(edit, ls)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "new-settings",
					Description: "comma-separated key-value pairs that you want updated (edit only)",
					Required:    false,
				},
			},
		},
		{
			Name:        "sandbox-settings",
			Description: "list current server's sandbox settings",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"server": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			// Get the operation
			option, _ := optionMap["operation"]
			operation := option.StringValue()
	
			acknowledge(s, i, fmt.Sprintf("Processing settings request operation: '%s'", operation))

			if (operation == "ping") {

				if (zomboid.IsServerActive()) {
					s.ChannelMessageSend(ChannelID, "Server is offline")
				} else {
					s.ChannelMessageSend(ChannelID, "Server is running")
				}

			} else if (operation == "start") {

				if (zomboid.IsServerActive()) {
					s.ChannelMessageSend(ChannelID, "Server is already running")
					return
				}

				zomboid.StartServer()

				s.ChannelMessageSend(ChannelID, "Server is starting, wait ~5 minutes")
			} else {
				s.ChannelMessageSend(ChannelID, fmt.Sprintf("Unknown operation: '%s'", operation))
			}
		},
		"server-settings": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			// Get the operation
			option, _ := optionMap["operation"]
			operation := option.StringValue()
	
			acknowledge(s, i, fmt.Sprintf("Processing settings request operation: '%s'", operation))

			if (operation == "ls") {
				config := zomboid.GetServerConfig("server-name")

				// create tmp dir if it doesn't exist
				os.MkdirAll("./tmp", os.ModePerm)
				out, e := os.Create("./tmp/server-config.txt")
				if e != nil {
					panic(e)
				}
				defer out.Close()
				// write config to tmp file
				out.WriteString(config)

				// create tmp file, attach it, send it
				file := &discordgo.File{
					ContentType: "text/plain",
					Name: "server-config.txt",
					Reader: strings.NewReader(config),
				}

				msgSend := &discordgo.MessageSend{
					Files: []*discordgo.File{},
				}
				msgSend.Files = append(msgSend.Files, file)

				s.ChannelMessageSendComplex(ChannelID, msgSend)
			} else if (operation == "edit") {

				newSettings := ""
				if option, ok := optionMap["new-settings"]; ok {
					newSettings = option.StringValue()
				} else {
					s.ChannelMessageSend(ChannelID, "Parameter `new-settings` is required for edit operations")
					return
				}

				zomboid.UpdateServerConfig("server-name", newSettings)

				s.ChannelMessageSend(ChannelID, "Saved new server config")
			} else {
				s.ChannelMessageSend(ChannelID, fmt.Sprintf("Unknown operation: '%s'", operation))
			}
		},
		"sandbox-settings": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			acknowledge(s, i, "Processing sandbox config list...")

			config := zomboid.GetSandboxConfig("server-name")

			// create tmp dir if it doesn't exist
			os.MkdirAll("./tmp", os.ModePerm)
			out, e := os.Create("./tmp/sandbox-config.txt")
			if e != nil {
				panic(e)
			}
			defer out.Close()
			// write config to tmp file
			out.WriteString(config)

			// create tmp file, attach it, send it
			file := &discordgo.File{
				ContentType: "text/plain",
				Name: "sandbox-config.txt",
				Reader: strings.NewReader(config),
			}

			msgSend := &discordgo.MessageSend{
				Files: []*discordgo.File{},
			}
			msgSend.Files = append(msgSend.Files, file)

			s.ChannelMessageSendComplex(ChannelID, msgSend)
		},
	}
)

func acknowledge(s *discordgo.Session, i *discordgo.InteractionCreate, text string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: text,
		},
	})
}