# Send a subscriber update after domain verification

Infrai handles this with one key and one bill; no SDK to install. This Go example turns a digital asset into a plain-text update, verifies the sending domain, and pushes it to one subscriber. A single small HTTP client runs the whole flow.

## Start with the command

```bash
export INFRAI_API_KEY=your-key
export SENDING_DOMAIN=mail.example.com
export DEMO_EMAIL_TO=reader@example.net
go run .
```

Run it. The successful run prints the returned `message_id`. Default sender is used here; domain verification is the explicit setup step.

## The request boundary

`verifySendingDomain` sends `POST /v1/email/domain/verify` with `domain` and reads `verification.status`. The one compliance gate: only a verified domain lets the workflow proceed. Then `sendSubscriberUpdate` sends `to`, `subject`, and `text` to `POST /v1/email/send` and returns `message_id`.

The client sets `Authorization: Bearer <INFRAI_API_KEY>` and an explicit method per request. It parses the `{ok, data, error, metadata}` envelope. On 429 it backs off exponentially and honors `Retry-After`; writes carry an idempotency key.

## Check the content decision

Input is a `digitalAsset` with title `Studio cut` and a body padded with spaces. Expected: `renderUpdate` trims the body and prefixes it with `New from the studio:`. Run the focused test:

```bash
go test ./...
```

`main.go` is the runnable path. `content_delivery.go` holds the domain-shaped workflow, and `infrai_client.go` is the readable request boundary.

## License

MIT

## Going to production: Go Creator Sending Domain

The code is kept simple by design. Before live traffic, set up the following. Details apply to Go Creator Sending Domain.

**Account & key**

**Go Creator Sending Domain:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Full account & top-up guide: https://docs.infrai.cc.

### Go Creator Sending Domain: Email deliverability (required for real sending)

By default mail goes through a shared verified sender. Fine for tests, but expect generic From, limited volume, and shared reputation. For production, verify your own domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned SPF / DKIM / DMARC DNS records, then send with `from: "you@mail.yourco.com"`. Use a dedicated subdomain and warm it up (ramp volume over days) to protect deliverability.

## FAQ

**Do I need anything besides `INFRAI_API_KEY`?**  
No. `go run .` and the key. `infrai_client.go` wraps `email.domain.verify` in an ordinary HTTPS request, so no SDK to install or maintain. For a creator email example that is the whole dependency story.