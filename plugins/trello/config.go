package trello

const (
	// ids can be found by going to the cacophony board, adding .json to the end
	// of url and searching for: "name": "*board name*", next to it is the id

	backlogBoardID = "5c7c490ac35ba5056fed77e9"
)

type Config struct {
	TrelloKey   string `envconfig:"TRELLO_KEY"`
	TrelloToken string `envconfig:"TRELLO_TOKEN"`
}
