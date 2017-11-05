-- DROP TABLE users CASCADE;
-- DROP TABLE contacts;
-- DROP TABLE checkins;

-- launchers table represents kid's credit card accounts.
CREATE TABLE launchers (
  id serial primary key,
  customer_id bigint NOT NULL,
  account_id bigint NOT NULL,
  first_name text NOT NULL,
  last_name text NOT NULL,
  credit_limit int NOT NULL,
  balance int default 0,
  due_date timestamp with time zone NOT NULL,
  minimum_payment int default 0,
  reward_balance int default 0, 
  created timestamp with time zone default current_timestamp,
  modified timestamp with time zone default NULL
);

--
CREATE TABLE transactions (
    id int64 primary key, -- From Capital One
    launchers_id int NOT NULL,
    type text NOT NULL,
    merchant text NOT NULL,
    amount float NOT NULL,
    purchase_date timestamp with time zone,
    CONSTRAINT launchers_fk
      FOREIGN KEY(launchers_id) REFERENCES launchers
      ON DELETE CASCADE
);
