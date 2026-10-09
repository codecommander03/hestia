# Hestia: Requirements

> **Name.** Hestia is the Greek goddess of the hearth and home. She received the first portion of every meal.
> **Status:** v0.1 draft, 2026-09-30. Edit and append freely. Keep the IDs stable, and when something is removed, mark it ~~struck~~ instead of renumbering, so references don't break.

## How to read this

- **FR** = functional requirement (what the system does). **NFR** = non-functional requirement (how well it does it; always measurable).
- IDs look like `FR-LST-3`. Design docs, API specs and tests should point back to these IDs.
- Priority: **[M]** Must (the MVP), **[S]** Should (next), **[C]** Could (nice to have). Won't-have items are listed in §8.
- Values in `{braces}` are tunables. Their defaults are listed in §9.
- The state lists in §5 are a first idea. The exact state machines go in the design doc.

---

## 1. Overview

**Problem.** Corporate cafeterias, college canteens, hostels, restaurants, event venues and households regularly have safe, surplus cooked or packaged food. NGOs and shelters that feed people in need could use it, but there's no quick, trusted way to connect the two before the food goes bad.

**Solution.** A mobile-first platform for this:
- Anyone can post surplus food.
- Verified NGOs claim it.
- Anyone eligible can carry it from the donor to the NGO.
- Each handoff is confirmed with an OTP and recorded in a tamper-evident log.

**Context.** A portfolio project for India. It's a demo, not a live service: all data and documents are fictional, and it runs at ₹0/month.

**Success looks like:**
1. The full flow works end to end on a real Android phone: post → claim → carry → handoff → accept.
2. The design holds up under a simulated city-scale load.
3. Every decision is written down and can be defended in an interview.

---

## 2. Actors and roles

| Actor | Who | Can |
|---|---|---|
| **User** (individual) | Anyone who signs up | Donate food personally, carry food (if eligible), create an org |
| **Org admin** (`ORG_ADMIN`) | Manages a donor org or NGO | Everything org staff can do, plus manage members, upload verification documents and set NGO policies |
| **Org staff** (`ORG_STAFF`) | Works at a donor org or NGO | Donor org: post listings on the org's behalf. NGO: claim listings and receive food |
| **Carrier** | Any user who is eligible for a particular transport task | Accept transport tasks and take part in handoffs. *This is a capability, not a role* (see §3) |
| **Admin** (`ADMIN`) | Platform operator | Verify orgs, moderate content, suspend users and orgs, view audit logs |
| **System** | Scheduled jobs | Expire listings, handle timeouts, widen the notification radius, send reminders |

Notes:
- There are two org types: `DONOR_ORG` (corporate, college, hostel, restaurant, caterer, event venue) and `NGO` (receiver). Each org has exactly one type.
- Only members of a **verified NGO** can receive food.
- A user can belong to several orgs. When posting or claiming, they choose whether they act **personally** or **on behalf of an org**.

---

## 3. Trust tiers

What a donor or carrier may do depends on their trust tier, not their role.

| Tier | Who | Donating | Carrying |
|---|---|---|---|
| **0: New** | Verified email and phone, little history | Up to `{T0_MAX_SERVINGS}` servings per listing and `{T0_MAX_LISTINGS_PER_DAY}` listings per day. Shown as "Unverified individual" | Allowed if the NGO's policy lets new carriers in |
| **1: Established** | Meets `{T1_CRITERIA}` | Higher limits (`{T1_*}`) | Counts as experienced under NGO policies |
| **2: Verified org** | Member of a verified donor org, acting for that org | Org limits (`{T2_MAX_SERVINGS}`) | Uses the person's own individual tier |

- A user's tier can go down after upheld reports or a high rate of rejected food.
- While a donor org's verification is pending, its members post at Tier 0 limits.

---

## 4. The main flow (example)

1. A college canteen (a verified donor org) has 60 servings of veg biryani left after lunch.
2. A staff member posts a listing: 60 servings, veg, cooked at 13:00, kept at room temperature. It's safe until 17:00, pickup is between 14:30 and 16:00, and a photo is attached.
3. NGOs within `{NOTIFY_RADIUS_KM}` get a notification. NGO A claims 40 servings and NGO B claims 20. A third NGO's claim fails cleanly with "nothing left".
4. NGO A opens its claim to carriers. Carriers whose chosen area covers the pickup point get a notification.
5. Carrier C accepts the task. Only now can C see the exact pickup location and the canteen's contact number.
6. At the canteen, C's app shows a **pickup code**. The canteen staff enter it, and the food is marked **picked up**.
7. At NGO A, the NGO's app shows a **drop code**. C enters it, and the food is marked **delivered**.
8. NGO A inspects the food and accepts all 40 servings.
9. Every step is recorded in the custody log. The canteen sees "40 servings delivered to NGO A at 15:42".

