package task

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"templrjs/pkg/duckdb"
	"time"

	log "github.com/sirupsen/logrus"
)

type InstagramResponse struct {
	Data struct {
		User struct {
			EdgeOwnerToTimelineMedia struct {
				Count int `json:"count"`
				Edges []struct {
					Node struct {
						DisplayURL         string `json:"display_url"`
						EdgeMediaToCaption struct {
							Edges []struct {
								Node struct {
									Text string `json:"text"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_caption"`
						Shortcode        string `json:"shortcode"`
						TakenAtTimestamp int64  `json:"taken_at_timestamp"`
						ThumbnailSrc     string `json:"thumbnail_src"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
		} `json:"user"`
	} `json:"data"`
}

type Post struct {
	Text          string `json:"text"`
	Thumbnail_src string `json:"thumbnail_src"`
	Display_url   string `json:"display_url"`
	Short_code    string `json:"short_code"`
	Base64        string `json:"base64"`
	Uploaded_at   string `json:"uploaded_at"`
}

func Scrape(profileID string, first int, entity string) {

	if first == 0 {
		first = 10 // default value
	}

	url := fmt.Sprintf("https://www.instagram.com/graphql/query/?query_id=17888483320059182&variables={\"id\":\"%s\",\"first\":%d,\"after\":null}", profileID, first)
	log.Info(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Error("Error fetching from Instagram: %s", err)
		return
	}
	log.Info("Response received")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Debug("Error reading response body: %s", err)
		log.Error("Error reading response body: %s", err)
		return
	}

	var instaResp InstagramResponse
	err = json.Unmarshal(body, &instaResp)
	if err != nil {
		log.Error("Error unmarshalling Instagram response: %s", err)
		return
	}
	log.Info("Response unmarshalled")
	log.Info(instaResp)
	log.Info("Count: %d", instaResp.Data.User.EdgeOwnerToTimelineMedia.Count)

	var posts []Post
	log.Info("Iterating over Instagram response")

	for i, edge := range instaResp.Data.User.EdgeOwnerToTimelineMedia.Edges {
		log.Info("Loop #%d", i)
		imgURL := edge.Node.DisplayURL
		imgResp, err := http.Get(imgURL)
		if err != nil {
			log.Error("Error fetching image: %s", err)
			continue
		}
		imgBytes, err := io.ReadAll(imgResp.Body)
		imgResp.Body.Close()
		if err != nil {
			log.Error("Error reading image bytes: %s", err)
			continue
		}
		imgBase64 := base64.StdEncoding.EncodeToString(imgBytes)
		post := Post{
			Text:          edge.Node.EdgeMediaToCaption.Edges[0].Node.Text,
			Thumbnail_src: edge.Node.ThumbnailSrc,
			Display_url:   edge.Node.DisplayURL,
			Short_code:    edge.Node.Shortcode,
			Base64:        imgBase64,
			Uploaded_at:   time.Unix(edge.Node.TakenAtTimestamp, 0).Format("2006-01-02"),
		}
		posts = append(posts, post)

		if entity != "" {
			log.Info("Creating entity map[string]interface{}")
			postMap := map[string]interface{}{
				"text":          post.Text,
				"thumbnail_src": post.Thumbnail_src,
				"display_url":   post.Display_url,
				"short_code":    post.Short_code,
				"base64":        post.Base64,
				"uploaded_at":   post.Uploaded_at,
			}
			log.Info("Creating entity")
			duckdb.CreateEntityCore(entity, postMap)
		}
	}

	// If restAPIURL is empty, save the results to a JSON file
	log.Info("Preparing to save data to file")
	if entity == "" {
		log.Info("Saving data to file")
		// Convert the posts slice to a JSON string and print
		jsonOutput, err := json.Marshal(posts)

		outputFilename := fmt.Sprintf("instagram_%s_%d.json", profileID, time.Now().Unix())

		err = os.WriteFile(outputFilename, jsonOutput, 0644)
		if err != nil {
			log.Error("Error writing JSON data to file: %s", err)
			return
		}
		log.Info("Data saved to %s\n", outputFilename)
	}
}
