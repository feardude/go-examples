-- name: insert-fx_rate
insert into fx_rates values ($1, $2, $3)
on conflict (code_cbr, date_time)
do update set value = $3
    where fx_rates.code_cbr = $1
        and fx_rates.date_time = $2;

-- name: select-last-date
select max(date_time)
from fx_rates
where code_cbr = $1;

-- name: select-currencies
select code_cbr, code_eng, name_rus, name_eng
from currencies;