---

## 5. Functional requirements

### 5.1 Accounts (ACC)
- **FR-ACC-1 [M]:** Users sign up and log in with a one-time code sent by email. There are no passwords.
- **FR-ACC-2 [M]:** Users verify their phone number with an OTP before they can donate or carry food. In the demo, a dev provider shows or logs the code instead of sending SMS (see §7).
- **FR-ACC-3 [M]:** At sign-up, users accept the terms and privacy notice. The system records which version they accepted and when.
- **FR-ACC-4 [M]:** Users can view and edit their profile: name, phone and an optional photo.
- **FR-ACC-5 [S]:** Users can log out of all their devices at once.
- **FR-ACC-6 [S]:** Users can delete their account. Their personal data is then deleted or anonymised as described in NFR-PRIV-4.
- **FR-ACC-7 [C]:** Users can sign in with Google.

### 5.2 Organisations (ORG)
- **FR-ORG-1 [M]:** Any user can create an org of type `DONOR_ORG` or `NGO`. The creator becomes its `ORG_ADMIN`.
- **FR-ORG-2 [M]:** An org admin can invite members by email, switch a member between `ORG_ADMIN` and `ORG_STAFF`, and remove members.
- **FR-ORG-3 [M]:** An org must always have at least one `ORG_ADMIN`. The last admin can't leave or be demoted.
- **FR-ORG-4 [M]:** An NGO sets its **drop location** (address and map pin), its **service radius** and the **diet types it accepts** (veg, egg, non-veg).
- **FR-ORG-5 [S]:** An org can have several **sites**, for example several campus cafeterias. Listings and drops refer to a site.
- **FR-ORG-6 [S]:** An NGO chooses a **carrier policy**:
  - anyone eligible (default);
  - only carriers with at least `{MIN_DELIVERIES}` successful deliveries;
  - only the NGO's own volunteers.
- **FR-ORG-7 [S]:** An NGO admin can mark users as the NGO's volunteers.
- **FR-ORG-8 [C]:** An NGO sets a daily capacity (maximum servings per day). The app blocks claims that would go over it.

### 5.3 Verification (VER)
- **FR-VER-1 [M]:** An org admin uploads the documents required for the org's type and enters each document's ID number and expiry date. The list of required documents is configurable. Defaults:
  - **NGO:** registration certificate (Trust, Society or Section 8 company), NGO Darpan ID, the organisation's PAN, and FSSAI registration. FSSAI's surplus-food regulations expect food recovery organisations to be registered; confirm the current rules.
  - **Donor org:** the FSSAI licence or registration of its kitchen or caterer, GSTIN or CIN, and an authorisation letter for the person uploading.
- **FR-VER-2 [M]:** The system automatically checks every upload:
  - the real file type, detected from the file's contents rather than its name;
  - the file size limit;
  - a malware scan;
  - the ID number's format: GSTIN (including its check digit), PAN pattern, 14-digit FSSAI number.

  If a check fails, the document is rejected automatically and the reason is shown.
- **FR-VER-3 [M]:** The system flags an ID number that another org has already used.
- **FR-VER-4 [M]:** Documents that pass the automated checks join an **admin review queue**, oldest first.
  - An admin picks an item, which locks it to them for `{REVIEW_LOCK_MIN}`.
  - Two admins can never review the same item at the same time.
- **FR-VER-5 [M]:** An admin approves or rejects each document, with a reason code and an optional note. The org becomes `VERIFIED` once every required document is approved.
- **FR-VER-6 [M]:** The org admin can see each document's status and rejection reason, and can resubmit rejected documents.
- **FR-VER-7 [M]:** The admin panel links to the public lookup portals (FoSCoS for FSSAI, NGO Darpan) so reviewers can check details by hand.
- **FR-VER-8 [S]:** Document expiry:
  - The org admin gets reminders `{DOC_EXPIRY_REMINDERS}` before a document expires.
  - When a document expires, the org becomes `EXPIRED`.
  - An expired NGO can't make new claims. An expired donor org falls back to Tier 0 limits.
- **FR-VER-9 [C]:** OCR reads the ID number and expiry date from the document and flags any mismatch with what the user typed.

Org statuses (first idea): `UNVERIFIED → PENDING → VERIFIED | REJECTED`, then `VERIFIED → EXPIRED | SUSPENDED`.

### 5.4 Listings (LST)
- **FR-LST-1 [M]:** A donor posts a listing with these fields:
  - title and food category (cooked meal, raw produce, packaged, bakery);
  - diet type (veg, egg, non-veg) and allergens;
  - quantity, in servings or in kg/units;
  - when it was prepared;
  - how it's stored (room temperature, kept hot, refrigerated);
  - how it's packed (sealed, packed, open);
  - safe-until time;
  - pickup time window;
  - pickup location (a site or a map pin);
  - at least one photo;
  - optional notes.
