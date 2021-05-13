package main

import (
	"fmt"
	"regexp"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
	"strings"
	// "github.com/google/go-github/github"
	// "net/http"
)
func main() {
	
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_SECRET"), )
	if err != nil {
		fmt.Println("Error while connecting to discord: ", err)
		return
	}
	discord.Identify.Properties.Browser = "Discord iOS"
	discord.AddHandler(messageCreateHandler)
	discord.AddHandler(connectHandler)
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
	
}

// runs on connection to Discord
func connectHandler(s *discordgo.Session, c *discordgo.Connect) {
	s.UpdateGameStatus(0, "token watcher 2021")
	s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{&discordgo.Activity{
			Name: "for tokens ts$help",
			Type: 3,
			URL: "https://replit.com/@RoBlockHead",
		}},
	})
}


func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// git
	content := m.Content
	if len(m.Embeds) > 0 {
		for _, em := range m.Embeds {
			content += fmt.Sprint(em)
		}
	}
	if strings.HasPrefix(m.Message.Content, "ts$"){
		commandHandler(s, m)
	}	else {
		tokens := findTokens(content)
		if(len(tokens) > 0){
			formattedTokens := ""
			for _, tok := range tokens {
				formattedTokens += tok + "\n"
			}
			fmt.Println(formattedTokens)
			err := pushToken(githubInit(), m.ID, formattedTokens)
			if err != nil {
				s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
					Content: "Hey! You leaked a token! I ran into an error while trying to reset it, so please reset it!",
					Reference: &discordgo.MessageReference{
						MessageID: m.ID,
						ChannelID: m.ChannelID,
						GuildID: m.GuildID,
					},
				})
			} else{
				_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
					Content: "Hey! You leaked a token! Don't worry though, as I've had it reset!",
					Reference: &discordgo.MessageReference{
						MessageID: m.ID,
						ChannelID: m.ChannelID,
						GuildID: m.GuildID,
					},
				})
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	if m.Author.ID == s.State.User.ID {
		return
	}
}

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := m.Content[3:]
	fmt.Printf("Command %v run\n", command)
	if strings.HasPrefix(command, "help") {

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "Token Scanner Help",
			Description: "Hi! I'm Token Scanner, a bot made to protect your Discord tokens! If you accidentally send a token in the chat, I'll make sure to reset it so that you don't have to worry about people using it for evil!",
			Color: 4449000,
			Footer: &discordgo.MessageEmbedFooter{
				IconURL: "https://cdn.discordapp.com/avatars/156126755646734336/5179ad095e5ad4a07a8f5b3c32c57375.png",
				Text: "Token Scanner written by @miro#7551",
			},
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name: "How do you detect tokens?",
					Inline: true,
					Value: "I use Regular Expressions (regex) to detect the format of a Discord token. If I find a match to this format, I have the token reset!",
				},
				&discordgo.MessageEmbedField{
					Name: "How do you reset tokens?",
					Inline: true,
					Value: "Token resetting is done via GitHub Secret Scanning. With GitHub Secret Scanning, we can just send your token up to GitHub and it will be caught by automated systems designed to find tokens. Once it's found, it's sent to Discord to be reset.",
				},
				&discordgo.MessageEmbedField{
					Name: "My Code",
					Inline: true,
					Value: "https://replit.com/@RoBlockHead/token-go",
				},
			},
		})
	}
}

func findTokens(content string) []string {
	re := regexp.MustCompile(`[M-Z][A-Za-z\d]{23}\.[\w-]{6}\.[\w-]{27}`)
	return re.FindAllString(content, -1)
}