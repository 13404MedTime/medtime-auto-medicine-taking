package function

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/spf13/cast"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	multipleUpdateUrl = "https://api.admin.u-code.io/v1/object/multiple-update/"

	appId     = "P-JV2nVIRUtgyPO5xRNeYll2mT4F5QG4bS"
	urlConst  = "https://api.admin.u-code.io"
	botToken  = "6799733885:AAHiknQxaSzSgeDXhZOVxTNajHoRyX_rhRE"
	chatID    = int64(-162256495)
	baseUrl   = "https://api.admin.u-code.io"
	notifSlug = "notifications"
)

// func main() {
// 	// data := `{"data":{"additional_parameters":[],"app_id":"P-JV2nVIRUtgyPO5xRNeYll2mT4F5QG4bS","method":"CREATE","object_data":{"app_id":"P-JV2nVIRUtgyPO5xRNeYll2mT4F5QG4bS","cleints_id":"8626c75c-1cf6-42dd-9531-78fa80f744ba","comment":"\u003cp\u003esd dsfsdf\u003c/p\u003e","guid":"e2a757f8-feb8-41d3-9d56-46e299b44f4d","json_body":"{\"type\":\"custom\",\"dayData\":[null],\"weekData\":[],\"customData\":{\"cycle_name\":\"week\",\"cycle_count\":1,\"time\":\"21:50:10\",\"dates\":[\"tuesday\",\"thursday\",\"saturday\"]},\"before_after_food\":\"after\",\"start_date\":\"2023-07-01\",\"end_date\":\"2023-07-29\"}","medicine_time":"","method":"CREATE","naznachenie_id":"d594dc27-9caf-4e84-9b9b-b748937dc189","preparati_id":"dcfb63b3-3db7-4689-bfaa-e86ac1cea3e0"},"object_ids":["e2a757f8-feb8-41d3-9d56-46e299b44f4d"],"table_slug":"medicine_taking","user_id":"c384727f-408a-4dd0-b242-db3cc9355edd"}}`

// 	data := `{
// 		"data": {
// 			"additional_parameters": [],
// 			"app_id": "P-JV2nVIRUtgyPO5xRNeYll2mT4F5QG4bS",
// 			"environment_id": "dcd76a3d-c71b-4998-9e5c-ab1e783264d0",
// 			"label_ru": "Автоматический добавления лекарств",
// 			"label_uz": "Автоматический добавления лекарств",
// 			"method": "CREATE",
// 			"object_data": {
// 				"cleints_id": "adda1b93-3a8b-4ff8-922c-bb4f13d559aa",
// 				"company_service_environment_id": "dcd76a3d-c71b-4998-9e5c-ab1e783264d0",
// 				"company_service_project_id": "a4dc1f1c-d20f-4c1a-abf5-b819076604bc",
// 				"current_amount": 6,
// 				"description": [
// 					"before_food"
// 				],
// 				"dosage": 2,
// 				"frequency": [
// 					"several_times_day"
// 				],
// 				"guid": "036bc6a3-a198-43b3-a50c-1920a4fabc8f",
// 				"is_from_patient": true,
// 				"json_body": "{\"hours_of_day\":[\"16:15:00\",\"16:20:00\",\"16:25:00\",\"16:30:00\",\"16:35:00\"]}",
// 				"number": 1,
// 				"preparati_id": "fe40982d-59e0-493b-9bc4-ff5c69a52396",
// 				"unit_of_measure_id": "ce083fb6-94cd-4d2d-81f3-c89f28643925",
// 				"week_days": [
// 					"0",
// 					"1",
// 					"2",
// 					"3",
// 					"4",
// 					"5",
// 					"6"
// 				]
// 			},
// 			"object_data_before_update": null,
// 			"object_ids": [
// 				"036bc6a3-a198-43b3-a50c-1920a4fabc8f"
// 			],
// 			"project_id": "a4dc1f1c-d20f-4c1a-abf5-b819076604bc",
// 			"table_slug": "medicine_taking",
// 			"user_id": ""
// 		}
// 	}`

// 	fmt.Println(Handle([]byte(data)))
// }

