ALTER TABLE reviews ADD CONSTRAINT reviews_unique_key UNIQUE (hotel_id, name, date, headline);
