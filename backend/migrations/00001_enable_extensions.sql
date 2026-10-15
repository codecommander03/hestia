-- +goose Up
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS citext;

-- +goose Down
-- DROP EXTENSION IF EXISTS postgis;
-- postgis is not dropped: the postgis/postgis image installs it (plus
-- topology and tiger_geocoder, which depend on it) when the database is created.
DROP EXTENSION IF EXISTS citext;