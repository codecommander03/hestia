# ADR 0002: Backend in Go

- **Status:** Accepted
- **Date:** 2026-10-12
- **Supersedes:** [ADR 0001](0001-backend-framework.md)
- **Requirements:** NFR-MNT-1, NFR-MNT-2, NFR-CON-1 to 7, NFR-SCALE-3, FR-ADM-1, requirements §7

## Context

ADR 0001 chose FastAPI, because requirements §7 fixed the backend language as
Python. §7 now says Go. This ADR records why the language changed and which Go
libraries replace the FastAPI stack.

The backend was about 50 lines when this was decided: settings, a database
engine, a health endpoint and two tests. Switching now costs almost nothing.
Switching later would mean rewriting real features.

**Why Go:**

- **Learning.** This is a solo project, and learning Go is one of its goals.
  That's the main reason for the switch, and it's stated plainly here.
- **The SQL stays visible.** sqlc generates typed Go functions from
  hand-written `.sql` files. Each concurrency guarantee in `data-model.md`
  is one query in one file, which fits ADR 0001's goal of keeping the SQL
  that has to be defended in plain view, and fits it better.
- **No async discipline.** Each request runs in its own goroutine, so one
  blocking call doesn't stall other requests. ADR 0001 listed this as a
  cost, and in Go it goes away.
- **Small footprint.** The backend ships as one static binary that uses
  little memory, which suits one small free-tier VM (₹0/month).

**Not a reason:** "Python doesn't scale". For this workload it isn't true.
Most requests wait on Postgres, object storage or push notifications, so
throughput depends on the database, not on the language. The backend is
stateless (NFR-SCALE-3), so it scales by running more copies in either
language.

ADR 0001's list of what the backend must support still applies unchanged:
a typed client generated from an accurate OpenAPI spec, exact control over
SQL, versioned migrations, stateless processes, integration tests against a
real Postgres, and one developer at ₹0/month.

## Decision

**Go**, with:

- **Huma** (v2) on the standard library's `net/http` router. Request and
  response types are Go structs, and Huma generates the OpenAPI 3.1 spec from
  them and validates requests against it. This replaces FastAPI + Pydantic.
- **sqlc** with **pgx** (v5) for database access. Queries are written in SQL,
  and sqlc generates typed Go functions for them. This replaces SQLAlchemy.
- **goose** for migrations, written as plain SQL files. This replaces Alembic.
- **PostGIS through plain SQL.** `ST_DWithin`, `ST_MakePoint` and the other
  functions are called directly in the sqlc queries. Points go into queries
  as longitude/latitude parameters and come back out as `ST_X`/`ST_Y`
  columns, so Go code never handles geometry types.

The invariant-critical queries (conditional `UPDATE`s, `SELECT … FOR UPDATE`,
`FOR UPDATE SKIP LOCKED`) each get their own named sqlc query, one per
guarantee, and each has an integration test against Postgres.

## Alternatives considered

### Stay on FastAPI (ADR 0001)
Still a sound choice technically, and the strongest option for generating
OpenAPI from code. Rejected only because of the learning goal above.

### Go with oapi-codegen (spec first)
The OpenAPI spec is written by hand, and the Go server interfaces are
generated from it. Rejected because the spec then has to be written and
reviewed by hand. NFR-MNT-1 wants the spec to come from the code, and Huma
does that.

### Go with GORM
Rejected for the same reason ADR 0001 rejected the Django ORM: it hides the
exact SQL behind the concurrency guarantees, and that SQL is what has to be
defended.

### Go with Gin or Echo (plus swaggo for OpenAPI)
swaggo builds the spec from code comments, which can drift away from the
real types. Since Go 1.22, the standard library's router handles methods
and path parameters, so a third-party router adds little.

## Consequences

**Good:**
- Every concurrency guarantee is a plain `.sql` file with a generated typed
  function and its own test.
- The OpenAPI spec still comes from the same types that validate requests,
  so the Expo app's client can still be generated in CI.
- No event-loop pitfalls. File checks and malware scans (FR-VER-2) can still
  go to the background worker, but the reason is time spent, not blocking.
- One binary to deploy.

**Costs, accepted:**
- **Learning curve.** Go is new to the developer, so early progress is slower.
- **More boilerplate.** Error handling is explicit, and there's no ORM for
  everyday reads and writes. sqlc generates most of that code.
- **PostGIS is manual.** No GeoAlchemy2 equivalent. Geometry is converted
  in SQL as described above.
- **Huma has a smaller community than FastAPI.** The OpenAPI spec is the
  one part this choice depends on most.
- Still true from ADR 0001: no admin site, auth or permissions out of the box.
  OTP login, sessions and refresh-token rotation (NFR-SEC-2/4) are built by
  hand.

**Left for later ADRs** (unchanged from ADR 0001):
- the background job runner (expiry, timeouts, notifications, document checks);
- the admin web app's frontend framework;
- the hosting and deployment setup.
