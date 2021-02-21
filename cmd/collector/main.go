package main

import (
	"log"
	"os"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/SevereCloud/vksdk/v2/api"

	"collector/internal/config"

	_ "github.com/lib/pq"
)

func main() {
	conf := config.GetProjectConfig()
	vk := api.NewVK(conf.Token)

	f := excelize.NewFile()
	f.SetSheetName("Sheet1", "пользователи")

	f.SetCellValue("пользователи", "A1", "ID")
	f.SetCellValue("пользователи", "B1", "Имя")
	f.SetCellValue("пользователи", "C1", "Фамилия")
	f.SetCellValue("пользователи", "D1", "Никнейм")
	f.SetCellValue("пользователи", "E1", "Дата рождения")
	f.SetCellValue("пользователи", "F1", "Город")
	f.SetCellValue("пользователи", "G1", "Количество друзей")
	f.SetCellValue("пользователи", "H1", "Количество подписчиков")
	f.SetCellValue("пользователи", "I1", "Количество подписок")

	row := 2
	for i, search := range conf.Search {
		users, err := vk.UsersSearch(api.Params{
			"fields":    conf.Fields,
			"count":     search.Count,
			"city":      search.City,
			"sex":       search.Sex,
			"status":    search.Status,
			"age_from":  search.AgeFrom,
			"age_to":    search.AgeTo,
			"has_photo": search.HasPhoto,
			"religion":  search.Religion,
		})
		if err != nil {
			log.Printf("search[%d]: %s", i, err)
			os.Exit(1)
		}

		for l, user := range users.Items {
			f.SetCellValue("пользователи", "A"+strconv.Itoa(row+l), user.ID)
			f.SetCellValue("пользователи", "B"+strconv.Itoa(row+l), user.FirstName)
			f.SetCellValue("пользователи", "C"+strconv.Itoa(row+l), user.LastName)
			f.SetCellValue("пользователи", "D"+strconv.Itoa(row+l), user.Nickname)
			f.SetCellValue("пользователи", "E"+strconv.Itoa(row+l), user.Bdate)
			f.SetCellValue("пользователи", "F"+strconv.Itoa(row+l), user.City.Title)
			//! friends
			friends, _ := vk.FriendsGet(api.Params{
				"user_id": user.ID,
			})
			f.SetCellValue("пользователи", "G"+strconv.Itoa(row+l), friends.Count)
			//! followers
			followers, _ := vk.UsersGetFollowers(api.Params{
				"user_id": user.ID,
			})
			f.SetCellValue("пользователи", "H"+strconv.Itoa(row+l), followers.Count)
			//! subscriptions
			subscriptions, _ := vk.UsersGetSubscriptions(api.Params{
				"user_id": user.ID,
			})
			f.SetCellValue("пользователи", "I"+strconv.Itoa(row+l), subscriptions.Groups.Count)
		}
		row += search.Count
		log.Printf("search[%d] done\n", i)
		if err := f.SaveAs("base.xlsx"); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	if err := f.SaveAs("base.xlsx"); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}
