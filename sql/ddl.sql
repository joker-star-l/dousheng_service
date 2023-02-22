create table if not exists public."user"
(
    id               bigint not null
    constraint user_pk
    primary key,
    name             text,
    password         text,
    created_at       timestamp(0),
    updated_at       timestamp(0),
    deleted_at       timestamp(0),
    avatar           text,
    background_image text,
    signature        text
    );

alter table public."user"
    owner to postgres;

create index if not exists user_name_index
    on public."user" (name);

create table if not exists public.user_follow
(
    id         bigint not null
    constraint user_follow_pk
    primary key,
    user_from  bigint not null,
    user_to    bigint not null,
    created_at timestamp(0),
    updated_at timestamp(0),
    deleted_at timestamp(0)
    );

alter table public.user_follow
    owner to postgres;

create index if not exists user_follow_user_from_index
    on public.user_follow (user_from);

create index if not exists user_follow_user_to_index
    on public.user_follow (user_to);

create table if not exists public.user_friend
(
    id         bigint not null
    constraint user_friend_pk
    primary key,
    user0      bigint not null,
    user1      bigint not null,
    created_at timestamp(0),
    updated_at timestamp(0),
    deleted_at timestamp(0)
    );

alter table public.user_friend
    owner to postgres;

create index if not exists user_friend_user0_index
    on public.user_friend (user0);

create index if not exists user_friend_user1_index
    on public.user_friend (user1);

create table if not exists public.message
(
    id         bigint not null
    constraint message_pk
    primary key,
    user_from  bigint not null,
    user_to    bigint not null,
    content    text,
    created_at timestamp(0),
    updated_at timestamp(0),
    deleted_at timestamp(0)
    );

alter table public.message
    owner to postgres;

create index if not exists message_user_from_index
    on public.message (user_from);

create index if not exists message_user_to_index
    on public.message (user_to);

create table if not exists public.video
(
    id         bigint not null
    constraint video_pk
    primary key,
    title      text,
    play_url   text,
    cover_url  text,
    created_at timestamp(0),
    updated_at timestamp(0),
    deleted_at timestamp(0),
    user_id    bigint
    );

alter table public.video
    owner to postgres;

create index if not exists video_user_id_index
    on public.video (user_id);

create index if not exists video_created_at_index
    on public.video (created_at);

create table if not exists public.video_favorite
(
    id            bigint not null
    constraint video_favorite_pk
    primary key,
    user_id       bigint not null,
    video_id      bigint not null,
    created_at    timestamp(0),
    updated_at    timestamp(0),
    deleted_at    timestamp(0),
    video_user_id bigint not null
    );

alter table public.video_favorite
    owner to postgres;

create index if not exists video_favorite_user_id_index
    on public.video_favorite (user_id);

create index if not exists video_favorite_video_id_index
    on public.video_favorite (video_id);

create index if not exists video_favorite_video_user_id_index
    on public.video_favorite (video_user_id);

create table if not exists public.video_comment
(
    id            bigint not null
    constraint video_comment_pk
    primary key,
    user_id       bigint not null,
    video_id      bigint not null,
    comment       text   not null,
    created_at    timestamp(0),
    updated_at    timestamp(0),
    deleted_at    timestamp(0),
    video_user_id bigint not null
    );

alter table public.video_comment
    owner to postgres;

create index if not exists video_comment_user_id_index
    on public.video_comment (user_id);

create index if not exists video_comment_video_id_index
    on public.video_comment (video_id);

create index if not exists video_comment_video_user_id_index
    on public.video_comment (video_user_id);

create index if not exists video_comment_created_at_index
    on public.video_comment (created_at);

