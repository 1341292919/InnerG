-- ============================================
-- InnerG 音乐平台测试数据
-- 生成时间: 2026-03-21
-- ============================================

-- 1. 插入歌手数据
INSERT INTO InnerG.singer (id, name, description, avatar_url, status, created_at, updated_at, deleted_at) VALUES
                                                                                                              (400001, '周杰伦', '华语乐坛天王，2000年出道，开创了华语流行音乐的新时代。代表作《七里香》《青花瓷》《稻香》等。', 'https://p3.hippopx.com/preview/622/628/bengal-tiger-close-up-tiger-eyes-fur-pattern-wildlife-photography-animal-portrait-tiger-face-striking-eyes-big-cat-nature-photography-thumbnail.jpg', 1, NOW(), NOW(), NULL),
                                                                                                              (400002, '陈奕迅', '香港著名歌手，被誉为"歌神"张学友的接班人。拥有独特的嗓音和极强的感染力。', 'https://p3.hippopx.com/preview/386/205/tiger-majestic-tiger-intense-gaze-jungle-natural-setting-fierce-tiger-powerful-gaze-wildlife-animal-nature-thumbnail.jpg', 1, NOW(), NOW(), NULL),
                                                                                                              (400003, '邓紫棋', '创作型女歌手，拥有独特的嗓音和创作才华。2008年出道，迅速成为华语乐坛新生代代表。', 'https://p3.hippopx.com/preview/585/12/adult-tiger-majestic-tiger-close-up-tiger-striped-tiger-fierce-tiger-regal-tiger-tiger-face-tiger-portrait-tiger-stripes-tiger-eyes-thumbnail.jpg', 1, NOW(), NOW(), NULL),
                                                                                                              (400004, '林俊杰', '新加坡创作型歌手，被誉为"行走的CD"。拥有完美的唱功和出色的创作能力。', 'https://p0.hippopx.com/preview/893/689/119/tiger-lily-tiger-lily-flower-tiger-lily-plant-asiatic-lily-royalty-free-thumbnail.jpg', 1, NOW(), NOW(), NULL),
                                                                                                              (400005, '王菲', '华语乐坛天后，以其空灵的嗓音和独特的风格著称。', 'https://p0.hippopx.com/preview/865/211/815/pets-cute-cat-dog-royalty-free-thumbnail.jpg', 1, NOW(), NOW(), NULL),
                                                                                                              (400006, '五月天', '中国台湾摇滚乐团，由阿信、怪兽、石头、玛莎、冠佑组成。', 'https://p3.hippopx.com/preview/887/452/yawning-cat-close-up-cat-cat-fangs-cat-whiskers-cute-cat-amusing-cat-cat-mouth-open-cat-tongue-cat-expression-cat-face-thumbnail.jpg', 1, NOW(), NOW(), NULL);

