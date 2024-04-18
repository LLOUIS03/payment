-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO merchant (id, name, email, password)
VALUES  ('b1b9b1b9-1b9b-1b9b-1b9b-1b9b1b9b1b9b', 'Wendy', 'jkirkness0@Wendy.com', '$2a$04$8yUumMzTc3YvTG0GU1bCueQ5ZiPKGtiDFQ0zzvOP5KVfV1wxRqhJy'),
        ('A4A3E811-F71C-4A40-BD1A-504FE43014D3', 'McDonalds', 'crawlence1@McDonalds.com', '$2a$04$6N.DNmmWuxJfCx5G4MUMBuAOjHlAy53ihksUm0uCL92K0hFK0FVq6'),
        ('4353B756-62FE-4EAE-A5FD-BCAF27402E85', 'BurgerKing', 'sexell2@merchant.com', '$2a$04$wRRGoW0lqTgyCY.yfjJIEOVX8rVDlAjbVZcuT8QUD2w08B8ZDmC8y');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
