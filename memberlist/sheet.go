package memberlist

import (
	"context"
	"encoding/base64"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// https://docs.google.com/spreadsheets/d/<SPREADSHEETID>/edit#gid=<SHEETID>
const (
	MEMBERLIST_SPREADSHEET_ID   = "10vC_oi6rgBmVqJKgymokWobIvXOiP8yLx9F4sgfT994"
	MEMBERLIST_SHEED_ID         = "0"
	MEMBERLIST_SHEET_READ_RANGE = "A2:G200"
)

var sheetInstance *sheets.Service

func init() {
	creds, err := base64.StdEncoding.DecodeString(os.Getenv("GOOGLE_KEY_JSON_BASE64"))
	if err != nil {
		panic(err)
	}

	config, err := google.JWTConfigFromJSON(creds, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		panic(err)
	}

	client := config.Client(context.TODO())
	sheetInstance, err = sheets.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		panic(err)
	}
}

func GetMemberlistSheet() (*sheets.ValueRange, error) {
	resp, err := sheetInstance.Spreadsheets.Values.Get(MEMBERLIST_SPREADSHEET_ID, MEMBERLIST_SHEET_READ_RANGE).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func UpdateMemberlistSheet(m *Memberlist) error {
	var sheetValues [][]interface{}
	for _, v := range m.Members {
		sheetValues = append(sheetValues,
			[]interface{}{
				v.Uuid,
				v.Name,
				v.DiscordID,
				v.TeamSpeakID,
				v.Accounts.XLPC,
				v.Accounts.LPC,
				v.Rank,
			})
	}

	_, err := sheetInstance.Spreadsheets.Values.Update(MEMBERLIST_SPREADSHEET_ID, MEMBERLIST_SHEET_READ_RANGE, &sheets.ValueRange{
		Values: sheetValues,
	}).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		return err
	}

	return nil
}