-- 2. 插入歌曲数据
INSERT INTO InnerG.songs (id, name, description, cover_url, status, singer_id, source_url, duration, play_count, tags, created_at, updated_at, deleted_at) VALUES
                                                                                                                                                               (100001, '七里香', '周杰伦经典作品，充满青春回忆', 'https://p3.hippopx.com/preview/622/628/bengal-tiger-close-up-tiger-eyes-fur-pattern-wildlife-photography-animal-portrait-tiger-face-striking-eyes-big-cat-nature-photography-thumbnail.jpg', 1, 400001, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 280, 1250000, '周杰伦,经典,青春', NOW(), NOW(), NULL),
                                                                                                                                                               (100002, '青花瓷', '中国风代表作，意境优美', 'https://p3.hippopx.com/preview/386/205/tiger-majestic-tiger-intense-gaze-jungle-natural-setting-fierce-tiger-powerful-gaze-wildlife-animal-nature-thumbnail.jpg', 1, 400001, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 245, 980000, '中国风,周杰伦,经典', NOW(), NOW(), NULL),
                                                                                                                                                               (100003, '稻香', '励志歌曲，温暖人心', 'https://p3.hippopx.com/preview/585/12/adult-tiger-majestic-tiger-close-up-tiger-striped-tiger-fierce-tiger-regal-tiger-tiger-face-tiger-portrait-tiger-stripes-tiger-eyes-thumbnail.jpg', 1, 400001, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 223, 1100000, '励志,温暖,周杰伦', NOW(), NOW(), NULL),
                                                                                                                                                               (100004, '富士山下', '陈奕迅经典粤语歌曲', 'https://p0.hippopx.com/preview/893/689/119/tiger-lily-tiger-lily-flower-tiger-lily-plant-asiatic-lily-royalty-free-thumbnail.jpg', 1, 400002, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 260, 890000, '粤语,陈奕迅,经典', NOW(), NOW(), NULL),
                                                                                                                                                               (100005, '十年', 'KTV必点歌曲，感人至深', 'https://p0.hippopx.com/preview/865/211/815/pets-cute-cat-dog-royalty-free-thumbnail.jpg', 1, 400002, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 198, 1520000, '陈奕迅,情歌,经典', NOW(), NOW(), NULL),
                                                                                                                                                               (100006, '浮夸', '展现唱功的经典作品', 'https://p3.hippopx.com/preview/887/452/yawning-cat-close-up-cat-cat-fangs-cat-whiskers-cute-cat-amusing-cat-cat-mouth-open-cat-tongue-cat-expression-cat-face-thumbnail.jpg', 1, 400002, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 285, 760000, '陈奕迅,高音,经典', NOW(), NOW(), NULL),
                                                                                                                                                               (100007, '光年之外', '邓紫棋代表作，旋律优美', 'https://p3.hippopx.com/preview/879/82/ginger-cat-tabby-cat-lounging-cat-scratching-post-indoor-cat-close-up-cat-looking-at-camera-ginger-tabby-cat-on-scratching-post-cat-indoors-thumbnail.jpg', 1, 400003, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 235, 1350000, '邓紫棋,流行,太空', NOW(), NOW(), NULL),
                                                                                                                                                               (100008, '泡沫', '邓紫棋成名曲', 'https://p3.hippopx.com/preview/728/236/orange-tabby-cat-hiding-face-cute-cat-sleeping-cat-paws-covering-face-cozy-scene-endearing-scene-tabby-cat-orange-cat-cat-sleeping-thumbnail.jpg', 1, 400003, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 260, 1580000, '邓紫棋,情歌,成名曲', NOW(), NOW(), NULL),
                                                                                                                                                               (100009, '喜欢你', '邓紫棋改编版，充满激情', 'https://p3.hippopx.com/preview/964/209/cat-tabby-cat-green-eyes-gazing-upward-animal-lovers-selective-focus-photography-charming-cat-upward-gaze-tabby-cat-eyes-cat-photography-thumbnail.jpg', 1, 400003, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 215, 920000, '邓紫棋,改编,粤语', NOW(), NOW(), NULL),
                                                                                                                                                               (100010, '江南', '林俊杰成名曲，中国风', 'https://p3.hippopx.com/preview/448/205/tiger-nature-animal-wildlife-black-and-white-wild-animal-wild-cat-feline-thumbnail.jpg', 1, 400004, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 245, 1450000, '林俊杰,中国风,经典', NOW(), NOW(), NULL),
                                                                                                                                                               (100011, '不为谁而作的歌', '林俊杰创作巅峰', 'https://tse2-mm.cn.bing.net/th/id/OIP-C.UHSGBwk9A8u4OFAuU7KUoAHaE8?w=236&h=180&c=7&r=0&o=7&pid=1.7&rm=3', 1, 400004, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 278, 860000, '林俊杰,创作,励志', NOW(), NOW(), NULL),
                                                                                                                                                               (100012, '修炼爱情', '林俊杰情歌代表作', 'https://tse3-mm.cn.bing.net/th/id/OIP-C.0raQShBJRN3df0_NP8revgHaE7?w=274&h=183&c=7&r=0&o=7&pid=1.7&rm=3', 1, 400004, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 260, 1050000, '林俊杰,情歌,修炼', NOW(), NOW(), NULL),
                                                                                                                                                               (100013, '传奇', '王菲经典作品，空灵嗓音', 'https://p3.hippopx.com/preview/622/628/bengal-tiger-close-up-tiger-eyes-fur-pattern-wildlife-photography-animal-portrait-tiger-face-striking-eyes-big-cat-nature-photography-thumbnail.jpg', 1, 400005, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 305, 780000, '王菲,空灵,传奇', NOW(), NOW(), NULL),
                                                                                                                                                               (100014, '红豆', '王菲经典情歌', 'https://p3.hippopx.com/preview/386/205/tiger-majestic-tiger-intense-gaze-jungle-natural-setting-fierce-tiger-powerful-gaze-wildlife-animal-nature-thumbnail.jpg', 1, 400005, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 250, 1120000, '王菲,情歌,红豆', NOW(), NOW(), NULL),
                                                                                                                                                               (100015, '倔强', '五月天励志歌曲', 'https://p3.hippopx.com/preview/585/12/adult-tiger-majestic-tiger-close-up-tiger-striped-tiger-fierce-tiger-regal-tiger-tiger-face-tiger-portrait-tiger-stripes-tiger-eyes-thumbnail.jpg', 1, 400006, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 280, 960000, '五月天,励志,摇滚', NOW(), NOW(), NULL),
                                                                                                                                                               (100016, '突然好想你', '五月天经典情歌', 'https://p0.hippopx.com/preview/893/689/119/tiger-lily-tiger-lily-flower-tiger-lily-plant-asiatic-lily-royalty-free-thumbnail.jpg', 1, 400006, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 265, 1280000, '五月天,情歌,思念', NOW(), NOW(), NULL),
                                                                                                                                                               (100017, '恋爱ing', '五月天欢快歌曲', 'https://p0.hippopx.com/preview/865/211/815/pets-cute-cat-dog-royalty-free-thumbnail.jpg', 1, 400006, 'http://tc6io7vt3.hn-bkt.clouddn.com/innerG/108001/1000.flac', 210, 890000, '五月天,欢快,恋爱', NOW(), NOW(), NULL);

