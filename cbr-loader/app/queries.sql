-- name: create-table-currencies
create table if not exists currencies (
    code_cbr varchar(10) primary key,
    code_eng varchar(3) not null,
    name_rus varchar(128),
    name_eng varchar(128)
);

-- name: create-table-fx_rates
create table if not exists fx_rates (
    code_cbr varchar references currencies(code_cbr),
    date_time timestamp not null,
    value real not null
);
create unique index if not exists fx_rate_unique_idx on fx_rates (code_cbr, date_time);

-- name: insert-currency
insert into currencies values ($1, $2, $3, $4)
on conflict do nothing;

-- name: insert-fx_rate
insert into fx_rates values ($1, $2, $3)
on conflict do nothing;

-- name: select-last-date
select max(date_time)
from fx_rates
where code_cbr = $1;