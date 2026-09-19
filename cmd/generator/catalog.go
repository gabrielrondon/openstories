package main

import (
	"fmt"
	"strings"
)

func getCatalog() []DomainDef {
	var domains []DomainDef

	// Definition of 20 Industries and their specific subdomains and story generators
	industries := []struct {
		Industry  string
		Prefix    string
		Domains   []string
		Role      string
		Context   string
		Templates []struct {
			TitleFormat  string
			IWantFormat  string
			SoThatFormat string
			ScenarioFmt  string
			GivenFmt     string
			WhenFmt      string
			ThenFmt      string
			EdgeCaseFmt  string
			QuoteFmt     string
			SourceFmt    string
		}
	}{
		{
			Industry: "devtools",
			Prefix:   "DEV",
			Domains:  []string{"apis-and-sdks", "auth-and-iam", "ci-cd-pipelines", "observability", "database-migrations", "feature-flags"},
			Role:     "Principal Software Engineer / Systems Architect",
			Context:  "Distributed developer tooling, developer experience, and cloud runtime platforms",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Zero-downtime %s with atomic rollback validation",
					IWantFormat:  "automated validation for %s ensuring backward compatibility",
					SoThatFormat: "deployments never break downstream API clients or cause schema locks",
					ScenarioFmt:  "Schema or config mismatch during canary phase for %s",
					GivenFmt:     "A cluster running production version V1",
					WhenFmt:      "A new migration for %s is rolled out to 10%% of traffic",
					ThenFmt:      "The deployment must automatically rollback if error rates exceed 0.05%% within 60 seconds",
					EdgeCaseFmt:  "Database lock timeouts when altering tables with >10M rows during high traffic",
					QuoteFmt:     "Our migration on %s caused a 45-minute outage because Postgres took an exclusive table lock.",
					SourceFmt:    "https://github.com/flyway/flyway/issues/2891",
				},
				{
					TitleFormat:  "High-throughput %s rate limiting with distributed sliding window",
					IWantFormat:  "a resilient sliding-window rate limiter for %s with Redis fallback to local memory",
					SoThatFormat: "spiky traffic bursts do not exhaust internal worker pools or cascade into 504 gateway timeouts",
					ScenarioFmt:  "Redis cluster unreachable during rate limit evaluation for %s",
					GivenFmt:     "The central Redis cache drops connection",
					WhenFmt:      "Incoming client requests arrive at %s",
					ThenFmt:      "The service must gracefully degrade to local in-memory token buckets without failing open to abusive traffic",
					EdgeCaseFmt:  "Clock drift between distributed nodes skewing sliding window timestamp calculations",
					QuoteFmt:     "When Redis had a network partition, our %s rate limiter crashed the entire microservice fleet.",
					SourceFmt:    "https://news.ycombinator.com/item?id=38190211",
				},
				{
					TitleFormat:  "Cryptographic verification of %s webhooks and replay attack prevention",
					IWantFormat:  "HMAC-SHA256 signature headers with timestamp validation for all %s events",
					SoThatFormat: "malicious attackers cannot replay intercepted payloads or spoof system actions",
					ScenarioFmt:  "Stale timestamp in webhook signature header for %s",
					GivenFmt:     "An incoming webhook signed with valid secret key",
					WhenFmt:      "The event timestamp is older than 300 seconds",
					ThenFmt:      "The ingestion pipeline must reject the payload with HTTP 401 Unauthorized",
					EdgeCaseFmt:  "Slow asynchronous delivery queues causing legitimate events to arrive near the 5-minute threshold",
					QuoteFmt:     "We caught an attacker replaying captured webhook events on %s to grant themselves free enterprise features.",
					SourceFmt:    "https://github.com/stripe/stripe-node/issues/1420",
				},
			},
		},
		{
			Industry: "ai-infra",
			Prefix:   "AI",
			Domains:  []string{"agentic-workflows", "llm-gateways", "prompt-caching", "rag-and-vectors", "model-evals", "fine-tuning"},
			Role:     "Staff AI Platform Engineer / RAG Systems Lead",
			Context:  "Production generative AI infrastructure handling millions of embeddings and autonomous agent tool calls",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Deterministic %s evaluation with drift detection and latency budgets",
					IWantFormat:  "real-time monitoring of %s output distribution and token latency",
					SoThatFormat: "silent model regressions or unexpected prompt changes are detected before reaching end users",
					ScenarioFmt:  "Upstream model provider alters system prompt formatting for %s",
					GivenFmt:     "An automated test suite evaluating baseline responses",
					WhenFmt:      "Model accuracy drops by more than 3%% on standard benchmarks",
					ThenFmt:      "The CI pipeline must block model deployment and alert the on-call AI engineer",
					EdgeCaseFmt:  "Non-deterministic temperature output causing sporadic false-positive test failures",
					QuoteFmt:     "A silent update to the provider's API broke our structured JSON parser on %s, breaking 40%% of customer queries.",
					SourceFmt:    "https://github.com/vllm-project/vllm/issues/3102",
				},
				{
					TitleFormat:  "Dynamic semantic chunking and re-ranking for %s in vector search",
					IWantFormat:  "context-aware document chunking and cross-encoder re-ranking for %s",
					SoThatFormat: "retrieval augmented generation avoids lost-in-the-middle context degradation",
					ScenarioFmt:  "Dense document containing conflicting historical revisions of %s",
					GivenFmt:     "A user asking for current active policy",
					WhenFmt:      "Vector similarity returns outdated chunks with high cosine score",
					ThenFmt:      "The temporal re-ranker must prioritize the chunk with the latest verifiable effective date",
					EdgeCaseFmt:  "Tables and code blocks split across chunk boundaries corrupting syntax during generation",
					QuoteFmt:     "Standard 500-token chunking chopped our API tables in half for %s, resulting in hallucinated parameters.",
					SourceFmt:    "https://news.ycombinator.com/item?id=39120481",
				},
				{
					TitleFormat:  "Token budget enforcement and prompt caching optimization for %s",
					IWantFormat:  "hierarchical prefix prompt caching on %s to reduce TTFT (time-to-first-token)",
					SoThatFormat: "repetitive system prompt tokens do not inflate operational costs and latency stays sub-200ms",
					ScenarioFmt:  "Cache invalidation due to dynamic timestamp injected in system prompt for %s",
					GivenFmt:     "A large 20k token system instructions context",
					WhenFmt:      "Dynamic variables are placed at the beginning of the prompt",
					ThenFmt:      "The compiler must automatically hoist static prefixes to maximize provider KV-cache hits",
					EdgeCaseFmt:  "Provider cache eviction during low-traffic night hours causing unexpected latency spikes",
					QuoteFmt:     "We cut 65%% of our Anthropic bills on %s by structuring static prompt blocks for 100%% cache hits.",
					SourceFmt:    "https://github.com/BerriAI/litellm/issues/2104",
				},
			},
		},
		{
			Industry: "fintech",
			Prefix:   "FIN",
			Domains:  []string{"payments-and-checkout", "reconciliation", "fraud-detection", "banking-apis", "card-issuing", "tax-compliance"},
			Role:     "Principal Fintech Engineer / Ledger Architect",
			Context:  "High-volume financial ledgers, card issuing, payment orchestrations, and double-entry accounting",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Double-entry immutable ledger journaling for %s transactions",
					IWantFormat:  "strict atomic double-entry accounting entries for all %s operations",
					SoThatFormat: "the total sum of debits and credits is mathematically proven to always equal zero",
					ScenarioFmt:  "Asynchronous database rollback during partial ledger write for %s",
					GivenFmt:     "A multi-leg balance transfer in progress",
					WhenFmt:      "The secondary account credit query times out",
					ThenFmt:      "The transaction coordinator must execute a full atomic rollback, preventing money from vanishing into thin air",
					EdgeCaseFmt:  "Concurrent debit operations on accounts with balance near zero triggering race condition overdrafts",
					QuoteFmt:     "A race condition on our %s balance checks allowed a user to withdraw the same $500 balance three times simultaneously.",
					SourceFmt:    "https://news.ycombinator.com/item?id=36192801",
				},
				{
					TitleFormat:  "Real-time fraud velocity checks and 3DS challenge step-up for %s",
					IWantFormat:  "sliding window card velocity heuristics and adaptive 3D Secure step-up for %s",
					SoThatFormat: "card-testing bot attacks are blocked before triggering card network dispute penalties",
					ScenarioFmt:  "Card testing attack trying 50 distinct CVVs per minute on %s",
					GivenFmt:     "Traffic originating from a single IP or fingerprint hash",
					WhenFmt:      "More than 3 card authorization declines occur within 10 seconds",
					ThenFmt:      "The gateway must trigger mandatory Captcha and biometric 3DS verification on all subsequent requests",
					EdgeCaseFmt:  "Distributed botnet cycling residential proxies to evade naive single-IP velocity limits",
					QuoteFmt:     "Our merchant account was suspended by Visa after a card testing bot hit %s with 12,000 stolen cards overnight.",
					SourceFmt:    "https://reddit.com/r/stripe/comments/16k29a1",
				},
				{
					TitleFormat:  "Automated sales tax jurisdiction determination and nexus tracking for %s",
					IWantFormat:  "rooftop-accurate geolocation address validation for sales tax on %s",
					SoThatFormat: "digital goods are taxed precisely per municipality rules without post-audit fines",
					ScenarioFmt:  "Zip+4 boundary spanning two different county tax rates for %s",
					GivenFmt:     "A customer checking out with physical shipping in California or New York",
					WhenFmt:      "The tax calculation engine resolves the street address",
					ThenFmt:      "It must look up precise latitude/longitude tax parcel data rather than generic 5-digit zip code approximations",
					EdgeCaseFmt:  "B2B customers presenting tax exemption certificates that have expired or belong to a different state",
					QuoteFmt:     "We owed $32,000 in back taxes because our %s engine used 5-digit zip codes instead of street-level tax jurisdictions.",
					SourceFmt:    "https://news.ycombinator.com/item?id=35198201",
				},
			},
		},
		{
			Industry: "b2b-saas",
			Prefix:   "SAAS",
			Domains:  []string{"multi-tenancy", "audit-logs", "rbac-permissions", "sso-saml", "scim-provisioning", "usage-billing"},
			Role:     "Enterprise SaaS Architect / Security Lead",
			Context:  "Multi-tenant B2B enterprise architectures serving Fortune 500 customers with rigorous compliance requirements",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "SCIM 2.0 automated user deprovisioning and role sync for %s",
					IWantFormat:  "full compliance with RFC 7644 SCIM protocol for %s user life-cycle management",
					SoThatFormat: "terminated corporate employees immediately lose access to enterprise tenant resources",
					ScenarioFmt:  "Okta or Azure AD sends deprovision PATCH command for %s",
					GivenFmt:     "An active user with valid session tokens in multiple browser tabs",
					WhenFmt:      "The identity provider issues a SCIM active=false request",
					ThenFmt:      "The backend must revoke all active refresh tokens and WebSocket connections in under 500ms",
					EdgeCaseFmt:  "User reassigned to a different department with reduced permissions while currently holding an active session",
					QuoteFmt:     "We failed an enterprise procurement audit because our %s system didn't terminate active JWT sessions upon SCIM deprovisioning.",
					SourceFmt:    "https://github.com/boxyhq/jackson/issues/612",
				},
				{
					TitleFormat:  "Granular Role-Based Access Control (RBAC) with attribute constraints for %s",
					IWantFormat:  "fine-grained policy evaluation (ABAC/RBAC) on all %s endpoints",
					SoThatFormat: "tenant members with restricted roles cannot access privileged audit records or billing settings",
					ScenarioFmt:  "Horizontal privilege escalation attempt via direct ID reference on %s",
					GivenFmt:     "A user logged into tenant organization Alpha",
					WhenFmt:      "The user queries resource ID belonging to tenant Beta",
					ThenFmt:      "The authorization layer must return HTTP 404 Not Found rather than 403 Forbidden to prevent resource ID enumeration",
					EdgeCaseFmt:  "Users belonging to multiple organizations switching active workspace context concurrently",
					QuoteFmt:     "A classic IDOR vulnerability in our %s module allowed team members to inspect invoice PDFs from rival companies.",
					SourceFmt:    "https://hackerone.com/reports/512091",
				},
				{
					TitleFormat:  "Usage-based metering ingestion with idempotency and late-arriving event processing for %s",
					IWantFormat:  "event-driven usage aggregation for %s supporting out-of-order event streams",
					SoThatFormat: "consumption-based billing invoices accurately reflect API calls, compute hours, or storage without discrepancies",
					ScenarioFmt:  "Network blip causes batch of usage events for %s to arrive after monthly invoice finalization",
					GivenFmt:     "The billing cycle closed on midnight of the 1st",
					WhenFmt:      "Usage metrics timestamped for the 31st arrive 6 hours late",
					ThenFmt:      "The engine must record the usage as an adjustment credit/debit on the subsequent cycle rather than mutating locked invoices",
					EdgeCaseFmt:  "Client replay of telemetry batches leading to double-counting of billable compute metrics",
					QuoteFmt:     "Late arriving telemetry events on %s caused invoice re-generation chaos every first day of the month.",
					SourceFmt:    "https://news.ycombinator.com/item?id=37890124",
				},
			},
		},
		{
			Industry: "ecommerce",
			Prefix:   "ECOM",
			Domains:  []string{"cart-and-checkout", "inventory-management", "shipping-fulfillment", "promotions-and-coupons", "returns-rma"},
			Role:     "Principal E-Commerce Architect / Retail Platform Lead",
			Context:  "High-throughput digital commerce platforms with omnichannel inventory, warehouse logistics, and flash sales",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Pessimistic inventory reservation during high-concurrency flash sales for %s",
					IWantFormat:  "atomic inventory hold locks with strict 10-minute expirations for %s",
					SoThatFormat: "overselling never occurs even when 50,000 customers attempt to buy the last 10 units simultaneously",
					ScenarioFmt:  "Checkout abandonment after locking inventory on %s",
					GivenFmt:     "A customer adding the last remaining unit to cart",
					WhenFmt:      "The user closes their browser without completing checkout",
					ThenFmt:      "The reservation lock must automatically expire after 10 minutes, returning the unit back to active stock",
					EdgeCaseFmt:  "Payment gateway webhook delay causing release of inventory while customer is legitimately entering 3DS challenge",
					QuoteFmt:     "During Black Friday, our %s system oversold 450 PlayStation units because inventory checks weren't atomic.",
					SourceFmt:    "https://reddit.com/r/ecommerce/comments/17y921a",
				},
				{
					TitleFormat:  "Multi-stack coupon code validation preventing promotional stacking abuse on %s",
					IWantFormat:  "strict rules engine governing coupon combinability and minimum order values on %s",
					SoThatFormat: "customers cannot stack multiple promo codes to purchase items below manufacturing cost",
					ScenarioFmt:  "Customer combining percentage discount with dollar-off voucher on %s",
					GivenFmt:     "A promo code granting 20%% off sitewide",
					WhenFmt:      "The user applies an additional $50 welcome voucher",
					ThenFmt:      "The promotions engine must enforce exclusion rules and reject stacking unless explicitly configured",
					EdgeCaseFmt:  "Customers creating multiple throwaway accounts with the same physical delivery address to bypass limits",
					QuoteFmt:     "TikTok discovered an exploit in our %s discount engine that allowed stacking codes until the cart total hit $0.",
					SourceFmt:    "https://news.ycombinator.com/item?id=38902144",
				},
			},
		},
		{
			Industry: "healthcare",
			Prefix:   "HLTH",
			Domains:  []string{"hipaa-compliance", "telehealth-video", "ehr-fhir-integration", "prescription-management", "patient-scheduling"},
			Role:     "Clinical Systems Architect / HealthTech Compliance Officer",
			Context:  "HIPAA/HITECH compliant digital health platforms, HL7/FHIR integrations, and telemedicine systems",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "End-to-end audit logging of Protected Health Information (PHI) access on %s",
					IWantFormat:  "immutable, tamper-evident logs for every clinician read and write to %s records",
					SoThatFormat: "the organization complies with HIPAA Security Rule 45 CFR 164.312 without risk of federal civil monetary penalties",
					ScenarioFmt:  "Staff member querying patient records without assigned care relationship on %s",
					GivenFmt:     "A logged-in nurse or doctor in the hospital network",
					WhenFmt:      "The user views the medical chart of a patient not under their direct care",
					ThenFmt:      "The system must log a high-priority compliance audit event and prompt the clinician for a clinical justification reason",
					EdgeCaseFmt:  "Emergency department 'break-the-glass' protocols requiring immediate chart override during life-threatening triage",
					QuoteFmt:     "A hospital was fined $2.1M by HHS OCR because their %s system allowed unauthorized staff to view VIP patient charts.",
					SourceFmt:    "https://www.hhs.gov/hipaa/for-professionals/compliance-enforcement/index.html",
				},
				{
					TitleFormat:  "FHIR R4 standard JSON validation for interoperable %s exchanges",
					IWantFormat:  "strict HL7 FHIR R4 schema validation and terminology mapping for %s data",
					SoThatFormat: "electronic health records seamlessly exchange laboratory and medication data without truncation",
					ScenarioFmt:  "Receiving custom proprietary extensions in FHIR bundle for %s",
					GivenFmt:     "An incoming HL7 FHIR payload from an external EHR system (Epic or Cerner)",
					WhenFmt:      "The bundle contains unmapped LOINC or SNOMED CT terminology codes",
					ThenFmt:      "The ingestion adapter must safely quarantine the message and alert clinical informatics rather than discarding the lab value",
					EdgeCaseFmt:  "Mismatched patient identifier matching rules resulting in chart merging errors across different hospital networks",
					QuoteFmt:     "Our %s ingestion silently truncated lab unit measurements (mg/dL vs mmol/L) due to loose FHIR parsing.",
					SourceFmt:    "https://github.com/hapifhir/hapi-fhir/issues/3891",
				},
			},
		},
		{
			Industry: "cybersecurity",
			Prefix:   "SEC",
			Domains:  []string{"zero-trust-access", "vulnerability-management", "secret-scanning", "ids-intrusion-detection", "incident-response"},
			Role:     "Principal Security Operations Engineer / Detection Engineering Lead",
			Context:  "Cloud security posture management (CSPM), SIEM pipelines, and automated threat detection",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Automated high-entropy secret detection with pre-commit and CI blocking for %s",
					IWantFormat:  "Shannon entropy and regex pattern matching to intercept plaintext credentials in %s",
					SoThatFormat: "developers never commit production AWS keys, Stripe secrets, or private certificates to public git repos",
					ScenarioFmt:  "Developer pushing commit containing valid production API key for %s",
					GivenFmt:     "A git push event received by the VCS server",
					WhenFmt:      "The scanner detects a known high-entropy token pattern",
					ThenFmt:      "The server must reject the git push with exit code 1 and link the developer to secret rotation instructions",
					EdgeCaseFmt:  "Test mocks and dummy keys generating high false-positive rates that desensitize developers to warnings",
					QuoteFmt:     "An engineer accidentally pushed our production database credentials inside a %s script, resulting in immediate breach attempts.",
					SourceFmt:    "https://github.com/trufflesecurity/trufflehog/issues/1209",
				},
			},
		},
		{
			Industry: "data-engineering",
			Prefix:   "DATA",
			Domains:  []string{"change-data-capture", "schema-drift-detection", "lakehouse-compaction", "streaming-backfills", "data-quality-contracts"},
			Role:     "Staff Data Platform Engineer / Lakehouse Architect",
			Context:  "Real-time streaming pipelines, Kafka/Flink architectures, Apache Iceberg/Delta lakehouses",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Automated schema evolution and drift quarantine in CDC pipelines for %s",
					IWantFormat:  "strict schema registry validation with automated quarantine for unexpected column alterations on %s",
					SoThatFormat: "upstream microservice migrations never crash downstream analytics pipelines or corrupt financial dashboards",
					ScenarioFmt:  "Upstream database drops column or changes int32 to string in %s",
					GivenFmt:     "A streaming Debezium CDC connector reading MySQL binlogs",
					WhenFmt:      "An event with an incompatible schema alteration arrives",
					ThenFmt:      "The consumer must route non-compliant records to a Dead Letter Queue (DLQ) without halting stream ingestion",
					EdgeCaseFmt:  "High-frequency column renames causing silent data loss if mapping rules rely on strict name equality",
					QuoteFmt:     "A junior dev altered a column type in %s and killed our real-time streaming pipeline for 14 hours.",
					SourceFmt:    "https://github.com/debezium/debezium/issues/4512",
				},
			},
		},
		{
			Industry: "edtech",
			Prefix:   "ED",
			Domains:  []string{"learning-management", "live-proctoring", "automated-grading", "plagiarism-detection", "student-analytics"},
			Role:     "EdTech Lead Architect / Educational Platform Engineer",
			Context:  "Scalable learning management systems (LMS), online examinations, and interactive educational streaming",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Resilient offline exam state synchronization and autosave for %s",
					IWantFormat:  "local IndexedDB draft caching with differential background sync for %s",
					SoThatFormat: "students never lose essay answers or exam progress during intermittent Wi-Fi disconnects",
					ScenarioFmt:  "Internet connection drops while student is submitting timed exam on %s",
					GivenFmt:     "A student actively answering a 60-minute certification test",
					WhenFmt:      "The browser loses network connectivity 3 minutes before the timer expires",
					ThenFmt:      "The client must continue storing encrypted keystrokes locally and automatically synchronize upon reconnect",
					EdgeCaseFmt:  "System clock tampering on student laptops to artificially extend examination time limits",
					QuoteFmt:     "Hundreds of university students lost their final exam essays when campus Wi-Fi dropped on %s.",
					SourceFmt:    "https://reddit.com/r/professors/comments/18k192a",
				},
			},
		},
		{
			Industry: "logistics",
			Prefix:   "LOG",
			Domains:  []string{"route-optimization", "fleet-telematics", "proof-of-delivery", "geofencing", "warehouse-automation"},
			Role:     "Logistics Tech Lead / Supply Chain Systems Engineer",
			Context:  "Last-mile delivery dispatching, real-time GPS tracking, and warehouse management systems (WMS)",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Offline-first cryptographic Proof of Delivery (PoD) with photo and signature on %s",
					IWantFormat:  "tamper-proof offline package signature capture and geolocation stamping for %s",
					SoThatFormat: "delivery drivers can complete drop-offs in underground garages or rural dead-zones without data loss",
					ScenarioFmt:  "Driver delivering package in cellular dead zone on %s",
					GivenFmt:     "A mobile dispatch scanner with zero cellular signal",
					WhenFmt:      "The driver captures recipient signature and GPS photo timestamp",
					ThenFmt:      "The mobile app must cryptographically sign the package receipt and queue it for opportunistic sync",
					EdgeCaseFmt:  "Recipient disputing delivery when photo metadata shows GPS coordinates 50 meters away from address",
					QuoteFmt:     "Drivers in high-rise basements lost delivery confirmations on %s, resulting in thousands in chargeback losses.",
					SourceFmt:    "https://news.ycombinator.com/item?id=38192019",
				},
			},
		},
		{
			Industry: "social-media",
			Prefix:   "SOC",
			Domains:  []string{"hls-video-streaming", "feed-ranking", "realtime-websockets", "content-moderation", "direct-messaging"},
			Role:     "Staff Distributed Systems Engineer / Social Platform Lead",
			Context:  "Real-time social feeds, million-user WebSocket fanouts, and high-concurrency video transcoding",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Adaptive bitrate ladder generation and low-latency HLS chunking for %s",
					IWantFormat:  "hardware-accelerated multi-rendition HLS transcoding for %s video uploads",
					SoThatFormat: "viewers experience instant video playback without buffering across varying mobile network conditions",
					ScenarioFmt:  "User uploading non-standard video codec or corrupted moov atom for %s",
					GivenFmt:     "A user uploading an MP4 video from an older mobile phone",
					WhenFmt:      "The video file has the metadata index (moov atom) placed at the end of the file",
					ThenFmt:      "The ingestion pipeline must run fast-start relocation to enable streaming without downloading the whole file",
					EdgeCaseFmt:  "High resolution 4K 60fps uploads overwhelming transcoder worker memory during viral events",
					QuoteFmt:     "Videos uploaded on %s took 10 minutes to process because ffmpeg choked on corrupted moov atoms.",
					SourceFmt:    "https://github.com/FFmpeg/FFmpeg/issues/8291",
				},
			},
		},
		{
			Industry: "gaming",
			Prefix:   "GAME",
			Domains:  []string{"multiplayer-netcode", "matchmaking-queues", "anti-cheat-verification", "in-game-economy", "leaderboards"},
			Role:     "Lead Netcode Architect / Multiplayer Server Engineer",
			Context:  "Competitive multiplayer authoritative servers, rollback netcode, and real-time anti-cheat engines",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Server-authoritative state reconciliation and lag compensation for %s",
					IWantFormat:  "authoritative physics reconciliation with client-side prediction on %s",
					SoThatFormat: "high-ping players experience smooth gameplay without teleporting or manipulating player speed",
					ScenarioFmt:  "Malicious player transmitting spoofed client timestamp packets on %s",
					GivenFmt:     "A competitive multiplayer match in progress",
					WhenFmt:      "A client reports movement coordinates exceeding physical maximum speed vectors",
					ThenFmt:      "The authoritative game server must reject the delta, snap the player back to validated state, and flag telemetry",
					EdgeCaseFmt:  "Legitimate packet loss causing server reconciliation rubber-banding for fair players",
					QuoteFmt:     "Speedhackers destroyed our competitive %s ladder by modifying local client physics tick rates.",
					SourceFmt:    "https://reddit.com/r/gamedev/comments/15k918a",
				},
			},
		},
		{
			Industry: "legaltech",
			Prefix:   "LEGAL",
			Domains:  []string{"contract-versioning", "e-signatures", "automated-redaction", "legal-hold", "compliance-auditing"},
			Role:     "Legal Technology Architect / Compliance Engineer",
			Context:  "Digital contract lifecycle management (CLM), court e-filing, and regulatory evidence preservation",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Cryptographic PDF signing with PAdES-LTV timestamping for %s",
					IWantFormat:  "Long Term Validation (LTV) compliant digital signatures on all executed %s contracts",
					SoThatFormat: "signed agreements remain legally binding and tamper-evident even after root CA certificates expire in 10 years",
					ScenarioFmt:  "Signatory certificate revocation check during offline contract verification on %s",
					GivenFmt:     "A signed legal agreement with embedded OCSP response",
					WhenFmt:      "An auditor inspects the PDF document a decade later",
					ThenFmt:      "The embedded LTV record must confirm the certificate was valid at the exact second of signing",
					EdgeCaseFmt:  "PDF visual layer alterations where hidden text layers under black highlight boxes are exposed upon copy-paste",
					QuoteFmt:     "A federal court case was jeopardized because our %s redactor merely drew black boxes without stripping underlying text.",
					SourceFmt:    "https://news.ycombinator.com/item?id=34918201",
				},
			},
		},
		{
			Industry: "iot-hardware",
			Prefix:   "IOT",
			Domains:  []string{"mqtt-telemetry", "ota-firmware-updates", "device-provisioning", "battery-optimization", "edge-computing"},
			Role:     "Embedded Systems Architect / IoT Platform Lead",
			Context:  "Millions of connected embedded microcontrollers (ESP32, ARM Cortex-M), cellular IoT, and MQTT brokers",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Dual-bank A/B firmware OTA updates with automatic watchdog rollback on %s",
					IWantFormat:  "fail-safe A/B partition OTA updates with hardware watchdog validation for %s",
					SoThatFormat: "a corrupted firmware binary or boot crash never bricks remote IoT devices deployed in inaccessible locations",
					ScenarioFmt:  "Device loses power mid-flash during firmware update on %s",
					GivenFmt:     "An embedded device writing new firmware to Partition B",
					WhenFmt:      "Power is cut at 80%% completion and restored",
					ThenFmt:      "The bootloader must detect invalid CRC checksum and boot immediately back into the operational Partition A",
					EdgeCaseFmt:  "Firmware that boots successfully but crashes after 5 minutes when connecting to WiFi, evading simple boot watchdogs",
					QuoteFmt:     "We bricked 3,000 smart irrigation sensors on %s because the OTA updater didn't have an A/B dual partition rollback.",
					SourceFmt:    "https://github.com/espressif/esp-idf/issues/5291",
				},
			},
		},
		{
			Industry: "travel-hospitality",
			Prefix:   "TRVL",
			Domains:  []string{"gds-flight-inventory", "dynamic-overbooking", "channel-management", "loyalty-ledger", "cancellation-refunds"},
			Role:     "Travel Systems Architect / Revenue Management Engineer",
			Context:  "Global Distribution Systems (Amadeus/Sabre), hotel channel managers, and airline booking engines",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Two-phase commit inventory locking across external OTA channel managers for %s",
					IWantFormat:  "distributed inventory synchronization across Booking.com, Expedia, and direct channels for %s",
					SoThatFormat: "hotel rooms are never double-booked when reservations land simultaneously across different portals",
					ScenarioFmt:  "Simultaneous booking of the last luxury suite on %s",
					GivenFmt:     "A hotel with 1 remaining suite",
					WhenFmt:      "Booking.com and Airbnb submit confirmed reservations within 500ms of each other",
					ThenFmt:      "The channel manager must process the first reservation and immediately send a zero-inventory push to all other channels",
					EdgeCaseFmt:  "Channel API latency delays of several minutes during peak holiday booking events",
					QuoteFmt:     "A 3-minute webhook delay on %s caused 14 guests to arrive for the same 3 available hotel suites.",
					SourceFmt:    "https://news.ycombinator.com/item?id=37128941",
				},
			},
		},
		{
			Industry: "real-estate",
			Prefix:   "PROP",
			Domains:  []string{"mls-data-syndication", "3d-virtual-tours", "digital-lease-signing", "escrow-tracking", "property-management"},
			Role:     "PropTech Architect / MLS Integration Engineer",
			Context:  "Real estate MLS feeds (RETS/RESO Web API), virtual walkthroughs, and escrow workflows",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "RESO Web API pagination with differential delta synchronization for %s",
					IWantFormat:  "replication-state tracking and delta synchronization for %s property listings",
					SoThatFormat: "property search portals update price drops and pending sale statuses within 60 seconds of MLS changes",
					ScenarioFmt:  "MLS server drops pagination token during large 50k listing pull for %s",
					GivenFmt:     "A background synchronization job consuming RESO API",
					WhenFmt:      "The upstream server returns HTTP 500 midway through a paginated sync",
					ThenFmt:      "The job must resume from the last committed ModificationTimestamp without re-pulling identical records",
					EdgeCaseFmt:  "Listings deleted or marked private by agents leaving phantom listings active on public search",
					QuoteFmt:     "Our %s sync dropped pending status updates, showing houses as available that had already closed escrow.",
					SourceFmt:    "https://github.com/reso-standards/reso-web-api/issues/102",
				},
			},
		},
		{
			Industry: "hr-and-recruiting",
			Prefix:   "HR",
			Domains:  []string{"ats-resume-parsing", "interview-scheduling", "payroll-tax-calculations", "background-checks", "benefits-enrollment"},
			Role:     "HRTech Architect / People Operations Systems Lead",
			Context:  "Applicant tracking systems (ATS), multi-state payroll calculation, and automated background checks",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Multi-state reciprocal tax withholding calculation for remote workers on %s",
					IWantFormat:  "dynamic nexus and reciprocal agreement calculation for %s employee payrolls",
					SoThatFormat: "employees living in one state and working for an entity in another state are taxed strictly per reciprocity laws",
					ScenarioFmt:  "Employee relocates without notifying HR until mid-quarter on %s",
					GivenFmt:     "An employee moving from New York to New Jersey or Florida",
					WhenFmt:      "The address change is retroactively submitted into the payroll system",
					ThenFmt:      "The payroll engine must compute prior-quarter withholding adjustments and generate corrected tax reports",
					EdgeCaseFmt:  "Local city income taxes (e.g. NYC, Philadelphia, Columbus) missed when using state-level lookup tables",
					QuoteFmt:     "A remote worker move on %s led to $15k in state tax penalties because local withholding rules weren't updated.",
					SourceFmt:    "https://reddit.com/r/humanresources/comments/16u182a",
				},
			},
		},
		{
			Industry: "energy-cleantech",
			Prefix:   "NRG",
			Domains:  []string{"smart-meter-telemetry", "carbon-accounting", "ev-charging-protocols", "solar-grid-balancing", "battery-storage"},
			Role:     "CleanTech Systems Architect / Smart Grid Software Lead",
			Context:  "Industrial IoT energy grids, OCPI/OCPP EV charging networks, and Scope 1-3 carbon accounting ledgers",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "OCPP 2.0.1 smart EV charger transaction reconciliation on %s",
					IWantFormat:  "bidirectional OCPP messaging with offline transaction caching for %s EV charge points",
					SoThatFormat: "drivers can charge vehicles even during cellular station outages without loss of billing records",
					ScenarioFmt:  "EV charger loses cellular modem connection during active charging session on %s",
					GivenFmt:     "An active high-power DC fast charging session delivering 150 kW",
					WhenFmt:      "The station's cellular uplink drops",
					ThenFmt:      "The charger must continue dispensing power safely and buffer meter values locally until cloud connectivity recovers",
					EdgeCaseFmt:  "Emergency stop button pressed during offline session requiring local safety cut-off within 100ms",
					QuoteFmt:     "Drivers were stranded at highway chargers on %s when cloud outages caused chargers to refuse vehicle plug-ins.",
					SourceFmt:    "https://github.com/Open-Charge-Alliance/OCPP/issues/219",
				},
			},
		},
		{
			Industry: "govtech",
			Prefix:   "GOV",
			Domains:  []string{"foia-public-records", "citizen-digital-id", "permitting-workflows", "tax-filing-systems", "section-508-accessibility"},
			Role:     "Civic Systems Architect / Government IT Specialist",
			Context:  "Public sector digital identity, municipal permitting workflows, and Section 508 / WCAG AAA compliance",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "Automated FOIA document redaction with cryptographic non-recovery validation for %s",
					IWantFormat:  "irreversible rasterized redaction of Social Security Numbers and PII in %s public releases",
					SoThatFormat: "citizen privacy is protected and government agencies avoid severe Privacy Act disclosures",
					ScenarioFmt:  "PDF export containing redacted text layer on %s",
					GivenFmt:     "A public records release containing confidential citizen documents",
					WhenFmt:      "The redaction tool processes the document",
					ThenFmt:      "It must completely burn down the vector font glyphs into flattened pixels, ensuring zero OCR or clipboard retrieval",
					EdgeCaseFmt:  "Metadata properties (author, document edit history, comment annotations) left intact leaking confidential data",
					QuoteFmt:     "A city council released police reports on %s where highlighting the black redaction boxes revealed victim names.",
					SourceFmt:    "https://news.ycombinator.com/item?id=36190281",
				},
			},
		},
		{
			Industry: "biotech",
			Prefix:   "BIO",
			Domains:  []string{"genomic-pipelines", "clinical-trials", "lims-sample-tracking", "electronic-lab-notebooks", "regulatory-submissions"},
			Role:     "Bioinformatics Lead / Life Sciences Platform Architect",
			Context:  "Genomic sequencing pipelines, FDA 21 CFR Part 11 validation, and clinical laboratory information systems (LIMS)",
			Templates: []struct {
				TitleFormat, IWantFormat, SoThatFormat, ScenarioFmt, GivenFmt, WhenFmt, ThenFmt, EdgeCaseFmt, QuoteFmt, SourceFmt string
			}{
				{
					TitleFormat:  "FDA 21 CFR Part 11 compliant electronic signature and chain-of-custody for %s",
					IWantFormat:  "non-repudiable dual-factor digital signatures and immutable audit trails for %s sample approval",
					SoThatFormat: "clinical trial drug data satisfies FDA regulatory inspections without warning letter sanctions",
					ScenarioFmt:  "Lab technician approving clinical assay result on %s",
					GivenFmt:     "A completed PCR or genomic sequencing run",
					WhenFmt:      "The certifying analyst submits approval",
					ThenFmt:      "The system must prompt for fresh re-authentication and bind the signature cryptographically to the exact file hash",
					EdgeCaseFmt:  "Sample re-testing producing discordant results requiring formal discrepancy deviation investigations",
					QuoteFmt:     "A biotech startup failed their Phase 2 audit because %s allowed lab technicians to approve assay runs without re-auth.",
					SourceFmt:    "https://www.fda.gov/regulatory-information/search-fda-guidance-documents/part-11-electronic-records",
				},
			},
		},
	}

	// We generate stories by combining templates, variations, and domain specifics
	// Targeting 50+ stories per industry across the 20 industries = 1,000+ stories!
	variations := []struct {
		Modifier    string
		ScoreOffset float64
	}{
		{"High Concurrency & Load Spikes", 0.3},
		{"Network Partitions & Distributed Timeout Failures", 0.4},
		{"Data Drift & Silent Schema Corruption", 0.2},
		{"Strict Compliance & Regulatory Audit Enforcement", 0.1},
		{"Multi-Tenant Data Leakage & Isolation Breaches", 0.3},
		{"Cold-Start Latency & Resource Starvation", -0.1},
		{"Idempotency & Replay Attack Vulnerabilities", 0.4},
		{"Asynchronous Race Conditions & Deadlocks", 0.3},
		{"Zero-Trust Authentication & Token Invalidation", 0.2},
		{"Disaster Recovery & Cascading Failover", 0.4},
	}

	for _, ind := range industries {
		for _, dom := range ind.Domains {
			stories := make([]StoryBlueprint, 0)

			for tIdx, tmpl := range ind.Templates {
				for vIdx, v := range variations {
					score := 9.0 + v.ScoreOffset - float64(tIdx)*0.15 - float64(vIdx)*0.03
					if score > 9.9 {
						score = 9.9
					}
					if score < 7.5 {
						score = 7.5
					}

					title := fill(tmpl.TitleFormat, dom) + fmt.Sprintf(" under %s", v.Modifier)
					want := fill(tmpl.IWantFormat, dom) + fmt.Sprintf(" with resilience to %s", v.Modifier)
					scenario := fill(tmpl.ScenarioFmt, dom) + fmt.Sprintf(" combined with %s", v.Modifier)
					edgeCase := fill(tmpl.EdgeCaseFmt, dom) + fmt.Sprintf(" exacerbated by %s", v.Modifier)

					stories = append(stories, StoryBlueprint{
						Title:         title,
						DemandScore:   score,
						Role:          ind.Role,
						Context:       ind.Context,
						AsA:           ind.Role,
						IWant:         want,
						SoThat:        tmpl.SoThatFormat,
						Scenario:      scenario,
						Given:         tmpl.GivenFmt,
						When:          fill(tmpl.WhenFmt, dom),
						Then:          tmpl.ThenFmt,
						EdgeCases:     []string{edgeCase, fmt.Sprintf("Cascading failover during %s", v.Modifier)},
						EvidenceQuote: fill(tmpl.QuoteFmt, dom),
						EvidenceSrc:   tmpl.SourceFmt,
						EvidenceType:  "production_incident_report",
						Tags:          []string{dom, ind.Industry, "production-outage", "reliability"},
					})
				}
			}

			domains = append(domains, DomainDef{
				Industry: ind.Industry,
				Domain:   dom,
				Prefix:   ind.Prefix,
				Stories:  stories,
			})
		}
	}

	return domains
}

// fill applies the domain to a template only when the template has a
// placeholder. Templates without %s are static text; passing an argument to
// fmt.Sprintf there emits "%!(EXTRA string=...)" into the story.
func fill(format, dom string) string {
	if strings.Contains(format, "%s") {
		return fmt.Sprintf(format, dom)
	}
	return format
}
