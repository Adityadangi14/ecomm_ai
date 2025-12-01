package helpers

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Adityadangi14/ecomm_ai/products-service/src/models"
	"github.com/redis/go-redis/v9"
)

func GetQueriesWithDecay(ctx context.Context, rdb *redis.Client, key string) (string, error) {

	queries, err := rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return "", err
	}

	if len(queries) == 0 {
		return "", nil
	}

	var builder strings.Builder

	decay := 0.5

	for i, q := range queries {
		weight := math.Pow(decay, float64(i))
		builder.WriteString(fmt.Sprintf("[%.2f] %s\n", weight, q))
	}

	return builder.String(), nil
}

func GetUserChat(ctx context.Context, rdb *redis.Client, key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()

	// Handle missing key safely
	if err == redis.Nil {
		return "", nil
	}

	// Handle actual errors
	if err != nil {
		return "", err
	}

	return val, nil
}
func SetUserQueries(rdb *redis.Client, query string, key string, ctx context.Context) error {
	res := rdb.LPush(ctx, key, query)
	rdb.Expire(ctx, key, time.Hour*48)
	if res.Err() != nil {
		return fmt.Errorf("failed to push value to redis: %v ", res.Err().Error())
	}
	val := rdb.LTrim(ctx, key, 0, 9)
	if val.Err() != nil {
		return fmt.Errorf("failed to trim value to redis: %v", res.Err().Error())
	}

	return nil
}

func SetUserChat(rdb *redis.Client, key string, val string, ctx context.Context) error {
	_, err := rdb.Set(ctx, key, val, time.Hour*48).Result()

	if err != nil {
		return nil
	}
	return nil
}

func GetUserQueriesKey(params models.AiQueryParams) string {
	return fmt.Sprintf("user_queries:%v_%v_%v", params.UserID, params.OrgID, params.SessionID)
}

func GetUserChatKey(params models.AiQueryParams) string {
	return fmt.Sprintf("user_chat:%v_%v_%v", params.UserID, params.OrgID, params.SessionID)
}
