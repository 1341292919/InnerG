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

INSERT INTO InnerG.user (id, account, username, email, avatar, gender, role_type, password_hash, status, created_at, updated_at, deleted_at)
VALUES (
           108000,
           'lbl102300218',
           '林柏林爱喝酒',
           '15759762783@163.com',
           'https://ts3.tc.mm.bing.net/th/id/OIP-C.MBqJCFa5AMwdgTr5_91G4gHaE7?rs=1&pid=ImgDetMain&o=7&rm=3',
           1,
           1,  -- role_type: 1-普通用户
           '$2a$12$sHsUsfhcJPri5.LO0wgovuG4KFozx85tc7CQivRPGCDTz.TUplc7q',
           0,
           '2026-03-18 15:02:46',
           '2026-03-18 23:03:31',
           NULL
       );
-- 插入管理员账号
INSERT INTO InnerG.user (id, account, username, email, avatar, gender, role_type, password_hash, status, created_at, updated_at, deleted_at)
VALUES (
           108001,
           'admin_zhang',
           '张知行',  -- 管理员姓名，寓意"知行合一"
           'zhangzhixing@qq.com',
           'https://ts3.tc.mm.bing.net/th/id/OIP-C.MBqJCFa5AMwdgTr5_91G4gHaE7?rs=1&pid=ImgDetMain&o=7&rm=3',  -- 头像不变
           1,  -- 性别：男
           0,  -- role_type: 0-管理员
           '$2a$12$sHsUsfhcJPri5.LO0wgovuG4KFozx85tc7CQivRPGCDTz.TUplc7q',  -- 密码不变
           1,  -- status: 1-正常
           '2026-03-20 10:00:00',  -- 创建时间
           '2026-03-20 10:00:00',  -- 更新时间
           NULL
       );

CREATE TABLE InnerG.songs
(
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(128) NOT NULL COMMENT '歌曲名称',
    description TEXT NULL COMMENT '歌曲描述',
    cover_url VARCHAR(512) NULL COMMENT '封面图片URL',
    status TINYINT DEFAULT 1 NOT NULL COMMENT '状态：0-下架，1-上架，2-审核中',
    singer_id BIGINT UNSIGNED NOT NULL COMMENT '歌手ID，关联user表',
    source_url VARCHAR(512) NOT NULL COMMENT '歌曲文件URL',
    duration INT DEFAULT 0 NOT NULL COMMENT '歌曲时长（秒）',
    play_count BIGINT DEFAULT 0 NOT NULL COMMENT '播放次数',
    tags VARCHAR(255) NULL COMMENT '标签，多个标签用逗号分隔',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    INDEX idx_singer_id (singer_id),
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

CREATE TABLE InnerG.singer
(
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(128) NOT NULL COMMENT '歌手名称',  -- 注释写的是"专辑名称"，但表是歌手表，建议改为"歌手名称"
    description TEXT NULL COMMENT '描述',
    avatar_url VARCHAR(512) NULL COMMENT '歌手头像URL',  -- 建议改为"歌手头像URL"
    status TINYINT DEFAULT 1 NOT NULL COMMENT '状态',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
) AUTO_INCREMENT = 400000 COMMENT '歌手表';