//UFID:88489948
//video link: https://www.youtube.com/watch?v=LjmHUTVEqbw
package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"

    "reddit_part2/internal/models"
)

type Client struct {
    BaseURL string
    UserID  string
}

func NewClient(baseURL, userID string) *Client {
    return &Client{BaseURL: baseURL, UserID: userID}
}

func (c *Client) Register() error {
    user := models.User{ID: c.UserID}
    return c.post("/register", user)
}

func (c *Client) CreateSubreddit(name string) error {
    sub := models.Subreddit{Name: name}
    return c.post("/subreddit", sub)
}

func (c *Client) CreatePost(content, subreddit string, isRepost bool) (int, error) {
    post := models.Post{Content: content, UserID: c.UserID, Subreddit: subreddit, IsRepost: isRepost}
    var response map[string]int
    err := c.post("/post", post, &response)
    return response["id"], err
}

func (c *Client) CreateComment(content string, parentID int) error {
    comment := models.Comment{Content: content, UserID: c.UserID, ParentID: parentID}
    return c.post("/comment", comment)
}

func (c *Client) Vote(itemID int, upvote bool) error {
    vote := struct {
        ItemID int  `json:"item_id"`
        Upvote bool `json:"upvote"`
    }{ItemID: itemID, Upvote: upvote}
    return c.post("/vote", vote)
}

func (c *Client) GetFeed() ([]*models.Post, error) {
    resp, err := http.Get(fmt.Sprintf("%s/feed?user_id=%s", c.BaseURL, c.UserID))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var feed []*models.Post
    err = json.NewDecoder(resp.Body).Decode(&feed)
    return feed, err
}

func (c *Client) SendDirectMessage(to, content string) error {
    msg := models.DirectMessage{From: c.UserID, To: to, Content: content}
    return c.post("/message", msg)
}

func (c *Client) post(path string, body interface{}, response ...interface{}) error {
    jsonBody, _ := json.Marshal(body)
    resp, err := http.Post(c.BaseURL+path, "application/json", bytes.NewBuffer(jsonBody))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status: %s", resp.Status)
    }

    if len(response) > 0 {
        return json.NewDecoder(resp.Body).Decode(response[0])
    }
    return nil
}