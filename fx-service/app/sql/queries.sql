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

-- name: select-rate
select c.code_eng, date_time, value
from fx_rates r
join currencies c on c.code_cbr = r.code_cbr
where c.code_eng = $1
      and r.date_time in (
        select max(date_time) as last_date
        from fx_rates r
        join currencies c on c.code_cbr = r.code_cbr
        where c.code_eng = $1
        and date_time <= $2
      );
