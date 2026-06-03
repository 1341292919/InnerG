CREATE TABLE InnerG.user
(
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    account VARCHAR(32) NULL,
    username VARCHAR(64) NOT NULL,
    email VARCHAR(128) NULL,
    avatar VARCHAR(512) NULL,
    status TINYINT DEFAULT 1 NOT NULL,
    role_type TINYINT NOT NULL DEFAULT 1 COMMENT '角色类型：0-管理员，1-普通用户',
    password_hash VARCHAR(255) NOT NULL,
    gender TINYINT DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    UNIQUE KEY (account),
    UNIQUE KEY (email)
) AUTO_INCREMENT = 108000;

CREATE TABLE InnerG.songs
(
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(128) NOT NULL COMMENT '歌曲名称',
    description TEXT NULL COMMENT '歌曲描述',
    cover_url VARCHAR(512) NULL COMMENT '封面图片URL',
    status TINYINT DEFAULT 1 NOT NULL COMMENT '状态：0-下架，1-上架，2-审核中',
    singer_name VARCHAR(128) NOT NULL COMMENT '歌手名称',
    album VARCHAR(256) NULL COMMENT '专辑名称',
    source_url VARCHAR(512) NOT NULL COMMENT '歌曲文件URL',
    duration INT DEFAULT 0 NOT NULL COMMENT '歌曲时长（秒）',
    play_count BIGINT DEFAULT 0 NOT NULL COMMENT '播放次数',
    tags VARCHAR(255) NULL COMMENT '标签，多个标签用逗号分隔',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    INDEX idx_singer_name (singer_name),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) AUTO_INCREMENT = 100000 COMMENT '歌曲表';

CREATE TABLE InnerG.playlist
(
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(128) NOT NULL COMMENT '歌单名称',
    description TEXT NULL COMMENT '歌单描述',
    cover_url VARCHAR(512) NULL COMMENT '歌单封面URL',
    status TINYINT DEFAULT 1 NOT NULL COMMENT '状态：0-私密，1-公开，2-删除',
    play_count BIGINT DEFAULT 0 NOT NULL COMMENT '播放次数',
    song_count INT DEFAULT 0 NOT NULL COMMENT '歌曲数量',
    tags VARCHAR(255) NULL COMMENT '标签，多个标签用逗号分隔',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)  -- 移除了不存在的字段索引
) AUTO_INCREMENT = 200000 COMMENT '歌单表';

CREATE TABLE InnerG.playlist_songs
(
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    playlist_id BIGINT UNSIGNED NOT NULL COMMENT '歌单ID',
    song_id BIGINT UNSIGNED NOT NULL COMMENT '歌曲ID',
    sort_order INT DEFAULT 0 NOT NULL COMMENT '歌曲在歌单中的排序',
    added_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    UNIQUE KEY unique_playlist_song (playlist_id, song_id),
    INDEX idx_playlist_id (playlist_id),
    INDEX idx_song_id (song_id)  -- 移除了多余的逗号
) AUTO_INCREMENT = 300000 COMMENT '歌单歌曲关联表';


CREATE TABLE InnerG.messages (
                                  id BIGINT PRIMARY KEY AUTO_INCREMENT,
                                  msg_id VARCHAR(64) NOT NULL,              -- 客户端生成的消息ID
                                 from_user BIGINT NOT NULL,                -- 发送者
                                 to_user BIGINT NOT NULL,                  -- 接收者
                                 content TEXT NOT NULL,                    -- 消息内容
                                 type TINYINT NOT NULL DEFAULT 1,          -- 消息类型
                                 status TINYINT NOT NULL DEFAULT 0,        -- 0=已推送 1=未推送 2=已撤回 4=用户已删除
                                 created_at BIGINT NOT NULL,               -- 消息时间戳
                                 updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                 deleted_at DATETIME NULL,
                                 INDEX idx_from_to (from_user, to_user),
                                 INDEX idx_to_from (to_user, from_user),
                                  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE InnerG.friend (
                               id BIGINT PRIMARY KEY COMMENT '雪花ID',
                               user_id BIGINT NOT NULL COMMENT '用户ID',
                               friend_id BIGINT NOT NULL COMMENT '好友用户ID',
                               status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-正常，2-已删除',
                               created_at BIGINT NOT NULL COMMENT '创建时间戳',
                               updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                               deleted_at DATETIME NULL,
                               UNIQUE KEY idx_user_friend (user_id, friend_id),
                               INDEX idx_user_status (user_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '好友关系表';

CREATE TABLE InnerG.friend_request (
                                       id BIGINT PRIMARY KEY COMMENT '雪花ID',
                                       from_user BIGINT NOT NULL COMMENT '申请人ID',
                                       to_user BIGINT NOT NULL COMMENT '接收人ID',
                                       status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-待处理，2-已接受，3-已拒绝，4-已取消',
                                       message VARCHAR(100) NOT NULL DEFAULT '' COMMENT '好友申请打招呼内容',
                                       created_at BIGINT NOT NULL COMMENT '创建时间戳',
                                       updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                       deleted_at DATETIME NULL,
                                       UNIQUE KEY idx_from_to (from_user, to_user),
                                       INDEX idx_from_status (from_user, status),
                                       INDEX idx_to_status (to_user, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '好友申请表';
