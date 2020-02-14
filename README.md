# NcaaBballStats
The source code for ncaabballstats.com, an API to grab information from sports-reference regarding college basketball. Scrapes sports-reference cbb section and returns data as JSON

## Background

I wanted a way to grab the data from sports-reference cbb section so that I could try my hand at data analytics / big data / ML / cool data stuff. As March Madness was about a month away at the start of this project, I thought it would give me plenty of time to practice and then seeing what I've learned by prediciting the winner of the tournament. 

However, there actually isn't any API that gives you the information that you need. The `roclark/sportsreference` package on github looks amazing but requires that you use python :(. So I decided to make a REST API that would allow for more parity.

## Setup

To run locally, just run `go run cmd/ncaabballstats/main.go` which will start up a local server at `localhost:8080`

## Notes

* The API currently only grabs the per_game stats for each player taken from the college home page

* You'll notice that when you're running locally, it creates a JSON and a CSV. The way the API currently works is that it transforms the HTML table to a CSV and then the CSV to a JSON and serves the JSON. The thinking behind this originally was to put the CSV data in storage / data lake to have a good idea of what teams people are currently interested in. As this currently runs on Cloud Run, there's no real issue of storage running out (that I'm aware of)

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

* Add test cases

* Create API Documentation

* Add more endpoints to grab other data from team page

* Add query parameters to filter by player

* Convert directly from HTML table to JSON data

* Add error codes for when user hits a non-existant URL route

* Consider moving from REST API implementation to a GraphQL one


