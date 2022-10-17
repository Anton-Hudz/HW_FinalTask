package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
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
	castingProblemErr           = errors.New("problem with casting")
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

	fmt.Print("Введите критерий, по которому нужно отсортировать поезда в результате. \nВалидные значения: price, arrival-time, departure-time: ")
	criteria, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	criteria = strings.TrimSpace(criteria)

	//result,
	result, err := FindTrains(departureStation, arrivalStation, criteria)
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
	const maximumNumberOfTrains = 3

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

	contentOfDataInByte, err := os.ReadFile("data.json")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(contentOfDataInByte, &unmarshaledJsonDatas); err != nil {
		return nil, err
	}

	for _, v := range unmarshaledJsonDatas {
		if v.DepartureStationID == departureStationToInt && v.ArrivalStationID == arrivalStationToInt {
			choosedTrains = append(choosedTrains, v)
		}
	}

	switch criteria {
	case "price":
		sort.SliceStable(choosedTrains, func(i, j int) bool {
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
	const timeStandart = "15:04:00"
	var jsonTypeToTaskType map[string]interface{}

	if err := json.Unmarshal(data, &jsonTypeToTaskType); err != nil {
		return err
	}

	for lineOfStruct, v := range jsonTypeToTaskType {
		switch lineOfStruct {
		case "trainId":
			trainId, ok := v.(float64)
			if !ok {
				return castingProblemErr
			}
			t.TrainID = int(trainId)
		case "departureStationId":
			departureStationID, ok := v.(float64)
			if !ok {
				return castingProblemErr
			}
			t.DepartureStationID = int(departureStationID)
		case "arrivalStationId":
			arrivalStationId, ok := v.(float64)
			if !ok {
				return castingProblemErr
			}
			t.ArrivalStationID = int(arrivalStationId)
		case "price":
			price, ok := v.(float64)
			if !ok {
				return castingProblemErr
			}
			t.Price = float32(price)
		case "arrivalTime":
			arrivalTime, ok := v.(string)
			if !ok {
				return castingProblemErr
			}
			time, err := time.Parse(timeStandart, arrivalTime)
			if err != nil {
				return err
			}
			t.ArrivalTime = time
		case "departureTime":
			departureTime, ok := v.(string)
			if !ok {
				return castingProblemErr
			}
			time, err := time.Parse(timeStandart, departureTime)
			if err != nil {
				return err
			}
			t.DepartureTime = time
		}
	}

	return nil
}
