CREATE TABLE review_reply_info (
    `id` bigint(32) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `create_by` varchar(48) NOT NULL DEFAULT '' COMMENT '创建⽅标识',
    `update_by` varchar(48) NOT NULL DEFAULT '' COMMENT '更新⽅标识',
    `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
    `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_at` timestamp COMMENT '逻辑删除标记',
    `version` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '乐观锁标记',
    `reply_id` bigint(32) NOT NULL DEFAULT '0' COMMENT '回复id',
    `review_id` bigint(32) NOT NULL DEFAULT '0' COMMENT '评价id',
    `store_id` bigint(32) NOT NULL DEFAULT '0' COMMENT '店铺id',
    `content` varchar(512) NOT NULL COMMENT '评价内容',
    `pic_info` varchar(1024) NOT NULL DEFAULT '' COMMENT '媒体信息：图⽚',
    `video_info` varchar(1024) NOT NULL DEFAULT '' COMMENT '媒体信息：视频',
    `ext_json` varchar(1024) NOT NULL DEFAULT '' COMMENT '信息扩展',
    `ctrl_json` varchar(1024) NOT NULL DEFAULT '' COMMENT '控制扩展'
    PRIMARY KEY (`id`),
    KEY `idx_delete_at` (`delete_at`) COMMENT '逻辑删除索引',
    UNIQUE KEY `uk_reply_id` (`reply_id`) COMMENT '回复id索引',
    KEY `idx_review_id` (`review_id`) COMMENT '评价id索引',
    KEY `idx_store_id` (`store_id`) COMMENT '店铺id索引'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT= '评价商家回复表';