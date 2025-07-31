ALTER TABLE hotels ADD COLUMN fts tsvector
  GENERATED ALWAYS AS (
    to_tsvector('simple',
      coalesce(hotel_name, '') || ' ' ||
      coalesce(city, '') || ' ' ||
      coalesce(country, '') || ' ' ||
      coalesce(description, '')
    )
  ) STORED;

CREATE INDEX idx_fts_hotels ON hotels USING gin (fts);

ALTER TABLE reviews ADD COLUMN fts tsvector
  GENERATED ALWAYS AS (
    to_tsvector('simple',
      coalesce(headline, '') || ' ' ||
      coalesce(pros, '') || ' ' ||
      coalesce(cons, '')
    )
  ) STORED;

CREATE INDEX idx_fts_reviews ON reviews USING gin (fts);