- **FR-LST-2 [M]:** The safe-until time can't be later than preparation time + `{MAX_SAFE_WINDOW[storage]}`. Sealed packaged food uses the expiry date printed on it.
- **FR-LST-3 [M]:** For every listing, the donor ticks a food-safety declaration: the food is fresh, stored as stated, and its allergens are listed.
- **FR-LST-4 [M]:** The donor's trust-tier limits are enforced: servings per listing and listings per day.
- **FR-LST-5 [M]:** Before anyone claims a listing, the donor can edit or cancel it. After a claim, the donor can only reduce the unclaimed quantity, or cancel with a reason. Cancelling notifies the NGOs that claimed it and counts against the donor's reliability.
- **FR-LST-6 [M]:** A listing expires automatically at safe-until or at the end of the pickup window, whichever comes first. After that, nobody can make new claims. Claims already picked up continue as normal.
- **FR-LST-7 [M]:** The donor can see each listing's current status and full history.
- **FR-LST-8 [S]:** The donor can offer to **deliver the food themselves** within `{SELF_DELIVERY_MAX_KM}`.
- **FR-LST-9 [S]:** The donor can record a **seal or tag number**, or a photo of the sealed package.
- **FR-LST-10 [C]:** The donor can repost a previous listing as a template.

Listing statuses (first idea): `ACTIVE ⇄ FULLY_CLAIMED` (it goes back to `ACTIVE` if a claim is cancelled), plus `EXPIRED` and `CANCELLED`. A listing that's only partly claimed stays `ACTIVE`.

### 5.5 Discovery and claims (CLM)
- **FR-CLM-1 [M]:** NGO members see active listings within their NGO's service radius, sorted by distance or by soonest expiry. Listings that don't match the NGO's diet settings are hidden.
- **FR-CLM-2 [M]:** Until a carrier is assigned, a listing shows only an **approximate location**, rounded to `{APPROX_LOCATION_M}` or shown as the area name.
- **FR-CLM-3 [M]:** Only members of a `VERIFIED` NGO that isn't suspended can claim.
- **FR-CLM-4 [M]:** An NGO claims all of a listing or part of it. The minimum claim is `{MIN_CLAIM}`, or whatever remains if that's less.
  - If several claims arrive at the same moment and together ask for more than is left, the extra ones fail.
  - A failed claim shows a clear message with how much is left (see NFR-CON-1).
- **FR-CLM-5 [M]:** When claiming, the NGO chooses how the food gets moved:
  - `NGO_PICKUP`: the NGO collects it itself;
  - `OPEN_TO_CARRIERS`: any eligible carrier can take it;
  - `DONOR_DELIVERY`: only if the donor offered it (FR-LST-8).
- **FR-CLM-6 [M]:** An NGO can cancel a claim before pickup, with a reason. If the listing is still valid, the claimed quantity goes back into it.
- **FR-CLM-7 [M]:** If the food isn't picked up by the end of the pickup window, the claim expires automatically.
- **FR-CLM-8 [S]:** A map view, plus filters for category, minimum quantity and donor tier.

Claim statuses (first idea): `ACTIVE → COMPLETED (with the acceptance result) | CANCELLED | EXPIRED`.

### 5.6 Transport (TRN)
- **FR-TRN-1 [M]:** Every claim creates exactly one **transport task**. Who carries the food depends on the transport mode:
  - `NGO_PICKUP`: an NGO member takes the task.
  - `OPEN_TO_CARRIERS`: any eligible carrier can take it (FR-TRN-3).
  - `DONOR_DELIVERY`: the donor carries it.
- **FR-TRN-2 [M]:** Carriers choose an **availability area** (a map pin and a radius). They're notified about open tasks whose pickup point falls inside it. The app doesn't track location in the background.
- **FR-TRN-3 [M]:** A carrier can accept an `OPEN_TO_CARRIERS` task only if **all** of these are true:
  - their phone number is verified;
  - they aren't suspended;
  - they meet the NGO's carrier policy;
  - they have fewer than `{CARRIER_MAX_ACTIVE_TASKS}` active tasks;
  - they didn't donate that listing.
- **FR-TRN-4 [M]:** Only one person can take a task. If several try at the same moment, exactly one succeeds (see NFR-CON-2).
- **FR-TRN-5 [M]:** Once a carrier accepts, they see:
  - the exact pickup and drop locations, as coordinates and addresses;
  - the straight-line distance;
  - a **Navigate** button that opens Google Maps with directions.
