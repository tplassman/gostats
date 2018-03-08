# Go Stats
This application is used to query the HubSpot API for the blog articles and display social media share information for analytics reporting.
Versions of this application show the progression and re-architecture to take advantage of Go language features such as concurrency and interface types

# Languages
- GO

# Versions

## V0
* Synchronous execution of network requests to fetch share count for each article for each social network

## V1
* Make network requests in parallel by launching share count functions in goroutines
  * Wrap share count functions in closures to contain concurrent processing in controller and maintain existing share count function signature
  * This works on the surface by will ocassionally panic due to concurrent processes acceessing shared memory

## V2
* Restructure application to return data from share count functions over channels rather than updating shared memeory
* Share count now updated in single for loop that will block while waiting to receive over the channel
* Add shared data structure to define information passed over channel

## V3
* Restructure application to use a shared interface which each social network share count type can implicitly satify by defining a "GetShareCount" method
* Use switch block on shared interface struct implentor type to determine social network to add share count for

## Final
* Add error handling between concurrent processes by utilizing dedicated error channel
* Utilize select block inside of for loop to handle data from multiple channels

