//UFID:88489948
//video link: https://www.youtube.com/watch?v=LjmHUTVEqbw

package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "math/rand"
    "time"
    "bufio"
    "os"
    "strings"

    "reddit_part2/internal/models"
)

type RESTClient struct {
    UserID  string
    BaseURL string
}

func NewRESTClient(userID, baseURL string) *RESTClient {
    return &RESTClient{
        UserID:  userID,
        BaseURL: baseURL,
    }
}

func (c *RESTClient) Register() error {
    user := models.User{ID: c.UserID}
    return c.post("/register", user)
}
func (c *RESTClient) GetDirectMessages() ([]models.DirectMessage, error) {
    resp, err := http.Get(fmt.Sprintf("%s/messages?user_id=%s", c.BaseURL, c.UserID))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var messages []models.DirectMessage
    err = json.NewDecoder(resp.Body).Decode(&messages)
    if err != nil {
        // If decoding as an array fails, try decoding as a single message
        var singleMessage models.DirectMessage
        err = json.NewDecoder(resp.Body).Decode(&singleMessage)
        if err != nil {
            return nil, err
        }
        messages = []models.DirectMessage{singleMessage}
    }

    return messages, nil
}


func (c *RESTClient) CreateSubreddit(name string) error {
    _, err := http.Post(c.BaseURL+"/subreddit", "application/json", bytes.NewBuffer([]byte(`{"name":"`+name+`"}`)))
    if err != nil {
        log.Printf("Could not create subreddit: %v", err)
        return err
    }
    log.Printf("User %s created subreddit %s", c.UserID, name)
    return nil
}

func (c *RESTClient) CreatePost(content, subreddit string, isRepost bool) (int, error) {
    post := struct {
        Content   string `json:"content"`
        UserID    string `json:"user_id"`
        Subreddit string `json:"subreddit"`
        IsRepost  bool   `json:"is_repost"`
    }{
        Content:   content,
        UserID:    c.UserID,
        Subreddit: subreddit,
        IsRepost:  isRepost,
    }
    var response struct {
        ID int `json:"id"`
    }
    err := c.post("/post", post, &response)
    if err != nil {
        return 0, err
    }
    if response.ID == 0 {
        return 0, fmt.Errorf("failed to create post: subreddit %s does not exist", subreddit)
    }
    return response.ID, nil
}

/**func (c *RESTClient) CreateComment(content string, parentID int) error {
    comment := models.Comment{Content: content, UserID: c.UserID, ParentID: parentID}
    return c.post("/comment", comment)
}**/
/**func (c *RESTClient) CreateComment(content string, parentID int) (int, error) {
    comment := models.Comment{Content: content, UserID: c.UserID, ParentID: parentID}
    var response struct {
        ID int `json:"id"`
    }
    err := c.post("/comment", comment, &response)
    if err != nil {
        return 0, err
    }
    return response.ID, nil
}**/
/**func (c *RESTClient) CreateComment(content string, parentID int) (int, error) {
    comment := models.Comment{Content: content, UserID: c.UserID, ParentID: parentID}
    var response struct {
        ID int `json:"id"`
    }
    err := c.post("/comment", comment, &response)
    if err != nil {
        return 0, err
    }
    return response.ID, nil
}
**/
func (c *RESTClient) CreateComment(content string, parentID int) (int, error) {
    comment := models.Comment{Content: content, UserID: c.UserID, ParentID: parentID}
    var response struct {
        ID int `json:"id"`
    }
    err := c.post("/comment", comment, &response)
    if err != nil {
        return 0, err
    }
    return response.ID, nil
}



func (c *RESTClient) Vote(itemID int, upvote bool) error {
    vote := struct {
        ItemID int  `json:"item_id"`
        Upvote bool `json:"upvote"`
    }{ItemID: itemID, Upvote: upvote}
    return c.post("/vote", vote)
}

/**func (c *RESTClient) GetFeed() ([]models.Post, error) {
    resp, err := http.Get(fmt.Sprintf("%s/feed?user_id=%s", c.BaseURL, c.UserID))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var feed []models.Post
    err = json.NewDecoder(resp.Body).Decode(&feed)
    return feed, nil
}
**/
func (c *RESTClient) GetFeed() ([]models.Post, error) {
    resp, err := http.Get(fmt.Sprintf("%s/feed?user_id=%s", c.BaseURL, c.UserID))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var feed []models.Post
    err = json.NewDecoder(resp.Body).Decode(&feed)
    return feed, nil
}