// Handle a serverless request
func Handle(req []byte) string {
	var (
		response Response
		request  NewRequestBody
	)

	// Send("req body medicine: QFIHQLHFLQELK  " + string(req))

	if err := json.Unmarshal(req, &request); err != nil {
		return Handler("error", err.Error())
	}

	var (
		// tableSlug  = "medicine_taking"
		objectData = request.Data["object_data"].(map[string]interface{})
		clientId   = ""
		medicine   Medicine
		body       = objectData["json_body"].(string)
	)

	if objectData["cleints_id"] != nil {
		clientId = objectData["cleints_id"].(string)
	}

	err := json.Unmarshal([]byte(body), &medicine)
	if err != nil {
		Handler("error", err.Error())
	}

	var json_data map[string]interface{}

	err = json.Unmarshal([]byte(body), &json_data)
	if err != nil {
		Handler("error", err.Error())
	}

	st := cast.ToSlice(json_data["hours_of_day"])

	if request.Data["method"].(string) == "CREATE" {

		isFromPatient, _ := objectData["is_from_patient"].(bool)

		if true {

			var (
				daysStr, _ = objectData["week_days"].([]interface{})
				days       = []int{}
			)

			for _, v := range daysStr {

				var (
					valstr, _ = v.(string)
					val, _    = strconv.Atoi(valstr)
				)

				days = append(days, int(val))
			}

			frequencyInt, _ := objectData["frequency"].([]interface{})
			if len(frequencyInt) < 1 {
				return Handler("error", err.Error())
			}

			frequencyList := []string{}

			for _, v := range frequencyInt {
				val, _ := v.(string)
				frequencyList = append(frequencyList, val)
			}

			frequencyType := frequencyList[0]
			if frequencyType == "several_times_day" || frequencyType == "always" {
				days = []int{0, 1, 2, 3, 4, 5, 6}
			}

			sort.Ints(days)

			var (
				timeString     = medicine.HoursOfDay
				sortedTimes, _ = sortHours(timeString)
				amount         float64
			)

			switch v := objectData["current_amount"].(type) {
			case string:
				amountInt, err := strconv.Atoi(v)
				if err != nil {
					response.Data = map[string]interface{}{"error while strconv.Atoi": err.Error()}
					response.Status = "error"
					responseByte, _ := json.Marshal(response)
					return string(responseByte)
				}
				amount = float64(amountInt)
			case float64:
				amount = v
			default:
				if frequencyType == "always" {
					amount = cast.ToFloat64(10 * len(st))
				} else {
					return Handler("error", "error while getting current amount"+" current amount is not string or float64")
				}
			}

			var dosage float64

			dosageStr, ok := objectData["dosage"].(string)
			if ok {
				dosageInt, err := strconv.Atoi(dosageStr)
				if err != nil {
					return Handler("error", err.Error())
				}

				dosage = float64(dosageInt)

			} else {

				dosage, ok = objectData["dosage"].(float64)
				if !ok {
					return Handler("error", "error while getting dosage, message")
				}
			}
			// Send(fmt.Sprintf("a m o u n t %d, d o s a g e %f", int(amount), dosage))
			// tableSlug = "patient_medication"

			var (
				currentTime   = time.Now()
				afterFoodInt  = objectData["description"].([]interface{})
				afterFoodList = []string{}
			)

			for _, v := range afterFoodInt {
				val, _ := v.(string)
				afterFoodList = append(afterFoodList, val)
			}

			var (
				boolAfterFood    = afterFoodList[0]
				preparatId, _    = objectData["preparati_id"].(string)
				appointmentId, _ = objectData["naznachenie_id"].(string)
				requests         = []map[string]interface{}{}
				notifRequests    = MultipleUpdateRequest{}
			)

			if frequencyType == "always" {
				amount = cast.ToFloat64(10*len(st)) * dosage
			}

			for amount > 0 {
				time := getNextDate(currentTime, days, sortedTimes)
				currentTime = time

				var (
					serverTime      = currentTime
					stringTime      = serverTime.Format("2006-01-02T15:04:05.000Z")
					preparatName, _ = objectData["medicine_name"].(string)
				)

				requests = append(requests, map[string]interface{}{
					"naznachenie_id":     appointmentId,
					"medicine_taking_id": objectData["guid"].(string),
					"time_take":          stringTime,
					"before_after_food":  boolAfterFood,
					"cleints_id":         clientId,
					"preparati_id":       preparatId,
					"is_from_patient":    isFromPatient,
					"count":              dosage,
					"preparat_name":      preparatName,
				},
				)

				// notifRequest := Request{
				// 	Data: map[string]interface{}{
				// 		"client_id":    clientId,
				// 		"title":        "Время принятия препарата!",
				// 		"body":         "Вам назначен препарат: ",
				// 		"title_uz":     "Preparatni qabul qilish vaqti bo'ldi!",
				// 		"body_uz":      "Sizga preparat tayinlangan: ",
				// 		"is_read":      false,
				// 		"preparati_id": preparatId,
				// 		"time_take":    stringTime,
				// 	},
				// }
				notifRequests.Data.Objects = append(notifRequests.Data.Objects, map[string]interface{}{
					"client_id":    clientId,
					"title":        "Время принятия препарата!",
					"body":         "Вам назначен препарат: ",
					"title_uz":     "Preparatni qabul qilish vaqti bo'ldi!",
					"body_uz":      "Sizga preparat tayinlangan: ",
					"is_read":      false,
					"preparati_id": preparatId,
					"time_take":    stringTime,
				})
				amount -= dosage
			}

			// req := Request{
			// 	Data: map[string]interface{}{
			// 		"objects": requests,
			// 	},
			// }
			body, err := MultipleUpdateObject(urlConst, "patient_medication", Request{Data: map[string]interface{}{"objects": requests}})
			if err != nil {
				return Handler("error", err.Error())
			}

			var bodyMap map[string]interface{}

			err = json.Unmarshal(body, &bodyMap)
			if err != nil {
				return Handler("error while unmarshal 1", err.Error())
			}

			m_object := cast.ToSlice(cast.ToStringMap(cast.ToStringMap(bodyMap["data"])["data"])["objects"])

			if len(m_object) > 0 {
				mId := cast.ToStringMap(cast.ToStringMap(cast.ToStringMap(m_object[len(m_object)-1]))["data"])["medicine_taking_id"]
				timeToTake := cast.ToStringMap(cast.ToStringMap(cast.ToStringMap(m_object[len(m_object)-1]))["data"])["time_take"]

				newM := map[string]interface{}{
					"guid":      mId,
					"last_time": timeToTake,
				}

				if frequencyType == "always" {
					newM["current_amount"] = 0
				}

				var data = map[string]interface{}{
					"data": newM,
				}

				_, err = DoRequest(baseUrl+"/v1/object/medicine_taking", "PUT", data, appId)
				if err != nil {
					return Handler("error while do update request 1", err.Error())
				}
			}

			// req = Request{
			// 	Data: map[string]interface{}{
			// 		"objects": notifRequests,
			// 	},
			// }
			// res, _ := json.Marshal(req)
			// Send("Multiple Update notif    " + string(res))
			// err = MultipleUpdateObject(urlConst, notifSlug, req)

			_, err = DoRequest(multipleUpdateUrl+"notifications", "PUT", notifRequests, appId)
			if err != nil {
				return Handler("error", err.Error())
			}
			if !isFromPatient {

				_, ok = objectData["preparati_id"].(string)
				_, ok2 := objectData["naznachenie_id"].(string)

				if ok && ok2 {
					manyReq := RequestMany2Many{
						IdFrom:    objectData["naznachenie_id"].(string),
						IdTo:      []string{objectData["preparati_id"].(string)},
						TableFrom: "naznachenie",
						TableTo:   "preparati",
					}

					err, _ = UpdateObjectMany2Many(urlConst, appId, manyReq)
					if err != nil {
						return Handler("error", err.Error())
					}
				}
			}
		} else {
			manyReq := RequestMany2Many{
				IdFrom:    objectData["naznachenie_id"].(string),
				IdTo:      []string{objectData["preparati_id"].(string)},
				TableFrom: "naznachenie",
				TableTo:   "preparati",
			}

			// testdat, _ := json.Marshal(manyReq)
			// Send(string(testdat))
			err, _ := UpdateObjectMany2Many(urlConst, appId, manyReq)
			if err != nil {
				return Handler("error", err.Error())
			}

			boolAfterFood := medicine.BeforeAfterFood

			// tableSlug = "patient_medication"

			// create patient medicine by day
			if medicine.Type == "day" {
				var (
					startDateStr = medicine.StartDate
					endDateStr   = medicine.EndDate
					hours        = medicine.DayData
					startDate, _ = time.Parse("2006-01-02", startDateStr)
					endDate, _   = time.Parse("2006-01-02", endDateStr)
				)
				// Loop over the range of dates
				for currentDate := startDate; currentDate.Before(endDate) || currentDate.Equal(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
					for _, hour := range hours {
						hourTime, _ := time.Parse("15:04:05", hour)
						combinedDateTime := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), hourTime.Hour(), hourTime.Minute(), hourTime.Second(), 0, currentDate.Location())
						combinedDateTime = combinedDateTime.Add(time.Hour * -5)

						// Create objects response example

						var (
							timeTake             = combinedDateTime.Format("2006-01-02T15:04:05.000Z")
							preparatId           = objectData["preparati_id"].(string)
							createtObjectRequest = Request{
								Data: map[string]interface{}{
									"naznachenie_id":     objectData["naznachenie_id"].(string),
									"medicine_taking_id": objectData["guid"].(string),
									"time_take":          timeTake,
									"before_after_food":  boolAfterFood,
									"cleints_id":         clientId,
									"preparati_id":       preparatId,
								},
							}
						)

						_, err, response = CreateObject(urlConst, "patient_medication", appId, createtObjectRequest)
						if err != nil {
							responseByte, _ := json.Marshal(response)
							return string(responseByte)
						}

						// add to notification table
						addToNotifTable(timeTake, preparatId, clientId, appId, notifSlug)
					}
				}

			}
			// create patient medicine by week
			if medicine.Type == "week" {
				var (
					startDateStr = medicine.StartDate
					endDateStr   = medicine.EndDate
					dayTimes     = medicine.WeekData
					startDate, _ = time.Parse("2006-01-02", startDateStr)
					endDate, _   = time.Parse("2006-01-02", endDateStr)
					datesInRange = getDatesInRange(startDate, endDate, dayTimes)
				)

				for _, date := range datesInRange {
					date = date.Add(time.Hour * -5)

					// tableSlug = "patient_medication"
					//create objects response example
					var (
						timeTake             = date.Format("2006-01-02T15:04:05.000Z")
						preparatId           = objectData["preparati_id"].(string)
						createtObjectRequest = Request{
							Data: map[string]interface{}{
								"naznachenie_id":     objectData["naznachenie_id"].(string),
								"medicine_taking_id": objectData["guid"].(string),
								"time_take":          timeTake,
								"before_after_food":  boolAfterFood,
								"cleints_id":         clientId,
								"preparati_id":       preparatId,
							},
						}
					)

					_, err, response = CreateObject(urlConst, "patient_medication", appId, createtObjectRequest)
					if err != nil {
						responseByte, _ := json.Marshal(response)
						return string(responseByte)
					}

					// add to notification table
					addToNotifTable(timeTake, preparatId, clientId, appId, notifSlug)
				}
			}
			// create patient medicine by month
			if medicine.Type == "month" {
				dateTimes := medicine.MonthData

				for _, dateTime := range dateTimes {
					var (
						inputDateStr     = dateTime.Date
						inputTimeStr     = dateTime.Time
						layout           = "2006-01-02 15:04:05" // Combine date and time format
						inputDateTimeStr = fmt.Sprintf("%s %s", inputDateStr, inputTimeStr)
						inputDateTime, _ = time.Parse(layout, inputDateTimeStr)
					)

					// Parse input date and time
					if err != nil {
						response.Data = map[string]interface{}{
							"data": "Error parsing date and time:",
						}
						responseByte, _ := json.Marshal(response)
						return string(responseByte)
					}

					// Convert to desired format "2006-01-02T15:04:05.000Z"
					resultDateTimeStr := inputDateTime.UTC().Format("2006-01-02T15:04:05.000Z")

					// tableSlug = "patient_medication"
					// create objects response example
					var (
						preparatId = objectData["preparati_id"].(string)
					)
					createtObjectRequest := Request{
						// some filters
						Data: map[string]interface{}{
							"naznachenie_id":     objectData["naznachenie_id"].(string),
							"medicine_taking_id": objectData["guid"].(string),
							"time_take":          resultDateTimeStr,
							"before_after_food":  boolAfterFood,
							"cleints_id":         clientId,
							"preparati_id":       preparatId,
						},
					}
					_, err, response = CreateObject(urlConst, "patient_medication", appId, createtObjectRequest)
					if err != nil {
						responseByte, _ := json.Marshal(response)
						return string(responseByte)
					}

					// add to notification table
					addToNotifTable(resultDateTimeStr, preparatId, clientId, appId, notifSlug)
				}

			}
			// // create patient medicine by custom
			if medicine.Type == "custom" {
				startDateStr := medicine.StartDate
				endDateStr := medicine.EndDate
				custom := medicine.CustomData

				if custom.CycleName == "day" || custom.CycleName == "month" {
					startDate, _ := time.Parse("2006-01-02", startDateStr)
					endDate, _ := time.Parse("2006-01-02", endDateStr)

					customData := CustomDataObj{
						CycleCount: custom.CycleCount,
						CycleName:  custom.CycleName,
						Time:       custom.Time,
					}
					pillDates := calculatePillDates(startDate, endDate, customData)

					// Print the pill dates
					for _, date := range pillDates {

						// tableSlug = "patient_medication"
						//create objects response example
						var (
							timeTake   = date.Format("2006-01-02T15:04:05.000Z")
							preparatId = objectData["preparati_id"].(string)
						)
						createtObjectRequest := Request{
							// some filters
							Data: map[string]interface{}{
								"naznachenie_id":     objectData["naznachenie_id"].(string),
								"medicine_taking_id": objectData["guid"].(string),
								"time_take":          timeTake,
								"before_after_food":  boolAfterFood,
								"cleints_id":         clientId,
								"preparati_id":       preparatId,
							},
						}
						_, err, response = CreateObject(urlConst, "patient_medication", appId, createtObjectRequest)
						if err != nil {
							responseByte, _ := json.Marshal(response)
							return string(responseByte)
						}

						// add to notification table
						addToNotifTable(timeTake, preparatId, clientId, appId, notifSlug)
					}

				} else if custom.CycleName == "week" {
					days := custom.Dates
					dayDates := getDatesInRangeWeek(startDateStr, endDateStr, days, custom.CycleCount, custom.Time)

					res := []string{}
					for _, obj := range dayDates {

						for i := 0; i < len(obj.Dates); i += custom.CycleCount {
							res = append(res, obj.Dates[i])
						}
					}

					for _, v := range res {
						// tableSlug = "patient_medication"
						//create objects response example
						var (
							preparatId = objectData["preparati_id"].(string)
						)
						createtObjectRequest := Request{
							// some filters
							Data: map[string]interface{}{
								"naznachenie_id":     objectData["naznachenie_id"].(string),
								"medicine_taking_id": objectData["guid"].(string),
								"time_take":          v,
								"before_after_food":  boolAfterFood,
								"cleints_id":         clientId,
								"preparati_id":       preparatId,
							},
						}
						_, err, response := CreateObject(urlConst, "patient_medication", appId, createtObjectRequest)
						if err != nil {
							responseByte, _ := json.Marshal(response)
							return string(responseByte)
						}

						// add to notification table
						addToNotifTable(v, preparatId, clientId, appId, notifSlug)
					}
				}
			}
		}
	}

	response.Data = map[string]interface{}{}
	response.Status = "done" //if all will be ok else "error"
	responseByte, _ := json.Marshal(response)

	return string(responseByte)
}

