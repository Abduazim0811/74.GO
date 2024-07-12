package main

import "74.Go/internal/redis"

func main() {
	rdb := redis.CreateClient()
	defer rdb.Close()

	userScores := map[string]float64{
		"user1": 100,
		"user2": 200,
		"user3": 150,
	}

	redis.AddUserScores(rdb, userScores)
	redis.DisplayLeaderboard(rdb, 3)

	redis.UpdateUserScores(rdb, "user1", 50)
	redis.DisplayLeaderboard(rdb, 3)
}
