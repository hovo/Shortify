[![GoDoc](https://godoc.org/github.com/khachikyan/Shortify?status.svg)](https://godoc.org/github.com/khachikyan/Shortify)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/khachikyan/Shortify/blob/master/LICENSE)

# Shortify

Yet another URL shortener API in Go.

## Running locally

Running this service locally requres Docker and Docker Compose. Please refer to the official [Docker website](https://docs.docker.com/docker-for-mac/install/) for instructions on how to install Docker and Docker compose.

Run the following command in your terminal to start the service: 

```sh
docker-compose up -d
```

## API
**POST** /short

Convert a long URL to short URL

Sample Request
```
curl -d '{"long_url": "https://apple.com"}' -X POST localhost:80/short
```

Sample Response
```json
{
  "short_url": "jCL9RStYD9"
}
```
## 
**GET** /{slug}

Redirects to destination (long) URL

Enter the following in your browser:
```
localhost:80/jCL9RStYD9
```

## 
**GET** /{slug}/clicks

This will return the click counts for a specified short URL.

Query Parameter: **unit = {day, week, all}**

Sample Request
```
curl -X GET 'localhost:80/jCSrXmVI1t/clicks?unit=day'
```

Sample Response
```json
{
  "url_visits": 3
}
```
## Design
### Assumptions
- If no query parameter is provided for analytics, default to last 24 hours
- The system does not report unique URL visits (not based on IP address)
- The maximum length of URL to be shortened is 255
- For simplicity, the system will not cache popular URLs

### URL Encoding
The initial iteration of the system would encode the URL by generating a random integer and converting it to base 62 (a-zA-Z0-9). The following approach has several issues:

1. Every time we generate a slug, we'd have to query the database to determine whether it is available
2. Not very scalable

After a quick search for distributed unique ID generator, I came across [sony/snowflake](https://github.com/sony/sonyflake) a Go implementation of [Twitter's Snowflake](https://developer.twitter.com/en/docs/basics/twitter-ids.html) (used to generate tweet IDs). To encode the URL, the system generates a unique 64-bit unsigned integer and performs a base 62 conversion, which yields a 10 character long slug. There are several advantages to this approach:

1. Generating IDs are uncoordinated, meaning machines generating ids do not have to coordinate with one another
2. Low probability of finding a working short URL from an existing one
3. Scales better

This simplicity of this design decision has its drawbacks, notably the length of the URL slug which is 10 characters long. Having a short slug gives provides a competitive advantage over other URL shortener services.

### Data Persistence
The system utilizes a relational database PostgreSQL to persist URL entries. The database schema consists of two tables, one responsible for mapping URL slug to destination URL, and the other analytics. 

Database schema for URLs:

id: int | slug: varchar(10) | destination: varchar(255)
------------ | ------------- | -------------
1 | jCDzk88W6l | https://stripe.com/docs/api/charges

Database schema for URL metrics:

id: int | slug: varchar(10) | visited_at: timestamp
------------ | ------------- | -------------
1 | jCDzk88W6l | 2019-09-06

This design approach has several issues:

1. If the system grows, database insertions for URL visit events would be pretty costly.
2. For popular URLs we'd still have to query the database to fetch the destination URL

#### Future Improvements
We may want to switch from relational database to NoSQL database for scalability and high accessibility. Something like Cassandra would be a good choice for the system, as it's designed for fast reads and writes.

Additionally, an in-memory cache layer should be added to the system for handling popular URLs. We could also use the least recently used (LRU) cache eviction policy once the cache capacity is reached.

### Analytics
#### Future Improvements
Instead of saving URL visit requests into the database, we can log successful short URL visit events to a system file on a separate container. The business logic for analytics would be handled on a separate service, which would run periodically every 24-hour intervals to tally the URL visits and store the results in both the database and in-memory cache such as Redis or Memcached. This would mean that our analytics data would be stale for 24 hours.

This approach would decouple the system, thus improving performance and maintenance.

### Scalability
The nature of microservices architecture enables the current system to scale horizontally. We can spin-up additional API services by invoking the following command:

```
docker-compose scale <service name> = <no of instances>
```

### Further Improvements
- Cache popular URLs into (Redis or memcached) for high availability, and use a LRU for cache eviction
- Remove the single point of failure from the system
  - Add a database cluster (master/slave) for data replication and failover
- Explore new approaches for generating shorter URL slugs

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
