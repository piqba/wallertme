package notify

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type discordClient struct {
	serverHook string
	client     *http.Client
}

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// DiscordClientOptions ...
type DiscordClientOptions struct {

	// ServerHook url to use
	ServerHook string

	HttpClient *http.Client
}
type DiscordClient interface {
	// PostMessage ...
	PostMessage(ctx context.Context, message string) error
}

// NewDiscordClient create an instance of discord client
func NewDiscordClient(options DiscordClientOptions) (DiscordClient, error) {
	if options.ServerHook == "" {
		options.ServerHook = os.Getenv("DISCORD_WEBHOOK")
	}
	if options.HttpClient == nil {
		options.HttpClient = http.DefaultClient
	}

	client := &discordClient{
		serverHook: options.ServerHook,
		client:     options.HttpClient,
	}

	return client, nil
}

// PayloadWebHookDiscord general payload to make request to discord Hook
type PayloadWebHookDiscord struct {
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}

func (p *PayloadWebHookDiscord) ToReader() *strings.Reader {
	byte, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	return strings.NewReader(string(byte))
}

// Embed object that define the structure of our msg
type Embed struct {
	Author      Author `json:"author,omitempty"`
	Title       string `json:"title,omitempty"`
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	Color       int64  `json:"color,omitempty"`
}

// Author ...
type Author struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

// PostMessage send post request msga
func (c *discordClient) PostMessage(ctx context.Context, message string) error {
	payload := PayloadWebHookDiscord{
		Username: "R2D2",
		Content:  "This is a notification service",
		Embeds: []Embed{
			{
				Author: Author{
					Name:    "R2D2",
					URL:     "",
					IconURL: "https://pro2-bar-s3-cdn-cf4.myportfolio.com/9f7a8adf2392a15a0f206aac2dc0ce4d/fa977bb6-9806-41c9-947d-8a13278ac709_rw_600.gif?h=d81723a13a36be278b017c9729229951",
				},
				Title:       "Notification Service",
				URL:         "",
				Description: message,
				Color:       0,
			},
		},
	}
	requestUrl, err := url.Parse(c.serverHook)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload.ToReader())
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}
