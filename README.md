# Reddit_Clone
Contents of this file: 

• Project overview 

• How to run? 

• Implementation Details 

• Key Components

Project Overview: 

This project implements a Reddit Clone using Go, featuring a core engine with REST API 
functionality and a client simulator. The system replicates key features of Reddit, allowing 
users to interact with content across various topics. 

How to run? 
First start the server by navigating to cmd/server directory and running the following 
command: 
>go run main.go 
Then start various clients by navigating to cmd/client directory and running the same 
command again, we can run multiple clients in multiple terminals. 
>go run main.go


Implementation Details: 

Account Registration 
• Users can create new accounts

Subreddit Management 
• Create subreddits 
• Join and leave subreddits 
• Post text content in subreddits 
• Comment on posts and other comments (hierarchical structure)

Voting 
• Upvote and downvote posts and comments 
• Karma computation based on votes 

User Feed 
• Generate feed

Direct Messaging 
• Send and receive direct messages 
• Reply to messages 

Key Components: 
• Engine: Core logic handler 
• Models: Data structures for users, posts, comments, etc. 
• API Handlers: REST endpoint implementations 
• Client: Command-line interface for interacting with the API
