package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"time"

	"cloud.google.com/go/civil"
	log "github.com/sirupsen/logrus"
)

type FlightSlide struct {
	// Set of flights/days that the user wants to track
	TrackedFlights []FlightAndDay
	// Flight currently being displayed/requested
	ActiveFlight FlightAndDay

	HttpHelper   *HttpHelper
	RedrawTicker *time.Ticker
	DisplayData  FlightDisplayData
}

type FlightAndDay struct {
	Id   string
	Date civil.Date
}

func NewFlightSlide(flights map[string]string) *FlightSlide {
	sl := new(FlightSlide)
	sl.HttpHelper = NewHttpHelper(HttpConfig{
		SlideId:            "FlightSlide",
		RefreshInterval:    10 * time.Minute,
		RequestUrlCallback: sl.BuildRequest,
		ParseCallback:      sl.Parse,
	})

	// Copy the flights from the input into a struct (to be used later)
	// Expected input is date (e.g. "2020-01-15") to flight (e.g AA 1234).
	for date, id := range flights {
		dateStruct, err := civil.ParseDate(date)
		if err != nil {
			log.WithFields(log.Fields{
				"date": date,
			}).Warn("Could not parse date input on flight slide.")
			continue
		}
		sl.TrackedFlights = append(sl.TrackedFlights, FlightAndDay{
			Id:   id,
			Date: dateStruct,
		})
	}

	return sl
}

func (sl *FlightSlide) Initialize() {
	// Get the flight to focus on
	flight, ok := sl.GetActiveFlight()
	log.WithFields(log.Fields{
		"ok":     ok,
		"flight": flight,
	}).Debug("Chose flight")
	// If no active flight, return immediately
	if !ok {
		return
	}
	sl.ActiveFlight = flight

	// For development, always read data from file instead of live API call
	/*f := "flight_sample.json"
	  data, err := ioutil.ReadFile(f)
	  if err != nil {
	      log.WithFields(log.Fields{
	          "file":  f,
	          "error": err,
	      }).Warn("Could not open test JSON.")
	      return
	  }
	  sl.Parse(data)*/

	// Start fetching data
	sl.HttpHelper.StartLoop()
}

func (sl *FlightSlide) Terminate() {
	// Fetching might already be stopped if flight is complete
	sl.HttpHelper.StopLoop()
}

func (sl *FlightSlide) StartDraw(d Display) {
	sl.RedrawTicker = DrawEverySecond(d, sl.Draw)
}

func (sl *FlightSlide) StopDraw() {
	sl.RedrawTicker.Stop()
}

func (sl *FlightSlide) IsEnabled() bool {
	// Slide should be enabled if there is an active flight
	_, ok := sl.GetActiveFlight()
	return ok
}

