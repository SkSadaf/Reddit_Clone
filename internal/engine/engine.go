//UFID:88489948
//video link: https://www.youtube.com/watch?v=LjmHUTVEqbw
package engine

import (
    "log"
    "sort"
    "sync"
    "fmt"

    "reddit_part2/internal/models"
)

type Engine struct {
    Users          map[string]*models.User
    Subreddits     map[string]*models.Subreddit
    Posts          map[int]*models.Post
    Comments       map[int]*models.Comment
    DirectMessages map[string][]*models.DirectMessage
    NextPostID     int
    NextCommentID  int
    Mu             sync.Mutex
}

func NewEngine() *Engine {
    return &Engine{
        Users:          make(map[string]*models.User),
        Subreddits:     make(map[string]*models.Subreddit),
        Posts:          make(map[int]*models.Post),
        Comments:       make(map[int]*models.Comment),
        DirectMessages: make(map[string][]*models.DirectMessage),
        NextPostID:     1,
        NextCommentID:  1,
    }
}

func (e *Engine) RegisterAccount(userID string) {
    e.Mu.Lock()
    defer e.Mu.Unlock()
    e.Users[userID] = &models.User{
        ID:            userID,
        Active:        true,
        Subscriptions: []string{},
    }
    log.Printf("User registered: %s", userID)
}

func (e *Engine) CreateSubreddit(name string) {
    e.Mu.Lock()
    defer e.Mu.Unlock()
    if _, exists := e.Subreddits[name]; !exists {
        e.Subreddits[name] = &models.Subreddit{Name: name}
        log.Printf("Subreddit created: %s", name)
    }
}

/**func (e *Engine) CreatePost(content, userID, subreddit string, isRepost bool) int {
    e.Mu.Lock()
    defer e.Mu.Unlock()
    if _, exists := e.Subreddits[subreddit]; !exists {
        log.Printf("Attempted to post in non-existent subreddit: %s", subreddit)
        return 0
    }
    post := &models.Post{
        ID:        e.NextPostID,
        Content:   content,
        UserID:    userID,
        Subreddit: subreddit,
        IsRepost:  isRepost,
    }
    e.Posts[e.NextPostID] = post
    e.Subreddits[subreddit].Posts = append(e.Subreddits[subreddit].Posts, post)
    e.NextPostID++
    log.Printf("Post created in %s by %s", subreddit, userID)
    return post.ID
}**/
func (e *Engine) CreatePost(content, userID, subreddit string, isRepost bool) int {
    e.Mu.Lock()
    defer e.Mu.Unlock()
    if _, exists := e.Subreddits[subreddit]; !exists {
        log.Printf("Attempted to post in non-existent subreddit: %s", subreddit)
        return 0
    }
    post := &models.Post{
        ID:        e.NextPostID,
        Content:   content,
        UserID:    userID,
        Subreddit: subreddit,
        IsRepost:  isRepost,
    }
    e.Posts[e.NextPostID] = post
    e.Subreddits[subreddit].Posts = append(e.Subreddits[subreddit].Posts, post)
    
    // Subscribe the user to the subreddit if not already subscribed
    user := e.Users[userID]
    if !contains(user.Subscriptions, subreddit) {
        user.Subscriptions = append(user.Subscriptions, subreddit)
    }
    
    e.NextPostID++
    log.Printf("Post created in %s by %s", subreddit, userID)
    return post.ID
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}


/**func (e *Engine) CreateComment(content, userID string, parentID int) {
    e.Mu.Lock()
    defer e.Mu.Unlock()
    comment := &models.Comment{
        ID:       e.NextCommentID,
        Content:  content,
        UserID:   userID,
        ParentID: parentID,
    }
    e.Comments[e.NextCommentID] = comment
    if post, ok := e.Posts[parentID]; ok {
        post.Comments = append(post.Comments, comment)
    } else if parentComment, ok := e.Comments[parentID]; ok {
        parentComment.Comments = append(parentComment.Comments, comment)
    }
    e.NextCommentID++
    log.Printf("Comment created by %s", userID)
}**/
func (e *Engine) CreateComment(content, userID string, parentID int) int {
    e.Mu.Lock()
    defer e.Mu.Unlock()
    comment := &models.Comment{
        ID:       e.NextCommentID,
        Content:  content,
        UserID:   userID,
        ParentID: parentID,
    }
    e.Comments[e.NextCommentID] = comment
    if post, ok := e.Posts[parentID]; ok {
        post.Comments = append(post.Comments, comment)
    } else if parentComment, ok := e.Comments[parentID]; ok {
        parentComment.Comments = append(parentComment.Comments, comment)
    }
    e.NextCommentID++
    log.Printf("Comment created by %s", userID)
    return comment.ID
}


