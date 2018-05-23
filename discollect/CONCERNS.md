Component Concerns
------

Metastore tracks state for individual scrapes, watching status / config / start info

When a scrape is completed, the *ScrapeStatus is written to the metastore. 
Active scrapes are checked by polling the Queue  

