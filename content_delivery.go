package main

import (
	"fmt"
	"strings"
)

type digitalAsset struct{ Title, Body string }
type subscriber struct{ Email string }

func renderUpdate(asset digitalAsset) string {
	return fmt.Sprintf("New from the studio: %s\n\n%s", asset.Title, strings.TrimSpace(asset.Body))
}

func sendSubscriberUpdate(c *client, domain string, asset digitalAsset, sub subscriber) (string, error) {
	verified, err := verifySendingDomain(c, domain)
	if err != nil {
		return "", err
	}
	if !verified {
		return "", fmt.Errorf("sending domain is not verified")
	}
	var reply struct {
		MessageID string `json:"message_id"`
	}
	err = c.call("POST", "/v1/email/send", map[string]string{"to": sub.Email, "subject": asset.Title, "body": renderUpdate(asset)}, &reply)
	return reply.MessageID, err
}