func MultipleUpdateObject(url, tableSlug string, request Request) ([]byte, error) {
	resp, err := DoRequest(url+"/v1/object/multiple-update/"+tableSlug, "PUT", request, appId)
	// fmt.Println("resp", string(resp), "err", err)
	if err != nil {
		return nil, errors.New("error while updating multiple objects" + err.Error())
	}
	return resp, nil
}

func sortHours(timeStrings []string) ([]time.Time, error) {
	// Parse the time strings into time.Time objects
	times := make([]time.Time, len(timeStrings))
	for i, str := range timeStrings {
		parsedTime, err := time.Parse("15:04:05", str)
		if err != nil {
			// fmt.Println("Error parsing time:", err)
			return nil, err
		}
		parsedTime = parsedTime.Add(time.Hour * -5)
		times[i] = parsedTime
	}

	// Sort the time.Time objects
	sort.Slice(times, func(i, j int) bool {
		return times[i].Before(times[j])
	})

	// // Format the sorted times as strings
	// sortedTimeStrings := make([]string, len(times))
	// for i, t := range times {
	// 	sortedTimeStrings[i] = t.Format("15:04:05")
	// }

	return times, nil
}

func getNextDate(current time.Time, days []int, times []time.Time) time.Time {
	nextDate := current

	// Get next hour
	var nextTime time.Time

	for _, t := range times {
		if t.Hour() == current.Hour() {
			if t.Minute() > current.Minute() {
				nextTime = t
				break
			}
		} else if t.Hour() > current.Hour() {
			nextTime = t
			break
		}
	}

	if nextTime == (time.Time{}) {
		nextTime = times[0]
		nextDate = nextDate.AddDate(0, 0, 1)
	}
	// current day of the week
	currentDay := int(nextDate.Weekday())

	// iterate days array and find next upcoming day
	addition := -1
	for _, day := range days {
		if day >= currentDay {
			addition = day - currentDay
			nextDate = nextDate.AddDate(0, 0, day-currentDay)
			break
		}
	}
	if addition == -1 {
		nextDate = nextDate.AddDate(0, 0, days[0]+7-currentDay)
	}

	// Combine the next date and time
	nextDateTime := time.Date(nextDate.Year(), nextDate.Month(), nextDate.Day(), nextTime.Hour(), nextTime.Minute(), nextTime.Second(), 0, nextDate.Location())

	return nextDateTime
}

