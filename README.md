# On Routers 

Lenrouter tries to guess which handler to forward your request to based on the length of the request's url path
(`r.URL.Path` in golang). If this guessing works we have a fast router. The article below explains how to get lucky.
Most routers use trie, tree, regular expressions. We will take another route. 


## Basic strategy
When a request comes in
1. Calculate `l := len(r.URL.Path)`, the length of URL's path.
2. Try to guess which handler to forward the request to based only on `l`.
3. If guessing fails, use brute force. Check every route to find the one that matches request. 

We want to become good at guessing so that we don't have to resort to brute force. 

## Details
TODO: Writup. For now, please read the code which also needs some rewriting, reorganization. 


## Performance comparison
    