-- 3. 插入歌单数据（5个歌单）
INSERT INTO InnerG.playlist (id, name, description, cover_url, status, play_count, song_count, tags, created_at, updated_at, deleted_at) VALUES
                                                                                                                                             (200001, '华语经典金曲', '收录华语乐坛最经典的歌曲，每一首都是回忆杀', 'https://p3.hippopx.com/preview/622/628/bengal-tiger-close-up-tiger-eyes-fur-pattern-wildlife-photography-animal-portrait-tiger-face-striking-eyes-big-cat-nature-photography-thumbnail.jpg', 1, 258000, 15, '经典,华语,怀旧', NOW(), NOW(), NULL),
                                                                                                                                             (200002, '治愈系音乐', '温暖心灵的歌曲，给你满满正能量', 'https://p3.hippopx.com/preview/386/205/tiger-majestic-tiger-intense-gaze-jungle-natural-setting-fierce-tiger-powerful-gaze-wildlife-animal-nature-thumbnail.jpg', 1, 189000, 15, '治愈,温暖,放松', NOW(), NOW(), NULL),
                                                                                                                                             (200003, 'KTV必点热歌', '聚会K歌必备歌曲，嗨翻全场', 'https://p3.hippopx.com/preview/585/12/adult-tiger-majestic-tiger-close-up-tiger-striped-tiger-fierce-tiger-regal-tiger-tiger-face-tiger-portrait-tiger-stripes-tiger-eyes-thumbnail.jpg', 1, 325000, 15, 'KTV,聚会,热门', NOW(), NOW(), NULL),
                                                                                                                                             (200004, '影视OST精选', '经典影视剧原声音乐，重温感动瞬间', 'https://p0.hippopx.com/preview/893/689/119/tiger-lily-tiger-lily-flower-tiger-lily-plant-asiatic-lily-royalty-free-thumbnail.jpg', 1, 156000, 15, 'OST,影视,原声', NOW(), NOW(), NULL),
                                                                                                                                             (200005, '深夜单曲循环', '适合深夜聆听的歌曲，陪伴你的孤独时光', 'https://p0.hippopx.com/preview/865/211/815/pets-cute-cat-dog-royalty-free-thumbnail.jpg', 1, 298000, 15, '深夜,安静,抒情', NOW(), NOW(), NULL);