- **FR-TRN-6 [M]:** The donor's phone number is shown to the carrier, and the carrier's to the donor and the NGO, **only while the task is active**. The numbers are hidden again once the task ends.
- **FR-TRN-7 [M]:** The first time someone accepts a task as a carrier, they accept a carrier declaration: carry the food hygienically, never open sealed packs, and deliver directly.
- **FR-TRN-8 [M]:** A carrier can back out before pickup, with a reason. The task then reopens.
- **FR-TRN-9 [M]:** If the carrier hasn't completed pickup within `{PICKUP_DEADLINE}`, it counts as a **no-show**. The task reopens if the pickup window allows, and the carrier's reliability score drops.
- **FR-TRN-10 [S]:** Show a rough arrival time, based on straight-line distance and an average speed.

Transport task statuses (first idea): `OPEN → ACCEPTED → PICKED_UP → DELIVERED`, plus `EXPIRED` and `CANCELLED`.
- If the carrier backs out or doesn't show up, the task goes from `ACCEPTED` back to `OPEN`.
- A `DONOR_DELIVERY` task starts at `PICKED_UP`, because the donor already has the food.

### 5.7 Handoffs and acceptance (HND)
- **FR-HND-1 [M]:** **Pickup:** the carrier's app shows a pickup code. The donor, or any staff member at the donor site, enters it. When the code is correct, the task becomes `PICKED_UP`.
- **FR-HND-2 [M]:** **Drop-off:** the NGO's app shows a drop code, and the carrier enters it. When the code is correct, the task becomes `DELIVERED`.
- **FR-HND-3 [M]:** The rule is always the same: the person **receiving** the food shows the code, and the person **handing it over** enters it. Codes appear only in the app and are never sent by SMS. NFR-SEC-3 covers code security.
- **FR-HND-4 [M]:** In `NGO_PICKUP` mode, the pickup code is exchanged between the donor and the NGO member. The drop-off step becomes a simple "received at facility" confirmation by an NGO member.
- **FR-HND-5 [S]:** In `DONOR_DELIVERY` mode there is no pickup step. The drop-off code is exchanged between the NGO and the donor.
- **FR-HND-6 [M]:** After drop-off, an NGO member records the **acceptance result** within `{ACCEPTANCE_WINDOW_MIN}`:
  - accept everything;
  - accept part, with the quantity accepted;
  - reject, with a reason: spoiled, unsafe temperature, tampered packaging, not as described, or other. A photo is optional.
- **FR-HND-7 [S]:** At each handoff, the app takes one GPS reading. If it's more than `{ARRIVAL_CHECK_M}` from the expected point, the handoff is **flagged but still allowed**, because GPS is unreliable indoors and can be faked.
- **FR-HND-8 [S]:** The carrier takes a photo at pickup.
- **FR-HND-9 [S]:** If the donor recorded a seal number, the NGO confirms at drop-off whether the seal is intact.
- **FR-HND-10 [S]:** **Offline handoff:** a handoff can be completed even when one phone has no network, and it syncs once the phone reconnects (NFR-AVL-3).
- **FR-HND-11 [S]:** If the food was picked up but safe-until passes before delivery, the task is flagged and the NGO is told to inspect it carefully.

### 5.8 Notifications (NTF)
- **FR-NTF-1 [M]:** Push notifications for these events:
  - new listing near the NGO;
  - listing claimed (to the donor);
  - task open near the carrier;
  - carrier assigned (to the donor and the NGO);
  - food picked up (to the NGO);
  - food delivered (to the donor);
  - acceptance result (to the donor and the carrier);
  - any timeout, cancellation or expiry (to the people affected).
- **FR-NTF-2 [M]:** Email for login codes, verification results and documents about to expire.
- **FR-NTF-3 [S]:** **Widening radius:** if a listing is still unclaimed after `{RADIUS_STEP_MIN}`, the notification radius grows by `{RADIUS_STEP_KM}`, up to `{MAX_RADIUS_KM}`.
- **FR-NTF-4 [S]:** Reminders:
  - to the donor, when an unclaimed listing is about to expire;
  - to the carrier, when the pickup deadline is near.
- **FR-NTF-5 [S]:** An in-app inbox of past notifications.
- **FR-NTF-6 [C]:** Each user can set notification preferences and quiet hours.

### 5.9 Reputation and moderation (REP)
- **FR-REP-1 [M]:** Any user can report another user or a listing, with a reason. Reports go to an admin queue.
- **FR-REP-2 [S]:** The system calculates scores automatically from what actually happened. There are no star ratings, which people could game.
  - **Donor quality:** how often their food is accepted.
  - **Carrier reliability:** on-time pickups compared with no-shows and backing out.
  - **NGO reliability:** how often they cancel claims.
