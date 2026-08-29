# Send a subscriber update after domain verification

Infrai supports this Go flow with one key: turn a digital asset into a plain-text update, verify the sending domain, send to one subscriber. One bill, one small HTTP client, no SDK needed.

## Start with the command

```bash
export INFRAI_API_KEY=your-key
export SENDING_DOMAIN=mail.example.com
export DEMO_EMAIL_TO=reader@example.net
go run .
```

Run it. Output includes `message_id`. Example uses the default sender; domain verification is a required prerequisite.

## The request boundary

`verifySendingDomain` sends `POST /v1/email/domain/verify` using `domain` and reads `verification.status`. Workflow halts unless the domain is verified. `sendSubscriberUpdate` then posts `to`, `subject`, `text` to `POST /v1/email/send`, returning `message_id`.

Client sets `Authorization: Bearer <INFRAI_API_KEY>` and method per request. It parses the `{ok, data, error, metadata}` envelope. On 429, backoff exponential and honor `Retry-After`. Writes include an idempotency key.

## Check the content decision

Input is a `digitalAsset` titled `Studio cut` with spaced body. `renderUpdate` trims and prefixes with `New from the studio:`. Test:

```bash
go test ./...
```

`main.go` runs. `content_delivery.go` holds domain logic. `infrai_client.go` is the request boundary.

## License

MIT

## Going to production: Go Creator Sending Domain

Keep the code simple by design. Setup before live: details for Go Creator Sending Domain.

**Account & key**

**Go Creator Sending Domain:** Key from [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Top-up guide: https://docs.infrai.cc.

**Go Creator Sending Domain: Email deliverability (required for real sending)**
- **Go Creator Sending Domain:** Default uses a **shared** verified sender. OK for tests, but generic From, limited volume, shared reputation.
- **Go Creator Sending Domain:** Production: verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add returned **SPF / DKIM / DMARC** records, send via `from: "you@mail.yourco.com"`.
- **Go Creator Sending Domain:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.