CREATE TABLE subscriptions (
    ID SERIAL PRIMARY KEY,
    ServiceName TEXT NOT NULL,
    Price INTEGER NOT NULL,
    UserID UUID NOT NULL,
    StartDate DATE NOT NULL
)