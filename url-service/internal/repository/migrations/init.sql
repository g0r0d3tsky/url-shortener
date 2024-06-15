CREATE TABLE url_data (
                          id uuid PRIMARY KEY,
                          long_url varchar(3000),
                          short_url varchar(255),
                          expires_at timestamp with time zone
);

CREATE TABLE url_keys (
                          id uuid PRIMARY KEY,
                          key_serial serial,
                          encode varchar(255),
                          url_id uuid,
                          FOREIGN KEY (url_id) REFERENCES url_data (id) ON DELETE SET NULL
);

CREATE SEQUENCE key_serial_seq OWNED BY url_keys.key_serial;