func (sl *FlightSlide) BuildRequest() (*http.Request, error) {
	url := fmt.Sprintf(
		"http://flightxml.flightaware.com/json/FlightXML3/FlightInfoStatus?ident=%s&howMany=5",
		sl.ActiveFlight.Id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(FLIGHTAWARE_USERNAME, FLIGHTAWARE_API_KEY)
	return req, nil
}

func (sl *FlightSlide) Parse(respBytes []byte) bool {
	var respData FlightInfoStatusResponse
	jsonErr := json.Unmarshal(respBytes, &respData)
	if jsonErr != nil {
		log.WithFields(log.Fields{
			"error": jsonErr,
			"data":  string(respBytes),
		}).Warn("Could not interpret flights JSON.")
		return false
	}

	displayData := FlightDisplayData{}

	// First we find the flight that was meant to depart today
	var targetFlight FlightInfoStatus
	for i := range respData.Result.Flights {
		f := respData.Result.Flights[i]

		depDate := civil.DateOf(time.Unix(f.FiledDepartureTime.LocalTime, 0))
		if depDate == civil.DateOf(time.Now()) {
			targetFlight = f
			break
		}
	}

	// If no flight was found, nothing to display
	if targetFlight == (FlightInfoStatus{}) {
		log.Info("No matching flight found in response")
		sl.DisplayData = displayData
		return false
	}

	displayData.Title = fmt.Sprintf("%s %s",
		targetFlight.AirlineIata, targetFlight.FlightNumber)
	// Use the IATA code since it's easier to read
	displayData.Origin = targetFlight.Origin.AlternateIdent
	displayData.Destination = targetFlight.Destination.AlternateIdent

	// Departure stats
	if targetFlight.EstimatedDepartureTime == (FlightInfoTime{}) {
		// If no estimated time, fall back to filed time (should always exist)
		displayData.DepartureTime = time.Unix(targetFlight.FiledDepartureTime.LocalTime, 0)
	} else if targetFlight.ActualDepartureTime == (FlightInfoTime{}) {
		// If no actual time, fall back to estimated time
		displayData.DepartureTime = time.Unix(targetFlight.EstimatedDepartureTime.LocalTime, 0)
	} else {
		// Otherwise the actual time should be safe to use
		displayData.DepartureTime = time.Unix(targetFlight.ActualDepartureTime.LocalTime, 0)
		displayData.HasDeparted = true
	}
	displayData.DepartureDelay = time.Duration(targetFlight.DepartureDelay) * time.Second

	// Arrival stats
	if targetFlight.EstimatedArrivalTime == (FlightInfoTime{}) {
		// If no estimated time, fall back to filed time (should always exist)
		displayData.ArrivalTime = time.Unix(targetFlight.FiledArrivalTime.LocalTime, 0)
	} else if targetFlight.ActualArrivalTime == (FlightInfoTime{}) {
		// If no actual time, fall back to estimated time
		displayData.ArrivalTime = time.Unix(targetFlight.EstimatedArrivalTime.LocalTime, 0)
	} else {
		// Otherwise the actual time should be safe to use
		displayData.ArrivalTime = time.Unix(targetFlight.ActualArrivalTime.LocalTime, 0)
		displayData.HasArrived = true

		// If the flight has arrived, stop requesting new data
		log.Info("Flight has arrived, stopping HTTP fetcher")
		sl.HttpHelper.StopLoop()
	}
	displayData.ArrivalDelay = time.Duration(targetFlight.ArrivalDelay) * time.Second

	log.WithFields(log.Fields{"data": displayData}).Debug("Parsed display data")

	sl.DisplayData = displayData
	return true
}

func (sl *FlightSlide) GetActiveFlight() (FlightAndDay, bool) {
	for i := range sl.TrackedFlights {
		if sl.TrackedFlights[i].Date == civil.DateOf(time.Now()) {
			return sl.TrackedFlights[i], true
		}
	}
	return FlightAndDay{}, false
}

func (sl *FlightSlide) GetDurationString(d time.Duration) string {
	if d.Hours() >= 1.0 {
		dm := d.Round(time.Minute)
		h := dm / time.Hour
		dm -= h * time.Hour
		m := dm / time.Minute
		return fmt.Sprintf("%d:%02d", h, m)
	}
	return fmt.Sprintf("%d Min", int(d.Minutes()))
}

func (sl *FlightSlide) Draw(img *image.RGBA) {
	if !sl.HttpHelper.LastFetchSuccess {
		DrawError(img, "Flight Status", "Connection error.")
		return
	}

	if sl.DisplayData == (FlightDisplayData{}) {
		DrawError(img, "Flight Status", "No data.")
		return
	}

	aqua := color.RGBA{0, 255, 255, 255}
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	// Show flight ID on top line
	WriteString(img, sl.DisplayData.Title, aqua, ALIGN_CENTER, 64, 0)

	// Draw origin/destination boxes on sides
	ow := GetDisplayWidth(sl.DisplayData.Origin)
	DrawBox(img, aqua, 0, 11, ow+4, 9)
	WriteString(img, sl.DisplayData.Origin, black, ALIGN_LEFT, 2, 12)

	dw := GetDisplayWidth(sl.DisplayData.Destination)
	DrawBox(img, aqua, 128-dw-4, 11, dw+4, 9)
	WriteString(img, sl.DisplayData.Destination, black, ALIGN_RIGHT, 125, 12)

	// Timing status
	status := "On Time"
	statusColor := color.RGBA{0, 255, 0, 255}
	if !sl.DisplayData.HasDeparted {
		if sl.DisplayData.DepartureDelay > 0 {
			status = fmt.Sprintf("%s Late", sl.GetDurationString(sl.DisplayData.DepartureDelay))
			statusColor = color.RGBA{255, 255, 0, 255}
		}
	} else if !sl.DisplayData.HasArrived {
		if sl.DisplayData.ArrivalDelay > 0 {
			status = fmt.Sprintf("%s Late", sl.GetDurationString(sl.DisplayData.ArrivalDelay))
			statusColor = color.RGBA{255, 255, 0, 255}
		}
	} else {
		status = "Arrived"
	}
	WriteString(img, status, statusColor, ALIGN_CENTER, 64, 8)

	// Departure
	depPrefix := "Dep. "
	if !sl.DisplayData.HasDeparted {
		depPrefix = "Est. Dep. "
	}
	WriteString(img, depPrefix+sl.DisplayData.DepartureTime.Format("3:04 PM"), white, ALIGN_CENTER, 64, 16)

	// Arrival
	arrPrefix := "Arr. "
	if !sl.DisplayData.HasArrived {
		arrPrefix = "Est. Arr. "
	}
	WriteString(img, arrPrefix+sl.DisplayData.ArrivalTime.Format("3:04 PM"), white, ALIGN_CENTER, 64, 24)
}

// Data structures used by the FlightAware v3 API
type FlightInfoStatusResponse struct {
	Result FlightInfoStatusResult `json:"FlightInfoStatusResult"`
}

type FlightInfoStatusResult struct {
	Flights []FlightInfoStatus `json:"flights"`
}

type FlightInfoStatus struct {
	Airline                string             `json:"airline"`
	AirlineIata            string             `json:"airline_iata"`
	FlightNumber           string             `json:"flightnumber"`
	Blocked                bool               `json:"blocked"`
	Diverted               bool               `json:"diverted"`
	Cancelled              bool               `json:"cancelled"`
	Origin                 FlightInfoLocation `json:"origin"`
	Destination            FlightInfoLocation `json:"destination"`
	FiledDepartureTime     FlightInfoTime     `json:"filed_departure_time"`
	EstimatedDepartureTime FlightInfoTime     `json:"estimated_departure_time"`
	ActualDepartureTime    FlightInfoTime     `json:"actual_departure_time"`
	DepartureDelay         int                `json:"departure_delay"`
	FiledArrivalTime       FlightInfoTime     `json:"filed_arrival_time"`
	EstimatedArrivalTime   FlightInfoTime     `json:"estimated_arrival_time"`
	ActualArrivalTime      FlightInfoTime     `json:"actual_arrival_time"`
	ArrivalDelay           int                `json:"arrival_delay"`
}

type FlightInfoLocation struct {
	Code           string `json:"code"`
	City           string `json:"city"`
	AlternateIdent string `json:"alternate_ident"`
	AirportName    string `json:"airport_name"`
}

type FlightInfoTime struct {
	LocalTime int64  `json:"epoch"`
	TimeZone  string `json:"tz"`
}

// Internal representation of what to draw on the slide
type FlightDisplayData struct {
	Title          string
	Origin         string
	Destination    string
	HasDeparted    bool
	HasArrived     bool
	DepartureTime  time.Time
	ArrivalTime    time.Time
	DepartureDelay time.Duration
	ArrivalDelay   time.Duration
}
