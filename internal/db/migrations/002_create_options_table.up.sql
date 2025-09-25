CREATE TABLE options
(
    id      INT GENERATED ALWAYS AS IDENTITY,
    text    TEXT NOT NULL,
    poll_id INT  NOT NULL,
    votes   INT DEFAULT 0,

    CONSTRAINT fk_poll
        FOREIGN KEY (poll_id)
            REFERENCES polls (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE
);

CREATE INDEX idx_options_poll_id on options(poll_id);