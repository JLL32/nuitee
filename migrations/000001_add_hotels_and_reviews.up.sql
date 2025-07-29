CREATE TABLE hotels (
    hotel_id INTEGER PRIMARY KEY,
    main_image_th TEXT,
    hotel_name TEXT NOT NULL,
    phone TEXT,
    email TEXT,
    address TEXT,
    city TEXT,
    state TEXT,
    country TEXT,
    postal_code TEXT,
    stars INTEGER,
    rating DECIMAL(3,2),
    review_count INTEGER,
    child_allowed BOOLEAN,
    pets_allowed BOOLEAN,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    hotel_id INTEGER NOT NULL,
    average_score INTEGER,
    country TEXT,
    type TEXT,
    name TEXT,
    date TIMESTAMP,
    headline TEXT,
    language CHAR(2),
    pros TEXT,
    cons TEXT,
    source TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (hotel_id) REFERENCES hotels(hotel_id) ON DELETE CASCADE
);
