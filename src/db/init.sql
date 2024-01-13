create table if not exists users
(
    id				serial primary key, 
	email			varchar(100) not null unique,
	password 		bytea not null,
	profile_name 	varchar(100) not null,
	balance 		int not null default 0,
	token_version 	bigint not null default 0
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
    id			int not null,
    bookId		int not null,
	name    	varchar(100) not null,
	price		int not null,
	content 	text not null,
	status		int not null default 0,
   	foreign key (bookId) REFERENCES books(id),
	primary key (id, bookId)
);

create table if not exists transactions
(
	id					serial primary key,
	bookid				int not null,
	chapterid    		int not null,
	receivinguserid		int not null,
	payinguserid 		int not null,
	amount 				int not null,
	foreign key (chapterid, bookid) references chapters(id, bookId),
	foreign key (bookid) references books (id),
	foreign key (payinguserid) references users(id),
	foreign key (receivinguserid) references users(id)
);

insert into users (email, password, profile_name, balance)
values ('test@test.com', '<bcrypt password hash 1>', 'Toni Tester', 1000),
       ('test', '<bcrypt password hash 1>', 'Test User', 1000);

insert into books (name, authorId, description)
values ('Book One', 1, 'A good book'),
       ('Book Two', 2, 'A bad book'),
       ('Book Three', 1, 'A mid book');

insert into chapters (id, bookId, name, price, content, status)
values (1, 1, 'The beginning', 0, 'Lorem Ipsum', 1),
       (2, 1,'The beginning 2: Electric Boogaloo', 100, 'Lorem Ipsum 2', 1),
       (3, 1, 'The beginning 3: My Enemy', 100, 'Lorem Ipsum 3', 1),
       (1, 2, 'A different book chapter 1', 0, 'Lorem Ipsum 4', 1),
       (2, 2, 'What came after', 100, 'Lorem Ipsum 5', 1),
	   (3, 2, 'What came after that', 500, 'Lorem Ipsum 6', 1),
	   (4, 2, 'And there after ', 400, 'Lorem Ipsum 7', 1),
	   (1, 3, 'The third book chapter 1', 750, 'Lorem Ipsum 8', 1),
	   (2, 3, 'The third book chapter 2', 800, 'Lorem Ipsum 9', 1),
	   (3, 3, 'The third book chapter 3', 900, 'Lorem Ipsum 10', 1);

insert into transactions (bookid, chapterid, receivinguserid, payinguserid, amount)
values (1, 1, 1, 2, 0),
       (1, 2, 1, 2, 100),
       (2, 1, 2, 1, 0),
	   (2, 4, 2, 1, 400),
	   (3, 1, 1, 2, 750),
	   (3, 2, 1, 2, 800);