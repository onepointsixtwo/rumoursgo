# rumoursgo

A blockchain implementation in Go (and a terrible Fleetwood Mac based pun that I hope I don't get sued for).


## The Chain

The actual blockchain uses the SHA256 hashing algorithm, and it uses this in a non-dependency injection based way at the moment. This is because as far as I understand it, you couldn't possibly change the hashing algorithm being used for a blockchain without invalidating the chain and having to start over, so there seemed little point in making that abstraction.


## Upcoming changes

Serialisation and deserialisation of the chain so it can be shared between computers as a blockchain must to be secured.
