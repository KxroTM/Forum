SQLite format 3   @     \   	   	                                                       \ .v�= ���
�	�	�	�                                                                                                                                                                                                                                                                                                                               �'�$	�tableposts_tmpposts_tmpCREATE TABLE `posts_tmp` (
						posts_id TEXT PRIMARY KEY NOT NULL,
						UUID TEXT NOT NULL,
                        categorie TEXT NOT NULL,
                        title TEXT NOT NULL,
						text TEXT NOT NULL,
						like INTEGER NOT NULL,
						liker TEXT NOT NULL,
						dislike INTEGER NOT NULL,
						retweet INTEGER NOT NULL,
						retweeter TEXT NOT NULL,
						date TEXT NOT NULL,
						report INTEGER NOT NULL, disliker TEXT NOT NULL,
						FOREIGN KEY (UUID) REFERENCES users (UUID)
�0	�?tablepostspostsCREATE TABLE "posts" (
						posts_id TEXT PRIMARY KEY NOT NULL,
						UUID TEXT NOT NULL,
                        categorie TEXT NOT NULL,
                        title TEXT NOT NULL,
						text TEXT NOT NULL,
						like INTEGER NOT NULL,
						liker TEXT NOT NULL,
						dislike INTEGER NOT NULL,
						retweet INTEGER NOT NULL,
						retweeter TEXT NOT NULL,
						date TEXT NOT NULL,
						report INTEGER NOT NULL, disliker TEXT NOT NULL, user_pfp TEXT NOT NULL,
						FOREIGN KEY (UUID) REFERENCES users (UUID)
                    )�6�KtableusersusersCREATE TABLE "users" (
						UUID TEXT PRIMARY KEY NOT NULL,
						role TEXT NOT NULL,
                        username TEXT NOT NULL,
                        email TEXT NOT NULL,
						password TEXT NOT NULL,
						created_at TEXT NOT NULL,
						updated_at TEXT NOT NULL,
						profilePicture TEXT NOT NULL,
						followers INTEGER NOT NULL,
						following INTEGER NOT NULL,
						bio TEXT NOT NULL,
						links TEXT NOT NULL,
						categoriesSub TEXT NOT NULL
                    , followersList TEXT NOT NULL, followingList TEXT NOT NULL)
D��gtablepostspostsCREATE TABLE posts (
						posts_id INTEGER PRIMARY KEY NOT NULL,
						UUID INTEGER NOT NULL,
                        categorie TEXT NOT NULL,
                        title TEXT NOT NULL,
						text TEXT NOT NULL,
						like INTEGER NOT NULL,
						liker TEXT NOT NULL,
						dislike INTEGER NOT NULL,
						retweet INTEGER NOT NULL,
	7K% indexsqlite_autoindex_comments_tmp_1comments_tmp
D 9    /C indexsqlite_autoindex_comments_1comments)
= indexsqlite_autoindex_posts_1posts ~�wtablecommentscommentsCREATE TABLE comments (
						comment_id INTEGER PRIMARY KEY NOT NULL,
						post_id INTEGER NOT NULL,
						UUID INTEGER NOT NULL,
                        text TEXT NOT NULL,
                        date TEXT NOT NULL,
						like INTEGER NOT NULL,
						dislike INTEGER NOT NULL,
						report INTEGER NOT NULL,
						FOREIGN KEY (UUID) REFERENCES us��stablecommentscommentsCREATE TABLE "comments" (
						comment_id TEXT PRIMARY KEY NOT NULL,
						post_id TEXT NOT NULL,
						UUID TEXT NOT NULL,
                        text TEXT NOT NULL,
                        date TEXT NOT NULL,
						like INTEGER NOT NULL,
						dislike INTEGER NOT NULL,
						report INTEGER NOT NULL, liker TEXT NOT NULL, disliker TEXT NOT NULL, user_pfp TEXT NOT NULL,
						FOREIGN KEY (UUID) REFERENCES users (UUID)
						FOREIGN KEY (post_id) REFERENCES posts (post_id)
                    ))= indexsqlite_autoindex_users_1users  ��WtableusersusersCREAT�]�Z�W%%�qtablecomments_tmpcomments_tmpCREATE TABLE `comments_tmp` (
						comment_id TEXT PRIMARY KEY NOT NULL,
						post_id TEXT NOT NULL,
						UUID TEXT NOT NULL,
                        text TEXT NOT NULL,
                        date TEXT NOT NULL,
						like INTEGER NOT NULL,
						dislike INTEGER NOT NULL,
						report INTEGER NOT NULL,
						FOREIGN KEY (UUID) REFERENCES users (UUID)
						FOREIGN KEY (post_id) REFERENCES posts (post_id)
                    ))                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    �                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           yUU-82cfb232-ec5b-4aff-b083-61a97b714d7f561f3361-7088-40a0-90f0-a947370dc607categorietitletext27-04-2024 21:39                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              � � �$����;�                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   �	U�!!ec6be46f-4559-4ed4-8b6d-12efeda8d673userFYTFGVEZU99ad61fdaae0d48bd0ded78985bc7cae3e87f3c50977c8f27a10dc51c37b77ce26-04-202426-04-2024�U�!!2c806f94-5f70-4fdc-ae38-2f31227a2c7buserFYTFGVEZU99ad61fdaae0�SU-�!!i04fa1879-e000-443a-8734-9b176768f80buserjuifdu77ootdb1@gmail.com081dcd1424b917e9f938eb3590fa5059bf220d79f7209239f137f829e991dbaa06-05-202406-05-2024../../style/media/default_avatar/avatar_03.png�`U);�!!i17d6e584-1800-4689-bf14-8b505f261e5cuseryoussef.ammariyoussef.ammari@ynov.comf979eb331d77599a3032eb54492020edb075c4da5f946101bfbd2f6674287b5206-05-202406-05-2024../../style/media/default_avatar/avatar_04.png�_U!?�!!ife8db749-8d27-4464-b698-8943813fe34cadminForumAdminforumprojetynov@gmail.comb4704bd614684a5f67c8419fb828803810c4cb1dc021885b234f2e2bf1e2d7a929-04-202429-04-2024../../style/media/default_avatar/avatar_01.png   �U;;�!!1dedc900-c6f9-4692-b3dc-ad42e14056c0useryoussef.ammari@ynov.comyoussef.ammari@ynov.com956d8bc6d6e2c740b34aa8497fbd�iU;;�!!i�YU!5�!!i1cfcbf4d-ea50-4ec2-a52b-09741fb9c2b1userendhayyyggendhayyygg@gmail.com986b3ca81817821cbf9ff837696046c1917d57daeeae55a1455d2be4841dd13406-05-202406-05-2024../../style/media/default_avatar/avatar_04.png
� ] ]��������                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      (Uec6be46f-4559-4ed4-8b6d-12efeda8d673	(U2c806f94-5f70-4fdc-ae38-2f31227a2c7b(Ua90fcb7f-42c5-49cb-99f5-36843c634e1a(U14b8dd2d-a094-4ce7-814d-b74e2b686e61(U04fa1879-e000-443a-8734-9b176768f80b(U17d6e584-1800-4689-bf14-8b505f261e5c(U1cfcbf4d-ea50-4ec2-a52b-09741fb9c2b1(Ufe8db749-8d27-4464-b698-8943813fe34c   (7e75f060-a076-43f0-bd82-999c335255ad
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             
      �                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              'U	82cfb232-ec5b-4aff-b083-61a97b714d7f                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              