func addToNotifTable(timeTake, preparatId, clientId, appId, tableSlug string) {

	notifRequest := Request{
		Data: map[string]interface{}{
			"client_id":    clientId,
			"title":        "Время принятия препарата!",
			"body":         "Вам назначен препарат: ",
			"title_uz":     "Preparatni qabul qilish vaqti bo'ldi!",
			"body_uz":      "Sizga preparat tayinlangan: ",
			"is_read":      false,
			"preparati_id": preparatId,
			"time_take":    timeTake,
		},
	}
	// res, _ := json.Marshal(notifRequest)
	// Send("addToNotifTable request   " + string(res))
	CreateObject(urlConst, tableSlug, appId, notifRequest)
}

func CreateObject(url, tableSlug, appId string, request Request) (Datas, error, Response) {
	response := Response{}

	var createdObject Datas
	createObjectResponseInByte, err := DoRequest(url+"/v1/object/"+tableSlug+"?from-ofs=true&project-id=a4dc1f1c-d20f-4c1a-abf5-b819076604bc", "POST", request, appId)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while creating object"}
		response.Status = "error"
		return Datas{}, errors.New("error"), response
	}
	err = json.Unmarshal(createObjectResponseInByte, &createdObject)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while unmarshalling create object object"}
		response.Status = "error"
		return Datas{}, errors.New("error"), response
	}
	return createdObject, nil, response
}

