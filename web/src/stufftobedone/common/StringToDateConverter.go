package common

import (
    "time"
    "strconv"
)

type ConversionResult struct {
    DateAsString string
    DateAsInt int
    Date time.Time
}

func ConvertIntToDisplay(dateAsInt int) string {
    dayAsString := strconv.Itoa(dateAsInt)
    
    return dayAsString
}

func ConvertStringToDates(dateAsString string) (ConversionResult, error) {
    result := ConversionResult{}
    result.DateAsString = dateAsString

    dateAsInt, err := strconv.Atoi(dateAsString)
    if err != nil {
        return result, nil
    }
    result.DateAsInt = dateAsInt

    date, err := time.Parse("20060102", dateAsString)
    if err != nil {
        return result, nil
    }
    result.Date = date

    return result, nil
}

func (d* ConversionResult) Tomorrow() (ConversionResult, error) {
    return ConvertStringToDates(d.Date.AddDate(0,0,1).Format("20060102"))
}

func (d* ConversionResult) Yesterday() (ConversionResult, error) {
    return ConvertStringToDates(d.Date.AddDate(0,0,-1).Format("20060102"))
}