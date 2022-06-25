package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Posts struct {
	UserId   int    `json:"userId"`
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Comments string
}

type Comments struct {
	PostId int    `json:"postId"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func main() {
	posts, err := GetPosts()
	if err != nil {
		log.Fatal(err)
	}

	comments, err := GetComments()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(posts); i++ {
		for j := 0; j < len(comments); j++ {
			if comments[j].PostId == posts[i].Id {
				posts[i].Comments = fmt.Sprintf("%s|", comments[j])
			}
		}
		log.Println(posts[i].Comments)
		posts[i].Comments = strings.Trim(posts[i].Comments, "|")
	}

	log.Println(len(posts))

	csvFile, err := os.Create("posts.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
		return
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	for _, v := range posts {
		log.Println(strings.Split(v.Comments, "|"))
		if err = csvWriter.Write(strings.Split(v.Comments, "|")); err != nil {
			log.Fatal(err)
			return
		}
	}
}

func GetPosts() ([]Posts, error) {
	b, err := Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		return nil, err
	}

	var jsonResp []Posts
	if err = json.Unmarshal(b, &jsonResp); err != nil {
		return nil, err
	}
	return jsonResp, nil
}

func GetComments() ([]Comments, error) {
	b, err := Get("https://jsonplaceholder.typicode.com/comments")
	if err != nil {
		return nil, err
	}
	var jsonResp []Comments
	if err = json.Unmarshal(b, &jsonResp); err != nil {
		return nil, err
	}
	return jsonResp, nil
}

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
