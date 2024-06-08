CREATE DATABASE IF NOT EXISTS `employee`;

/* TODO: position can be a seperate table with a maping to this table */
CREATE TABLE IF NOT EXISTS employee (
    id tinyint unsigned auto_increment primary key,
    name varchar(56) not null,
    position varchar(56) default '',
    created_on datetime not null default current_timestamp,
    updated_on datetime not null default current_timestamp,
    deleted_on datetime
) engine = innodb;
insert into employee (name, position) values
  ('Sachin Prasad','Senior Software Engineer'),
  ('Nadim Ayaz','Software Engineeri'),
  ('Ritesh Banerjee','Staff Engineer')
;
