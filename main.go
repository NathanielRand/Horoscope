package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const prefix string = "!hs"
const version string = "1.1.3"

// Result struct contains shortened URL data.
type apiResponse struct {
	Color         string `json:"color"`
	Compatibility string `json:"compatibility"`
	CurrentDate   string `json:"current_date"`
	DateRange     string `json:"date_range"`
	Description   string `json:"description"`
	LuckyNumber   string `json:"lucky_number"`
	LuckyTime     string `json:"lucky_time"`
	Mood          string `json:"mood"`
}

func goDotEnvVariable(key string) string {
	// Load .env file.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Return value from key provided.
	return os.Getenv(key)
}

func main() {
	// Grab bot token env var.
	botToken := goDotEnvVariable("BOT_TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// guildID := m.Message.GuildID

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Grab message content from guild.
	content := m.Content

	// Trim bot command from string to grab User tagged
	trimmedPrefixCommand := strings.TrimPrefix(content, prefix)
	trimmedSpaceCommand := strings.Title(strings.TrimSpace(trimmedPrefixCommand))

	if strings.Contains(content, prefix) && trimmedSpaceCommand == "" {
		// Build start vote message
		author := m.Author.Username
		message := "Yo " + author + "..." + "looks like you forgot to add a sign. (example: !hs Aquarius). Give it another try, you got this."

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+"help") {
		// Build help message
		author := m.Author.Username

		// Title
		commandHelpTitle := "Looks like you need a hand. Check out my goodies below... \n \n"

		// Notes
		note1 := "- Bot will return a Horoscope based on cosmic events. \n"
		note2 := "- Commands are case-sensitive. They must be in lower-case (except the sign name, that is optional) :) \n"
		note3 := "- Dev: Narsiq#5638. DM me for requests/questions/love. \n"

		// Commands
		commandHelp := "❔  " + prefix + "help : Provides a list of my commands. \n"
		commandHoroscope := "🦶🏽  " + prefix + " <Sign> : Return your Horoscope based on cosmic events. Do not include '<>' in the command. \n"
		commandInvite := "🔗  " + prefix + "invite : Invite link for the Horoscope Bot. \n"
		commandSite := "🔗  " + prefix + "site : Link to the Horoscope website. \n"
		commandSupport := "✨  " + prefix + "support : Link to the Horoscope Patreon. \n"
		commandStats := "📊  " + prefix + "stats : Check out Horoscope stats. \n"
		commandVersion := "🤖  " + prefix + "version : Current Horoscope version. \n"

		// Build sub messages
		notesMessage := note1 + note2 + note3
		commandsMessage := commandHelp + commandHoroscope
		othersMessage := commandInvite + commandSite + commandSupport + commandStats + commandVersion

		// Build full message
		message := "Whats up " + author + "\n \n" + commandHelpTitle + "NOTES: \n \n" + notesMessage + "\n" + "COMMANDS: \n \n" + commandsMessage + "\n" + "OTHER: \n \n" + othersMessage + "\n \n" + "https://www.patreon.com/BotVoteTo"

		// Reply to help request with build message above.
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+"site") {
		// Build start vote message
		author := m.Author.Username
		message := "Here ya go " + author + "..." + "\n" + "https://discordbots.dev/"

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+"support") {
		// Build start vote message
		author := m.Author.Username
		message := "Thanks for thinking of me " + author + " 💖." + "\n" + "https://www.patreon.com/BotVoteTo"

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+"version") {
		// Build start vote message
		message := "Horoscope is currently running version " + version

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+"stats") {
		// TODO: This will need to be updated to iterate through
		// all shards once the bot joins 1,000 servers.
		guilds := s.State.Ready.Guilds
		fmt.Println(len(guilds))
		guildCount := len(guilds)

		guildCountStr := strconv.Itoa(guildCount)

		// // Build start vote message
		message := "Horoscope is currently on " + guildCountStr + " servers. Such wow!"

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+"invite") {
		author := m.Author.Username

		// // Build start vote message
		message := "Wow! Such nice " + author + ". Thanks for spreading the 💖. Here is an invite link made just for you... \n \n" + "https://discord.com/api/oauth2/authorize?client_id=921254599913013320&permissions=274878036032&scope=bot"

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+trimmedPrefixCommand) && trimmedSpaceCommand != "" {

		// Check is command matches a symbol, if not, return
		symbol := getSymbol(trimmedSpaceCommand)
		if symbol == "?" {
			return
		}

		// Make api call while passing in sign command
		responseAPIBody := callAPI(trimmedSpaceCommand)
		var resultObject apiResponse
		err := json.Unmarshal(responseAPIBody, &resultObject)
		if err != nil {
			returnErrorMessage(s, m)
			return
		}

		// Store gathered data in a object and assign to var
		data := apiResponse{
			Color:         resultObject.Color,
			Compatibility: resultObject.Compatibility,
			CurrentDate:   resultObject.CurrentDate,
			DateRange:     resultObject.DateRange,
			Description:   resultObject.Description,
			LuckyNumber:   resultObject.LuckyNumber,
			LuckyTime:     resultObject.LuckyTime,
			Mood:          resultObject.Mood,
		}

		// Grab author
		author := m.Author.Username

		messageGreeting := author + ", \n \n"
		messageTitle := symbol + " " + trimmedSpaceCommand + " \n"
		messageColor := "🎨 Color: " + data.Color + "\n"
		messageCompatibility := "❤️ Compatibility: " + data.Compatibility + "\n"
		messageCurrentDate := data.CurrentDate + "\n \n"
		messageDateRange := "📆 Date Range: " + data.DateRange
		messageDescription := data.Description + "\n \n"
		messageLuckyNumber := "🍀 Lucky Number: " + data.LuckyNumber + "\n"
		messageLuckyTime := "🕰️ Lucky Time: " + data.LuckyTime + "\n \n"
		messageMood := "🧠 Mood: " + data.Mood + "\n"

		messageFull := "``` \n" + messageGreeting + messageTitle + messageCurrentDate + messageDescription + messageMood + messageColor + messageCompatibility + messageLuckyNumber + messageLuckyTime + messageDateRange + "\n ```"

		// Send start vote message
		_, err = s.ChannelMessageSendReply(m.ChannelID, messageFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}
}

func callAPI(sign string) []byte {
	url := "https://sameer-kumar-aztro-v1.p.rapidapi.com/?sign=" + sign + "&day=today"

	req, _ := http.NewRequest("POST", url, nil)

	req.Header.Add("x-rapidapi-host", "sameer-kumar-aztro-v1.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "a1fc98aff4mshe36beba570421c9p12a438jsnd30142a7efd7")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func getSymbol(command string) string {

	var symbol string

	switch command {
	case "Aries":
		symbol = "♈︎"
	case "Taurus":
		symbol = "♉︎"
	case "Gemini":
		symbol = "♊︎"
	case "Cancer":
		symbol = "♋︎"
	case "Leo":
		symbol = "♌︎"
	case "Virgo":
		symbol = "♍︎"
	case "Libra":
		symbol = "♎︎"
	case "Scorpio":
		symbol = "♏︎"
	case "Sagittarius":
		symbol = "♐︎"
	case "Capricorn":
		symbol = "♑︎"
	case "Aquarius":
		symbol = "♒︎"
	case "Pisces":
		symbol = "♓︎"
	default:
		symbol = "?"
	}

	return symbol
}

func returnErrorMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Build start vote message
	message := "Horoscope bot encountered an unknown error. Development team has been notified. Sorry we suck..."

	// Send start vote message
	_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
	if err != nil {
		fmt.Println(err)
	}
}
