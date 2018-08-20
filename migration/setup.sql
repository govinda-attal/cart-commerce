CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS cartItem;
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS inventory;

DROP TYPE IF EXISTS cart_state;

CREATE TYPE cart_state AS ENUM ('INPROGRESS', 'CANCELLED', 'SETTLED');

CREATE TABLE inventory (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  itemName text unique not null,
  stock int DEFAULT 0,
  avail int DEFAULT 0,
  price DECIMAL(13, 2)
);

CREATE TABLE cart (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  state cart_state not null
);

CREATE TABLE cartItem (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  cartId uuid references cart not null,
  itemName text references inventory(itemName),
  quantity int DEFAULT 0
);



INSERT INTO inventory (itemName, stock, avail, price) VALUES ('BELTS', 10, 10, 20.00);
INSERT INTO inventory (itemName, stock, avail, price) VALUES ('SHIRTS', 5, 5, 60.00);
INSERT INTO inventory (itemName, stock, avail, price) VALUES ('SUITS', 2, 2, 300.00);
INSERT INTO inventory (itemName, stock, avail, price) VALUES ('TROUSERS', 4, 4, 70.00);
INSERT INTO inventory (itemName, stock, avail, price) VALUES ('SHOES', 1, 1, 120.00);
INSERT INTO inventory (itemName, stock, avail, price) VALUES ('TIES', 8, 8, 20.00);

CREATE OR REPLACE FUNCTION updatecart(
 argCartID text, 
 argItemName text,
 argQuantity NUMERIC) 
RETURNS text AS $$
declare item_found NUMERIC;
declare prevQuantity NUMERIC;
declare availDiff NUMERIC;
declare varCartID UUID;
BEGIN
  
  IF (argCartID = '') IS TRUE THEN
    INSERT INTO cart (state) VALUES ('INPROGRESS') RETURNING id INTO varCartID;
  ELSE
    varCartID := argCartID::UUID;
  END IF;  
  
  
  SELECT count(*) INTO item_found FROM cartItem CI WHERE CI.cartID = varCartID AND CI.itemName = argItemName;

  IF item_found = 0 THEN
    prevQuantity := 0;
  ELSE
    SELECT CI.quantity INTO prevQuantity FROM cartItem CI WHERE CI.cartID = varCartID AND CI.itemName = argItemName;
  END IF;

  availDiff := prevQuantity - argQuantity;

  IF argQuantity > 0 THEN
    IF item_found = 0 THEN
      INSERT INTO cartItem (cartID, itemName, quantity) VALUES (varCartID, argItemName, argQuantity);
    ELSE    
      UPDATE cartItem CI SET quantity = argQuantity WHERE CI.cartId = varCartID AND CI.itemName = argItemName;
    END IF;
  ELSE
    DELETE FROM cartItem CI WHERE CI.cartId = varCartID AND CI.itemName = argItemName;
  END IF;

  UPDATE inventory SET avail = avail + availDiff WHERE itemName = argItemName;

  RETURN varCartID::text;

END; $$
 
LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION updateCartState(
   argCartID text,
   argState text
) 
RETURNS BOOLEAN AS $$
DECLARE
    rec RECORD;
BEGIN
    FOR rec IN SELECT CI.itemName, CI.quantity FROM cartItem CI WHERE CI.cartID = argCartID::UUID
    LOOP 
      CASE argState
      WHEN 'SETTLED' THEN
        UPDATE inventory SET stock = stock - rec.quantity WHERE itemName = rec.itemName; 
      WHEN 'CANCELLED' THEN
        UPDATE inventory SET avail = avail + rec.quantity WHERE itemName = rec.itemName;
      END CASE;
    END LOOP;
    UPDATE cart C SET state = argState::cart_state WHERE C.ID = argCartID::UUID;
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;