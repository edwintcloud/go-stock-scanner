# objectives
1. show stocks with pre-market gap -4% in an easy to understand dashboard
2. also show relevant metrics for each stock in the dashboard table *maybe*

# how to accomplish
1. each scanner update, save the minute bar for all symbols (that are not etf)
2. minute bars go in a hash bucket (slice) based on number of minutes in a day so they can be efficiently retrieved and reset each day
3. use ticker snapshots for prev day info and calculating rvol