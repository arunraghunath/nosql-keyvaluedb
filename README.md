# nosql-keyvaluedb

This is a no-sql key value db implementation using BTree to store data.

The database is a file with an extension ".db", which contains pages to store data. 

Each pagesize is defaulted to 4096 bytes --> 4kb

The page 0 is named "meta" and is reserved to store the DB metadata information. This is represented by the struct meta.

Another page is named "freespace" represented by the struct freespace.
This holds information of the last page number which has been utilised to create a page, and the list of pages which were freed up when data was deleted. This helps to reutilise the freed up pages which would otherwise lead to fragmentation.

Data is serialized and converted into binary using bigendian and then stored in the DB. This is done so as to efficiently parse the data on retrieve.

DB Components -- 
1.  Data Access Layer -- Low level access to DB data via a struct named     "accesslayer"