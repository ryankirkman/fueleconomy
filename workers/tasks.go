package workers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/teasherm/fueleconomy/global"
	"github.com/teasherm/fueleconomy/models"
)

type WorkRequest struct {
	Target  string
	Fetcher Fetcher
	Action  func(Fetcher) error
}

func (w *WorkRequest) DoWork() error {
	return w.Action(w.Fetcher)
}

func GenerateWorkRequest(target string) (WorkRequest, error) {
	switch target {
	case "vehicles":
		return WorkRequest{
			Target:  target,
			Fetcher: FileFetcher{},
			Action:  IngestVehicles}, nil
	case "fuelprices":
		return WorkRequest{
			Target:  target,
			Fetcher: RestFetcher{},
			Action:  IngestFuelPrices}, nil
	default:
		return WorkRequest{}, errors.New(fmt.Sprintf("Ingestion target %s not valid", target))
	}
}

func IngestFuelPrices(f Fetcher) error {
	data, err := f.Fetch("https://www.fueleconomy.gov/ws/rest/fuelprices")
	if err != nil {
		return err
	}

	fp := models.FuelPrices{}
	xml.Unmarshal(data, &fp)

	insertedId, err := global.Db.InsertOne("fuel_prices", &fp)
	if err != nil {
		return err
	}
	global.Logger.Println("Fuel Prices Inserted With ID:", insertedId)

	return nil
}

func IngestVehicles(f Fetcher) error {
	data, err := f.Fetch("vehicles")
	if err != nil {
		return err
	}

	var rvo models.RawVehiclesOuter
	err = xml.Unmarshal(data, &rvo)
	if err != nil {
		return err
	}

	insertedIds := make([]int, 0)
	count := 0
	for _, vehicle := range rvo.RawVehicles {
		fv, err := models.NewVehicleFromRaw(&vehicle)
		if err != nil {
			return err
		}
		insertedId, err := global.Db.UpsertOne("vehicles", "epa_id", fv)
		if err != nil {
			return err
		}
		if insertedId > 0 {
			insertedIds = append(insertedIds, insertedId)
		}
		count++
	}
	inserted := len(insertedIds)
	updated := count - len(insertedIds)
	global.Logger.Println("Vehicle Updates:", updated)
	global.Logger.Println("Vehicle Inserts:", inserted)
	err = ingestEmissionsInfo(f)
	if err != nil {
		return err
	}
	return nil
}

func ingestEmissionsInfo(f Fetcher) error {
	data, err := f.Fetch("emissions")
	if err != nil {
		return err
	}

	var rvo models.RawEmissionsInfoOuter
	err = xml.Unmarshal(data, &rvo)
	if err != nil {
		return err
	}

	global.Db.DeleteAll("emissions_info")
	insertedIds := make([]int, 0)
	violatingIds := make(map[int]int)
	for _, emissionsInfo := range rvo.RawEmissionsInfoes {
		ei, err := models.NewEmissionsInfoFromRaw(&emissionsInfo)
		insertedId, err := global.Db.InsertOne("emissions_info", ei)
		if err != nil {
			if strings.Contains(err.Error(), "violates foreign key constraint") {
				if count, ok := violatingIds[ei.EpaID]; ok {
					violatingIds[ei.EpaID] = count + 1
				} else {
					violatingIds[ei.EpaID] = 1
				}
				continue
			} else {
				return err
			}
		}
		if insertedId > 0 {
			insertedIds = append(insertedIds, insertedId)
		}
	}
	inserted := len(insertedIds)
	global.Logger.Println("Emissions Info Inserts:", inserted)
	global.Logger.Println("Emissions Info Foreign Key Violations[id:count]:", violatingIds)

	return nil
}
