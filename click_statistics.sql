CREATE TABLE `click_statistics` (
                                    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                                    `surl` varchar(11) NOT NULL COMMENT '短链接标识',
                                    `click_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '点击时间',
                                    `ip` varchar(50) DEFAULT NULL COMMENT '访问IP',
                                    `user_agent` text COMMENT '用户代理',
                                    `referer` varchar(255) DEFAULT NULL COMMENT '来源页面',
                                    PRIMARY KEY (`id`),
                                    KEY `idx_surl` (`surl`),
                                    KEY `idx_click_time` (`click_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建视图用于数据统计分析
CREATE VIEW `daily_click_summary` AS
SELECT
    surl,
        DATE(click_time) as click_date,
        COUNT(*) as total_clicks,
        COUNT(DISTINCT ip) as unique_visitors
        FROM click_statistics
        GROUP BY surl, DATE(click_time);