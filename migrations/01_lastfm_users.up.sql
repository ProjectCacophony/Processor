CREATE TABLE lastfm_users (
  user_id          VARCHAR                   NOT NULL,
  lastfm_username  VARCHAR                   NOT NULL,
  created_at       TIMESTAMP WITH TIME ZONE  NOT NULL  DEFAULT NOW(),
  updated_at       TIMESTAMP WITH TIME ZONE  NOT NULL  DEFAULT NOW(),
  CONSTRAINT user_id PRIMARY KEY(user_id)
)
;
