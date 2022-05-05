package bot

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/Harlow-CS/zomboidBot/zomboid"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", os.Getenv("bot_oauth"), "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	ChannelID = os.Getenv("server_channel_id")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func isSanitary(fileName string) bool {
	sanitary, _ := regexp.MatchString(`[^a-zA-Z0-9_\-\.]`, fileName)
	return sanitary
}

func acknowledge(s *discordgo.Session, i *discordgo.InteractionCreate, text string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: text,
		},
	})
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "start-server",
			Description: "Starts up a server",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "save-file",
					Description: "name of the save file (exclude extension)",
					Required:    true,
				},
			},
		},
		{
			Name:        "stop-server",
			Description: "Stops the current running server",
		},
		{
			Name:        "save-file",
			Description: "Create a save file or list the currently hosted save files",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "operation",
					Description: "(create or ls)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "name of save file you are creating (without extension) (ignored by ls)",
					Required:    false,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"start-server": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
			//check to see if a server is already running
			if (factorio.Server != nil) {
				// alert user that they're bad
				acknowledge(s, i, "Server is already running")
				return
			}
			
			acknowledge(s, i, "Starting server...")

			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			var fileName string = ""

			// Get the value from the option map.
			if option, ok := optionMap["save-file"]; ok {
				fileName = option.StringValue()
			}

			if (!isSanitary(fileName)) {
				s.ChannelMessageSend(ChannelID, "save-file name is invalid")
				return
			}

			if (!factorio.SaveFileExists(fileName)) {
				s.ChannelMessageSend(ChannelID, "The given save file does not exist")
				return
			}

			// if not, start the server
			factorio.StartServer(fileName)

			s.ChannelMessageSend(ChannelID, "Server is running!")

		},
		"stop-server": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
			//check to see if a server is running
			if (factorio.Server == nil) {
				// alert user that they're bad
				acknowledge(s, i, "No running server detected")
				return
			}
			
			acknowledge(s, i, "Shutting down server...")

			// if not, start the server
			factorio.StopServer()

			s.ChannelMessageSend(ChannelID, "Server is shutdown")

		},
		"save-file": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			var operation = ""

			// Get the value from the option map.
			if option, ok := optionMap["operation"]; ok {
				operation = option.StringValue()
			}

			if (operation == "create") {

				acknowledge(s, i, "Creating save file...")

				var saveName = ""
				if option, ok := optionMap["name"]; ok {
					saveName = option.StringValue()
				}

				if (!isSanitary(saveName)) {
					s.ChannelMessageSend(ChannelID, "File name is invalid")
					return
				}

				if (factorio.SaveFileExists(saveName)) {
					s.ChannelMessageSend(ChannelID, "Save file of the same name already exists")
					return
				}

				factorio.CreateSaveFile(saveName)

				s.ChannelMessageSend(ChannelID, "Save created")

			} else if (operation == "ls") {

				var saveFiles []string = factorio.ListSaveFiles()
				var messageString = strings.Join(saveFiles,"\n")
				acknowledge(s, i, "```\n" + messageString + "\n```\n")
			} else {
				// unknown op
				acknowledge(s, i, "Unknown operation")
			}

		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func Start() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		if err != nil {
			log.Fatalf("Could not fetch registered commands: %v", err)
		}

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}