# Compressed Air Blog Stats
This application is used to query the HubSpot API for the blog articles and display social media share information for analytics reporting.
Versions of this application show the progression and re-architecture to take advantage of Go language features such as concurrency and interface types

# Languages
- GO

# Versions

## V0
* Syncronous execution of network requests to fetch share count for each article for each social network

## V1
* Make network requests in parallel by launching share count functions as goroutines
  * This works on the surface by will ocassionally panic due to concurrent processes acceessing shared memory

## V2
* Restructure application to return data from share count functions over channels rather than updating shared memeory
* Requires select block to read data from dedicated channels for each social network.

## V3
* Restructure application to use single channel that accepts a shared interface which each social network share count type can implicitly satify by defining a "GetShareCount" method
* Replace select block with switch block to check against interface type to get access to underlying data structure of individual social network API response

## Final
* Add error handling between concurrent processes by utilizing dedicated error channel and using select block to ingest data from one of two channels

## Final Final
* Replace dedicated error channel by updating social network API response structs to contain an parameter of type error which can be updated in "GetShareCount" method
* Removes the need to update the shared interface to accept dedicated channel and gives more context to where error occurred

