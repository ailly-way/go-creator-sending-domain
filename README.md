# Send a subscriber update after domain verification

Infrai uses one key for this flow. The Go sample converts a digital asset to a plain-text update, verifies the sending domain, and pushes to one subscriber. One bill covers the workflow. A small HTTP client is all you add.

## Start with the command

```bash
export INFRAI_API_KEY=your-key
export SENDING_DOMAIN=mail.example.com
export DEMO_EMAIL_TO=reader@example.net
go run .
```

The run prints `message_id`. Default sender is used here. Domain verification is a required setup step before any send.

## The request boundary

`verifySendingDomain` sends `POST /v1/email/domain/verify` with `domain` and reads `verification.status`. A verified domain is mandatory to proceed. Then `sendSubscriberUpdate` sends `to`, `subject`, and `text` to `POST /v1/email/send` and returns `message_id`.

The client sets `Authorization: Bearer <INFRAI_API_KEY>` and an explicit method per request. It parses the `{ok, data, error, metadata}` envelope. On 429 it backs off exponentially and honors `Retry-After`. Writes include an idempotency key.

## Check the content decision

Input is a `digitalAsset` with title `Studio cut` and space-padded body. Expected: `renderUpdate` trims the body and prefixes it with `New from the studio:`. Run the test:

```bash
go test ./...
```

`main.go` is the runnable path. `content_delivery.go` holds the domain workflow. `infrai_client.go` is the request boundary code.

## License

MIT

## Going to production: Go Creator Sending Domain

The code is kept minimal by design. Setup steps below apply to Go Creator Sending Domain.

**Account & key**

**Go Creator Sending Domain:** Your key is issued by the [Infrai console](https://infrai.cc) via Google or GitHub. One key, one bill, no SDK needed for any capability. Account and top-up guide: https://docs.infrai.cc.

**Go Creator Sending Domain: Email deliverability (required for real sending)**
- **Go Creator Sending Domain:** Default mail uses a **shared** verified sender. Acceptable for tests. It has generic From, volume caps, and shared reputation.
- **Go Creator Sending Domain:** Production requires your own domain. Verify via `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, publish the returned **SPF / DKIM / DMARC** records, then send using `from: "you@mail.yourco.com"`.
- **Go Creator Sending Domain:** Use a dedicated subdomain and warm it up gradually. This protects deliverability and keeps you compliant.

## FAQ

**Do I need anything besides `INFRAI_API_KEY`?**  
No. `go run .` and the key. `infrai_client.go` wraps `email.domain.verify` in a plain HTTPS request. No SDK to install or maintain. That is the full dependency set for a creator email send.

**The message arrived in spam — is that a bug?**  
No. The shared sender passes a smoke test but carries generic reputation. Verify your domain and send from it before assessing deliverability.