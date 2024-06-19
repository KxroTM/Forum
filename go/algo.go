package forum

import (
	"database/sql"
	"math/rand"
	"strings"
)

func ForYouPageAlgorithm(db *sql.DB, user_id string) []Post {
	user := GetAccountById(db, user_id)

	var posts []Post
	var AllPosts = GetAllPosts(db)

	for _, post := range AllPosts {
		if strings.Contains(user.CategorieSub, post.Categorie) && user.CategorieSub != "" {
			if !contains(posts, post) {
				if strings.Contains(post.Liker, UserSession.Username) {
					post.IsLike = true
				}
				if strings.Contains(post.Disliker, UserSession.Username) {
					post.IsDislike = true
				}
				if strings.Contains(post.Retweeter, UserSession.Username) {
					post.IsRetweet = true
				}
				posts = append(posts, post)
			}
		}
	}

	for _, post := range AllPosts {
		if strings.Contains(user.FollowingList, post.User_id) {
			if !contains(posts, post) {
				if strings.Contains(post.Liker, UserSession.Username) {
					post.IsLike = true
				}
				if strings.Contains(post.Disliker, UserSession.Username) {
					post.IsDislike = true
				}
				if strings.Contains(post.Retweeter, UserSession.Username) {
					post.IsRetweet = true
				}
				posts = append(posts, post)
			}
		}
	}

	return posts
}

func RecommendedUsers(db *sql.DB, user_id string) RecommendedUser {
	accounts := RecommendUsersAlgoByCommonFollowings(db, user_id)
	algo := RecommendUserAlgorithmByCategorie(db, user_id)

	for i := 0; i < len(algo.RecommendedUsers); i++ {
		if !containsUser(accounts.RecommendedUsers, algo.RecommendedUsers[i]) {
			accounts.RecommendedUsers = append(accounts.RecommendedUsers, algo.RecommendedUsers[i])
			accounts.Reason = append(accounts.Reason, algo.Reason[i])
		}
	}

	return accounts
}

func RecommendUsersAlgoByCommonFollowings(db *sql.DB, user_id string) RecommendedUser {
	user := GetAccountById(db, user_id)

	var accounts RecommendedUser
	var AllUsers = GetAllUsers(db)

	if len(user.FollowingList) != 0 {
		followings := strings.Split(user.FollowingList, ",")
		for i := 0; i < len(followings)-1; i++ {
			user2 := GetAccountByUsername(db, followings[i])
			user3 := GetAccountByUsername(db, followings[i+1])
			followings2 := strings.Split(user2.FollowingList, ",")
			followings3 := strings.Split(user3.FollowingList, ",")
			for j := 0; j < len(followings2); j++ {
				for k := 0; k < len(followings3); k++ {
					if followings2[j] == followings3[k] {
						usertemp := GetAccountByUsername(db, followings2[j])
						if !containsUser(accounts.RecommendedUsers, usertemp) && usertemp != user && !strings.Contains(user.FollowingList, usertemp.Username) {
							accounts.RecommendedUsers = append(accounts.RecommendedUsers, usertemp)
							accounts.Reason = append(accounts.Reason, "Amis en commun")
						}
					}
				}
			}
		}
	}

	if len(AllUsers) < 3 {
		if len(accounts.RecommendedUsers) < 3 {
			randomUsers := rand.Perm(len(AllUsers))[:len(AllUsers)]
			for _, i := range randomUsers {
				randomUser := AllUsers[i]
				if !containsUser(accounts.RecommendedUsers, randomUser) && randomUser.User_id != user_id {
					accounts.RecommendedUsers = append(accounts.RecommendedUsers, randomUser)
					accounts.Reason = append(accounts.Reason, "Suggérer par le site")
				}
			}
		}
	} else {
		if len(accounts.RecommendedUsers) < 3 {
			randomUsers := rand.Perm(len(AllUsers))[:3]
			for _, i := range randomUsers {
				randomUser := AllUsers[i]
				if !containsUser(accounts.RecommendedUsers, randomUser) && randomUser.User_id != user_id {
					accounts.RecommendedUsers = append(accounts.RecommendedUsers, randomUser)
					accounts.Reason = append(accounts.Reason, "Suggérer par le site")
				}
			}
		}
	}

	return accounts
}

func RecommendUserAlgorithmByCategorie(db *sql.DB, user_id string) RecommendedUser {
	user := GetAccountById(db, user_id)

	var accounts RecommendedUser
	var AllUsers = GetAllUsers(db)

	if len(user.CategorieSub) != 0 {
		for len(accounts.RecommendedUsers) < 3 {
			for _, account := range AllUsers {
				if strings.Contains(user.CategorieSub, account.CategorieSub) {
					accounts.RecommendedUsers = append(accounts.RecommendedUsers, account)
				}
			}
		}
	}

	return accounts
}

func SearchPageAlgorithm(db *sql.DB, search string) []Post {
	var posts []Post
	AllPosts := GetAllPosts(db)

	for _, post := range AllPosts {
		if strings.Contains(post.Title, search) || strings.Contains(post.Text, search) || strings.Contains(post.Categorie, search) {
			posts = append(posts, post)
		}
	}

	return posts
}

func contains(posts []Post, post Post) bool {
	for _, p := range posts {
		if p == post {
			return true
		}
	}
	return false
}

func containsUser(accounts []User, account User) bool {
	for _, a := range accounts {
		if a == account {
			return true
		}
	}
	return false
}

func RealRecommendUsers(db *sql.DB, user_id string) RecommendedUser {
	user := GetAccountById(db, user_id)

	var accounts RecommendedUser
	var AllUsers = GetAllUsers(db)

	if len(user.FollowingList) != 0 {
		followings := strings.Split(user.FollowingList, ",")
		for i := 0; i < len(followings)-1; i++ {
			user2 := GetAccountByUsername(db, followings[i])
			user3 := GetAccountByUsername(db, followings[i+1])
			followings2 := strings.Split(user2.FollowingList, ",")
			followings3 := strings.Split(user3.FollowingList, ",")
			for j := 0; j < len(followings2); j++ {
				for k := 0; k < len(followings3); k++ {
					if followings2[j] == followings3[k] {
						usertemp := GetAccountByUsername(db, followings2[j])
						if !containsUser(accounts.RecommendedUsers, usertemp) && usertemp != user && !strings.Contains(user.FollowingList, usertemp.Username) {
							accounts.RecommendedUsers = append(accounts.RecommendedUsers, usertemp)
							accounts.Reason = append(accounts.Reason, "Amis en commun")
						}
					}
				}
			}
		}
	}

	if len(AllUsers) < 5 {
		randomUsers := rand.Perm(len(AllUsers))[:len(AllUsers)]
		for _, i := range randomUsers {
			randomUser := AllUsers[i]
			if !containsUser(accounts.RecommendedUsers, randomUser) && randomUser.User_id != user_id {
				accounts.RecommendedUsers = append(accounts.RecommendedUsers, randomUser)
				accounts.Reason = append(accounts.Reason, "Suggérer par le site")
			}
		}
	} else {
		randomUsers := rand.Perm(len(AllUsers))[:5]
		for _, i := range randomUsers {
			randomUser := AllUsers[i]
			if !containsUser(accounts.RecommendedUsers, randomUser) && randomUser.User_id != user_id {
				accounts.RecommendedUsers = append(accounts.RecommendedUsers, randomUser)
				accounts.Reason = append(accounts.Reason, "Suggérer par le site")
			}
		}
	}
	return accounts
}
