CREATE TABLE `t_product_category` (
    `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `parent_id` INT(11)  NOT NULL DEFAULT '0' COMMENT '父ID',
    `category_name` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '类目名称',
    `category_name_en` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '类名名称(英文)',
    `image` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '图片',
    `detail` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '详情',
    `detail_en` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '详情(英文)',
    `is_deleted` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY (`parent_id`)
) ENGINE=InnoDB CHARSET=utf8mb4 COMMENT='产品类目表';
