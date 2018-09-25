# rumoursgo

A blockchain implementation in Go (and a terrible Fleetwood Mac based pun that I hope I don't get sued for).


## The Chain

The actual blockchain uses the SHA256 hashing algorithm, and it uses this in a non-dependency injection based way at the moment. This is because as far as I understand it, you couldn't possibly change the hashing algorithm being used for a blockchain without invalidating the chain and having to start over, so there seemed little point in making that abstraction.


## Understanding Blockchain in Real World Scenarios

This is more a less a section of my notes from my understanding of Blockchain reading, and reading up on things like file structures for storing large chunked data to disk. I considered that maybe the blockchain would need B+ trees for storage like a database but since there is only a single 'index' which is time or to put it another way the direction of new blocks being appended, there didn't seem to be much point. 

However, when it comes to storing the blockchain to disk it seems there is a tradeoff. We have two options with different implications:

1) Fixed-size blocks

This makes scanning through blocks much much faster, because we know the exact size, so to get to a particular block we just go to file offset (blocksize * (x - 1)) and read the chunk of size blocksize (where x is block # we want). This is also true for say, calculation of length of chain.

2) Non-fixed-size blocks

This would involve the data section of the block having a header as with most file types on disk which stated the type of data and the size of the block's data to determine where a block ended and began. This means we effectively have to go through each block to get the boundaries or possibly have some kind of file header for the entire chain which defines the boundaries for blocks. It could however, make the blockchain much more flexible into the future and able to store large chunks if required. I'm not sure if this is desirable really - it seems like it could be, but at the same time in some ways the restricted size means that the importance of each chunk is maintained and we don't waste time with crap in our blockchain. Or worse, people storing illegal content of some kind.

The implementation so far I've chosen is fixed size blocks. The overhead is much lower in terms of doing something like grabbing a block from the chain.
