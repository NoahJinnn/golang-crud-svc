package mqttSvc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/models"
)

type UserIDPassword struct {
	UserId     string `json:"user_id"`
	RfidPass   string `json:"rfid_pw"`
	KeypadPass string `json:"keypad_pw"`
}

func ServerCreateDoorlockPayload(doorlock *models.Doorlock) string {
	msg := fmt.Sprintf(`{"doorlock_id":"%s"}`, doorlock.DoorSerialID)
	return PayloadWithGatewayId(doorlock.GatewayID, msg)
}

func ServerUpdateDoorlockPayload(doorlock *models.Doorlock) string {
	msg := fmt.Sprintf(`{"doorlock_id":"%s","description":"%s","location":"%s", "state":"%s"}`,
		doorlock.DoorSerialID, doorlock.Description, doorlock.Location, doorlock.State)
	return PayloadWithGatewayId(doorlock.GatewayID, msg)
}

func ServerDeleteDoorlockPayload(doorlock *models.DoorlockDelete) string {
	msg := fmt.Sprintf(`{"doorlock_id":"%s"}`, doorlock.DoorSerialID)
	return PayloadWithGatewayId(doorlock.GatewayID, msg)
}

func ServerCmdDoorlockPayload(gwId string, doorSerialId string, cmd string) string {
	msg := fmt.Sprintf(`{"doorlock_id":"%s","action":"%s"}`, doorSerialId, cmd)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerUpdateGatewayPayload(gw *models.Gateway) string {
	return fmt.Sprintf(`{"gateway_id":"%d","area_id":"%d","name":"%s", "state":"%s"}`,
		gw.ID, gw.AreaID, gw.Name, gw.State)
}

func ServerCreateRegisterPayload(
	gwId string,
	doorId string,
	sche *models.Scheduler,
	uP *UserIDPassword,
) string {

	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	startDmySlice := getDayMonthYearSlice(sche.StartDate)
	start := time.Date(startDmySlice[2], time.Month(startDmySlice[1]), startDmySlice[0], 0, 0, 0, 0, loc).Unix()
	endDmySlice := getDayMonthYearSlice(sche.EndDate)
	end := time.Date(endDmySlice[2], time.Month(endDmySlice[1]), endDmySlice[0], 23, 59, 59, 0, loc).Unix()

	msg := fmt.Sprintf(`{"register_id":"%d",
	"user_id":"%s",
	"doorlock_id":"%s",
	"rfid_pw":"%s",
	"keypad_pw":"%s",
	"start_date":"%d",
	"end_date":"%d",
	"week_day":"%d",
	"start_class":"%d",
	"end_class":"%d"}`,
		sche.ID, uP.UserId, doorId, uP.RfidPass, uP.KeypadPass,
		start, end, sche.WeekDay, sche.StartClassTime, sche.EndClassTime)

	return PayloadWithGatewayId(gwId, msg)
}

func ServerUpdateRegisterPayload(gwId string, uSche *models.UpdateScheduler) string {
	sche := uSche.Scheduler
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	startDmySlice := getDayMonthYearSlice(sche.StartDate)
	start := time.Date(startDmySlice[2], time.Month(startDmySlice[1]), startDmySlice[0], 0, 0, 0, 0, loc).Unix()
	endDmySlice := getDayMonthYearSlice(sche.EndDate)

	end := time.Date(endDmySlice[2], time.Month(endDmySlice[1]), endDmySlice[0], 23, 59, 59, 0, loc).Unix()
	msg := fmt.Sprintf(`{"register_id":"%d",
	"user_id":"%s",
	"doorlock_id":"%s",
	"start_date":"%d",
	"end_date":"%d",
	"week_day":"%d",
	"start_class":"%d",
	"end_class":"%d"}`,
		sche.ID, uSche.UserID, uSche.DoorSerialID,
		start, end, sche.WeekDay, sche.StartClassTime, sche.EndClassTime)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerDeleteRegisterPayload(gwId string, registerId uint) string {
	msg := fmt.Sprintf(`{"register_id":"%d"}`, registerId)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerBootuptHPEmployeePayload(gwId string, emps []models.Employee) string {
	bootupEmps := []UserIDPassword{}
	for _, emp := range emps {
		buEmp := UserIDPassword{
			UserId:     emp.MSNV,
			RfidPass:   emp.RfidPass,
			KeypadPass: emp.KeypadPass,
		}
		bootupEmps = append(bootupEmps, buEmp)
	}
	bootupEmpsJson, _ := json.Marshal(bootupEmps)
	return PayloadWithGatewayId(gwId, string(bootupEmpsJson))
}

func ServerUpdateUserPayload(gwId string, userId string, rfidPw string, keypadPw string) string {
	msg := fmt.Sprintf(`{"user_id":"%s","rfid_pw":"%s", "keypad_pw":"%s"}`,
		userId, rfidPw, keypadPw)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerDeleteUserPayload(gwId string, msnv string) string {
	msg := fmt.Sprintf(`{"user_id":"%s"}`, msnv)
	return PayloadWithGatewayId(gwId, msg)
}

func PayloadWithGatewayId(gwId string, msg string) string {
	return fmt.Sprintf(`{"gateway_id":"%s","message":%s}`, gwId, msg)
}

func getDayMonthYearSlice(str string) []int {
	strs := strings.Split(str, "/")
	var dmySlice = []int{}
	for _, s := range strs {
		number, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			panic(err)
		}
		dmySlice = append(dmySlice, int(number))
	}
	return dmySlice
}