- **FR-REP-3 [S]:** Users move up and down trust tiers automatically, based on `{T1_CRITERIA}`.
- **FR-REP-4 [S]:** Automatic moderation:
  - a listing is hidden after `{REPORTS_TO_HIDE}` reports;
  - a user is suspended until an admin reviews them, once `{REPORTERS_TO_SUSPEND}` different users have reported them.

### 5.10 Admin panel (ADM)
- **FR-ADM-1 [M]:** The admin panel is a **separate web app** and is not included in the public Android app. Admin accounts can only be created with a server-side command, never through sign-up.
- **FR-ADM-2 [M]:** Admins can search users and orgs, and suspend or reinstate them with a reason. Suspending an org cancels its claims that haven't been picked up yet (see Q2 in §10).
- **FR-ADM-3 [M]:** Admins work through the verification queue (FR-VER-4/5) and the reports queue (FR-REP-1).
- **FR-ADM-4 [M]:** Admins can see the full history of any listing, claim or task, built from the custody log.
- **FR-ADM-5 [S]:** Admins must use two-factor authentication with an authenticator app (TOTP).
- **FR-ADM-6 [S]:** Admins can run the custody log integrity check (FR-AUD-4) and see the result.
- **FR-ADM-7 [S]:** Admins can view and change tunables without redeploying.
- **FR-ADM-8 [C]:** A dashboard of platform stats: listings, deliveries, rejection rate, and the average time from posting to delivery.

### 5.11 Audit and custody log (AUD)
- **FR-AUD-1 [M]:** Every status change and handoff is recorded as a **custody event**, containing:
  - who did it, and in which role;
  - what happened;
  - when, by both the server's clock and the phone's clock;
  - where, if a GPS reading was taken;
  - hashes of any evidence, such as photos.
- **FR-AUD-2 [M]:** Custody events can only be added, never changed. Each event is linked to the previous one by its hash, so any change or deletion can be detected (NFR-INT-1).
- **FR-AUD-3 [M]:** Every admin action is recorded, including viewing a document.
- **FR-AUD-4 [S]:** An integrity check recalculates the hash chain and reports the first event that doesn't match.
- **FR-AUD-5 [M]:** Donors and NGOs can see the custody history of their own donations.

### 5.12 Impact (IMP)
- **FR-IMP-1 [S]:** Every user and org sees how much food they've helped deliver: number of donations, servings, and an estimate in kg (`{KG_PER_SERVING}`).
- **FR-IMP-2 [C]:** These numbers can be filtered by date range and exported to CSV for CSR reports.
- **FR-IMP-3 [C]:** A public page with platform-wide totals.

### 5.13 Demo and testing (DEMO)
- **FR-DEMO-1 [M]:** A seed script creates made-up orgs and users, plus **fake documents watermarked "SAMPLE"**. No real people, organisations or documents are ever used.
- **FR-DEMO-2 [M]:** A dev OTP provider shows or logs codes, so the demo never needs real SMS.
- **FR-DEMO-3 [S]:** Ready-made demo accounts for each role, listed in the README, so reviewers can try the app.
- **FR-DEMO-4 [S]:** A **simulator** generates realistic activity at a configurable rate:
  - donors post, with peaks after lunch and dinner;
  - NGOs claim;
  - carriers accept;
  - handoffs complete;
  - some food gets rejected.

  The same simulator drives the load tests (NFR-SCALE).
- **FR-DEMO-5 [C]:** Demo data resets every night.

---

## 6. Non-functional requirements

Unless stated otherwise, targets are measured on the demo server under the load in NFR-SCALE-2. These are first guesses; update them after the first load test.

### 6.1 Scale (SCALE)
- **NFR-SCALE-1: design target.** One large city:
  - 10,000 registered users, 1,000 active each day;
  - 300 donor orgs and 100 NGOs;
  - 1,000 listings a day, with half posted in two 2-hour peaks (after lunch and after dinner).
- **NFR-SCALE-2: load test target.** The server handles **100 requests/second** of mixed simulated traffic, meets the latency targets in 6.2, and returns errors for fewer than 0.1% of requests. This is several times the expected peak; the calculation goes in the design doc.
- **NFR-SCALE-3:** The API and worker processes keep no state in memory between requests, so more copies can be run later. The design doc explains what would break first and how to scale past it.

### 6.2 Performance (PERF)
- **NFR-PERF-1:** Response times, measured on the server:
  - 95% of read requests finish in under **300 ms**;
  - 95% of write requests finish in under **500 ms**;
  - 99% of all requests finish in under **1 s**.
- **NFR-PERF-2:** For 95% of listings, NGOs in the first radius are notified within **60 s** of posting.
- **NFR-PERF-3:** 95% of background jobs (expiry, timeouts, reopening tasks) run within **60 s** of when they're due.
- **NFR-PERF-4:** 95% of automated document checks finish within **5 minutes** of upload.
- **NFR-PERF-5:** On a low-end phone (NFR-USE-1), the app opens in under **4 s** and the listings screen shows saved data straight away.

