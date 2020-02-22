**Per Game Stats**
----
  Returns json data containing all per game stats for players that played at the school that year.

* **URL**

  /:team/:year/pergame

* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   `team=[string]`
   `year=[integer]`

* **Data Params**

  None

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** 
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
 
* **Error Response:**

  * **Code:** 404 NOT FOUND <br />
    **Content:** 
    ```json
                  {
                    "StatusCode": 404,
                    "Response": "This endpoint does not exist on this API. Please refer to the API documentation."
                  }
    ```

  OR

  * **Code:** 500 INTERNAL ERROR MESSAGE <br />
    **Content:** 
    ```json
                  {
                    "StatusCode": 500,
                    "Response": "Internal Error. Please contact administrator for more details."
                  }
    ```

  OR

  * **Code:** 400 BAD REQUEST <br />
    **Content:** 
    ```json
                  {
                    "StatusCode": 400,
                    "Response": "Please input a 4 digit year. Refer to the API documentation if you need more assistance."
                  }
    ```

  OR

  * **Code:** 404 NOT FOUND <br />
    **Content:** 
    ```json
                  {
                    "StatusCode": 404,
                    "Response": "There are no stats for the team, year combination given. Please check and make sure team and year were input correctly."
                  }
    ```



* **Sample Call:**

  ```javascript
    $.ajax({
      url: "/virginia/2018/pergame",
      dataType: "json",
      type : "GET",
      success : function(r) {
        console.log(r);
      }
    });
  ```
