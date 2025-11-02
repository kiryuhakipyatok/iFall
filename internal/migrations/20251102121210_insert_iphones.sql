-- +goose Up
-- +goose StatementBegin
INSERT INTO iphones (id, name, price, change, color) VALUES 
('apple_iphone_17_256gb_chernyy', 'Apple iPhone 17 256GB (черный)', 3191, 0, '353839'),
('apple_iphone_17_256gb_belyy', 'Apple iPhone 17 256GB (белый)', 3191.3, 0, 'F5F5F5'),
('apple_iphone_17_256gb_goluboy', 'iPhone 17 256GB (голубой)', 3191.3, 0, '96AED1'),
('apple_iphone_17_256gb_sirenevyy', 'Apple iPhone 17 256GB (сиреневый)', 3414.45, 0, 'DFCEEA'),
('apple_iphone_17_256gb_zelenyy', 'Apple iPhone 17 256GB (зеленый)', 3316.6, 0, 'A9B689'),
('apple_iphone_17_dual_esim_256gb_chernyy', 'Apple iPhone 17 Dual eSIM 256GB (черный)', 3132.34, 0, '353839'),
('apple_iphone_17_dual_esim_256gb_goluboy', 'Телефон Apple iPhone 17 Dual eSIM 256GB (голубой)', 3203.07, 0, '96AED1'),
('apple_iphone_17_dual_esim_256gb_sirenevyy', 'Телефон Apple iPhone 17 Dual eSIM 256GB (сиреневый)', 3283.9, 0, 'DFCEEA')
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM iphones
-- +goose StatementEnd