### 6.3 Consistency and invariants (CON)
These rules must **always** hold, even when requests arrive at the same moment or are retried. Each one gets its own automated concurrency test.
- **NFR-CON-1:** For every listing, the claimed quantity (active plus completed claims) never exceeds the listed quantity.
- **NFR-CON-2:** A transport task never has more than one accepted carrier.
- **NFR-CON-3:** A carrier never has more than `{CARRIER_MAX_ACTIVE_TASKS}` active tasks.
- **NFR-CON-4:** Each handoff code can be used successfully only once.
- **NFR-CON-5:** Statuses change only along the allowed transitions. An invalid change is refused completely and never half-applied.
- **NFR-CON-6:** Every org always has at least one `ORG_ADMIN`.
- **NFR-CON-7:** Every request that changes data accepts an **idempotency key** (a unique ID for that request). Retrying the same request never creates a second listing, claim or handoff.
- **NFR-CON-8:** Some data may be slightly out of date, where that does no harm:
  - the nearby-listings feed, by up to **10 s**;
  - reputation scores and trust tiers, by up to **5 minutes**;
  - impact numbers, by up to **1 hour**.

### 6.4 Availability (AVL)
- **NFR-AVL-1:** The demo aims to be up **99% of each month**. The single server is a known single point of failure, and this is written down rather than hidden.
- **NFR-AVL-2:** When the server can't be reached, the app says so clearly. Actions waiting to be sent are retried with the same idempotency key.
- **NFR-AVL-3 (with FR-HND-10):** An offline handoff syncs within **5 minutes** of either phone reconnecting.

### 6.5 Durability and reliability (DUR)
- **NFR-DUR-1:** Once the API has confirmed a write, it's never lost. A status change and its custody event are saved in the same database transaction.
- **NFR-DUR-2:** Every notification is sent **at least once**, and duplicates are filtered out so users never see the same alert twice.
- **NFR-DUR-3:** Backups:
  - The database is backed up to object storage every day, and backups are kept for 7 days.
  - At most **24 hours** of data can be lost, and the system can be restored within **4 hours**.
  - A restore has actually been tested and the steps written down.

### 6.6 Security (SEC)
- **NFR-SEC-1: deny by default.** Every endpoint checks both **what the user may do** (their role or capability) and **which records they may touch**. Tests prove that nobody can read or change another user's or org's data by changing an ID in the request.
- **NFR-SEC-2: login and phone OTPs.**
  - 6 digits, valid for 10 minutes, with 5 attempts allowed.
  - At most 5 codes sent per hour to the same email or phone, and from the same IP address.
  - Stored hashed, never in plain text.
- **NFR-SEC-3: handoff OTPs.**
  - 6 digits, single use, tied to one task and one handoff step.
  - Valid only while that step is active.
  - After 5 wrong attempts, the code is locked for `{OTP_LOCK_MIN}` and both people are notified.
  - Stored hashed, and compared in constant time, so response timing doesn't leak which digits were right.
- **NFR-SEC-4: sessions.**
  - Access tokens last **15 minutes**.
  - Refresh tokens last **30 days**, are replaced every time they're used, and can be revoked.
- **NFR-SEC-5: rate limits.** Writes and login attempts are limited per user and per IP address. Requests over the limit get a `429` response with a `Retry-After` header.
- **NFR-SEC-6: uploads.**
  - Documents: at most 10 MB, PDF, JPEG or PNG only.
  - Photos: at most `{PHOTO_MAX_KB}` after the phone compresses them.
  - Files are stored in a **private** bucket and can't be viewed until the malware scan passes.
- **NFR-SEC-7: transport and secrets.** HTTPS only, with HSTS. No secrets in the repository; they come from environment variables or a secret store.
- **NFR-SEC-8: checklist.** The design is checked against the OWASP API Security Top 10. **[S]** CI also scans dependencies for known vulnerabilities.

### 6.7 Privacy (PRIV): in line with the DPDP Act 2023
- **NFR-PRIV-1: collect as little as possible.** No Aadhaar or other government ID numbers from individuals. Only orgs upload documents.
- **NFR-PRIV-2: location.** Only the assigned carrier, the donor and admins can see the exact pickup location. Everyone else sees the approximate one (FR-CLM-2).
- **NFR-PRIV-3: contact details.** Phone numbers are only shared between the people on an active task (FR-TRN-6).
- **NFR-PRIV-4: how long data is kept.**

  | Data | Kept for |
  |---|---|
  | Rejected verification documents | 30 days after rejection, unless resubmitted |
  | Approved verification documents | While the org is active, plus 90 days |
  | Listing photos | 30 days after the listing closes (90 days if the food was rejected or disputed). Photo **hashes** stay in the custody log permanently |
  | Handoff GPS readings | 90 days |
  | Server logs | 14 days |
  | Deleted accounts | Personal data deleted or anonymised within **30 days** |

