INSERT INTO `gorent`.`rents` (
    `owner_id`, 
    `renter_id`, 
    `product_id`, 
    `restriction_id`,
    `processed`, 
    `total_cost`, 
    `duration`, 
    `start_date`, 
    `end_date`, 
    `created_at`, 
    `updated_at`
) VALUES (
    '1', 
    '2', 
    '1', 
    '1',
    true, 
    '100', 
    '1', 
    '2021/06/20', 
    '2021/06/20', 
    '2021/06/05', 
    '2021/06/05'
),(
    '1', 
    '3', 
    '1', 
    '1',
    true, 
    '100', 
    '2', 
    '2021/06/23', 
    '2021/06/24', 
    '2021/06/05', 
    '2021/06/05'
),(
    '1', 
    '2', 
    '1', 
    '1',
    true, 
    '100', 
    '1', 
    '2021/06/27', 
    '2021/06/27', 
    '2021/06/05', 
    '2021/06/05'
),(
    '1',
    '2',
    '3',
    '1',
    false, 
    '10', 
    '1', 
    '2021-06-15 00:00:00', 
    '2021-06-15 00:00:00', 
    '2021-06-15 00:00:00',
    '2021-06-15 00:00:00'
);