func UpdateObjectMany2Many(url, appId string, request RequestMany2Many) (error, Response) {
	response := Response{}

	_, err := DoRequest(url+"/v1/many-to-many/?project-id=a4dc1f1c-d20f-4c1a-abf5-b819076604bc", "PUT", request, appId)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while updating object"}
		response.Status = "error"
		return errors.New("error"), response
	}
	return nil, response
}

func convertToWeekday(day string) time.Weekday {
	switch day {
	case "sunday":
		return time.Sunday
	case "monday":
		return time.Monday
	case "tuesday":
		return time.Tuesday
	case "wednesday":
		return time.Wednesday
	case "thursday":
		return time.Thursday
	case "friday":
		return time.Friday
	case "saturday":
		return time.Saturday
	default:
		return time.Sunday // Return a default value (Sunday) in case of an invalid day.
	}
}

func getDatesInRange(startDate, endDate time.Time, dayTimes []DayTime) []time.Time {
	var datesInRange []time.Time

	for current := startDate; !current.After(endDate); current = current.AddDate(0, 0, 1) {
		for _, dt := range dayTimes {
			day := convertToWeekday(dt.Day)
			if current.Weekday() == day {
				t, _ := time.Parse("15:04:05", dt.Time)
				dateTime := current.Add(time.Hour*time.Duration(t.Hour()) + time.Minute*time.Duration(t.Minute()))
				datesInRange = append(datesInRange, dateTime)
			}
		}
	}

	return datesInRange
}

