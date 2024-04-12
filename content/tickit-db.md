---
title: Mastering SQL Queries with the Tickit Database on DuckDB
description: Dive into the SQL with the Tickit database, a comprehensive dataset for a fictional ticket sales company. This guide will explore the nuances of SQL operations, utilizing DuckDB's efficient analytical data management system. From querying event and venue data to analyzing ticket sales and customer interactions, enhance your SQL skills through practical scenarios and examples.
date: 2024-04-11T12:00:00.000Z
published: true
coverimage: 
ogImage: 
author: Senthilnathan Karuppaiah
avatar: 
type: Blog
tags: [DuckDB, SQL, Data Analysis, Query Optimization]
---

# TickitDB

The ***Tickit*** database contains tables related to a fictional ticket-selling company, originally created by AWS for ***Amazon Redshift*** to demonstrate the capabilities of Redshift for educational purposes. I have ported it to DuckDB to facilitate learning. This database includes data about events, venues, ticket sales, and customer information. It is designed to showcase various SQL operations and query scenarios.

By integrating this dataset with DuckDB Studio, complete with pre-configured queries and an immediately usable database, we aim to streamline the learning process, making it more efficient and engaging for budding data engineers and SQL enthusiasts who are looking to sharpen their data manipulation skills.


![TickitDB](/docs/tickitdb.png)

## Duckdb Practice SQL Queries

### Top 5 States with the Most Event Venues

This query ranks the top five states with the highest number of event venues, indicating their popularity for hosting events. It retrieves and counts venue entries from the Tickit database, grouped by the state in which each venue is located. The results are ordered by the total number of venues descending, highlighting the states that are most active in the event industry. This analysis helps identify key regions for event planning and potential market expansion.

```sql
SELECT
  "venue"."venuestate" AS "venuestate",
  COUNT("venue"."id") AS "total"
FROM "tickit"."venue" AS "venue"
GROUP BY
  "venue"."venuestate"
ORDER BY
  "total" DESC,
  "venuestate"
LIMIT 5
```

### Top 10 Least Popular Events Based on Ticket Sales

This query identifies the ten events with the lowest ticket sales, revealing insights into event popularity. It aggregates total tickets sold by event and category from the Tickit database, using correct joins between the event, category, and sale tables. By grouping the results by event and category names, and ordering by the sum of tickets sold in ascending order, the query highlights events that might require additional marketing efforts or review of their appeal to increase future attendance.

```sql
SELECT
  "event"."eventname",
  "category"."catname",
  SUM("sale"."qtysold") AS "total_tickets"
FROM "tickit"."event"
JOIN "tickit"."category"
  ON "category"."id" = "event"."id"
JOIN "tickit"."sale"
  ON "event"."id" = "sale"."id"
GROUP BY
  "eventname",
  "catname"
ORDER BY
  "total_tickets" DESC
LIMIT 10
```

### Inventory Aging Analysis for Event Listings

This query analyzes the time elapsed between the listing and sale of tickets for events, referred to as inventory aging. It utilizes the DuckDB date formatting and calculations. The query extracts data from the venue, category, event, sale, and listing tables within the Tickit database. Specific attention is paid to converting date and time strings into a standardized format using STRPTIME() function, facilitating precise calculations of the duration (in hours) between the list time and sale time. Results include venue name, category name, event name, precise listing and sale times, and the aging of the inventory in hours. Sorted to emphasize listings with the longest durations before sale, this analysis aids in identifying slow-moving inventory and improving event management strategies. The output presents the top 20 listings with the highest inventory aging, providing valuable insights for optimizing sales efficiency.

```sql
SELECT
  "venue"."venuename" AS "venuename",
  "category"."catname" AS "catname",
  "event"."eventname" AS "eventname",
  STRPTIME("tickit"."sale"."saletime", '%m/%d/%Y %H:%M:%S') AS "saletime",
  STRPTIME("tickit"."listing"."listtime", '%Y-%m-%d %H:%M:%S') AS "listtime",
  STRPTIME("tickit"."sale"."saletime", '%m/%d/%Y %H:%M:%S') - STRPTIME("tickit"."listing"."listtime", '%Y-%m-%d %H:%M:%S') AS "sale_age_hrs"
FROM "tickit"."sale" AS "sale"
JOIN "tickit"."event" AS "event"
  ON "event"."id" = "sale"."id"
JOIN "tickit"."listing" AS "listing"
  ON "listing"."id" = "sale"."listid"
JOIN "tickit"."category" AS "category"
  ON "category"."id" = "event"."catid"
JOIN "tickit"."venue" AS "venue"
  ON "event"."venueid" = "venue"."id"
ORDER BY
  "sale_age_hrs" DESC,
  "sale"."saletime" DESC,
  "listing"."listtime" DESC,
  "venue"."venuename",
  "category"."catname",
  "event"."eventname"
LIMIT 20
```

