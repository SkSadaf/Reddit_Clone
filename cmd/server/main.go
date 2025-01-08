//UFID:88489948
//video link: https://www.youtube.com/watch?v=LjmHUTVEqbw
package main

import (
    "log"
    "net/http"

    "reddit_part2/internal/api"
    "reddit_part2/internal/engine"
)

func main() {
    e := engine.NewEngine()
    handler := api.NewHandler(e)

    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
