CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Create tables that have no dependencies
CREATE TABLE IF NOT EXISTS "user" (
  id UUID PRIMARY KEY,
  name VARCHAR NOT NULL,
  email VARCHAR UNIQUE NOT NULL,
  phone VARCHAR UNIQUE NOT NULL,
  password VARCHAR NOT NULL,
  role VARCHAR NOT NULL CHECK (role IN ('user', 'courier', 'admin')),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "category" (
  id UUID PRIMARY KEY,
  name VARCHAR NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "product" (
  id UUID PRIMARY KEY,
  category_id UUID NOT NULL REFERENCES "category"(id),
  name VARCHAR NOT NULL,
  description TEXT,
  price DECIMAL NOT NULL,
  image_url VARCHAR,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "orderiteam" (
  id UUID PRIMARY KEY,
  order_id UUID NOT NULL REFERENCES "order"(id),
  product_id UUID NOT NULL REFERENCES "product"(id),
  total DECIMAL(10, 2) NOT NULL,
  quantity INT NOT NULL,
  price DECIMAL NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

-- 2. Create tables that reference the above tables
CREATE TABLE IF NOT EXISTS "order" (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES "user"(id),
  total_price DECIMAL NOT NULL,
  status VARCHAR NOT NULL CHECK (status IN ('pending', 'confirmed', 'picked_up', 'delivered')) DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "payment" (
  id UUID PRIMARY KEY DEFAULT,
  user_id UUID NOT NULL REFERENCES "user"(id),
  order_id UUID NOT NULL REFERENCES "order"(id),
  is_paid BOOLEAN NOT NULL DEFAULT false,
  payment_method VARCHAR NOT NULL CHECK (payment_method IN ('click', 'payme', 'naxt pul')),
  created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "courierassignment" (
  id UUID PRIMARY KEY DEFAULT,
  order_id UUID NOT NULL REFERENCES "order"(id),
  courier_id UUID NOT NULL REFERENCES "user"(id),
  status VARCHAR NOT NULL CHECK (status IN ('assigned', 'picked_up', 'en_route', 'delivered', 'payment_collected')),
  assigned_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "notification" (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES "user"(id),
  message TEXT NOT NULL,
  is_read BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "deliveryhistory" (
  id UUID PRIMARY KEY DEFAULT,
  courier_id UUID NOT NULL REFERENCES "user"(id),
  order_id UUID NOT NULL REFERENCES "order"(id),
  earnings DECIMAL NOT NULL,
  delivered_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "locations" (
  id UUID PRIMARY KEY DEFAULT,
  user_id UUID REFERENCES "user"(id),
  address VARCHAR,
  latitude DECIMAL(10, 8),
  longitude DECIMAL(11, 8),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "combos" (
  id UUID PRIMARY KEY,
  name VARCHAR,
  description TEXT,
  price DECIMAL(10, 2),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);


CREATE TABLE IF NOT EXISTS "banner" (
  image_url VARCHAR,
  created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "branch" (
  id UUID PRIMARY KEY DEFAULT,
  name VARCHAR NOT NULL,
  address VARCHAR NOT NULL,
  latitude DECIMAL(10, 8),
  longitude DECIMAL(11, 8),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);
