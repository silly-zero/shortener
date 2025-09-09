
CREATE TABLE `short_url_map` (
                                 `id` bigint(20) unsigned not null auto_increment comment '主键',
                                 `create_at` datetime not null default current_timestamp comment '创建时间',
                                 `create_by` varchar(64) not null default '' comment '创建者',
                                 `is_del` tinyint unsigned not null default '0' comment '是否删除:0正常1删除',

                                 `lurl` varchar(2048) default null comment '长链接',
                                 `md5` char(32) default null comment '长链接md5',
                                 `surl` varchar(11) default null comment '短链接',
                                 primary key (`id`),
                                 index(`is_del`),
                                 unique(`md5`),
                                 unique(`surl`)
)engine=InnoDB default charset=utf8 comment='长短链映射表';