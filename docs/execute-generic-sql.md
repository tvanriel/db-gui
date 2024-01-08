Executing "generic" sql is non-trivial.  Go's implementation of the SQL interface doesn't allow for "checking" whether a Database response should be read with Query or with Exec before executing the query.

### github.com/xwb1989/sqlparser
This package implements a SQL parser which allows us to programmatically determine whether a query should be executed with Query or Execute by parsing whatever is _in_ the query.

On top of that I want to see a progress bar for the SQL scriptm to show how far we've come.  This requires executing the queries all one by one.

Since this package returns an iterator of Statements that are to be executed.

I think the best option would be to first parse the whole SQL file, and then start submitting a progress amount with the amount of queries that were executed.

To report on the progress bar, we can return a Websocket upgrade.


### http://arlimus.github.io/articles/gin.and.gorilla/