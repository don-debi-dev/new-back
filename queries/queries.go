package queries

import (
	"database/sql"
	"log"
	"new-back/utils"

	"github.com/aws/aws-lambda-go/events"
)

type UserProfile struct {
	UserId   int64  `json:"id"`
	UserName string `json:"userName"`
}

type LoginDetails struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type tempProfile struct {
	userId   int64
	userName sql.NullString
}

func GetUserProfileFromId(conn *sql.DB, id int64) (UserProfile, error) {
	var u tempProfile

	if err := conn.QueryRow(`SELECT * FROM mybackdatabase.users WHERE id = ?`, id).Scan(&u.userId, &u.userName); err != nil {
		if err == sql.ErrNoRows {
			return UserProfile{}, nil
		} else {
			return UserProfile{}, err
		}
	}

	return UserProfile{
		UserId:   u.userId,
		UserName: utils.ParseNullString(u.userName),
	}, nil
}

func DeleteUserProfileById(conn *sql.DB, id int64) (UserProfile, error) {
	tx, err := conn.Begin()
	if err != nil {
		log.Printf("Error connecting to database Error: %s\n", err.Error())
		return UserProfile{}, err
	}
	defer tx.Rollback()

	deleteUser, err := tx.Prepare(`DELETE FROM mybackdatabase.users WHERE id = ?`)
	if err != nil {
		log.Printf("Error preparing delete user profile statement Error: %s\n", err.Error())
		return UserProfile{}, err
	}
	defer deleteUser.Close()

	_, err = deleteUser.Exec(id)
	if err != nil {
		log.Printf("Error executing delete user profile statement Error: %s\n", err.Error())
		return UserProfile{}, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error commiting delete user profile statement Error: %s\n", err.Error())
		return UserProfile{}, err
	}

	return UserProfile{
		UserId: id,
	}, nil
}

func CreateUserProfile(conn *sql.DB, userName string) (UserProfile, error) {
	tx, err := conn.Begin()
	if err != nil {
		log.Printf("Error connecting to database Error: %s\n", err.Error())
		return UserProfile{}, err
	}
	defer tx.Rollback()

	createUser, err := tx.Prepare(`INSERT INTO mybackdatabase.users (username) VALUES (?)`)
	if err != nil {
		log.Printf("Error preparing user profile statement Error: %s\n", err.Error())
		return UserProfile{}, err
	}
	defer createUser.Close()

	_, err = createUser.Exec(userName)
	if err != nil {
		log.Printf("Error executing user profile statement Error: %s\n", err.Error())
		return UserProfile{}, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error commiting user profile statement Error: %s\n", err.Error())
		return UserProfile{}, err
	}

	return UserProfile{}, nil
}

func GetUserProfileFromName(conn *sql.DB, userName string) (UserProfile, error) {
	var u tempProfile
	if err := conn.QueryRow(`SELECT * FROM mybackdatabase.users WHERE username = ?`, userName).Scan(&u.userId, &u.userName); err != nil {
		if err == sql.ErrNoRows {
			// Handle case when no rows were found
			return UserProfile{}, nil
		} else {
			// Handle other errors
			return UserProfile{}, err
		}
	}

	return UserProfile{
		UserId:   u.userId,
		UserName: utils.ParseNullString(u.userName),
	}, nil
}

func GetUserIdFromNameAndPassword(conn *sql.DB, login LoginDetails) (UserProfile, error) {
	var u tempProfile
	if err := conn.QueryRow(`SELECT id FROM mybackdatabase.users WHERE username = ? AND password = ?`, login.UserName, login.Password).Scan(&u.userId, &u.userName); err != nil {
		if err == sql.ErrNoRows {
			// Handle case when no rows were found
			return UserProfile{}, nil
		} else {
			// Handle other errors
			return UserProfile{}, err
		}
	}

	return UserProfile{
		UserId:   u.userId,
	}, nil
}

func SetUserSession(conn *sql.DB, userName string) (UserProfile, error) {
	var u tempProfile
	if err := conn.QueryRow(`SELECT * FROM mybackdatabase.users WHERE username = ?`, userName).Scan(&u.userId, &u.userName); err != nil {
		if err == sql.ErrNoRows {
			// Handle case when no rows were found
			return UserProfile{}, nil
		} else {
			// Handle other errors
			return UserProfile{}, err
		}
	}

	return UserProfile{
		UserId:   u.userId,
		UserName: utils.ParseNullString(u.userName),
	}, nil
}

func GetLambdaResponse(response CustomResponse, body ...string) (events.APIGatewayProxyResponse, error) {
	head := map[string]string{"Content-Type": "application/json", "Access-Control-Allow-Origin": "*"}

	if body != nil {
		return events.APIGatewayProxyResponse{
			Headers:    head,
			StatusCode: int(response.Status),
			Body:       string(body[0]),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Headers:    head,
		StatusCode: int(response.Status),
		Body:       string(response.Message),
	}, nil
}

func GetAllUserProfiles(conn *sql.DB) ([]UserProfile, error) {
	users := []UserProfile{}

	rows, err := conn.Query(`SELECT * FROM mybackdatabase.users ORDER BY id`)

	if err != nil {
		if err == sql.ErrNoRows {
			return users, nil
		} else {
			return users, err
		}
	}

	for rows.Next() {
		var u tempProfile

		if err := rows.Scan(&u.userId, &u.userName); err != nil {
			log.Printf("Error Iterating Users %s\n", err.Error())
			return users, err
		}

		user := UserProfile{
			UserId:   u.userId,
			UserName: utils.ParseNullString(u.userName),
		}

		users = append(users, user)
	}

	return users, nil
}
