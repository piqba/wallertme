package notify

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type webHookClient struct {
	serverHook string
	client     *http.Client
}

type WebHookClientOptions struct {

	// ServerHook url to use
	ServerHook string

	HttpClient *http.Client
}

type WebHookClient interface {
	// PostMessage ...
	PostMessage(ctx context.Context, message string) error
}

func NewWebHookClient(options WebHookClientOptions) (WebHookClient, error) {
	if options.ServerHook == "" {
		options.ServerHook = os.Getenv("CUSTOM_WEBHOOK")
	}
	if options.HttpClient == nil {
		options.HttpClient = http.DefaultClient
	}

	client := &webHookClient{
		serverHook: options.ServerHook,
		client:     options.HttpClient,
	}

	return client, nil
}

type PayloadWebHook struct {
	Username string `json:"username,omitempty"`
	Content  string `json:"content,omitempty"`
	Message  string `json:"message,omitempty"`
}

func (p *PayloadWebHook) ToReader() *strings.Reader {
	byte, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	return strings.NewReader(string(byte))
}

// PostMessage send post request msga
func (c *webHookClient) PostMessage(ctx context.Context, message string) error {
	_, span := otel.Tracer(nameNotifierWebHook).Start(ctx, "webhook.PostMessage")
	defer span.End()
	payload := PayloadWebHook{
		Username: "R2D2",
		Content:  "This is a notification service",
		Message:  message,
	}
	requestUrl, err := url.Parse(c.serverHook)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), payload.ToReader())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	defer res.Body.Close()
	span.SetAttributes(attribute.String("notifier.webhook.client", res.Status))
	return nil
}
