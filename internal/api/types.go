package api

// Board represents a Trello board.
type Board struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Desc           string     `json:"desc"`
	Closed         bool       `json:"closed"`
	IDOrganization string     `json:"idOrganization"`
	URL            string     `json:"url"`
	ShortURL       string     `json:"shortUrl"`
	ShortLink      string     `json:"shortLink"`
	DateLastActivity string   `json:"dateLastActivity"`
	Prefs          BoardPrefs `json:"prefs"`
	LabelNames     LabelNames `json:"labelNames"`
}

// BoardPrefs holds board preferences.
type BoardPrefs struct {
	PermissionLevel    string `json:"permissionLevel"`
	Voting             string `json:"voting"`
	Comments           string `json:"comments"`
	Background         string `json:"background"`
	BackgroundColor    string `json:"backgroundColor"`
	BackgroundImage    string `json:"backgroundImage"`
	SelfJoin           bool   `json:"selfJoin"`
	CardCovers         bool   `json:"cardCovers"`
	IsTemplate         bool   `json:"isTemplate"`
	CardAging          string `json:"cardAging"`
}

// LabelNames maps label colors to custom names.
type LabelNames struct {
	Black  string `json:"black"`
	Blue   string `json:"blue"`
	Green  string `json:"green"`
	Lime   string `json:"lime"`
	Orange string `json:"orange"`
	Pink   string `json:"pink"`
	Purple string `json:"purple"`
	Red    string `json:"red"`
	Sky    string `json:"sky"`
	Yellow string `json:"yellow"`
}

// TrelloList represents a Trello list (column).
type TrelloList struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Closed     bool    `json:"closed"`
	IDBoard    string  `json:"idBoard"`
	Pos        float64 `json:"pos"`
	Subscribed bool    `json:"subscribed"`
}

// Card represents a Trello card.
type Card struct {
	ID              string   `json:"id"`
	IDShort         int      `json:"idShort"`
	Name            string   `json:"name"`
	Desc            string   `json:"desc"`
	Closed          bool     `json:"closed"`
	IDBoard         string   `json:"idBoard"`
	IDList          string   `json:"idList"`
	IDMembers       []string `json:"idMembers"`
	IDLabels        []string `json:"idLabels"`
	Labels          []Label  `json:"labels"`
	Due             *string  `json:"due"`
	DueComplete     bool     `json:"dueComplete"`
	Start           *string  `json:"start"`
	Pos             float64  `json:"pos"`
	ShortLink       string   `json:"shortLink"`
	ShortURL        string   `json:"shortUrl"`
	URL             string   `json:"url"`
	Subscribed      bool     `json:"subscribed"`
	DateLastActivity string  `json:"dateLastActivity"`
	Badges          CardBadges `json:"badges"`
}

// CardBadges holds summary counts for a card.
type CardBadges struct {
	Attachments       int    `json:"attachments"`
	CheckItems        int    `json:"checkItems"`
	CheckItemsChecked int    `json:"checkItemsChecked"`
	Comments          int    `json:"comments"`
	Description       bool   `json:"description"`
	Due               *string `json:"due"`
	DueComplete       bool   `json:"dueComplete"`
	Subscribed        bool   `json:"subscribed"`
	Votes             int    `json:"votes"`
}

// Label represents a card label.
type Label struct {
	ID      string `json:"id"`
	IDBoard string `json:"idBoard"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

// Member represents a Trello member.
type Member struct {
	ID          string `json:"id"`
	FullName    string `json:"fullName"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatarUrl"`
	URL         string `json:"url"`
	IDBoards    []string `json:"idBoards"`
	MemberType  string `json:"memberType"`
	Confirmed   bool   `json:"confirmed"`
}

// Organization represents a Trello workspace/organization.
type Organization struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Desc        string `json:"desc"`
	URL         string `json:"url"`
	IDBoards    []string `json:"idBoards"`
}

// Checklist represents a card checklist.
type Checklist struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	IDBoard    string      `json:"idBoard"`
	IDCard     string      `json:"idCard"`
	Pos        float64     `json:"pos"`
	CheckItems []CheckItem `json:"checkItems"`
}

// CheckItem represents an item within a checklist.
type CheckItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	State       string  `json:"state"` // "complete" or "incomplete"
	IDChecklist string  `json:"idChecklist"`
	IDCard      string  `json:"idCard"`
	Pos         float64 `json:"pos"`
	Due         *string `json:"due"`
}

// Attachment represents a file or link attached to a card.
type Attachment struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	MimeType string `json:"mimeType"`
	Bytes    int64  `json:"bytes"`
	Date     string `json:"date"`
	IsUpload bool   `json:"isUpload"`
}

// Action represents a Trello activity/audit log entry.
type Action struct {
	ID              string      `json:"id"`
	IDMemberCreator string      `json:"idMemberCreator"`
	Type            string      `json:"type"`
	Date            string      `json:"date"`
	Data            ActionData  `json:"data"`
	MemberCreator   *Member     `json:"memberCreator,omitempty"`
}

// ActionData holds context data for an action.
type ActionData struct {
	Text  string `json:"text"`
	Board *struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		ShortLink string `json:"shortLink"`
	} `json:"board,omitempty"`
	Card *struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		ShortLink string `json:"shortLink"`
		IDShort   int    `json:"idShort"`
	} `json:"card,omitempty"`
	List *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"list,omitempty"`
}

// Webhook represents a Trello webhook subscription.
type Webhook struct {
	ID                    string `json:"id"`
	Description           string `json:"description"`
	IDModel               string `json:"idModel"`
	CallbackURL           string `json:"callbackURL"`
	Active                bool   `json:"active"`
	ConsecutiveFailures   int    `json:"consecutiveFailures"`
	FirstConsecutiveFailDate *string `json:"firstConsecutiveFailDate"`
}

// SearchResult holds the result of a search query.
type SearchResult struct {
	Cards   []Card   `json:"cards"`
	Boards  []Board  `json:"boards"`
	Members []Member `json:"members"`
	Options struct {
		Terms    []string `json:"terms"`
		Modifiers []string `json:"modifiers"`
	} `json:"options"`
}

// TrelloError is returned when the API responds with an error.
type TrelloError struct {
	StatusCode int
	Message    string
}

func (e *TrelloError) Error() string {
	return e.Message
}