- **NFR-PRIV-5: deleting data from a log that can't change.** The custody log only stores **pseudonymous IDs**, never names, phone numbers or addresses. Deleting an account deletes the link between the ID and the person, so the hash chain can still be checked.
- **NFR-PRIV-6:** Every time an admin views a document, it's logged (FR-AUD-3). Logs never contain OTPs, tokens, document contents or exact home locations.

### 6.8 Integrity (INT)
- **NFR-INT-1:** Any change to, or deletion of, a past custody event is caught by the integrity check (FR-AUD-4).
- **NFR-INT-2:** The database account the app uses can't update or delete custody events. The **database itself** enforces this, not just the app code.
- **NFR-INT-3 [S]:** The integrity check runs every day. A mismatch triggers an alert.

### 6.9 Usability (USE)
- **NFR-USE-1:** The app works well on low-end Android phones: 2 GB RAM, 720p screen, the oldest Android version the chosen Expo SDK supports, and slow 3G-like connections.
- **NFR-USE-2:** A returning donor can post a listing in under **60 seconds**, using saved sites and default values.
- **NFR-USE-3:** English for the MVP. All app text lives in translation files from day one, so other languages can be added. Hindi is **[S]**.
- **NFR-USE-4: accessibility.**
  - Tap targets are at least 48 dp.
  - Text follows the phone's font-size setting.
  - Colour is never the only signal. For example, veg and non-veg also get an icon and a label.

### 6.10 Cost (COST)
- **NFR-COST-1:** Running costs are **₹0/month**. Every third-party service stays within its free tier.
- **NFR-COST-2:** Photos are the biggest storage risk. At the design scale, uncompressed photos would fill the free storage within weeks. To prevent that:
  - the phone compresses each photo to at most `{PHOTO_MAX_KB}`;
  - each listing requires exactly one photo;
  - photos are deleted on the NFR-PRIV-4 schedule.

  Check: 1,000 listings/day × 200 KB × 30 days ≈ 6 GB, which is under the free limit.

### 6.11 Maintainability (MNT)
- **NFR-MNT-1:** The API contract is an OpenAPI spec generated from the backend code. The app's typed API client is generated from that spec.
- **NFR-MNT-2:** Every database schema change is a versioned migration. No manual changes to the database.
- **NFR-MNT-3:** Tests:
  - unit tests for the business logic;
  - integration tests that run against a real Postgres database;
  - automated tests for every status change and every rule in §6.3.
- **NFR-MNT-4:** CI runs the linter, strict type checks and all tests on every push. Nothing is merged unless CI passes.
- **NFR-MNT-5:** One command starts the whole system locally.
- **NFR-MNT-6:** Every major design decision is written up as an ADR (architecture decision record).

### 6.12 Observability (OBS)
- **NFR-OBS-1:** Logs are structured, and each request carries an ID that shows up in both the API's and the worker's logs.
- **NFR-OBS-2:** A dashboard shows:
  - request rate and response times (p50/p95/p99);
  - error rate;
  - job queue size and how far behind the jobs are;
  - failed notifications;
  - flagged handoffs.
- **NFR-OBS-3:** An alert fires if more than 2% of requests fail for 5 minutes, or if jobs fall more than 5 minutes behind.
- **NFR-OBS-4:** Errors are tracked for both the backend and the mobile app.

---

## 7. Constraints and decisions already made

- **Budget:** ₹0/month, free tiers only (NFR-COST-1).
- **Platforms:**
  - Android app (React Native + Expo) for all users. **No iOS app.**
  - Separate web app for admins.
- **Backend:** Python. The framework choice and the reasons for it go in an ADR.
- **Hosting:** one small cloud VM for the demo.
- **No blockchain.** Handoffs are confirmed with OTPs, and records are protected by the hash-chained custody log. The ADR explains why:
  - blockchain can't prove the physical food wasn't tampered with (the oracle problem);
  - cost;
  - privacy.
- **No real SMS.** Sending SMS in India requires TRAI DLT registration. The demo uses email codes plus a dev phone-OTP provider. A real launch would need DLT registration or WhatsApp authentication templates.
- **No live tracking.** Only pickup and drop coordinates, plus a Navigate link.
- **Solo developer.**
- **Legal context:** India, the DPDP Act 2023 and FSSAI's surplus-food rules. All demo data is fictional.

---

## 8. Out of scope (v1)

