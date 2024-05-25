-- Insert Users
INSERT INTO Users (username, email, password_hash) VALUES
('user1', 'user1@gmail.com', 'ZKn4vbb1fmXfWsQiRuqSnA==:99bd3597109a3131997a2377feff8d30615a7a459dd5d1f689fb1c2a5b9ee6f06abc64b5323021073392e4b9823405dee85ab03495f7c65dd8668413c2515150'),
('user2', 'user2@gmail.com', 'PgSH8f4Y0ZxmTwgVqVHrcQ==:18de3e67e0dd57adc011d9e463d5c2f83c4fd2a48197fe646be27526d00f34b706d14e5f963b13caa6a783f17ca007c613ae2737f457b2974a6663478487437a'),
('user3', 'user3@gmail.com', 'BV0BPdVbbdMzvJOrZ0/I/g==:163ca5bb5cbeaaab9d5334f6d1f1f98e028e2985565262fcf729edd529c084a3093a056de890406e6fb29d0824a75b2d465ee9c72f6a37336f51a5928bd9e02d'),
('user4', 'user4@gmail.com', '9qO/2JTJbROm5kslNt0hLw==:02ee7357a203936f25f00f061bb11f2a1279a84f57ac999bb4a3aace33c635fb7f5f2ddbb647cf5a6e2bb7829053cafd06e072dfdee1e5270d84045e8708b761'),
('victoroliveirab', 'eu@victoroliveira.com.br', 'z4QVaSQT3u0Q/EQ33vbHFw==:ZWLxUg0P/xeBTSItizg1Zq7P1vTy3cg4kxEK83oK2CA=');

-- Insert Groups
INSERT INTO Groups (name, description) VALUES
('Group A', 'Description for Group A'),
('Group B', 'Description for Group B');

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
