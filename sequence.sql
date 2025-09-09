
CREATE TABLE `sequence` (
	`id` bigint(20) unsigned not null auto_increment ,
	`stub` varchar(1) not null,
    `timestamp` timestamp not null default current_timestamp ,
	primary key (`id`),
    unique key `idx_uniq_stub` (`stub`)
) engine=MyISAM default charset = utf8 comment = '序号表';

