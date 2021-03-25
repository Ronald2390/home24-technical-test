create table public."user"
(
	"id" serial not null,
	"name" varchar(100) not null,
	"email" varchar(100) not null,
	"address" varchar(255) not null,
	"password" varchar(255) null,
    "createdBy" int not null,
	"createdAt" timestamptz NOT NULL,
    "updatedBy" int not null,
	"updatedAt" timestamptz NOT NULL,
    "deletedBy" int,
	"deletedAt" timestamptz,
	constraint user_pkey primary key ("id")
);