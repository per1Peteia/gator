# gator

a rss blog aggregator 

## TODO

- might have to refactor my file handling
    - i am using simple os.ReadFile() os.WriteFile() and no decoders (i am (un)marshaling)
    - could lead to significant memory overhead, but i don't know what kind of operation logic the config will have to handle... could be fine