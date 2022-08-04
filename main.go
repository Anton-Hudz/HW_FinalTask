package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	// "io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	unsupportedCriteriaErr      = errors.New("unsupported criteria")
	emptyDepartureStationErr    = errors.New("empty departure station")
	emptyArrivalStationErr      = errors.New("empty arrival station")
	badArrivalStationInputErr   = errors.New("bad arrival station input")
	badDepartureStationInputErr = errors.New("bad departure station input")
)

type Trains []Train
type Train struct {
	TrainID            int       `json:"trainId"`
	DepartureStationID int       `json:"departureStationId"`
	ArrivalStationID   int       `json:"arrivalStationId"`
	Price              float32   `json:"price"`
	ArrivalTime        time.Time `json:"arrivalTime"`
	DepartureTime      time.Time `json:"departureTime"`
}

func main() {
	var (
		departureStation string
		arrivalStation   string
		criteria         string
	)

	//	... запит даних від користувача
	fmt.Print("Введите станцию отправления: ")
	departureStation, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	departureStation = strings.TrimSpace(departureStation)

	fmt.Print("Введите станцию прибытия: ")
	arrivalStation, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	arrivalStation = strings.TrimSpace(arrivalStation)

	fmt.Print("Введите критерий, по которому нужно отсортировать поезда в результате. Валидные значения: price, arrival-time, departure-time: ")
	criteria, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	criteria = strings.TrimSpace(criteria)

	//result, err добавить ошибку
	result, err := FindTrains(departureStation, arrivalStation, criteria)

	//	... обробка помилки
	if err != nil {
		fmt.Println(err)
		return
	}

	//	... друк result
	for _, i := range result {
		fmt.Printf("%+v\n", i)
	}
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	// ... код
	var (
		choosedTrains        Trains
		unmarshaledJsonDatas Trains
	)
	maximumNumberOfTrains := 3

	if departureStation == "" {
		return nil, emptyDepartureStationErr
	}

	departureStationToInt, err := strconv.Atoi(departureStation)
	if err != nil {
		return nil, badDepartureStationInputErr
	}
	if arrivalStation == "" {
		return nil, emptyArrivalStationErr
	}

	arrivalStationToInt, err := strconv.Atoi(arrivalStation)
	if err != nil {
		return nil, badArrivalStationInputErr
	}

	contentOfDataInByte, _ := os.ReadFile("data.json")
	err = json.Unmarshal(contentOfDataInByte, &unmarshaledJsonDatas)
	if err != nil {
		return nil, err
	}

	for _, v := range unmarshaledJsonDatas {
		if v.DepartureStationID == departureStationToInt && v.ArrivalStationID == arrivalStationToInt {
			choosedTrains = append(choosedTrains, v)
		}
	}

	switch criteria {
	case "price":
		sort.Slice(choosedTrains, func(i, j int) bool {
			return choosedTrains[j].Price > choosedTrains[i].Price
		})
	case "arrival-time":
		sort.SliceStable(choosedTrains, func(i, j int) bool {
			return choosedTrains[j].ArrivalTime.After(choosedTrains[i].ArrivalTime)
		})
	case "departure-time":
		sort.SliceStable(choosedTrains, func(i, j int) bool {
			return choosedTrains[j].DepartureTime.After(choosedTrains[i].DepartureTime)
		})
	default:
		return nil, unsupportedCriteriaErr
	}

	if len(choosedTrains) > maximumNumberOfTrains {
		choosedTrains = choosedTrains[:maximumNumberOfTrains]
	}

	if len(choosedTrains) == 0 {
		return nil, nil
	}
	return choosedTrains, nil
}

func (t *Train) UnmarshalJSON(data []byte) error {
	var jsonTypeToTaskType map[string]interface{}
	err := json.Unmarshal(data, &jsonTypeToTaskType)

	if err != nil {
		return err
	}

	for line, v := range jsonTypeToTaskType {
		switch line {
		case "trainId":
			t.TrainID = int(v.(float64))
		case "departureStationId":
			t.DepartureStationID = int(v.(float64))
		case "arrivalStationId":
			t.ArrivalStationID = int(v.(float64))
		case "price":
			t.Price = float32(v.(float64))
		case "arrivalTime":
			time, err := time.Parse("15:04:00", v.(string))
			if err != nil {
				return err
			}
			t.ArrivalTime = time
		case "departureTime":
			time, err := time.Parse("15:04:00", v.(string))
			if err != nil {
				return err
			}
			t.DepartureTime = time
		}
	}
	return nil
}
