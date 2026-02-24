package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const apiBase = "https://api.trello.com/1"

// Client is an authenticated Trello API client.
type Client struct {
	apiKey     string
	apiToken   string
	httpClient *http.Client
}

// NewClient creates a new authenticated Client.
func NewClient(apiKey, apiToken string) *Client {
	return &Client{
		apiKey:   apiKey,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// authParams returns the base auth query params added to every request.
func (c *Client) authParams() url.Values {
	p := url.Values{}
	p.Set("key", c.apiKey)
	p.Set("token", c.apiToken)
	return p
}

// buildURL constructs a full API URL merging auth params with caller params.
func (c *Client) buildURL(path string, params url.Values) string {
	u, _ := url.Parse(apiBase + path)
	q := c.authParams()
	for k, vs := range params {
		for _, v := range vs {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// doRequest executes an HTTP request and returns the body bytes.
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		msg := string(body)
		return nil, &TrelloError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, msg)}
	}

	return body, nil
}

// Get makes a GET request to path with the given extra params.
func (c *Client) Get(path string, params url.Values) ([]byte, error) {
	if params == nil {
		params = url.Values{}
	}
	req, err := http.NewRequest(http.MethodGet, c.buildURL(path, params), nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

// Post makes a POST request to path with a JSON body.
func (c *Client) Post(path string, params url.Values, payload any) ([]byte, error) {
	if params == nil {
		params = url.Values{}
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encoding request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, c.buildURL(path, params), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.doRequest(req)
}

// Put makes a PUT request to path with a JSON body.
func (c *Client) Put(path string, params url.Values, payload any) ([]byte, error) {
	if params == nil {
		params = url.Values{}
	}
	var bodyReader io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("encoding request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(http.MethodPut, c.buildURL(path, params), bodyReader)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.doRequest(req)
}

// Delete makes a DELETE request to path.
func (c *Client) Delete(path string, params url.Values) ([]byte, error) {
	if params == nil {
		params = url.Values{}
	}
	req, err := http.NewRequest(http.MethodDelete, c.buildURL(path, params), nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

// ---- Boards ----

// GetBoard returns a board by ID.
func (c *Client) GetBoard(id string, params url.Values) (*Board, error) {
	body, err := c.Get("/boards/"+id, params)
	if err != nil {
		return nil, err
	}
	var b Board
	return &b, json.Unmarshal(body, &b)
}

// GetMyBoards returns all boards for the authenticated member.
func (c *Client) GetMyBoards(filter string) ([]Board, error) {
	params := url.Values{}
	if filter != "" {
		params.Set("filter", filter)
	}
	body, err := c.Get("/members/me/boards", params)
	if err != nil {
		return nil, err
	}
	var boards []Board
	return boards, json.Unmarshal(body, &boards)
}

// CreateBoard creates a new board.
func (c *Client) CreateBoard(name, desc, idOrganization string, prefs url.Values) (*Board, error) {
	params := url.Values{}
	params.Set("name", name)
	if desc != "" {
		params.Set("desc", desc)
	}
	if idOrganization != "" {
		params.Set("idOrganization", idOrganization)
	}
	for k, vs := range prefs {
		for _, v := range vs {
			params.Set(k, v)
		}
	}
	body, err := c.Post("/boards", params, nil)
	if err != nil {
		return nil, err
	}
	var b Board
	return &b, json.Unmarshal(body, &b)
}

// UpdateBoard updates a board.
func (c *Client) UpdateBoard(id string, params url.Values) (*Board, error) {
	body, err := c.Put("/boards/"+id, params, nil)
	if err != nil {
		return nil, err
	}
	var b Board
	return &b, json.Unmarshal(body, &b)
}

// DeleteBoard deletes (closes) a board.
func (c *Client) DeleteBoard(id string) error {
	_, err := c.Delete("/boards/"+id, nil)
	return err
}

// GetBoardLists returns all lists for a board.
func (c *Client) GetBoardLists(boardID, filter string) ([]TrelloList, error) {
	params := url.Values{}
	if filter != "" {
		params.Set("filter", filter)
	}
	body, err := c.Get("/boards/"+boardID+"/lists", params)
	if err != nil {
		return nil, err
	}
	var lists []TrelloList
	return lists, json.Unmarshal(body, &lists)
}

// GetBoardCards returns all cards for a board.
func (c *Client) GetBoardCards(boardID, filter string) ([]Card, error) {
	params := url.Values{}
	if filter != "" {
		params.Set("filter", filter)
	}
	body, err := c.Get("/boards/"+boardID+"/cards", params)
	if err != nil {
		return nil, err
	}
	var cards []Card
	return cards, json.Unmarshal(body, &cards)
}

// GetBoardMembers returns all members of a board.
func (c *Client) GetBoardMembers(boardID string) ([]Member, error) {
	body, err := c.Get("/boards/"+boardID+"/members", nil)
	if err != nil {
		return nil, err
	}
	var members []Member
	return members, json.Unmarshal(body, &members)
}

// GetBoardLabels returns all labels on a board.
func (c *Client) GetBoardLabels(boardID string) ([]Label, error) {
	body, err := c.Get("/boards/"+boardID+"/labels", nil)
	if err != nil {
		return nil, err
	}
	var labels []Label
	return labels, json.Unmarshal(body, &labels)
}

// ---- Lists ----

// GetList returns a list by ID.
func (c *Client) GetList(id string) (*TrelloList, error) {
	body, err := c.Get("/lists/"+id, nil)
	if err != nil {
		return nil, err
	}
	var l TrelloList
	return &l, json.Unmarshal(body, &l)
}

// CreateList creates a new list on a board.
func (c *Client) CreateList(name, idBoard string, pos string) (*TrelloList, error) {
	params := url.Values{}
	params.Set("name", name)
	params.Set("idBoard", idBoard)
	if pos != "" {
		params.Set("pos", pos)
	}
	body, err := c.Post("/lists", params, nil)
	if err != nil {
		return nil, err
	}
	var l TrelloList
	return &l, json.Unmarshal(body, &l)
}

// UpdateList updates a list.
func (c *Client) UpdateList(id string, params url.Values) (*TrelloList, error) {
	body, err := c.Put("/lists/"+id, params, nil)
	if err != nil {
		return nil, err
	}
	var l TrelloList
	return &l, json.Unmarshal(body, &l)
}

// ArchiveList archives (closes) a list.
func (c *Client) ArchiveList(id string, archive bool) (*TrelloList, error) {
	params := url.Values{}
	if archive {
		params.Set("value", "true")
	} else {
		params.Set("value", "false")
	}
	body, err := c.Put("/lists/"+id+"/closed", params, nil)
	if err != nil {
		return nil, err
	}
	var l TrelloList
	return &l, json.Unmarshal(body, &l)
}

// GetListCards returns all cards in a list.
func (c *Client) GetListCards(listID, filter string) ([]Card, error) {
	params := url.Values{}
	if filter != "" {
		params.Set("filter", filter)
	}
	body, err := c.Get("/lists/"+listID+"/cards", params)
	if err != nil {
		return nil, err
	}
	var cards []Card
	return cards, json.Unmarshal(body, &cards)
}

// ---- Cards ----

// GetCard returns a card by ID.
func (c *Client) GetCard(id string, params url.Values) (*Card, error) {
	body, err := c.Get("/cards/"+id, params)
	if err != nil {
		return nil, err
	}
	var card Card
	return &card, json.Unmarshal(body, &card)
}

// CreateCard creates a new card.
func (c *Client) CreateCard(idList, name, desc string, params url.Values) (*Card, error) {
	p := url.Values{}
	p.Set("idList", idList)
	p.Set("name", name)
	if desc != "" {
		p.Set("desc", desc)
	}
	for k, vs := range params {
		for _, v := range vs {
			p.Set(k, v)
		}
	}
	body, err := c.Post("/cards", p, nil)
	if err != nil {
		return nil, err
	}
	var card Card
	return &card, json.Unmarshal(body, &card)
}

// UpdateCard updates a card.
func (c *Client) UpdateCard(id string, params url.Values) (*Card, error) {
	body, err := c.Put("/cards/"+id, params, nil)
	if err != nil {
		return nil, err
	}
	var card Card
	return &card, json.Unmarshal(body, &card)
}

// DeleteCard deletes a card.
func (c *Client) DeleteCard(id string) error {
	_, err := c.Delete("/cards/"+id, nil)
	return err
}

// MoveCard moves a card to a different list (and optionally board).
func (c *Client) MoveCard(id, idList, idBoard string) (*Card, error) {
	params := url.Values{}
	params.Set("idList", idList)
	if idBoard != "" {
		params.Set("idBoard", idBoard)
	}
	return c.UpdateCard(id, params)
}

// GetCardChecklists returns all checklists for a card.
func (c *Client) GetCardChecklists(cardID string) ([]Checklist, error) {
	body, err := c.Get("/cards/"+cardID+"/checklists", nil)
	if err != nil {
		return nil, err
	}
	var checklists []Checklist
	return checklists, json.Unmarshal(body, &checklists)
}

// GetCardAttachments returns all attachments for a card.
func (c *Client) GetCardAttachments(cardID string) ([]Attachment, error) {
	body, err := c.Get("/cards/"+cardID+"/attachments", nil)
	if err != nil {
		return nil, err
	}
	var attachments []Attachment
	return attachments, json.Unmarshal(body, &attachments)
}

// AddComment adds a comment to a card.
func (c *Client) AddComment(cardID, text string) (*Action, error) {
	params := url.Values{}
	params.Set("text", text)
	body, err := c.Post("/cards/"+cardID+"/actions/comments", params, nil)
	if err != nil {
		return nil, err
	}
	var action Action
	return &action, json.Unmarshal(body, &action)
}

// AddLabelToCard adds a label to a card.
func (c *Client) AddLabelToCard(cardID, labelID string) error {
	params := url.Values{}
	params.Set("value", labelID)
	_, err := c.Post("/cards/"+cardID+"/idLabels", params, nil)
	return err
}

// RemoveLabelFromCard removes a label from a card.
func (c *Client) RemoveLabelFromCard(cardID, labelID string) error {
	_, err := c.Delete("/cards/"+cardID+"/idLabels/"+labelID, nil)
	return err
}

// AddMemberToCard assigns a member to a card.
func (c *Client) AddMemberToCard(cardID, memberID string) error {
	params := url.Values{}
	params.Set("value", memberID)
	_, err := c.Post("/cards/"+cardID+"/idMembers", params, nil)
	return err
}

// RemoveMemberFromCard removes a member from a card.
func (c *Client) RemoveMemberFromCard(cardID, memberID string) error {
	_, err := c.Delete("/cards/"+cardID+"/idMembers/"+memberID, nil)
	return err
}

// ---- Members ----

// GetMember returns a member by ID or username (use "me" for self).
func (c *Client) GetMember(idOrUsername string, params url.Values) (*Member, error) {
	body, err := c.Get("/members/"+idOrUsername, params)
	if err != nil {
		return nil, err
	}
	var m Member
	return &m, json.Unmarshal(body, &m)
}

// GetMemberBoards returns all boards for a member.
func (c *Client) GetMemberBoards(idOrUsername, filter string) ([]Board, error) {
	params := url.Values{}
	if filter != "" {
		params.Set("filter", filter)
	}
	body, err := c.Get("/members/"+idOrUsername+"/boards", params)
	if err != nil {
		return nil, err
	}
	var boards []Board
	return boards, json.Unmarshal(body, &boards)
}

// GetMemberCards returns all cards assigned to a member.
func (c *Client) GetMemberCards(idOrUsername, filter string) ([]Card, error) {
	params := url.Values{}
	if filter != "" {
		params.Set("filter", filter)
	}
	body, err := c.Get("/members/"+idOrUsername+"/cards", params)
	if err != nil {
		return nil, err
	}
	var cards []Card
	return cards, json.Unmarshal(body, &cards)
}

// GetMemberOrganizations returns all organizations/workspaces for a member.
func (c *Client) GetMemberOrganizations(idOrUsername string) ([]Organization, error) {
	body, err := c.Get("/members/"+idOrUsername+"/organizations", nil)
	if err != nil {
		return nil, err
	}
	var orgs []Organization
	return orgs, json.Unmarshal(body, &orgs)
}

// ---- Checklists ----

// GetChecklist returns a checklist by ID.
func (c *Client) GetChecklist(id string) (*Checklist, error) {
	body, err := c.Get("/checklists/"+id, nil)
	if err != nil {
		return nil, err
	}
	var cl Checklist
	return &cl, json.Unmarshal(body, &cl)
}

// CreateChecklist creates a new checklist on a card.
func (c *Client) CreateChecklist(idCard, name string) (*Checklist, error) {
	params := url.Values{}
	params.Set("idCard", idCard)
	params.Set("name", name)
	body, err := c.Post("/checklists", params, nil)
	if err != nil {
		return nil, err
	}
	var cl Checklist
	return &cl, json.Unmarshal(body, &cl)
}

// DeleteChecklist deletes a checklist.
func (c *Client) DeleteChecklist(id string) error {
	_, err := c.Delete("/checklists/"+id, nil)
	return err
}

// CreateCheckItem adds an item to a checklist.
func (c *Client) CreateCheckItem(checklistID, name string) (*CheckItem, error) {
	params := url.Values{}
	params.Set("name", name)
	body, err := c.Post("/checklists/"+checklistID+"/checkItems", params, nil)
	if err != nil {
		return nil, err
	}
	var item CheckItem
	return &item, json.Unmarshal(body, &item)
}

// UpdateCheckItem updates the state of a check item on a card.
func (c *Client) UpdateCheckItem(cardID, checklistID, checkItemID, state string) (*CheckItem, error) {
	params := url.Values{}
	params.Set("state", state)
	params.Set("idChecklist", checklistID)
	body, err := c.Put("/cards/"+cardID+"/checklist/"+checklistID+"/checkItem/"+checkItemID, params, nil)
	if err != nil {
		return nil, err
	}
	var item CheckItem
	return &item, json.Unmarshal(body, &item)
}

// ---- Search ----

// Search performs a global search across Trello.
func (c *Client) Search(query string, modelTypes []string, limit int) (*SearchResult, error) {
	params := url.Values{}
	params.Set("query", query)
	if len(modelTypes) > 0 {
		for _, t := range modelTypes {
			params.Add("modelTypes", t)
		}
	} else {
		params.Set("modelTypes", "all")
	}
	if limit > 0 {
		params.Set("cards_limit", fmt.Sprintf("%d", limit))
		params.Set("boards_limit", fmt.Sprintf("%d", limit))
		params.Set("members_limit", fmt.Sprintf("%d", limit))
	}
	params.Set("card_fields", "id,name,idBoard,idList,shortUrl,labels,due,dueComplete")
	params.Set("board_fields", "id,name,shortUrl,closed")

	body, err := c.Get("/search", params)
	if err != nil {
		return nil, err
	}
	var result SearchResult
	return &result, json.Unmarshal(body, &result)
}

// ---- Labels ----

// GetLabel returns a label by ID.
func (c *Client) GetLabel(id string) (*Label, error) {
	body, err := c.Get("/labels/"+id, nil)
	if err != nil {
		return nil, err
	}
	var l Label
	return &l, json.Unmarshal(body, &l)
}

// CreateLabel creates a new label on a board.
func (c *Client) CreateLabel(idBoard, name, color string) (*Label, error) {
	params := url.Values{}
	params.Set("idBoard", idBoard)
	params.Set("name", name)
	params.Set("color", color)
	body, err := c.Post("/labels", params, nil)
	if err != nil {
		return nil, err
	}
	var l Label
	return &l, json.Unmarshal(body, &l)
}

// DeleteLabel deletes a label.
func (c *Client) DeleteLabel(id string) error {
	_, err := c.Delete("/labels/"+id, nil)
	return err
}
