CREATE TABLE IF NOT EXISTS url
  (
     shortenedid    VARCHAR NOT NULL UNIQUE,
     originalurl    VARCHAR NOT NULL,
     apikey         VARCHAR NOT NULL,
     creationtime   TIMESTAMP NOT NULL,
     expirationdate TIMESTAMP NOT NULL
  );

CREATE TABLE IF NOT EXISTS stats
  (
     shortenedid VARCHAR NOT NULL UNIQUE,
     redirects   INTEGER DEFAULT 0,
     visitdate   TIMESTAMP
  );