- Live GPS tracking, route planning, and grouping several pickups into one trip.
- An iOS app.
- Money donations, payments, and 80G tax receipts.
- Individuals or unverified orgs receiving food.
- Automatic checks against government databases (FSSAI, GST and NGO Darpan APIs).
- Sending real SMS or WhatsApp messages.
- In-app chat and masked phone calls.
- More than one city, or languages other than English and Hindi.
- NGOs posting what food they need ("reverse listings"). This is a good idea for v2.
- Predicting demand, or matching donations with ML.
- Blockchain.

---

## 9. Tunables (defaults)

| Tunable | Default | Used by |
|---|---|---|
| `NOTIFY_RADIUS_KM` | 3 km | FR-NTF-1 |
| `RADIUS_STEP_MIN` / `RADIUS_STEP_KM` / `MAX_RADIUS_KM` | 10 min / +2 km / 10 km | FR-NTF-3 |
| `MAX_SAFE_WINDOW` | Room temperature 4 h · kept hot 6 h · refrigerated 24 h · sealed packaged: printed expiry. **These are placeholders; check them against FSSAI guidance** | FR-LST-2 |
| `T0_MAX_SERVINGS` / `T0_MAX_LISTINGS_PER_DAY` | 20 / 2 | §3 |
| `T1_CRITERIA` | ≥ 5 delivered donations, rejection rate < 20%, no upheld reports | §3, FR-REP-3 |
| `T1_MAX_SERVINGS` / `T1_MAX_LISTINGS_PER_DAY` | 100 / 5 | §3 |
| `T2_MAX_SERVINGS` | 1,000 (safety cap) | §3 |
| `MIN_CLAIM` | 5 servings | FR-CLM-4 |
| `SELF_DELIVERY_MAX_KM` | 5 km | FR-LST-8 |
| `MIN_DELIVERIES` (NGO carrier policy) | 3 | FR-ORG-6 |
| `CARRIER_MAX_ACTIVE_TASKS` | 1 | FR-TRN-3, NFR-CON-3 |
| `PICKUP_DEADLINE` | 60 min after accepting, or the end of the pickup window, whichever comes first | FR-TRN-9 |
| `ACCEPTANCE_WINDOW_MIN` | 30 min | FR-HND-6 |
| `ARRIVAL_CHECK_M` | 200 m | FR-HND-7 |
| `APPROX_LOCATION_M` | 500 m | FR-CLM-2 |
| `OTP_LOCK_MIN` | 15 min | NFR-SEC-3 |
| `REVIEW_LOCK_MIN` | 30 min | FR-VER-4 |
| `DOC_EXPIRY_REMINDERS` | 30 days and 7 days before | FR-VER-8 |
| `REPORTS_TO_HIDE` / `REPORTERS_TO_SUSPEND` | 3 / 3 | FR-REP-4 |
| `PHOTO_MAX_KB` | 200 KB | NFR-SEC-6, NFR-COST-2 |
| `KG_PER_SERVING` | 0.4 kg (placeholder) | FR-IMP-1 |

---

## 10. Assumptions and open questions

1. **How do the carrier and donor contact each other?** For now, phone numbers are shown while the task is active. Alternatives: in-app chat, or masked calls (paid).
2. **What happens when an NGO is suspended or expires while it has active claims?** Proposal: cancel the claims that haven't been picked up, and let food already on its way be delivered.
3. **Can an org be both a donor and a receiver?** One example is a temple kitchen that also feeds people. Assumed no: each org has one type.
4. **Should Tier 0 individuals only be allowed to donate sealed or packaged food?** Assumed no: cooked food is allowed, with lower limits.
5. **The safe-time limits** need checking against FSSAI guidance before they're treated as real.
6. **Throwing away rejected food:** assumed out of scope. The NGO or carrier disposes of it according to local rules.
7. **One hash chain per donation, or one for the whole system?** To be decided in the design doc. It affects how many writes compete for the chain and how verification works.
8. **Can a carrier belong to the receiving NGO?** Assumed yes, since this is the normal volunteer case.

---

## 11. Glossary

- **Listing:** food a donor has posted.
- **Claim:** an NGO reserving all or part of a listing.
- **Transport task:** the job of moving one claim from pickup to drop-off.
- **Carrier:** the person doing a transport task.
- **Leg:** one half of the trip: pickup (donor → carrier) or drop-off (carrier → NGO).
- **Handoff OTP:** the one-time code proving a leg happened between the right people.
- **Custody event:** one entry in the add-only, hash-chained log.
- **Trust tier:** a level (0–2) that sets how much a user can donate or carry.
- **Safe-until:** the latest time the food is considered safe to eat.
- **Site:** one physical location of an org.
- **Tunable:** a setting whose value can be changed; see §9.
