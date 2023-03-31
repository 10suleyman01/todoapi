create table todo (
    id uuid primary key default gen_random_uuid(),
    title varchar(100) not null,
    user_id uuid NOT NULL,
    constraint user_fk foreign key (user_id) references public.users(id)
)