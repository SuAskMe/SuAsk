-- ----------------------------
-- Records of users
-- ----------------------------

DELETE FROM users;

INSERT INTO `users` VALUES (1, 'root', '1@suask.me', 'iqjbenzvfz', '6917ef8afa3ffeb3cb02643b9feb2a46', 'student', 'root', '# 这里是root', NULL, 0, '2024-12-19 03:25:28', '2024-12-19 03:25:42', NULL);
INSERT INTO `users` VALUES (2, 'teacher', '2@suask.me', 'dDQCQ5zaNL', '26e6d44c96358cdfc6ea25d15fa2442a', 'student', 'teacher', '# 我是老师', NULL, 0, '2024-12-19 03:27:03', '2024-12-19 03:27:03', NULL);
INSERT INTO `users` VALUES (1000, 'student', '3@suask.me', 'Ettykmu1Zc', '4fd5724ce6fdbe4ac6df655eb1dc31dc', 'student', 'student', '# 我是学生', NULL, 0, '2024-12-19 03:26:25', '2024-12-19 03:26:25', NULL);

-- ----------------------------
-- Records of questions
-- ----------------------------

DELETE FROM questions;

INSERT INTO `questions` VALUES (1,1,NULL,'你认为是什么导致人们对旅行的热爱？','不同的文化、美丽的风景，以及逃离日常生活的压力。',_binary '\0','2016-04-04 14:51:12',540,652);
INSERT INTO `questions` VALUES (2,1,NULL,'如果可以选择任何一个历史人物共进晚餐，你会选择谁？','达·芬奇，想听听他的创意和思想。',_binary '\0','2009-02-07 03:34:10',169,195);
INSERT INTO `questions` VALUES (3,1,NULL,'你最喜欢的书是什么，为什么？','《1984》，因为它对权力和自由的深刻思考令人警醒。',_binary '\0','2014-01-08 18:45:44',500,885);
INSERT INTO `questions` VALUES (4,1,NULL,'有哪些你认为值得一试的爱好或活动？','尝试陶艺或摄影，可以提升创造力。',_binary '\0','2001-06-22 13:21:38',852,427);
INSERT INTO `questions` VALUES (5,1,NULL,'你认为科技会如何改变未来的工作方式？','更多的远程工作、自动化和人工智能将改变传统的工作模式。',_binary '\0','2020-02-08 14:57:07',502,99);
INSERT INTO `questions` VALUES (6,1,NULL,'你最难忘的假期是去哪里，为什么？','去了一次日本，丰富的文化和美味的食物给我留下深刻印象。',_binary '\0','2008-03-26 17:43:25',878,833);
INSERT INTO `questions` VALUES (7,1,NULL,'如果能学会一种乐器，你希望学习什么？','钢琴，因为它的音色丰富多样。',_binary '\0','2012-01-02 13:50:47',216,255);
INSERT INTO `questions` VALUES (8,1,NULL,'你对环保有什么建议或想法？','大力推广可再生能源和减少塑料使用是关键。',_binary '\0','2006-06-17 13:53:20',756,811);
INSERT INTO `questions` VALUES (9,1,NULL,'在你的生活中，有没有什么事情让你感到特别感恩？','感恩家人和朋友的支持，他们在我困难时给予了我力量。',_binary '\0','2002-12-07 01:39:24',476,398);
INSERT INTO `questions` VALUES (10,1,NULL,'你觉得友情和爱情之间最大的区别是什么？','友情更倾向于无条件的支持，而爱情通常伴随着深度的情感和浪漫。',_binary '\0','2011-02-13 00:29:18',783,273);
INSERT INTO `questions` VALUES (11,1,NULL,'你认为幸福的定义是什么？','幸福是内心的满足和平静，能与人分享美好的时刻。',_binary '\0','2012-07-28 05:38:14',412,160);
INSERT INTO `questions` VALUES (12,1,NULL,'你最喜欢的电影是什么，有什么特别的理由吗？','《肖申克的救赎》，因为它传递了希望和坚持的力量。',_binary '\0','2019-11-16 10:18:27',108,20);
INSERT INTO `questions` VALUES (13,1,NULL,'如果有一天可以拥有超能力，你希望拥有哪种？','弗雷迪·默丘里，想听听他的音乐创作理念。',_binary '\0','2013-03-13 16:21:04',911,571);
INSERT INTO `questions` VALUES (14,1,NULL,'你如何看待素食主义的影响？','对环境和健康有积极影响，但选择要灵活，适合个人情况。',_binary '\0','2014-03-23 06:19:31',738,146);
INSERT INTO `questions` VALUES (15,1,NULL,'对于你的人生目标，有哪些计划或梦想？','希望能在职业和个人生活中找到平衡，同时实现自我成长。',_binary '\0','2024-08-18 20:44:58',558,687);
INSERT INTO `questions` VALUES (16,1,NULL,'你认为心理健康的重要性体现在哪些方面？','心理健康影响情绪、决策和人际关系，良好的心理状态是全面健康的重要部分。',_binary '\0','2011-11-27 03:08:41',402,712);
INSERT INTO `questions` VALUES (17,1,NULL,'如果复活任何一位已故的音乐家或艺术家，你会选择谁？','学习到新技能和帮助他人成功，让我感到自豪。',_binary '\0','2012-07-03 21:44:12',728,768);
INSERT INTO `questions` VALUES (18,1,NULL,'有哪些事情是你希望未来能改善的？','持续学习和自我反思是关键，可以帮助认清自己的目标和进步。',_binary '\0','2022-02-09 10:33:18',66,664);
INSERT INTO `questions` VALUES (19,1,NULL,'你觉得社交媒体对人际关系的影响是什么？','帮助保持联系，但也会导致表面化的关系和误解。',_binary '\0','2015-08-06 05:36:43',372,675);
INSERT INTO `questions` VALUES (20,1,NULL,'你对未来的科技有什么期待？','期待人工智能进一步改善人类生活，如医疗、交通等行业的发展。',_binary '\0','2013-01-13 17:43:06',312,243);
INSERT INTO `questions` VALUES (21,1,NULL,'你有哪个特别，不为人知的才能或爱好吗？','喜欢写作，偶尔会写短篇故事。',_binary '\0','2005-07-03 11:53:54',465,256);
INSERT INTO `questions` VALUES (22,1,NULL,'在你的生活中，有哪些人对你产生过重大影响？','我的父母和一些老师，他们的支持和教导对我的发展至关重要。',_binary '\0','2018-10-22 18:46:41',469,935);
INSERT INTO `questions` VALUES (23,1,NULL,'你最喜欢的季节是什么，又是什么原因？','秋天，天气宜人，树叶变色也很美。',_binary '\0','2005-02-18 13:35:29',867,154);
INSERT INTO `questions` VALUES (24,1,NULL,'有哪些书籍或电影改变了你的观点或价值观？','《小王子》让我思考人际关系的深度和纯真。',_binary '\0','2008-04-22 02:01:55',488,500);
INSERT INTO `questions` VALUES (25,1,NULL,'你怎么看待教育系统的变革？','需要更关注学生的个性化发展和批判性思维的培养。',_binary '\0','2000-12-05 06:09:54',35,156);
INSERT INTO `questions` VALUES (26,1,NULL,'如果未来可以住在任何地方，你希望选择哪里？','希望住在一个风景优美、生活节奏适中的城市，如阿姆斯特丹。',_binary '\0','2023-12-08 04:48:53',898,349);
INSERT INTO `questions` VALUES (27,2,NULL,'你认为自我提升的最有效方法是什么？','设定明确目标并坚持行动，持续学习和反思。',_binary '\0','2024-10-16 20:18:33',589,234);
INSERT INTO `questions` VALUES (28,2,NULL,'过去一年中，有哪些事情让你感到特别骄傲？','成功完成某个项目或帮助他人实现梦想，让我感到自豪。',_binary '\0','2002-02-13 14:53:04',36,218);
INSERT INTO `questions` VALUES (29,2,NULL,'你对未来的环保措施有什么看法？','需要更多的支持和行动，以减少污染并保护生态系统。',_binary '\0','2019-01-16 18:59:07',981,348);
INSERT INTO `questions` VALUES (30,2,NULL,'你认为人类未来面临的最大挑战是什么？','应对气候变化和资源短缺是人类未来最严峻的挑战。',_binary '\0','2023-03-06 15:32:49',749,825);
INSERT INTO `questions` VALUES (31,2,NULL,'如果你能选择一个超能力，你希望是什么，为什么？','能隐形，因为可以秘密探索。',_binary '\0','2015-10-28 15:01:41',934,870);
INSERT INTO `questions` VALUES (32,2,NULL,'你最喜欢的书是什么，它对你有什么影响？','《活着》，让我更珍惜生命。',_binary '\0','2004-06-02 22:06:32',523,236);
INSERT INTO `questions` VALUES (33,2,NULL,'如果你可以与历史上的任何人共进晚餐，你会选择谁？','爱因斯坦，想听他的见解。',_binary '\0','2011-06-20 07:52:55',197,400);
INSERT INTO `questions` VALUES (34,2,NULL,'描述一个你最难忘的旅行经历。','去日本赏樱花。',_binary '\0','2007-07-22 03:18:18',182,163);
INSERT INTO `questions` VALUES (35,2,NULL,'你认为人工智能在未来将如何改变我们的生活？','会使生活更便利。',_binary '\0','2023-12-11 01:15:56',809,577);
INSERT INTO `questions` VALUES (36,2,NULL,'你有没有一个人生座右铭？是什么？','\"活在当下\"。',_binary '\0','2021-11-18 01:49:00',978,169);
INSERT INTO `questions` VALUES (37,2,NULL,'在你心中，成功的定义是什么？','能实现自己的价值。',_binary '\0','2013-07-12 19:10:01',953,117);
INSERT INTO `questions` VALUES (38,2,NULL,'你觉得科技进步更多是帮助了人类，还是创造了新的问题？','带来了便利，但也带来了依赖。',_binary '\0','2011-12-15 06:48:33',879,374);
INSERT INTO `questions` VALUES (39,2,NULL,'如果你可以住在任何地方，你希望住在哪里，为什么？','纽约，充满活力和文化。',_binary '\0','2010-08-14 11:03:21',399,667);
INSERT INTO `questions` VALUES (40,2,NULL,'你最喜欢的季节是什么，是什么让你喜欢它？','秋天，喜欢它的色彩和气候。',_binary '\0','2007-01-03 09:18:51',939,809);
INSERT INTO `questions` VALUES (41,2,NULL,'有哪部电影或电视剧对你产生了深刻的影响？','《海上钢琴师》，让我思考人生选择。',_binary '\0','2021-06-01 18:09:04',494,248);
INSERT INTO `questions` VALUES (42,2,NULL,'如果时间旅行成为可能，你想去哪个时代，为什么？','维多利亚时代，想看历史的变化。',_binary '\0','2007-03-16 09:29:21',304,428);
INSERT INTO `questions` VALUES (43,2,NULL,'你在生活中遇到的最大挑战是什么，你是如何克服的？','克服焦虑，学会自我调节。',_binary '\0','2001-06-04 12:12:38',796,940);
INSERT INTO `questions` VALUES (44,2,NULL,'你觉得友谊最重要的品质是什么？','诚实和支持。',_binary '\0','2005-07-26 00:08:30',443,679);
INSERT INTO `questions` VALUES (45,2,NULL,'如果你可以学习任何技能而不需时间，你想学什么？','能说流利的多种语言。',_binary '\0','2020-10-19 16:28:02',494,237);
INSERT INTO `questions` VALUES (46,2,NULL,'你最喜欢的食物是什么，有什么特别的回忆吗？','披萨，小时候的美好回忆。',_binary '\0','2019-07-06 12:16:56',84,591);
INSERT INTO `questions` VALUES (47,2,NULL,'对你来说，家庭意味着什么？','给予支持和爱的地方。',_binary '\0','2002-04-11 05:47:27',439,495);
INSERT INTO `questions` VALUES (48,2,NULL,'你有没有梦寐以求的工作，是什么？','作家，能创造故事。',_binary '\0','2022-10-14 23:24:21',273,43);
INSERT INTO `questions` VALUES (49,2,NULL,'如果你能改变一件历史事件，你希望改变什么？','改变二战，避免无辜生命的牺牲。',_binary '\0','2000-07-25 14:08:27',763,758);
INSERT INTO `questions` VALUES (50,2,NULL,'你认为人们对幸福的定义是一个主观的概念吗？','是的，每个人的经历不同。',_binary '\0','2007-02-03 23:10:02',731,694);
INSERT INTO `questions` VALUES (51,2,NULL,'描述一下你的理想周末是怎样的。','放松、看电影、与朋友聚会。',_binary '\0','2016-09-15 10:01:20',84,195);
INSERT INTO `questions` VALUES (52,2,NULL,'你认为教育的未来会是什么样子？','更加个性化和在线化。',_binary '\0','2020-10-11 21:12:34',289,630);
INSERT INTO `questions` VALUES (53,2,NULL,'你有没有过任何让你转变观念的经历？','遇到挫折后重新审视目标。',_binary '\0','2019-03-09 04:11:41',161,669);
INSERT INTO `questions` VALUES (54,2,NULL,'你最喜欢的音乐类型是什么，为什么？','摇滚乐，充满激情。',_binary '\0','2000-10-16 08:38:52',551,618);
INSERT INTO `questions` VALUES (55,2,NULL,'有哪件事情是你一直想做但还没有实现的？','学习弹吉他。',_binary '\0','2009-09-23 17:48:14',413,174);
INSERT INTO `questions` VALUES (56,2,NULL,'你认为社交媒体对人际关系的影响是积极还是消极？','有积极和消极的影响。',_binary '\0','2001-01-27 03:20:36',719,694);
INSERT INTO `questions` VALUES (57,2,NULL,'你是否相信运气在生活中的作用？为什么？','相信，常常能碰到好运。',_binary '\0','2024-09-16 04:21:19',993,542);
INSERT INTO `questions` VALUES (58,2,NULL,'如果你能发明一个新的节日，你希望它是关于什么的？','感恩节，专门感谢身边的人。',_binary '\0','2009-07-04 19:54:52',462,685);
INSERT INTO `questions` VALUES (59,2,NULL,'哪些品质让一个人值得信任？','诚实和一贯性。',_binary '\0','2013-07-25 00:11:41',582,115);
INSERT INTO `questions` VALUES (60,2,NULL,'你觉得善良在现代社会中仍然重要吗？','是的，它让社会更温暖。',_binary '\0','2000-05-30 13:43:02',657,819);
INSERT INTO `questions` VALUES (61,2,NULL,'你最喜欢的假期是什么时候，为什么？','圣诞节，因为有家人团聚。',_binary '\0','2023-01-26 05:15:51',544,825);
INSERT INTO `questions` VALUES (62,2,NULL,'如果你可以访问任何一个国家，你会选择哪个？','日本，因为文化丰富。',_binary '\0','2015-12-08 09:25:45',805,498);
INSERT INTO `questions` VALUES (63,2,NULL,'对于你而言，什么是成功？','完成自己设定的目标。',_binary '\0','2012-08-30 06:54:01',731,338);
INSERT INTO `questions` VALUES (64,2,NULL,'你最喜欢的书是什么，为什么？','《活着》，因为反映了人性的坚韧。',_binary '\0','2013-11-17 12:22:28',615,675);
INSERT INTO `questions` VALUES (65,2,NULL,'有什么事情是你一直想尝试但还没做的？','学习一门外语。',_binary '\0','2005-09-05 18:28:05',179,181);
INSERT INTO `questions` VALUES (66,2,NULL,'描述一次改变你人生的经历。','大学毕业典礼。',_binary '\0','2005-01-20 15:24:27',225,144);
INSERT INTO `questions` VALUES (67,2,NULL,'如果你可以选择任何一种超能力，你会选择什么？','隐形。',_binary '\0','2016-10-23 22:17:03',971,755);
INSERT INTO `questions` VALUES (68,2,NULL,'对你来说，家庭的意义是什么？','互相支持和爱护。',_binary '\0','2007-01-26 05:52:59',860,325);
INSERT INTO `questions` VALUES (69,2,NULL,'你最喜欢的音乐类型是什么，为什么？','摇滚乐，因为能激发能量。',_binary '\0','2004-07-16 04:58:59',842,694);
INSERT INTO `questions` VALUES (70,2,NULL,'有没有一部电影让你感动得流泪？是哪一部？','《海上钢琴师》。',_binary '\0','2012-06-18 04:39:44',274,135);
INSERT INTO `questions` VALUES (71,2,NULL,'你如何面对压力？','深呼吸和运动。',_binary '\0','2023-10-08 00:06:41',124,572);
INSERT INTO `questions` VALUES (72,2,NULL,'谈谈一次让你感到特别自豪的成就。','完成马拉松比赛。',_binary '\0','2004-05-04 15:05:46',348,752);
INSERT INTO `questions` VALUES (73,2,NULL,'你觉得科技在生活中发挥了怎样的作用？','改善我们的生活质量。',_binary '\0','2003-11-17 16:48:12',764,590);
INSERT INTO `questions` VALUES (74,2,NULL,'什么样的食物让你觉得幸福？','巧克力蛋糕。',_binary '\0','2015-11-21 17:26:43',535,693);
INSERT INTO `questions` VALUES (75,2,NULL,'如果可以回到过去，你最想改变什么？','做得更好一些。',_binary '\0','2015-01-02 02:15:21',302,894);
INSERT INTO `questions` VALUES (76,2,NULL,'你认为友情有什么重要性？','朋友能提供情感支持。',_binary '\0','2019-05-05 11:22:00',593,803);
INSERT INTO `questions` VALUES (77,2,NULL,'你有没有什么特别的爱好？是什么？','17.摄影。',_binary '\0','2019-05-10 01:26:28',702,74);
INSERT INTO `questions` VALUES (78,2,NULL,'你最喜欢的季节是什么，为什么？','春天，因为万物复苏。',_binary '\0','2005-04-25 21:41:25',791,262);
INSERT INTO `questions` VALUES (79,2,NULL,'谈谈你心目中的理想工作。','设计师/艺术家。',_binary '\0','2023-08-03 14:46:37',520,470);
INSERT INTO `questions` VALUES (80,2,NULL,'有没有一位名人对你产生过影响？是谁？','马丁·路德·金。',_binary '\0','2008-10-24 23:24:06',393,298);
INSERT INTO `questions` VALUES (81,2,NULL,'你对当前的社会发展有什么看法？','有正面影响，也有负面影响。',_binary '\0','2021-01-12 01:17:47',632,816);
INSERT INTO `questions` VALUES (82,2,NULL,'你如何定义幸福？','能够内心平和。',_binary '\0','2024-09-19 12:45:31',276,880);
INSERT INTO `questions` VALUES (83,2,NULL,'描述一次你最难忘的旅行经历。','去巴黎的旅行。',_binary '\0','2019-02-28 03:21:36',999,358);
INSERT INTO `questions` VALUES (84,2,NULL,'你喜欢通过什么方式放松自己？','听音乐和散步。',_binary '\0','2009-01-11 04:18:37',32,701);
INSERT INTO `questions` VALUES (85,2,NULL,'有什么事情是你希望了解的，但没机会学习？','心理学。',_binary '\0','2020-05-11 19:27:54',2,602);
INSERT INTO `questions` VALUES (86,2,NULL,'如果你能与任何人共进晚餐，你会选择谁？','爱因斯坦。',_binary '\0','2001-11-01 09:17:58',578,535);
INSERT INTO `questions` VALUES (87,2,NULL,'你觉得自己在生活中最大的挑战是什么？','工作与生活的平衡。',_binary '\0','2001-08-23 02:05:09',408,709);
INSERT INTO `questions` VALUES (88,2,NULL,'你会如何向别人描述你的性格？','开朗和乐于助人。',_binary '\0','2019-10-11 20:11:42',113,348);
INSERT INTO `questions` VALUES (89,2,NULL,'如果有人要写一部关于你的书，你希望它的主题是什么？','追求梦想与勇气。',_binary '\0','2001-08-22 19:11:01',439,513);
INSERT INTO `questions` VALUES (90,2,NULL,'描述一下你最喜欢的天气和活动。','晴天，适合户外活动。',_binary '\0','2012-02-24 04:10:49',936,515);
INSERT INTO `questions` VALUES (91,2,NULL,'你觉得教育在个人发展中扮演什么角色？','提升个体的思维和能力。',_binary '\0','2020-08-10 19:35:45',843,313);
INSERT INTO `questions` VALUES (92,2,NULL,'如果你可以改变世界上的一件事，那会是什么？','消除饥饿与贫困。',_binary '\0','2021-12-14 16:15:43',154,71);
INSERT INTO `questions` VALUES (93,2,NULL,'你有没有一个特别的梦想？是什么？','开一家自己的咖啡馆。',_binary '\0','2020-03-07 16:05:36',992,303);
INSERT INTO `questions` VALUES (94,2,NULL,'你认为成年人应该有怎样的责任？','关心社会与他人。',_binary '\0','2008-11-17 07:35:08',781,551);
INSERT INTO `questions` VALUES (95,2,NULL,'描述一下理想中的居住环境。','安静、和平、靠近自然。',_binary '\0','2023-06-25 22:04:48',467,397);
INSERT INTO `questions` VALUES (96,2,NULL,'有没有一本书或一部电影让你改变了看法？是什么？','《百年孤独》，让我对时间有新理解。',_binary '\0','2019-11-20 09:41:49',116,114);
INSERT INTO `questions` VALUES (97,2,NULL,'你认为与朋友的关系如何影响你的生活？','影响很大，能给生活带来乐趣。',_binary '\0','2004-10-15 23:38:22',43,875);
INSERT INTO `questions` VALUES (98,2,NULL,'什么样的事情能激励到你？','看到他人的努力和决心。',_binary '\0','2002-11-18 01:07:58',289,850);
INSERT INTO `questions` VALUES (99,2,NULL,'你在生活中最感激的事情是什么？','有爱我的家人和朋友。',_binary '\0','2003-08-27 01:31:37',439,770);
INSERT INTO `questions` VALUES (100,2,NULL,'你最喜欢的运动是什么，为什么？','足球，因为可以锻炼体力。',_binary '\0','2002-05-04 15:42:32',944,296);
INSERT INTO `questions` VALUES (101,2,NULL,'如果你可以学会任何一种乐器，你想学什么？','吉他。',_binary '\0','2018-10-20 13:45:49',304,480);

