create table if not exists users
(
    id				serial primary key, 
	email			varchar(100) not null unique,
	password 		bytea not null,
	profile_name 	varchar(100) not null,
	balance 		int not null default 0
);

create table if not exists books
(
    id			serial primary key,
	name    	varchar(100) not null,
	authorId	int not null,
	description text not null,
   	foreign key (authorId) REFERENCES users(id)
);

create table if not exists chapters
(
    id			serial primary key,
    bookId		int not null,
	name    	varchar(100) not null,
	price		int not null,
	content 	text not null,
   	foreign key (bookId) REFERENCES books(id)
);

create table if not exists transactions
(
	id					serial primary key,
	bookid				int not null,
	chapterid    		int not null,
	receivinguserid		int not null,
	payinguserid 		int not null,
	amount 				int not null,
	foreign key (chapterid) references chapters(id),
	foreign key (bookid) references books (id),
	foreign key (payinguserid) references users(id),
	foreign key (receivinguserid) references users(id)
);

insert into users (email, password, profile_name, balance)
values ('test@test.com', '<bcrypt password hash 1>', 'Toni Tester', 1000),
       ('max.mustermann@gmail.com', '<bcrypt password hash 2>', 'Max Mustermann', 1000);

insert into books (name, authorId, description)
values ('Book One', 1, 'A good book'),
       ('Book Two', 2, 'A bad book'),
       ('Book Three', 1, 'A mid book');

insert into chapters (bookId, name, price, content)
values (1, 'The beginning', 0, 'Lorem Ipsum'),
       (1, 'The beginning 2: Electric Boogaloo', 100, 'Lorem Ipsum 2'),
       (1, 'The beginning 3: My Enemy', 100, 'Lorem Ipsum 3'),
       (2, 'A different book chapter 1', 0, 'LorIp 4'),
       (2, 'What came after', 100, 'Lorem Ipsum 5');

insert into transactions (bookid, chapterid, receivinguserid, payinguserid, amount)
values (1, 1, 1, 2, 0),
       (1, 2, 1, 2, 100),
       (2, 4, 2, 1, 0);