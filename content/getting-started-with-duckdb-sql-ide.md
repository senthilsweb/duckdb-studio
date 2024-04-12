---
title: Mastering SQL Queries with the Tickit Database on DuckDB
description: Dive into the SQL with the Tickit database, a comprehensive dataset for a fictional ticket sales company. This guide will explore the nuances of SQL operations, utilizing DuckDB's efficient analytical data management system. From querying event and venue data to analyzing ticket sales and customer interactions, enhance your SQL skills through practical scenarios and examples.
date: 2024-04-11T12:00:00.000Z
published: true
coverimage: 
ogImage: 
author: Senthilnathan Karuppaiah
avatar: 
type: Documentation
tags: [DuckDB, SQL, Data Analysis, Query Optimization]
---
# Introduction

The ***Tickit*** database typically contains tables related to a fictional ticket-selling company, including data about events, venues, ticket sales, and customer information. It's designed to showcase various SQL operations and query scenarios.


## Blog Tables Setup

### Create the `id_sequence` table


```sql
CREATE TABLE IF NOT EXISTS id_sequence (entity_name VARCHAR, next_id INTEGER);
```

### Create the `blogs` table

```sql
CREATE TABLE IF NOT EXISTS blogs (
    id INT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT current_timestamp,
    modified_at TIMESTAMP DEFAULT current_timestamp,
    update_count INT DEFAULT 0,
    title VARCHAR,
    content TEXT
);
```

## Duckdb Practice SQL Queries from tickitdb

### Most Popular States

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
### Least Popular Events

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

### Inventory Aging

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

### event_location_and_time

```sql
SELECT  events.eventname, events.starttime, venues.venuename
FROM events, venues 
where events.venueid=venues.venueid 
order by events.eventname asc
```

### people_sold_tickets_to_events

```sql
SELECT  events.eventname, users.email as seller
FROM events, users, sales
where
sales.eventid=events.eventid
and
users.userid=sales.sellerid
```

### big_spenders

```sql
SELECT users.firstname, sum(pricepaid) as total_spent
FROM users, sales
WHERE users.userid=sales.buyerid
GROUP BY users.firstname
ORDER BY total_spent desc
LIMIT 100
```

### sellers who made more than 1 sale in a holiday

```sql
select d.caldate, u.firstname, u.lastname, count(*) as sale_count from dates d
 inner join sales s on s.dateid  = d.dateid
 inner join users u on u.userid = s.sellerid
where d.holiday = true
group by 1,2,3having count(*)>1
```

### top 10 sellers made most on holidays or weekends in 2008

```sql
select u.firstname, u.lastname, sum(s.commission) as holiday_weekend_commission
 from dates d
 inner join sales s on s.dateid  = d.dateid 
 inner join users u on s.sellerid = u.userid
 where (d.holiday = true or d.day in ('SU','SA')) and d.YEAR = 2008 group by 1,2order by  holiday_weekend_commission desc limit 10
 ```
 
### Shares of Amount Paid by Music Category among Buyers who like musicals

```sql
select c.catname, sum(s.pricepaid) pricepaid, sum(s.qtysold) as qtysold
 from categories  c
 inner join events e on e.catid = c.catid
 inner join sales s on s.eventid = e.eventid
 inner join users u on u.userid = s.buyerid
where  u.likemusicals = true
group by 1
```




