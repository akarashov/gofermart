BEGIN;
create table if not exists users (
    id SERIAL primary key,
    login VARCHAR(255) unique not null,
    password_hash VARCHAR(255) not null,
    constraint check_login_length check (LENGTH(login) between 3 and 255)
);

create type order_status_type as enum (
    'NEW',
    'PROCESSING', 
    'INVALID',
    'PROCESSED'
);

create table if not exists orders (
    id SERIAL primary key,
    user_id INT not null references users(id) on
delete
	cascade,
	order_number VARCHAR not null unique,
	-- leadering zeros are allowed
	order_status order_status_type not null default 'NEW',
	accrual numeric(10, 2),
	-- sum of accruals for the order (RUB,CENTS) it can be NULL if order is NEW or PROCESSING
	created_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
	-- RFC3339 format
	uploaded_at TIMESTAMPTZ,
	constraint unique_user_order unique(user_id, order_number),
	constraint valid_timeline check (uploaded_at >= created_at),
	constraint order_number_digits check (order_number ~ '^[0-9]+$'),
	constraint check_accrual_non_negative check (accrual >= 0)
);

create table if not exists balance (
    id SERIAL primary key,
    user_id INT not null references users(id) on delete cascade,
	current numeric(10, 2) default 0,
	-- sum of current balance loyalty points (RUB,CENTS)
	withdrawn numeric(10, 2) default 0,
	-- sum of withdrawn loyalty points (RUB,CENTS)
	created_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
	-- RFC3339 format
	uploaded_at TIMESTAMPTZ,
	constraint valid_timeline check (uploaded_at >= created_at),
	constraint check_withdrawn_non_negative check (withdrawn >= 0),
	constraint check_current_non_negative check (current >= 0)
);

create table if not exists withdrawals (
    id SERIAL primary key,
    user_id INT not null references users(id) on
delete
	cascade,
	order_id INT not null references orders(id) on
	delete
		cascade,
		sum numeric(10, 2) not null,
		processed_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
		-- RFC3339 format
    constraint unique_withdrawal_order unique(order_id),
		constraint check_withdrawal_sum_positive check (sum > 0)
);

create or replace
function update_uploaded_at_column()
returns trigger as $$
begin
    NEW.uploaded_at = CURRENT_TIMESTAMP;

return new;
end;

$$ language 'plpgsql';

create or replace function withdraw(p_user_id INT, p_order_number VARCHAR, p_sum numeric(10, 2))
returns void
language plpgsql
as $$
declare
    v_balance_id INT;
	v_order_id INT;

begin
    select
	id
into
	v_balance_id
from
	balance
where
	user_id = p_user_id
    for
update;
	
	update
	balance
set
	current = current - p_sum,
	withdrawn = withdrawn + p_sum,
	uploaded_at = CURRENT_TIMESTAMP
where
	id = v_balance_id;

insert
	into
	orders (
        user_id,
	order_number
    )
values (
        p_user_id,
        p_order_number
    )
    returning id
into
	v_order_id;

insert
	into
	withdrawals (
        user_id,
	order_id,
	sum,
	processed_at
    )
values (
        p_user_id,
        v_order_id,
        p_sum,
        CURRENT_TIMESTAMP
    );
end;

$$;


create trigger update_orders_uploaded_at
    before
update or insert
	on
	orders 
    for each row
    execute procedure update_uploaded_at_column();

create trigger update_balances_uploaded_at
    before
update
	on
	balance
    for each row
    execute procedure update_uploaded_at_column();

create index if not exists idx_users_login on
users(login);

create index if not exists idx_orders_order_number on
orders(order_number);

create index if not exists idx_orders_user_id on
orders(user_id);

create index if not exists idx_orders_status on
orders(order_status)
where
order_status in ('NEW', 'PROCESSING');

create index if not exists idx_orders_uploaded_at on
orders(uploaded_at desc);

create index if not exists idx_withdrawals_user_id on
withdrawals(user_id);

create index if not exists idx_withdrawals_processed_at on
withdrawals(processed_at desc);

create index if not exists idx_balances_user_id on
balance(user_id);

COMMIT;