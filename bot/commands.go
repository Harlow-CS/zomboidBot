package bot

import (
	"os"
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
	"github.com/Harlow-CS/zomboidBot/zomboid"
)

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
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
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"server-settings": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			operation := ""

			// Get the value from the option map.
			if option, ok := optionMap["operation"]; ok {
				operation = option.StringValue()
			} else {
				acknowledge(s, i, fmt.Sprintf("Processing settings request"))
				return
			}
	
			acknowledge(s, i, fmt.Sprintf("Processing settings request for '%s'", operation))

			config := zomboid.GetCurrentConfig("server-name")

			fmt.Println(config)

			// ChannelMessageSendComplex

			// create tmp dir if it doesn't exist
			os.MkdirAll("./tmp", os.ModePerm)
			out, e := os.Create("./tmp/out.txt")
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