package bot

var (
	bots map[string]Bot
)

func Register(name string, bot Bot) {
	if bots == nil {
		bots = make(map[string]Bot)
	}
	bots[name] = bot
}

func ListBots() map[string]Bot {
	return bots
}

func GetBot(name string) Bot {
	return bots[name]
}
