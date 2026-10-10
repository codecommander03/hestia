# ADR 0001: Backend framework

- **Status:** Accepted
- **Date:** 2026-10-10
- **Requirements:** NFR-MNT-1, NFR-MNT-2, NFR-CON-1 to 7, NFR-SCALE-3, FR-ADM-1, requirements §7

## Context

Requirements §7 fixes the backend language as Python and leaves the framework
open. The backend is a **JSON API only**: the Android app (React Native + Expo)
and the separate admin web app are its only clients. Nothing is rendered on the
server.

What the framework has to support well:

1. **A typed client for the app.** NFR-MNT-1: the OpenAPI spec is generated
   from the backend code, and the app's typed client is generated from that
   spec. The spec has to be accurate without being kept in sync by hand.
2. **Exact control over SQL.** The core guarantees in `data-model.md` are not
   app-level checks. They are specific queries:
   - conditional `UPDATE … WHERE quantity_claimed + $q <= quantity_total`
     (NFR-CON-1);
   - `UPDATE … WHERE status = 'OPEN'` to take a task (NFR-CON-2);
   - `SELECT … FOR UPDATE` row locks (NFR-CON-3, NFR-CON-6);
   - `FOR UPDATE SKIP LOCKED` for the review queue (FR-VER-4);
   - PostGIS `ST_DWithin` for every radius query.
3. **Versioned migrations** (NFR-MNT-2).
4. **Stateless processes** (NFR-SCALE-3), so more copies can be run later.
5. **Integration tests against a real Postgres** (NFR-MNT-3).
6. **One developer, ₹0/month.** Less framework to learn and run is a real
   benefit.

## Decision

**FastAPI**, with:

- **Pydantic** for request and response models. These are what the OpenAPI
  spec is generated from.
- **SQLAlchemy 2.x** (Core and ORM, async with `asyncpg`) for database access.
- **GeoAlchemy2** for PostGIS types and functions.
- **Alembic** for migrations.

The invariant-critical queries are written explicitly in SQLAlchemy Core, one
function each, so each guarantee lives in one readable place that a test can
target. The ORM is used for ordinary reads and writes.

## Alternatives considered

### Django + Django REST Framework
The most complete option: ORM, migrations, auth and a built-in admin site.

Rejected because:
- **The admin site doesn't fit FR-ADM.** The admin panel needs a review queue
  with time-limited locks, TOTP, and an audit entry for every document *view*
  (FR-AUD-3). That means heavy customisation of Django admin, and it would run
  inside the same process as the public API. FR-ADM-1 asks for a separate app.
- **The ORM hides the queries that matter.** Conditional updates are possible
  with `F()` expressions and `select_for_update`, but the actual SQL is less
  visible, and it's the actual SQL that has to be defended.
- **OpenAPI is an add-on** (drf-spectacular), and serializers and models are
  defined separately, so the spec can drift from the code.
- **Most of what it includes is unused:** templates, forms, server-side
  sessions.

### Flask
Rejected: no built-in request validation or OpenAPI. Getting there means
choosing and combining several extensions, which is FastAPI with extra steps.

### Litestar
Technically close to FastAPI and arguably cleaner in places. Rejected because it
has a much smaller community, fewer examples and fewer answers to search for.
That cost lands entirely on a solo developer.

## Consequences

**Good:**
- The OpenAPI spec comes from the same Pydantic models that validate requests,
  so it can't drift. The Expo app's client is generated from it in CI.
- Every concurrency guarantee is a visible SQL statement with its own test.
- Async I/O suits this workload: most requests wait on Postgres, object storage
  or push notifications, not on the CPU.

**Costs, accepted:**
- **No admin site, auth or permissions out of the box.** Login by OTP, sessions
  and refresh-token rotation (NFR-SEC-2/4) are built by hand. They're
  specified closely enough in the requirements that this is a bounded amount
  of work, and it's work worth being able to explain.
- **More decisions are left open:** project layout, dependency injection,
  transaction handling. They get written down as conventions in Stage 0.
- **Async needs discipline.** A blocking call (a sync library, CPU-heavy work)
  inside a request stalls every request on that worker. File checks and
  malware scans (FR-VER-2) therefore run in the background worker, not in the
  request.

**Left for later ADRs:**
- the background job runner (expiry, timeouts, notifications, document checks);
- the admin web app's frontend framework;
- the hosting and deployment setup.
