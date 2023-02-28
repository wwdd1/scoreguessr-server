package repo

import (
	"fmt"
	"log"
	"wc22/tools"
	"wc22/types"
)

func CreateOrUpdateUser(user *types.User) error {
	db := tools.GetDb()
	query := fmt.Sprintf(`
		INSERT INTO public."users"(uid, name, email, profile_picture, last_login)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (uid)
		DO UPDATE
		SET profile_picture=$4, last_login=$5
	`)
	_, err := db.Query(query, user.Uid, user.Name, user.Email, user.ProfilePicture, user.LastLogin)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	return nil
}

func GetUser(uid string) (*types.User, error) {
	db := tools.GetDb()
	query := fmt.Sprintf(`
		SELECT uid, name, email, profile_picture, last_login
		FROM public."users"
		WHERE uid = $1
	`)
	result, err := db.Query(query, uid)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}
	fmt.Print(result)
	return nil, nil
}