func calculatePillDates(startDate, endDate time.Time, customData CustomDataObj) []time.Time {
	var pillDates []time.Time

	// Parse time from customData.Time
	customTime, _ := time.Parse("15:04:05", customData.Time)

	// Define the duration based on the cycle_name and cycle_count
	var duration time.Duration
	switch customData.CycleName {
	case "day":
		duration = time.Duration(customData.CycleCount) * 24 * time.Hour
	case "month":
		duration = time.Duration(customData.CycleCount) * 30 * 24 * time.Hour
	default:
		return pillDates
	}

	// Initialize the current date as the start date
	currentDate := startDate

	// Loop to calculate pill dates between startDate and endDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		// Combine currentDate and customTime to get the pillDate
		pillDate := currentDate.Add(time.Duration(customTime.Hour()) * time.Hour)
		pillDate = pillDate.Add(time.Duration(customTime.Minute()) * time.Minute)
		pillDate = pillDate.Add(time.Hour * -5)
		// Add the pillDate to the result
		pillDates = append(pillDates, pillDate)

		// Move currentDate to the next cycle based on the duration
		currentDate = currentDate.Add(duration)
	}

	return pillDates
}

// week logics
func DoRequest(url string, method string, body interface{}, appId string) ([]byte, error) {
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	// fmt.Println("dataaaaa", string(data))
	client := &http.Client{
		Timeout: time.Duration(25 * time.Second),
	}
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	request.Header.Add("authorization", "API-KEY")
	request.Header.Add("X-API-KEY", appId)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respByte, nil
}