### Event Locations and Start Times

This query provides details on the location and start time for each event. It links the event and venue tables within the Tickit database, extracting the event name, start time, and venue name. The information is sorted by event name, offering an organized view that helps attendees and organizers quickly reference the schedule and locations of upcoming events.

```sql
SELECT
  "event"."eventname" AS "eventname",
  "event"."starttime" AS "starttime",
  "venue"."venuename" AS "venuename"
FROM "tickit"."event" AS "event"
JOIN "tickit"."venue" AS "venue"
  ON "event"."venueid" = "venue"."id"
ORDER BY
  "event"."eventname"
```

### Sellers for Each Event

This query lists each event along with the emails of users who sold tickets for those events. It utilizes joins across the event, sale, and user tables within the Tickit database to connect event details with the corresponding sellers. The output will show the event name paired with the seller's email, providing a clear view of who facilitated ticket sales for specific events.

```sql
SELECT
  "event"."eventname" AS "eventname",
  "user"."email" AS "seller"
FROM "tickit"."event" AS "event"
JOIN "tickit"."sale" AS "sale"
  ON "event"."id" = "sale"."eventid"
JOIN "tickit"."user" AS "user"
  ON "sale"."sellerid" = "user"."id"
```

### Top 100 Big Spenders

This query identifies the top 100 users who have spent the most on purchases, ranked by their total expenditure. It aggregates the total amount paid by each buyer, pulling from the user and sale tables within the Tickit database. Results display the first names of users and their corresponding total spent, providing insights into the highest spending customers.

```sql
SELECT
  "user"."firstname" AS "firstname",
  SUM("sale"."pricepaid") AS "total_spent"
FROM "tickit"."user" AS "user"
JOIN "tickit"."sale" AS "sale"
  ON "sale"."buyerid" = "user"."id"
GROUP BY
  "user"."firstname"
ORDER BY
  "total_spent" DESC
LIMIT 100
```

### Sellers with Multiple Sales on Holidays

This query lists sellers who achieved more than one sale on holidays, detailing their names and the specific dates of the sales. It extracts data from the date, sale, and user tables within the Tickit database, focusing solely on transactions that occurred on designated holiday dates. The output is grouped by the calendar date and seller names, highlighting those who successfully made multiple sales on each holiday.

```sql
SELECT
  "d"."caldate" AS "caldate",
  "u"."firstname" AS "firstname",
  "u"."lastname" AS "lastname",
  COUNT("s"."id") AS "sale_count"
FROM "tickit"."date" AS "d"
JOIN "tickit"."sale" AS "s"
  ON "d"."id" = "s"."dateid"
JOIN "tickit"."user" AS "u"
  ON "s"."sellerid" = "u"."id"
WHERE
  "d"."holiday" = 1
GROUP BY
  "d"."caldate",
  "u"."firstname",
  "u"."lastname"
HAVING
  COUNT(*) > 1
```

### Top 10 Highest Earning Sellers on Holidays and Weekends in 2020

This query identifies the top 10 sellers based on their commission earnings during holidays and weekends in the year 2020. It filters data by combining information from the date, sale, and user tables in the Tickit database, focusing on specific days that fall on weekends (Saturday and Sunday) or are designated holidays. The result highlights the sellers who maximized their earnings on these peak days, ranked by their total commissions.

```sql
SELECT
  "u"."firstname" AS "firstname",
  "u"."lastname" AS "lastname",
  SUM("s"."commission") AS "holiday_weekend_commission"
FROM "tickit"."date" AS "d"
JOIN "tickit"."sale" AS "s"
  ON "d"."id" = "s"."dateid"
JOIN "tickit"."user" AS "u"
  ON "s"."sellerid" = "u"."id"
WHERE
  (
    "d"."day" IN ('SU', 'SA') OR "d"."holiday" = TRUE
  ) AND "d"."year" = 2020
GROUP BY
  "u"."firstname",
  "u"."lastname"
ORDER BY
  "holiday_weekend_commission" DESC
LIMIT 10
 ```
 
### Spending Patterns in Music Categories Among Musical Enthusiasts

This query aggregates the total amount spent and the quantity of tickets sold by music category for users who have a preference for musicals. It combines data from the categories, events, sales, and user tables in the Tickit database to provide a clear view of spending behaviors among musical fans.

```sql
SELECT
  "c"."catname" AS "catname",
  SUM("s"."pricepaid") AS "pricepaid",
  SUM("s"."qtysold") AS "qtysold"
FROM "tickit"."category" AS "c"
JOIN "tickit"."event" AS "e"
  ON "c"."id" = "e"."catid"
JOIN "tickit"."sale" AS "s"
  ON "e"."id" = "s"."eventid"
JOIN "tickit"."user" AS "u"
  ON "s"."buyerid" = "u"."id" AND "u"."likemusicals" = TRUE
GROUP BY
  "c"."catname"
```