-- ----------------------------
-- Records of favorites
-- ----------------------------

DELETE FROM favorites;

INSERT INTO `favorites` VALUES (1,1,1,'2024-12-18 13:25:35','默认收藏夹');
INSERT INTO `favorites` VALUES (2,1,2,'2024-12-15 09:36:32','默认收藏夹');
INSERT INTO `favorites` VALUES (3,1,3,'2024-12-15 08:10:24','默认收藏夹');
INSERT INTO `favorites` VALUES (4,1,4,'2024-12-16 02:16:41','默认收藏夹');
INSERT INTO `favorites` VALUES (5,1,5,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (6,1,17,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (7,1,18,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (8,1,19,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (9,1,20,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (10,1,21,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (11,1,22,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (12,1,23,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (13,1,24,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (14,1,25,'2024-12-16 07:15:58','默认收藏夹');
INSERT INTO `favorites` VALUES (15,1,26,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (16,1,27,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (17,1,28,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (18,1,29,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (19,1,30,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (20,1,31,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (21,1,32,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (22,1,33,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (23,1,34,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (24,1,35,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (25,1,36,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (26,1,37,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (27,1,38,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (28,1,39,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (29,1,40,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (30,1,41,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (31,1,42,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (32,1,43,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (33,1,44,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (34,1,45,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (35,1,46,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (36,1,47,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (37,1,48,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (38,1,49,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (39,1,50,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (40,1,51,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (41,1,52,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (42,1,53,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (43,1,54,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (44,1,55,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (45,1,56,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (46,1,57,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (47,1,58,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (48,1,59,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (49,1,60,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (50,1,61,'2024-12-18 13:23:19','默认收藏夹');
INSERT INTO `favorites` VALUES (51,1,62,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (52,1,63,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (53,1,64,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (54,1,65,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (55,1,66,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (56,1,67,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (57,1,68,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (58,1,69,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (59,1,70,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (60,1,71,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (61,1,72,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (62,1,73,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (63,1,74,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (64,1,75,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (65,1,76,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (66,1,77,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (67,1,78,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (68,1,79,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (69,1,80,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (70,1,81,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (71,1,82,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (72,1,83,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (73,1,84,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (74,1,85,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (75,1,86,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (76,1,87,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (77,1,88,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (78,1,89,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (79,1,90,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (80,1,91,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (81,1,92,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (82,1,93,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (83,1,94,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (84,1,95,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (85,1,96,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (86,1,97,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (87,1,98,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (88,1,99,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (89,1,100,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (90,1,101,'2024-12-18 13:25:08','默认收藏夹');

-- ----------------------------
-- Records of files
-- ----------------------------

DELETE FROM files;

INSERT INTO `files` VALUES (1, '1.png', 0x2CBDFB377C60CC77EB810E7134BCD01E1C3FC1ACD3EAA9BE815863B997428D36, 1, '2024-12-19 13:14:24');
INSERT INTO `files` VALUES (2, '2.png', 0xDAAA509C33DC6CA956ED3DD802AD0C6D851351D0D6C392BC81E1B4A56D1B2707, 1, '2024-12-19 13:14:34');
INSERT INTO `files` VALUES (3, '3.png', 0x2677AB5DEFC59056446EBDD43006F2DE27603660EFCE7BCAF01F2001C5CFBFD4, 1, '2024-12-19 13:15:55');
INSERT INTO `files` VALUES (4, '4.png', 0xAEB7BC6B3353AC89CE865DCAA4EFA8B2B5948DE57223537C8FBD1B8C295CC5A7, 1, '2024-12-19 13:16:12');
INSERT INTO `files` VALUES (5, '5.png', 0x5F36045C9D8C2FDF15B22DBB0F3ABAD11C1CCD3D8D1291587D54DC101BEE07B3, 1, '2024-12-19 13:16:17');
INSERT INTO `files` VALUES (6, '6.png', 0x7F43D95221A94AABE063D8F5435E069ACC23840F8DE6D8899BA46B56B8F6168D, 1, '2024-12-19 13:16:23');
INSERT INTO `files` VALUES (7, '7.png', 0x716A3C19E9BB451AC30AB010A509C8873D447E2AD5A940A53B42CBA1F61D0C1F, 1, '2024-12-19 13:16:31');
INSERT INTO `files` VALUES (8, '8.png', 0x47D6B5BA427F1DDB3EB885E406598E2BA568F2E1F53CF07DC86948DAD07C1C8C, 1, '2024-12-19 13:16:36');
INSERT INTO `files` VALUES (9, '9.png', 0xF770232CEF955B7BAAFB009B899D5775E991D7EB7EDA43D2D59B004CCCE473A3, 1, '2024-12-19 13:16:41');
INSERT INTO `files` VALUES (10, '10.png', 0x007045B2F2001DF842E0FD448587FA00CCCD04807BC26EEC7D01252C299EA137, 1, '2024-12-19 13:16:46');
