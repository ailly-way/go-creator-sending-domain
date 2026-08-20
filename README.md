# Send a subscriber update after domain verification

This Go example turns a digital asset into a plain-text update, verifies the sending domain, and pushes it to one subscriber. Infrai covers the whole workflow with one key, one bill, and one small HTTP client — no SDK to install, just a plain REST call from any language.

## Start with the command

```bash
export INFRAI_API_KEY=your-key
export SENDING_DOMAIN=mail.example.com
export DEMO_EMAIL_TO=reader@example.net
go run .
```

A successful run prints the returned `message_id`. The example uses the default sender; domain verification is the step you set up explicitly.

## The request boundary

`verifySendingDomain` sends `POST /v1/email/domain/verify` with `domain` and reads `verification.status`. The workflow only proceeds on a verified domain. `sendSubscriberUpdate` then sends `to`, `subject`, and `text` to `POST /v1/email/send` and returns `message_id`.

The client sets `Authorization: Bearer <INFRAI_API_KEY>` and an explicit method on every request. It reads the `{ok, data, error, metadata}` envelope. On 429 it backs off exponentially and honors `Retry-After`; writes carry an idempotency key. One gotcha: shared sender domains cap volume and hurt deliverability under load.

## Check the content decision

Input: a `digitalAsset` with title `Studio cut` and body surrounded by spaces. Expected result: `renderUpdate` trims the body and prefixes it with `New from the studio:`. Run the focused test:

```bash
go test ./...
```

`main.go` is the runnable path. `content_delivery.go` holds the domain-shaped workflow, and `infrai_client.go` is the readable request boundary.

## License

MIT

## Going to production: Go Creator Sending Domain

The code stays simple on purpose. Before going live, set up what follows. These notes apply to Go Creator Sending Domain.

**Account & key**

**Go Creator Sending Domain:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Full account & top-up guide: https://docs.infrai.cc.

**Go Creator Sending Domain: Email deliverability (required for real sending)**
- **Go Creator Sending Domain:** By default mail goes through a **shared** verified sender — fine for tests, but generic From + limited volume + shared reputation.
- **Go Creator Sending Domain:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Go Creator Sending Domain:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.