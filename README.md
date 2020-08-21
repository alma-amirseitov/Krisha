# Krisha Platform Crawler
Krisha Platform Crawler makes continuous 
collection of advertisements from 
Krisha.kz and writes it to DB for future use. 
## Installation
```bash
$ git clone https://github.com/alma-amirseitov/Krisha
```
## QUICK START
To run this project use Docker-compose
```bash
docker-compose up --scale links_parser=4
```
## Architecture

```mermaid
  links_scraper --> RabbitMq(queue links);
  RabbitMq(queue links) --> advertisement_scraper;
  advertisement_scraper --> RabbitMq(queue advertisements);
  RabbitMq(queue advertisements) --> Saver;
  Saver --> ElasticSearch(index krisha_ads);
```

## Used Technologies
- RabbitMq
- Docker
- ElasticSearch

## Project status
Finalization of project and working wit bags