func (e *Engine) Vote(itemID int, upvote bool) {
    e.Mu.Lock()
    defer e.Mu.Unlock()
    if post, ok := e.Posts[itemID]; ok {
        if upvote {
            post.Votes++
            e.Users[post.UserID].Karma++
        } else {
            post.Votes--
            e.Users[post.UserID].Karma--
        }
    } else if comment, ok := e.Comments[itemID]; ok {
        if upvote {
            comment.Votes++
            e.Users[comment.UserID].Karma++
        } else {
            comment.Votes--
            e.Users[comment.UserID].Karma--
        }
    }
}

func (e *Engine) GetFeed(userID string) []*models.Post {
    e.Mu.Lock()
    defer e.Mu.Unlock()
    user, exists := e.Users[userID]
    if !exists {
        return []*models.Post{}
    }
    var feed []*models.Post
    for _, subreddit := range user.Subscriptions {
        if sub, exists := e.Subreddits[subreddit]; exists {
            feed = append(feed, sub.Posts...)
        }
    }
    sort.Slice(feed, func(i, j int) bool {
        return feed[i].Votes > feed[j].Votes
    })
    return feed
}

func (e *Engine) SendDirectMessage(from, to, content string) {
    e.Mu.Lock()
    defer e.Mu.Unlock()
    message := &models.DirectMessage{From: from, To: to, Content: content}
    e.DirectMessages[to] = append(e.DirectMessages[to], message)
    log.Printf("DM sent from %s to %s", from, to)
}

/**func (e *Engine) GetDirectMessages(userID string) []*models.DirectMessage {
    e.Mu.Lock()
    defer e.Mu.Unlock()
    messages, exists := e.DirectMessages[userID]
    if !exists {
        return []*models.DirectMessage{}
    }
    // Clear the messages after retrieving them
    delete(e.DirectMessages, userID)
    return messages
}
**/
func (e *Engine) GetDirectMessages(userID string) []*models.DirectMessage {
    e.Mu.Lock()
    defer e.Mu.Unlock()
    messages, exists := e.DirectMessages[userID]
    if !exists {
        return []*models.DirectMessage{}
    }
    // Clear the messages after retrieving them
    delete(e.DirectMessages, userID)
    return messages
}

func (e *Engine) JoinSubreddit(userID, subredditName string) error {
    e.Mu.Lock()
    defer e.Mu.Unlock()
    
    user, exists := e.Users[userID]
    if !exists {
        return fmt.Errorf("user not found")
    }
    
    if _, exists := e.Subreddits[subredditName]; !exists {
        return fmt.Errorf("subreddit not found")
    }
    
    if !contains(user.Subscriptions, subredditName) {
        user.Subscriptions = append(user.Subscriptions, subredditName)
        log.Printf("User %s joined subreddit %s", userID, subredditName)
    }
    
    return nil
}

func (e *Engine) LeaveSubreddit(userID, subredditName string) error {
    e.Mu.Lock()
    defer e.Mu.Unlock()
    
    user, exists := e.Users[userID]
    if !exists {
        return fmt.Errorf("user not found")
    }
    
    for i, sub := range user.Subscriptions {
        if sub == subredditName {
            user.Subscriptions = append(user.Subscriptions[:i], user.Subscriptions[i+1:]...)
            log.Printf("User %s left subreddit %s", userID, subredditName)
            return nil
        }
    }
    
    return fmt.Errorf("user is not subscribed to this subreddit")
}
