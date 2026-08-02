package seeder

import (
	"database/sql"
	"fmt"
	"math/rand/v2"

	"golang.org/x/crypto/bcrypt"
)

type Seed struct {
	db *sql.DB
}

func NewSeed(db *sql.DB) Seed {
	return Seed{
		db: db,
	}
}

func (s Seed) Run() error {
	adminID, err := s.initADMIN()
	if err != nil {
		return err
	}

	if adminID == 0 {
		return nil
	}

	userIDs, err := s.seedUsers()
	if err != nil {
		return err
	}

	postIDs, err := s.seedPosts(userIDs)
	if err != nil {
		return err
	}

	return s.seedComments(userIDs, postIDs)
}

func (s Seed) initADMIN() (int64, error) {
	var exists bool

	err := s.db.QueryRow(
		`
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email = $1
		)
		`,
		"admin@gmail.com",
	).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("check admin existence: %w", err)
	}

	if exists {
		return 0, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("123456"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return 0, fmt.Errorf("bcrypt error: %w", err)
	}

	var id int64

	err = s.db.QueryRow(
		`
		INSERT INTO users(username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
		`,
		"admin",
		"admin@gmail.com",
		string(hashedPassword),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("seed admin error: %w", err)
	}

	return id, nil
}

func (s Seed) seedUsers() ([]int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("123456"),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, fmt.Errorf("bcrypt error: [%w]", err)
	}

	userIDs := make([]int64, 0, 10)

	for i := 1; i <= 10; i++ {
		username := fmt.Sprintf("user_%d", i)
		email := fmt.Sprintf("email%d@gmail.com", i)

		var id int64

		err := s.db.QueryRow(
			`
			INSERT INTO users(username, email, password)
			VALUES ($1, $2, $3)
			RETURNING id
			`,
			username,
			email,
			string(hashedPassword),
		).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("seed user insert error: idx [%d] - [%w]", i, err)
		}

		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}

func (s Seed) seedPosts(userIDs []int64) ([]int64, error) {
	postIDs := make([]int64, 0, len(userIDs)*3)

	tags := []string{
		"go", "backend", "microservice", "clean-architecture", "ddd", "tdd",
	}

	for _, userID := range userIDs {
		for i := 1; i <= 3; i++ {
			title := fmt.Sprintf("Post %d", i)
			content := fmt.Sprintf("Content of post: %d", i)

			var postID int64

			err := s.db.QueryRow(`
				INSERT INTO posts(title, content, user_id, tags)
				VALUES ($1, $2, $3, $4) 
				RETURNING id
				`,
				title,
				content,
				userID,
				tags[:rand.IntN(len(tags))+1],
			).Scan(&postID)

			if err != nil {
				return nil, fmt.Errorf("seed posts insert error: idx [%d] - [%w]", i, err)
			}

			postIDs = append(postIDs, postID)
		}
	}

	return postIDs, nil
}

func (s Seed) seedComments(userIDs, postIDs []int64) error {
	comments := []string{
		"Hay quá!",
		"Cảm ơn đã chia sẻ.",
		"Rất hữu ích.",
		"Đồng ý với quan điểm này.",
		"Mình học được khá nhiều.",
		"Cho mình hỏi thêm được không?",
		"Ví dụ rất dễ hiểu.",
		"Go thật sự rất hay.",
		"Backend thú vị đấy.",
		"Like!",
	}

	for _, postID := range postIDs {

		for _, userID := range userIDs {

			total := rand.IntN(2) + 1 // 1-2 comments/user

			for i := 0; i < total; i++ {

				content := comments[rand.IntN(len(comments))]

				_, err := s.db.Exec(
					`
					INSERT INTO comments(post_id, user_id, content)
					VALUES ($1, $2, $3)
					`,
					postID,
					userID,
					content,
				)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
