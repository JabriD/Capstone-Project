# Capstone-Project
GOAL

Building a Database: 2022 NBA Roster

This project involves creating a Cient that will make API calls to the NBA stat server and present that data back to the user in their desired format.
The NBA already offers a variety of data filters, but this executable would allow the User to filter/manipulate data in a more customized manner.
First, the server will tap a given NBA Statistics endpoint, retrieve the JSON response, and unmarshal it.
Next, the retrieved data will be used to populate a SQL database made of two tables (Teams/Players).
Players will be associated with their respective team's via their Team ID.

METHODOLOGY

Using the two active endpoints (standings/players) to compose a full active 2022 NBA Roster.
Standings includes: team names, cities, and ID's.
Players includes: Player name, height, weight, position, DOB, jersey number, player ID, Team ID, draft info, and other general career/professional details.
The JSON retrieved from these two endpoints will be consolidated into a single database, and made available for query.

USER STORIES

A Coach can use this roster to review and inform gameplan decisions.
A Journalist can use this roster to get player or team data to inform/develop articles and storylines.
An Admin can use this roster as a starting point to add more data in the future such as individual player statistics.
