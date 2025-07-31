CREATE INDEX IF NOT EXISTS idx_fts_hotels
ON hotels
USING gin (
  to_tsvector('simple',
    coalesce(hotel_name, '') || ' ' ||
    coalesce(city, '') || ' ' ||
    coalesce(country, '') || ' ' ||
    coalesce(description, '')
  )
);

CREATE INDEX IF NOT EXISTS idx_fts_reviews
ON reviews
USING gin (
  to_tsvector('simple',
    coalesce(headline, '') || ' ' ||
    coalesce(pros, '') || ' ' ||
    coalesce(cons, '')
  )
);
