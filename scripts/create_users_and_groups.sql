-- Insert Users
INSERT INTO Users (username, name, email, password_hash) VALUES
('user1', 'John Doe', 'user1@gmail.com', 'ZKn4vbb1fmXfWsQiRuqSnA==:99bd3597109a3131997a2377feff8d30615a7a459dd5d1f689fb1c2a5b9ee6f06abc64b5323021073392e4b9823405dee85ab03495f7c65dd8668413c2515150'),
('user2', 'John Doe', 'user2@gmail.com', 'PgSH8f4Y0ZxmTwgVqVHrcQ==:18de3e67e0dd57adc011d9e463d5c2f83c4fd2a48197fe646be27526d00f34b706d14e5f963b13caa6a783f17ca007c613ae2737f457b2974a6663478487437a'),
('user3', 'John Doe', 'user3@gmail.com', 'BV0BPdVbbdMzvJOrZ0/I/g==:163ca5bb5cbeaaab9d5334f6d1f1f98e028e2985565262fcf729edd529c084a3093a056de890406e6fb29d0824a75b2d465ee9c72f6a37336f51a5928bd9e02d'),
('user4', 'John Doe', 'user4@gmail.com', '9qO/2JTJbROm5kslNt0hLw==:02ee7357a203936f25f00f061bb11f2a1279a84f57ac999bb4a3aace33c635fb7f5f2ddbb647cf5a6e2bb7829053cafd06e072dfdee1e5270d84045e8708b761'),
('victoroliveirab', 'Victor Oliveira', 'eu@victoroliveira.com.br', 'z4QVaSQT3u0Q/EQ33vbHFw==:ZWLxUg0P/xeBTSItizg1Zq7P1vTy3cg4kxEK83oK2CA=');

-- Insert Groups
INSERT INTO Groups (admin, name, description, points_table, ranking) VALUES
(1, 'Group A', 'Description for Group A', '{"Perfect":18,"Opposite":-8,"DiffPlusWinner":15,"DiffPlusOpposite":0,"WinnerPlusWinnerGoals":12,"WinnerPlusLoserGoals":11,"Winner":10,"Draw":15,"OneGoalButDraw":0,"None":0}', '{"1":0,"2":0,"5":0}'),
(1, 'Group B', 'Description for Group B', '{"Perfect":20,"Opposite":-10,"DiffPlusWinner":15,"DiffPlusOpposite":-5,"WinnerPlusWinnerGoals":12,"WinnerPlusLoserGoals":11,"Winner":10,"Draw":15,"OneGoalButDraw":4,"None":0}', '{"1":0,"3":0}'),
(5, 'My Private Group', 'Group for me to do my private guesses', '{"Perfect":20,"Opposite":-10,"DiffPlusWinner":15,"DiffPlusOpposite":-5,"WinnerPlusWinnerGoals":12,"WinnerPlusLoserGoals":11,"Winner":10,"Draw":15,"OneGoalButDraw":4,"None":0}', '{"5":0}');

-- Assign Users to Groups
-- User 1 is part of both Group A and Group B
INSERT INTO User_Groups (user_id, group_id) VALUES
((SELECT id FROM Users WHERE username = 'user1'), (SELECT id FROM Groups WHERE name = 'Group A')),
((SELECT id FROM Users WHERE username = 'user1'), (SELECT id FROM Groups WHERE name = 'Group B'));

-- User 2 is part of Group A
INSERT INTO User_Groups (user_id, group_id) VALUES
((SELECT id FROM Users WHERE username = 'user2'), (SELECT id FROM Groups WHERE name = 'Group A'));

-- User 3 is part of Group B
INSERT INTO User_Groups (user_id, group_id) VALUES
((SELECT id FROM Users WHERE username = 'user3'), (SELECT id FROM Groups WHERE name = 'Group B'));

-- User 4 is not part of any group, so no entry in User_Groups for user4

-- Insert me as part of Group A
INSERT INTO User_Groups (user_id, group_id) VALUES
((SELECT id FROM Users WHERE username = 'victoroliveirab'), (SELECT id FROM Groups WHERE name = 'Group A')),
((SELECT id FROM Users WHERE username = 'victoroliveirab'), (SELECT id FROM Groups WHERE name = 'My Private Group'));

