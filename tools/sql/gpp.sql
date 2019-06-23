# 创建数据库
create database gpp;

# 创建doc基本信息表
create table t_doc_info (
	id bigint auto_increment,
	docid varchar(128) not null,
	title varchar(256) not null, 
	src varchar(32) not null,
	type varchar(32),
	publish_time timestamp default '1970-01-01 00:00:00',
	create_time timestamp not null default current_timestamp,
	update_time timestamp not null default current_timestamp on update current_timestamp,
	primary key(id),
	constraint uqdocid unique(src, docid);
) default charset = utf8;

# 创建doc动态信息表
create table t_doc_info (
	id bigint auto_increment,
	docid varchar(128) not null,
	src varchar(32) not null,
	views int default 0,
	loves int default 0,
	publish_time timestamp default '1970-01-01 00:00:00',
	create_time timestamp not null default current_timestamp,
	update_time timestamp not null default current_timestamp on update current_timestamp,
	primary key(id),
	constraint uqdocid unique(src, docid);
);
