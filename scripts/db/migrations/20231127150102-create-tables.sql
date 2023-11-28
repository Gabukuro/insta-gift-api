-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products (
  id UUID DEFAULT uuid_generate_v4(),
  name TEXT,
  main_category TEXT,
  sub_category TEXT,
  image TEXT,
  link TEXT,
  ratings FLOAT,
  no_of_ratings FLOAT,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY (id)
);

CREATE TABLE predictions (
  id UUID NOT NULL,
  username TEXT,
  feedback_rate INTEGER,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY (id)
);

CREATE TABLE prediction_products (
  prediction_id UUID NOT NULL,
  product_id UUID NOT NULL,
  rank_position INTEGER,
  PRIMARY KEY (prediction_id, product_id)
);

-- +migrate Down
DROP TABLE prediction_products;

DROP TABLE predictions;

DROP TABLE products;
