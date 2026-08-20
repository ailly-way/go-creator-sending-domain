package main

type verification struct {
	Status string `json:"status"`
}
type domainReply struct {
	Verification verification `json:"verification"`
}

func verifySendingDomain(c *client, domain string) (bool, error) {
	var reply domainReply
	if err := c.call("POST", "/v1/email/domain/verify", map[string]string{"domain": domain}, &reply); err != nil {
		return false, err
	}
	return reply.Verification.Status == "verified", nil
}
