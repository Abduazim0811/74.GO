package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func CreateClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	return rdb
}

func AddUserScores(rdb *redis.Client, userScores map[string]float64) {
	for user, score := range userScores {
		err := rdb.ZAdd(ctx, "leaderboard", redis.Z{
			Score:  score,
			Member: user,
		}).Err()
		if err != nil {
			log.Fatalf("Could not add user score: %v", err)
		}
	}
}

func GetUserScore(rdb *redis.Client, user string) (float64, error) {
	score, err := rdb.ZScore(ctx, "leaderboard", user).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("user %s does not exist", user)
	}
	return score, err
}

func DisplayLeaderboard(rdb *redis.Client, topN int64) {
	leaders, err := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, topN-1).Result()
	if err != nil {
		log.Fatalf("Could not get leaderboard: %v", err)
	}

	fmt.Println("Leaderboard:")
	for i, leader := range leaders {
		fmt.Printf("%d. %s: %.2f\n", i+1, leader.Member, leader.Score)
	}
}

func UpdateUserScores(rdb *redis.Client, user string, scoreDelta float64) {
	pipe := rdb.TxPipeline()
	pipe.ZIncrBy(ctx, "leaderboard", scoreDelta, user)
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Could not update user score: %v", err)
	}
}
