package helpers

import "time"

const defaultDateLayout = "2006-01-02"

func ParseStringToDate(value string) (time.Time, error){

	parsedDate, err := time.Parse(defaultDateLayout, value)
	if err != nil{
		return time.Time{}, err
	}

	return parsedDate, nil
}

func ParseDateToString(date time.Time) (string, error){
	return date.Format(defaultDateLayout), nil
}