create table currencies (
    code_cbr varchar(10) primary key,
    code_eng varchar(3) not null,
    name_rus varchar(128),
    name_eng varchar(128)
);

create table fx_rates (
    code_cbr varchar references currencies(code_cbr),
    date_time timestamp not null,
    value real not null
);
create unique index fx_rate_unique_idx on fx_rates (code_cbr, date_time);

insert into currencies values
('R01010', 'AUD', 'Австралийский доллар', 'Australian Dollar'),
('R01020A', 'AZN', 'Азербайджанский манат', 'Azerbaijan Manat'),
('R01035', 'GBP', 'Фунт стерлингов Соединенного королевства', 'British Pound Sterling'),
('R01060', 'AMD', 'Армянский драм', 'Armenia Dram'),
('R01090B', 'BYN', 'Белорусский рубль', 'Belarussian Ruble'),
('R01100', 'BGN', 'Болгарский лев', 'Bulgarian lev'),
('R01115', 'BRL', 'Бразильский реал', 'Brazil Real'),
('R01135', 'HUF', 'Венгерский форинт', 'Hungarian Forint'),
('R01200', 'HKD', 'Гонконгский доллар', 'Hong Kong Dollar'),
('R01215', 'DKK', 'Датская крона', 'Danish Krone'),
('R01235', 'USD', 'Доллар США', 'US Dollar'),
('R01239', 'EUR', 'Евро', 'Euro'),
('R01270', 'INR', 'Индийская рупия', 'Indian Rupee'),
('R01335', 'KZT', 'Казахстанский тенге', 'Kazakhstan Tenge'),
('R01350', 'CAD', 'Канадский доллар', 'Canadian Dollar'),
('R01370', 'KGS', 'Киргизский сом', 'Kyrgyzstan Som'),
('R01375', 'CNY', 'Китайский юань', 'China Yuan'),
('R01500', 'MDL', 'Молдавский лей', 'Moldova Lei'),
('R01535', 'NOK', 'Норвежская крона', 'Norwegian Krone'),
('R01565', 'PLN', 'Польский злотый', 'Polish Zloty'),
('R01585F', 'RON', 'Румынский лей', 'Romanian Leu'),
('R01589', 'XDR', 'СДР (специальные права заимствования)', 'SDR'),
('R01625', 'SGD', 'Сингапурский доллар', 'Singapore Dollar'),
('R01670', 'TJS', 'Таджикский сомони', 'Tajikistan Ruble'),
('R01700J', 'TRY', 'Турецкая лира', 'Turkish Lira'),
('R01710A', 'TMT', 'Новый туркменский манат', 'New Turkmenistan Manat'),
('R01717', 'UZS', 'Узбекский сум', 'Uzbekistan Sum'),
('R01720', 'UAH', 'Украинская гривна', 'Ukrainian Hryvnia'),
('R01760', 'CZK', 'Чешская крона', 'Czech Koruna'),
('R01770', 'SEK', 'Шведская крона', 'Swedish Krona'),
('R01775', 'CHF', 'Швейцарский франк', 'Swiss Franc'),
('R01810', 'ZAR', 'Южноафриканский рэнд', 'S.African Rand'),
('R01815', 'KRW', 'Вон Республики Корея', 'South Korean Won'),
('R01820', 'JPY', 'Японская иена', 'Japanese Yen');