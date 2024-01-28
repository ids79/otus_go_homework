create database calendar;
create user otus with encrypted password 'otus';
\c calendar
grant usage, create on schema public to otus;


