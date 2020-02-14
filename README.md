# NcaaBballStats
The source code for ncaabballstats.com, an API to grab information from sports-reference regarding college basketball. Scrapes sports-reference cbb section and returns data as JSON

## Setup

To run locally, just run `go run main.go` which will start up a local server at `localhost:8080`

## Notes

The API currently only grabs the per_game stats for each player taken from the college home page.

## Example

In order to get team data, team name must match what's given in sports-reference i.e. University of Texas at Austin is just texas but University of Nevada - Las Vegas is nevada-las-vegas. 

Once you have the name just pass in the year and the string `pergame` to get data.

Example URL : `localhost:8080/maryland-baltimore-county/2018/pergame`

Example Retun: 
```json
{
  "StatusCode": 200,
  "Response": [
    {
      "Rank": 1,
      "Player": "Jairus Lyles",
      "Games": 33,
      "Games Started": 32,
      "Minutes Played Per Game": 34.9,
      "Field Goals Per Game": 6.8,
      "Field Goal Attempts Per Game": 15.5,
      "Field Goal Percentage": 0.439,
      "2-Point Field Goals Per Game": 4.5,
      "2-Point Field Goal Attempts Per Game": 9.5,
      "2-Point Field Goal Percentage": 0.471,
      "3-Point Field Goals Per Game": 2.4,
      "3-Point Field Goal Attempts Per Game": 6.1,
      "3-Point Field Goal Percentage": 0.390,
      "Free Throws Per Game": 4.2,
      "Free Throw Attempts Per Game": 5.2,
      "Free Throw Percentage": 0.792,
      "Offensive Rebounds Per Game": 0.7,
      "Defensive Rebounds Per Game": 4.8,
      "Total Rebounds Per Game": 5.5,
      "Assists Per Game": 3.5,
      "Steals Per Game": 2.1,
      "Blocks Per Game": 0.2,
      "Turnovers Per Game": 3.1,
      "Personal Fouls Per Game": 2.0,
      "Points Per Game": 20.2
    }
  ]
}
```

## Credits

Was originally using `nivrrex/excel` package to go from xls to csv but no longer using

Took `Ahmad-Magdy/CSV-To-JSON-Converter` and turned it into a pacakge

## TODO

* Create API Documentation

* Add more endpoints to grab other data from team page

* Add query parameters to filter by player

* Convert directly from HTML table to JSON data

* Add error codes for when user hits a non-existant URL route