func Send(text string) {
	bot, _ := tgbotapi.NewBotAPI("6799733885:AAHiknQxaSzSgeDXhZOVxTNajHoRyX_rhRE")

	msg := tgbotapi.NewMessage(162256495, text)

	bot.Send(msg)
}

func Handler(status, message string) string {
	var (
		response Response
		Message  = make(map[string]interface{})
	)

	// sendMessage("automedicinetaking", status, message)
	response.Status = status
	data := Request{
		Data: map[string]interface{}{
			"data": message,
		},
	}
	response.Data = data.Data
	Message["message"] = message
	respByte, _ := json.Marshal(response)
	return string(respByte)
}

func sendMessage(functionName, errorStatus string, message interface{}) {
	bot, _ := tgbotapi.NewBotAPI("6799733885:AAHiknQxaSzSgeDXhZOVxTNajHoRyX_rhRE")

	chatID := int64(162256495)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("message from %s function: %s\n%s", functionName, errorStatus, message))

	bot.Send(msg)
}

func getDatesInRangeWeek(startDate, endDate string, days []string, cycleCount int, hourStr string) []DateObject {
	layout := "2006-01-02"
	start, err := time.Parse(layout, startDate)
	if err != nil {
		panic(err)
	}
	end, err := time.Parse(layout, endDate)
	if err != nil {
		panic(err)
	}

	currentTime := time.Now()
	if hourStr == "" {
		hourStr = currentTime.Format("15:04:05")
	} else {
		parsedTime, err := time.Parse("15:04:05", hourStr)
		if err != nil {
			// fmt.Println("Error parsing time:", err)
			panic(err)
		}
		subtractedTime := parsedTime.Add(-5 * time.Hour)
		outputLayout := "15:04:05"

		hourStr = subtractedTime.Format(outputLayout)
	}

	dayDates := make([]DateObject, 0)

	for !start.After(end) {
		for _, day := range days {
			weekday := getWeekday(day)
			if start.Weekday() == weekday {
				dateStr := start.Format(layout)
				foundDay := false

				// Check if the day already exists in the result slice
				for i, obj := range dayDates {
					if obj.Day == day {
						dayDates[i].Dates = append(dayDates[i].Dates, dateStr)
						foundDay = true
						break
					}
				}

				// If the day doesn't exist in the result slice, add a new DateObject
				if !foundDay {
					dayDates = append(dayDates, DateObject{
						Day:   day,
						Hour:  hourStr, // You can set your desired hour here
						Dates: []string{dateStr},
					})
				}
			}
		}
		start = start.AddDate(0, 0, 1)
	}

	// Add the hour to the dates and format them as "2023-07-11T14:00:00.000Z"
	for i, obj := range dayDates {
		for j, dateStr := range obj.Dates {
			dayDates[i].Dates[j] = dateStr + "T" + obj.Hour + ".000Z"
		}
	}

	return dayDates
}

