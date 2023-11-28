-- +migrate Up

INSERT INTO predictions("id", "username", "status", "feedback_rate", "created_at", "updated_at")
VALUES
  ('1ed1e4ec-c73f-6f67-8389-f60bb01150bc', 'test', 'completed', 5, now(), now());

-- +migrate Down