func (c *RESTClient) SendDirectMessage(to, content string) error {
    msg := models.DirectMessage{From: c.UserID, To: to, Content: content}
    return c.post("/message", msg)
}

func (c *RESTClient) post(path string, body interface{}, response ...interface{}) error {
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
func (c *RESTClient) JoinSubreddit(subreddit string) error {
    req := struct {
        UserID    string `json:"user_id"`
        Subreddit string `json:"subreddit"`
    }{
        UserID:    c.UserID,
        Subreddit: subreddit,
    }
    return c.post("/join", req)
}

func (c *RESTClient) LeaveSubreddit(subreddit string) error {
    req := struct {
        UserID    string `json:"user_id"`
        Subreddit string `json:"subreddit"`
    }{
        UserID:    c.UserID,
        Subreddit: subreddit,
    }
    return c.post("/leave", req)
}


func main() {
    rand.Seed(time.Now().UnixNano())
    client := NewRESTClient(fmt.Sprintf("user%d", rand.Intn(1000)), "http://localhost:8080")
    err := client.Register()
    if err != nil {
        log.Fatalf("Failed to register: %v", err)
    }

    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Println("\nChoose an action:")
        fmt.Println("1. Create Subreddit")
        fmt.Println("2. Create Post")
        fmt.Println("3. Create Comment")
        fmt.Println("4. Vote")
        fmt.Println("5. Get Feed")
        fmt.Println("6. Send Direct Message")
        fmt.Println("7. Get Direct Messages")
        fmt.Println("8. Join Subreddit")
        fmt.Println("9. Leave Subreddit")
        fmt.Println("10. Exit")
        //fmt.Println("8. Exit")

        choice, _ := reader.ReadString('\n')
        choice = strings.TrimSpace(choice)

        switch choice {
        case "1":
            fmt.Print("Enter subreddit name: ")
            name, _ := reader.ReadString('\n')
            name = strings.TrimSpace(name)
            err := client.CreateSubreddit(name)
            if err != nil {
                fmt.Printf("Error creating subreddit: %v\n", err)
            } else {
                fmt.Printf("Subreddit '%s' created successfully\n", name)
            }

        case "2":
            fmt.Print("Enter post content: ")
            content, _ := reader.ReadString('\n')
            content = strings.TrimSpace(content)
            fmt.Print("Enter subreddit name: ")
            subreddit, _ := reader.ReadString('\n')
            subreddit = strings.TrimSpace(subreddit)
            fmt.Print("Is this a repost? (true/false): ")
            isRepostStr, _ := reader.ReadString('\n')
            isRepostStr = strings.TrimSpace(isRepostStr)
            isRepost := isRepostStr == "true"
            postID, err := client.CreatePost(content, subreddit, isRepost)
            if err != nil {
                fmt.Printf("Error creating post: %v\n", err)
            } else {
                fmt.Printf("Post created successfully with ID: %d\n", postID)
            }

        /**case "3":
            fmt.Print("Enter comment content: ")
            content, _ := reader.ReadString('\n')
            content = strings.TrimSpace(content)
            fmt.Print("Enter parent post/comment ID: ")
            parentIDStr, _ := reader.ReadString('\n')
            parentIDStr = strings.TrimSpace(parentIDStr)
            parentID := 0
            fmt.Sscanf(parentIDStr, "%d", &parentID)
            err := client.CreateComment(content, parentID)
            if err != nil {
                fmt.Printf("Error creating comment: %v\n", err)
            } else {
                fmt.Println("Comment created successfully")
            }**/
        /**case "3":
            fmt.Print("Enter comment content: ")
            content, _ := reader.ReadString('\n')
            content = strings.TrimSpace(content)
            fmt.Print("Enter parent post/comment ID: ")
            parentIDStr, _ := reader.ReadString('\n')
            parentIDStr = strings.TrimSpace(parentIDStr)
            parentID := 0
            fmt.Sscanf(parentIDStr, "%d", &parentID)
            commentID, err := client.CreateComment(content, parentID)
            if err != nil {
                fmt.Printf("Error creating comment: %v\n", err)
            } else {
                fmt.Printf("Comment created successfully with ID: %d\n", commentID)
            }
            fmt.Printf("Comment created successfully with ID: %d\n", commentID)**/
        case "3":
            fmt.Print("Enter comment content: ")
            content, _ := reader.ReadString('\n')
            content = strings.TrimSpace(content)
            fmt.Print("Enter parent post/comment ID: ")
            parentIDStr, _ := reader.ReadString('\n')
            parentIDStr = strings.TrimSpace(parentIDStr)
            parentID := 0
            fmt.Sscanf(parentIDStr, "%d", &parentID)
            commentID, err := client.CreateComment(content, parentID)
            if err != nil {
                fmt.Printf("Error creating comment: %v\n", err)
            } else {
                fmt.Printf("Comment created successfully with ID: %d\n", commentID)
            }
        
        
        

        case "4":
            fmt.Print("Enter item ID to vote on: ")
            itemIDStr, _ := reader.ReadString('\n')
            itemIDStr = strings.TrimSpace(itemIDStr)
            itemID := 0
            fmt.Sscanf(itemIDStr, "%d", &itemID)
            fmt.Print("Upvote? (true/false): ")
            upvoteStr, _ := reader.ReadString('\n')
            upvoteStr = strings.TrimSpace(upvoteStr)
            upvote := upvoteStr == "true"
            err := client.Vote(itemID, upvote)
            if err != nil {
                fmt.Printf("Error voting: %v\n", err)
            } else {
                fmt.Println("Vote recorded successfully")
            }

        /**case "5":
            feed, err := client.GetFeed()
            if err != nil {
                fmt.Printf("Error getting feed: %v\n", err)
            } else if len(feed) == 0 {
                fmt.Println("Your feed is empty.")
            } else {
                fmt.Println("Your feed:")
                for _, post := range feed {
                    fmt.Printf("- [%d] %s (in %s, votes: %d)\n", post.ID, post.Content, post.Subreddit, post.Votes)
                }
            }
        **/
    case "5":
        feed, err := client.GetFeed()
        if err != nil {
            fmt.Printf("Error getting feed: %v\n", err)
        } else if len(feed) == 0 {
            fmt.Println("Your feed is empty.")
        } else {
            fmt.Println("Your feed:")
            for _, post := range feed {
                fmt.Printf("- [%d] %s (in %s, votes: %d)\n", post.ID, post.Content, post.Subreddit, post.Votes)
                if len(post.Comments) > 0 {
                    fmt.Println("  Comments:")
                    for _, comment := range post.Comments {
                        fmt.Printf("    - [%d] %s (votes: %d)\n", comment.ID, comment.Content, comment.Votes)
                    }
                }
            }
        }
    

        case "6":
            fmt.Print("Enter recipient user ID: ")
            to, _ := reader.ReadString('\n')
            to = strings.TrimSpace(to)
            fmt.Print("Enter message content: ")
            content, _ := reader.ReadString('\n')
            content = strings.TrimSpace(content)
            err := client.SendDirectMessage(to, content)
            if err != nil {
                fmt.Printf("Error sending message: %v\n", err)
            } else {
                fmt.Println("Message sent successfully")
            }

        case "7":
            messages, err := client.GetDirectMessages()
            if err != nil {
                fmt.Printf("Error getting messages: %v\n", err)
            } else if len(messages) == 0 {
                fmt.Println("You have no new messages.")
            } else {
                fmt.Println("Your messages:")
                for _, msg := range messages {
                    fmt.Printf("From %s: %s\n", msg.From, msg.Content)
                }
            }
        case "8":
            fmt.Print("Enter subreddit name to join: ")
            subreddit, _ := reader.ReadString('\n')
            subreddit = strings.TrimSpace(subreddit)
            err := client.JoinSubreddit(subreddit)
            if err != nil {
                fmt.Printf("Error joining subreddit: %v\n", err)
            } else {
                fmt.Printf("Successfully joined subreddit: %s\n", subreddit)
            }

        case "9":
            fmt.Print("Enter subreddit name to leave: ")
            subreddit, _ := reader.ReadString('\n')
            subreddit = strings.TrimSpace(subreddit)
            err := client.LeaveSubreddit(subreddit)
            if err != nil {
                fmt.Printf("Error leaving subreddit: %v\n", err)
            } else {
                fmt.Printf("Successfully left subreddit: %s\n", subreddit)
            }
        case "10":
            fmt.Println("Exiting...")
            return
        default:
            fmt.Println("Invalid choice. Please try again.")
        }
    }
}