CREATE TABLE IF NOT EXISTS launchers (
  id serial primary key,
  customer_id bigint NOT NULL UNIQUE,
  account_id bigint NOT NULL,
  first_name text NOT NULL,
  last_name text NOT NULL,
  interest_rate float NOT NULL,
  credit_limit bigint NOT NULL,
  balance float default 0,
  due_date timestamp with time zone NOT NULL,
  minimum_payment float default 0,
  reward_balance float default 0,
  created timestamp with time zone default current_timestamp,
  modified timestamp with time zone default NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id bigint primary key, -- From Capital One
    launchers_id int NOT NULL,
    type text NOT NULL,
    merchant text NOT NULL,
    amount float NOT NULL,
    purchase_date timestamp with time zone,
    CONSTRAINT launchers_fk
      FOREIGN KEY(launchers_id) REFERENCES launchers
      ON DELETE CASCADE
);
