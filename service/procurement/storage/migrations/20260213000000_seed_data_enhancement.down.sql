-- Reverse seed data enhancement: remove cancelled purchase order and related records.
BEGIN;

DELETE FROM purchase_order_fee WHERE uuid = '61000000-0000-0000-0000-000000000007';
DELETE FROM purchase_order_line WHERE uuid IN (
    '60000000-0000-0000-0000-000000000020',
    '60000000-0000-0000-0000-000000000021'
);
DELETE FROM purchase_order WHERE uuid = '50000000-0000-0000-0000-000000000008';

COMMIT;