func getWeekday(day string) time.Weekday {
	switch day {
	case "monday":
		return time.Monday
	case "tuesday":
		return time.Tuesday
	case "wednesday":
		return time.Wednesday
	case "thursday":
		return time.Thursday
	case "friday":
		return time.Friday
	case "saturday":
		return time.Saturday
	case "sunday":
		return time.Sunday
	default:
		panic("Invalid weekday")
	}
}

// Datas This is response struct from create
type Datas struct {
	Data struct {
		Data struct {
			Data map[string]interface{} `json:"data"`
		} `json:"data"`
	} `json:"data"`
}

// ClientApiResponse This is get single api response
type ClientApiResponse struct {
	Data ClientApiData `json:"data"`
}

type ClientApiData struct {
	Data ClientApiResp `json:"data"`
}

type ClientApiResp struct {
	Response map[string]interface{} `json:"response"`
}

type Response struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

// NewRequestBody's Data (map) field will be in this structure
//.   fields
// objects_ids []string
// table_slug string
// object_data map[string]interface
// method string
// app_id string

// but all field will be an interface, you must do type assertion
type DateObject struct {
	Day   string
	Hour  string
	Dates []string
}
type HttpRequest struct {
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Headers http.Header `json:"headers"`
	Params  url.Values  `json:"params"`
	Body    []byte      `json:"body"`
}

type AuthData struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type NewRequestBody struct {
	RequestData HttpRequest            `json:"request_data"`
	Auth        AuthData               `json:"auth"`
	Data        map[string]interface{} `json:"data"`
}
type Request struct {
	Data map[string]interface{} `json:"data"`
}

type RequestMany2Many struct {
	IdFrom    string   `json:"id_from"`
	IdTo      []string `json:"id_to"`
	TableFrom string   `json:"table_from"`
	TableTo   string   `json:"table_to"`
}

type MultipleUpdateRequest struct {
	Data struct {
		Objects []map[string]interface{} `json:"objects"`
	} `json:"data"`
}

// GetListClientApiResponse This is get list api response
type GetListClientApiResponse struct {
	Data GetListClientApiData `json:"data"`
}

type GetListClientApiData struct {
	Data GetListClientApiResp `json:"data"`
}

type GetListClientApiResp struct {
	Response []map[string]interface{} `json:"response"`
}

type CustomDataObj struct {
	CycleName  string   `json:"cycle_name"`
	CycleCount int      `json:"cycle_count"`
	Time       string   `json:"time"`
	Dates      []string `json:"dates"`
}

type DayTime struct {
	Day  string `json:"day"`
	Time string `json:"time"`
}
type DateTime struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

type Medicine struct {
	Type            string        `json:"type"`
	DayData         []string      `json:"dayData"`
	CustomData      CustomDataObj `json:"customData"`
	WeekData        []DayTime     `json:"weekData"`
	MonthData       []DateTime    `json:"monthData"`
	BeforeAfterFood string        `json:"before_after_food"`
	StartDate       string        `json:"start_date"`
	EndDate         string        `json:"end_date"`
	CurrentAmount   int           `json:"current_amount"`
	DaysOfWeek      []int         `json:"days_of_week"`
	HoursOfDay      []string      `json:"hours_of_day"`
	WithoutBreak    bool          `json:"without_break"`
}
