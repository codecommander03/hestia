# Hestia

> Hestia is the Greek goddess of the hearth and home. She received the first portion of every meal.

A mobile-first platform that gets surplus food from cafeterias, canteens, hostels, restaurants and event venues to verified NGOs before it goes bad.

- **Donors** post surplus food.
- **Verified NGOs** claim all or part of it.
- **Carriers** move it from the donor to the NGO.
- Every handoff is confirmed with an **in-app OTP** and recorded in a **tamper-evident, hash-chained custody log**.

## Status

**Design phase.** No code yet. The requirements and data model are written down first, so that every decision can be traced and defended.

| Document | What's in it |
|---|---|
| [Requirements](docs/requirements.md) | Actors, trust tiers, functional and non-functional requirements, tunables, open questions |
| [Data model](docs/data-model.md) | Entities, relationships, how the schema enforces the consistency rules, ER diagram |

## Planned stack

- **Mobile app:** React Native + Expo (Android only)
- **Admin panel:** separate web app
- **Backend:** Python
- **Database:** Postgres + PostGIS
- **Hosting:** one small cloud VM, ₹0/month on free tiers

## Context

A portfolio project set in India (DPDP Act 2023, FSSAI surplus-food rules). It's a demo, not a live service: **all data, organisations and documents are fictional.**

## License

[MIT](LICENSE)
