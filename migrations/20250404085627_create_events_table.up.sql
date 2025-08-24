CREATE TABLE events(
    id SERIAL PRIMARY KEY,
    user_id int not null,
    title text not null,
    time timestamptz not null,
    timeToNotify timestamp,
    is_sent BOOLEAN DEFAULT FALSE
)