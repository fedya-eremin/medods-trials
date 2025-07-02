create extension if not exists "uuid-ossp";
create schema auth;

create table auth.user (
    id UUID primary key
);

create table auth.refresh_token (
    id UUID primary key,
    user_id UUID references auth.user(id),
    token_hash text,
    used boolean,
    user_agent text,
    ip_addr text,
    expires_at timestamp
);
