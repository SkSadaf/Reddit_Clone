//UFID:88489948
//video link: https://www.youtube.com/watch?v=LjmHUTVEqbw
package models

type User struct {
    ID         string   `json:"id"`
    Karma      int      `json:"karma"`
    Active     bool     `json:"active"`
    Subreddits []string `json:"subreddits"`
	Subscriptions []string
}

type Subreddit struct {
    Name    string  `json:"name"`
    Members int     `json:"members"`
    Posts   []*Post `json:"posts"`
}

type Post struct {
    ID        int        `json:"id"`
    Content   string     `json:"content"`
    UserID    string     `json:"user_id"`
    Subreddit string     `json:"subreddit"`
    Votes     int        `json:"votes"`
    Comments  []*Comment `json:"comments"`
    IsRepost  bool       `json:"is_repost"`
}

type Comment struct {
    ID       int        `json:"id"`
    Content  string     `json:"content"`
    UserID   string     `json:"user_id"`
    Votes    int        `json:"votes"`
    ParentID int        `json:"parent_id"`
    Comments []*Comment `json:"comments"`
}

type DirectMessage struct {
    From    string `json:"from"`
    To      string `json:"to"`
    Content string `json:"content"`
}