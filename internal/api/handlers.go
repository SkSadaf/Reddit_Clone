//UFID:88489948
//video link: https://www.youtube.com/watch?v=LjmHUTVEqbw
package api

import (
    "encoding/json"
    "net/http"
    //"strconv"

    "reddit_part2/internal/engine"
    "reddit_part2/internal/models"
)

type Handler struct {
    engine *engine.Engine
}

func NewHandler(e *engine.Engine) *Handler {
    return &Handler{engine: e}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/register":
        h.register(w, r)
    case "/subreddit":
        h.subreddit(w, r)
    case "/post":
        h.post(w, r)
    case "/comment":
        h.comment(w, r)
    case "/vote":
        h.vote(w, r)
    case "/feed":
        h.feed(w, r)
    case "/message":
        h.message(w, r)
    case "/messages":
        if r.Method == http.MethodGet {
            h.getDirectMessages(w, r)
        } else if r.Method == http.MethodPost {
            h.message(w, r)
        }
    case "/join":
        h.joinSubreddit(w, r)
    case "/leave":
        h.leaveSubreddit(w, r)
    default:
        http.NotFound(w, r)
    }
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
    var user models.User
    json.NewDecoder(r.Body).Decode(&user)
    h.engine.RegisterAccount(user.ID)
    w.WriteHeader(http.StatusCreated)
}

func (h *Handler) subreddit(w http.ResponseWriter, r *http.Request) {
    var sub models.Subreddit
    json.NewDecoder(r.Body).Decode(&sub)
    h.engine.CreateSubreddit(sub.Name)
    w.WriteHeader(http.StatusCreated)
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
    var post models.Post
    json.NewDecoder(r.Body).Decode(&post)
    id := h.engine.CreatePost(post.Content, post.UserID, post.Subreddit, post.IsRepost)
    json.NewEncoder(w).Encode(map[string]int{"id": id})
}

/**func (h *Handler) comment(w http.ResponseWriter, r *http.Request) {
    var comment models.Comment
    json.NewDecoder(r.Body).Decode(&comment)
    h.engine.CreateComment(comment.Content, comment.UserID, comment.ParentID)
    w.WriteHeader(http.StatusCreated)
}**/
func (h *Handler) comment(w http.ResponseWriter, r *http.Request) {
    var comment models.Comment
    json.NewDecoder(r.Body).Decode(&comment)
    id := h.engine.CreateComment(comment.Content, comment.UserID, comment.ParentID)
    json.NewEncoder(w).Encode(map[string]int{"id": id})
}


func (h *Handler) vote(w http.ResponseWriter, r *http.Request) {
    var vote struct {
        ItemID int  `json:"item_id"`
        Upvote bool `json:"upvote"`
    }
    json.NewDecoder(r.Body).Decode(&vote)
    h.engine.Vote(vote.ItemID, vote.Upvote)
    w.WriteHeader(http.StatusOK)
}

func (h *Handler) feed(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    feed := h.engine.GetFeed(userID)
    json.NewEncoder(w).Encode(feed)
}

func (h *Handler) message(w http.ResponseWriter, r *http.Request) {
    var msg models.DirectMessage
    json.NewDecoder(r.Body).Decode(&msg)
    h.engine.SendDirectMessage(msg.From, msg.To, msg.Content)
    w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getFeed(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    feed := h.engine.GetFeed(userID)
    json.NewEncoder(w).Encode(feed)
}

func (h *Handler) sendDirectMessage(w http.ResponseWriter, r *http.Request) {
    var msg models.DirectMessage
    json.NewDecoder(r.Body).Decode(&msg)
    h.engine.SendDirectMessage(msg.From, msg.To, msg.Content)
    w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getDirectMessages(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    messages := h.engine.GetDirectMessages(userID)
    json.NewEncoder(w).Encode(messages)
}

func (h *Handler) joinSubreddit(w http.ResponseWriter, r *http.Request) {
    var req struct {
        UserID    string `json:"user_id"`
        Subreddit string `json:"subreddit"`
    }
    json.NewDecoder(r.Body).Decode(&req)
    err := h.engine.JoinSubreddit(req.UserID, req.Subreddit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (h *Handler) leaveSubreddit(w http.ResponseWriter, r *http.Request) {
    var req struct {
        UserID    string `json:"user_id"`
        Subreddit string `json:"subreddit"`
    }
    json.NewDecoder(r.Body).Decode(&req)
    err := h.engine.LeaveSubreddit(req.UserID, req.Subreddit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusOK)
}