-- 4. 插入歌单歌曲关联数据（每个歌单15首歌曲）
-- 歌单1: 华语经典金曲 (歌曲ID: 100001-100015)
INSERT INTO InnerG.playlist_songs (playlist_id, song_id, sort_order, added_at, created_at, updated_at, deleted_at) VALUES
                                                                                                                       (200001, 100001, 1, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100002, 2, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100003, 3, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100004, 4, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100005, 5, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100006, 6, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100007, 7, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100008, 8, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100009, 9, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100010, 10, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100011, 11, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100012, 12, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100013, 13, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100014, 14, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200001, 100015, 15, NOW(), NOW(), NOW(), NULL);

-- 歌单2: 治愈系音乐 (歌曲ID: 100003,100005,100007,100008,100009,100010,100011,100012,100013,100014,100015,100016,100017, 再加上两个)
INSERT INTO InnerG.playlist_songs (playlist_id, song_id, sort_order, added_at, created_at, updated_at, deleted_at) VALUES
                                                                                                                       (200002, 100003, 1, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100005, 2, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100007, 3, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100008, 4, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100009, 5, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100010, 6, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100011, 7, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100012, 8, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100013, 9, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100014, 10, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100015, 11, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100016, 12, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100017, 13, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100001, 14, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200002, 100002, 15, NOW(), NOW(), NOW(), NULL);

-- 歌单3: KTV必点热歌 (歌曲ID: 100001,100005,100006,100008,100010,100012,100015,100016,100017,100007,100009,100011,100013,100014,100004)
INSERT INTO InnerG.playlist_songs (playlist_id, song_id, sort_order, added_at, created_at, updated_at, deleted_at) VALUES
                                                                                                                       (200003, 100001, 1, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100005, 2, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100006, 3, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100008, 4, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100010, 5, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100012, 6, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100015, 7, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100016, 8, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100017, 9, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100007, 10, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100009, 11, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100011, 12, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100013, 13, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100014, 14, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200003, 100004, 15, NOW(), NOW(), NOW(), NULL);

-- 歌单4: 影视OST精选 (歌曲ID: 100001,100002,100004,100005,100007,100008,100010,100012,100013,100014,100015,100016,100003,100006,100009)
INSERT INTO InnerG.playlist_songs (playlist_id, song_id, sort_order, added_at, created_at, updated_at, deleted_at) VALUES
                                                                                                                       (200004, 100001, 1, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100002, 2, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100004, 3, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100005, 4, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100007, 5, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100008, 6, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100010, 7, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100012, 8, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100013, 9, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100014, 10, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100015, 11, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100016, 12, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100003, 13, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100006, 14, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200004, 100009, 15, NOW(), NOW(), NOW(), NULL);

-- 歌单5: 深夜单曲循环 (歌曲ID: 100002,100004,100005,100006,100008,100010,100011,100012,100013,100014,100015,100016,100001,100003,100007)
INSERT INTO InnerG.playlist_songs (playlist_id, song_id, sort_order, added_at, created_at, updated_at, deleted_at) VALUES
                                                                                                                       (200005, 100002, 1, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100004, 2, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100005, 3, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100006, 4, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100008, 5, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100010, 6, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100011, 7, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100012, 8, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100013, 9, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100014, 10, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100015, 11, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100016, 12, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100001, 13, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100003, 14, NOW(), NOW(), NOW(), NULL),
                                                                                                                       (200005, 100007, 15, NOW(), NOW(), NOW(), NULL);

-- 验证数据插入
SELECT '歌手数量:' AS '统计项', COUNT(*) AS '数量' FROM InnerG.singer
UNION ALL
SELECT '歌曲数量:', COUNT(*) FROM InnerG.songs
UNION ALL
SELECT '歌单数量:', COUNT(*) FROM InnerG.playlist
UNION ALL
SELECT '歌单歌曲关联数量:', COUNT(*) FROM InnerG.playlist_songs;