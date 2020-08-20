# Introduction

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

## Performance 

Number next to `BenchmarkLenrouter_*` in bracket on the left shows its rank in the benchmark. 

```
goos: darwin
goarch: amd64
pkg: github.com/julienschmidt/go-http-routing-benchmark
    BenchmarkAero_Param                   	27906673	        41.6 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_Param        	40726311	        28.9 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_Param                  	10314051	       116 ns/op	      32 B/op	       1 allocs/op
    BenchmarkGoji_Param                   	 2481369	       461 ns/op	     336 B/op	       2 allocs/op
    BenchmarkGojiv2_Param                 	  676768	      1777 ns/op	    1328 B/op	      11 allocs/op
    BenchmarkGorillaMux_Param             	  600684	      1971 ns/op	    1280 B/op	      10 allocs/op
    BenchmarkHttpRouter_Param             	13411254	        84.2 ns/op	      32 B/op	       1 allocs/op
(3) BenchmarkLenrouter_Param              	26824592	        48.3 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_Param5                  	17041465	        69.8 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_Param5       	12720734	        93.9 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_Param5                 	 3645216	       315 ns/op	     160 B/op	       1 allocs/op
    BenchmarkGoji_Param5                  	 1870831	       647 ns/op	     336 B/op	       2 allocs/op
    BenchmarkGojiv2_Param5                	  551001	      2061 ns/op	    1392 B/op	      11 allocs/op
    BenchmarkGorillaMux_Param5            	  422702	      2751 ns/op	    1344 B/op	      10 allocs/op
    BenchmarkHttpRouter_Param5            	 4724962	       254 ns/op	     160 B/op	       1 allocs/op
(1) BenchmarkLenrouter_Param5             	22314896	        55.7 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_Param20                 	41312287	        28.9 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_Param20      	 3168570	       377 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_Param20                	 1000000	      1093 ns/op	     640 B/op	       1 allocs/op
    BenchmarkGoji_Param20                 	  503769	      2052 ns/op	    1247 B/op	       2 allocs/op
    BenchmarkGojiv2_Param20               	  424921	      2661 ns/op	    1632 B/op	      11 allocs/op
    BenchmarkGorillaMux_Param20           	  183615	      6519 ns/op	    3453 B/op	      12 allocs/op
    BenchmarkHttpRouter_Param20           	 1463616	       805 ns/op	     640 B/op	       1 allocs/op
(2)* BenchmarkLenrouter_Param20            	 9937328	       123 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_ParamWrite              	16782452	        69.8 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_ParamWrite   	40697214	        29.2 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_ParamWrite             	 7735794	       148 ns/op	      32 B/op	       1 allocs/op
    BenchmarkGoji_ParamWrite              	 2410459	       507 ns/op	     336 B/op	       2 allocs/op
    BenchmarkGojiv2_ParamWrite            	  561968	      1945 ns/op	    1360 B/op	      13 allocs/op
    BenchmarkGorillaMux_ParamWrite        	  550056	      2054 ns/op	    1280 B/op	      10 allocs/op
    BenchmarkHttpRouter_ParamWrite        	 9528392	       118 ns/op	      32 B/op	       1 allocs/op
(3) BenchmarkLenrouter_ParamWrite         	15707262	        79.6 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_GithubStatic            	29153634	        39.9 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_GithubStatic 	26446176	        45.2 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_GithubStatic           	38772614	        31.0 ns/op	       0 B/op	       0 allocs/op
    BenchmarkGoji_GithubStatic            	 7615426	       157 ns/op	       0 B/op	       0 allocs/op
    BenchmarkGojiv2_GithubStatic          	  691842	      1738 ns/op	    1312 B/op	      10 allocs/op
    BenchmarkGorillaMux_GithubStatic      	  240483	      4435 ns/op	     976 B/op	       9 allocs/op
    BenchmarkHttpRouter_GithubStatic      	32366912	        36.6 ns/op	       0 B/op	       0 allocs/op
(1) BenchmarkLenrouter_GithubStatic       	38886072	        28.6 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_GithubParam             	13382012	        88.1 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_GithubParam  	 9139987	       131 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_GithubParam            	 4151624	       288 ns/op	     128 B/op	       1 allocs/op
    BenchmarkGoji_GithubParam             	 1602334	       822 ns/op	     336 B/op	       2 allocs/op
    BenchmarkGojiv2_GithubParam           	  495052	      2355 ns/op	    1408 B/op	      13 allocs/op
    BenchmarkGorillaMux_GithubParam       	  177650	      6378 ns/op	    1296 B/op	      10 allocs/op
    BenchmarkHttpRouter_GithubParam       	 5917116	       208 ns/op	      96 B/op	       1 allocs/op
(1) BenchmarkLenrouter_GithubParam        	18571885	        65.0 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_GithubAll               	   67444	     18003 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_GithubAll    	   59752	     19595 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_GithubAll              	   20631	     58208 ns/op	   20224 B/op	     167 allocs/op
    BenchmarkGoji_GithubAll               	    3493	    349016 ns/op	   56113 B/op	     334 allocs/op
    BenchmarkGojiv2_GithubAll             	    1581	    735225 ns/op	  352720 B/op	    4321 allocs/op
    BenchmarkGorillaMux_GithubAll         	     363	   3269097 ns/op	  251655 B/op	    1994 allocs/op
    BenchmarkHttpRouter_GithubAll         	   28476	     43224 ns/op	   13792 B/op	     167 allocs/op
(1) BenchmarkLenrouter_GithubAll          	   74970	     15904 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_GPlusStatic             	36163489	        31.6 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_GPlusStatic  	44769922	        26.7 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_GPlusStatic            	64841548	        18.3 ns/op	       0 B/op	       0 allocs/op
    BenchmarkGoji_GPlusStatic             	10739841	       112 ns/op	       0 B/op	       0 allocs/op
    BenchmarkGojiv2_GPlusStatic           	  575906	      2014 ns/op	    1312 B/op	      10 allocs/op
    BenchmarkGorillaMux_GPlusStatic       	  889083	      1728 ns/op	     976 B/op	       9 allocs/op
    BenchmarkHttpRouter_GPlusStatic       	46614662	        22.1 ns/op	       0 B/op	       0 allocs/op
(3) BenchmarkLenrouter_GPlusStatic        	44706841	        26.1 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_GPlusParam              	20786140	        57.4 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_GPlusParam   	27321030	        44.8 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_GPlusParam             	 6780992	       174 ns/op	      64 B/op	       1 allocs/op
    BenchmarkGoji_GPlusParam              	 2311022	       511 ns/op	     336 B/op	       2 allocs/op
    BenchmarkGojiv2_GPlusParam            	  568128	      1944 ns/op	    1328 B/op	      11 allocs/op
    BenchmarkGorillaMux_GPlusParam        	  454093	      2574 ns/op	    1280 B/op	      10 allocs/op
    BenchmarkHttpRouter_GPlusParam        	 8842765	       134 ns/op	      64 B/op	       1 allocs/op
(3) BenchmarkLenrouter_GPlusParam         	21670611	        54.4 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_GPlus2Params            	13847505	        87.0 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_GPlus2Params 	16481263	        71.8 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_GPlus2Params           	 5137830	       229 ns/op	      64 B/op	       1 allocs/op
    BenchmarkGoji_GPlus2Params            	 1693290	       703 ns/op	     336 B/op	       2 allocs/op
    BenchmarkGojiv2_GPlus2Params          	  485748	      2438 ns/op	    1408 B/op	      14 allocs/op
    BenchmarkGorillaMux_GPlus2Params      	  247010	      4958 ns/op	    1296 B/op	      10 allocs/op
    BenchmarkHttpRouter_GPlus2Params      	 7223353	       163 ns/op	      64 B/op	       1 allocs/op
(1) BenchmarkLenrouter_GPlus2Params       	20465742	        57.9 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_GPlusAll                	 1449496	       812 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_GPlusAll     	 1757090	       685 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_GPlusAll               	  426133	      2592 ns/op	     672 B/op	      11 allocs/op
    BenchmarkGoji_GPlusAll                	  154389	      7571 ns/op	    3696 B/op	      22 allocs/op
    BenchmarkGojiv2_GPlusAll              	   45004	     26875 ns/op	   17616 B/op	     154 allocs/op
    BenchmarkGorillaMux_GPlusAll          	   28720	     40591 ns/op	   16112 B/op	     128 allocs/op
    BenchmarkHttpRouter_GPlusAll          	  655842	      2004 ns/op	     640 B/op	      11 allocs/op
(2) BenchmarkLenrouter_GPlusAll           	 1627490	       728 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_ParseStatic             	33310864	        35.7 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_ParseStatic  	45130071	        26.4 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_ParseStatic            	56506653	        21.2 ns/op	       0 B/op	       0 allocs/op
    BenchmarkGoji_ParseStatic             	 8455958	       142 ns/op	       0 B/op	       0 allocs/op
    BenchmarkGojiv2_ParseStatic           	  596034	      1678 ns/op	    1312 B/op	      10 allocs/op
    BenchmarkGorillaMux_ParseStatic       	  658522	      1817 ns/op	     976 B/op	       9 allocs/op
    BenchmarkHttpRouter_ParseStatic       	54996940	        21.5 ns/op	       0 B/op	       0 allocs/op
(4) BenchmarkLenrouter_ParseStatic        	39056044	        28.7 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_ParseParam              	24300880	        47.7 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_ParseParam   	34384299	        34.8 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_ParseParam             	 7175407	       164 ns/op	      64 B/op	       1 allocs/op
    BenchmarkGoji_ParseParam              	 2105846	       580 ns/op	     336 B/op	       2 allocs/op
    BenchmarkGojiv2_ParseParam            	  605451	      1872 ns/op	    1360 B/op	      12 allocs/op
    BenchmarkGorillaMux_ParseParam        	  539708	      2071 ns/op	    1280 B/op	      10 allocs/op
    BenchmarkHttpRouter_ParseParam        	 9774423	       120 ns/op	      64 B/op	       1 allocs/op
(3) BenchmarkLenrouter_ParseParam         	22240795	        53.4 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_Parse2Params            	19694250	        60.6 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_Parse2Params 	23051820	        52.1 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_Parse2Params           	 5777247	       202 ns/op	      64 B/op	       1 allocs/op
    BenchmarkGoji_Parse2Params            	 2051320	       566 ns/op	     336 B/op	       2 allocs/op
    BenchmarkGojiv2_Parse2Params          	  545546	      1946 ns/op	    1344 B/op	      11 allocs/op
    BenchmarkGorillaMux_Parse2Params      	  488367	      2641 ns/op	    1296 B/op	      10 allocs/op
    BenchmarkHttpRouter_Parse2Params      	 7544499	       155 ns/op	      64 B/op	       1 allocs/op
(3) BenchmarkLenrouter_Parse2Params       	19148420	        61.3 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_ParseAll                	  845348	      1400 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_ParseAll     	 1000000	      1041 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_ParseAll               	  305773	      3840 ns/op	     928 B/op	      16 allocs/op
    BenchmarkGoji_ParseAll                	   92605	     14080 ns/op	    5376 B/op	      32 allocs/op
    BenchmarkGojiv2_ParseAll              	   23923	     52831 ns/op	   34448 B/op	     277 allocs/op
    BenchmarkGorillaMux_ParseAll          	   14480	     86450 ns/op	   30288 B/op	     250 allocs/op
    BenchmarkHttpRouter_ParseAll          	  540916	      2616 ns/op	     640 B/op	      16 allocs/op
(2) BenchmarkLenrouter_ParseAll           	  877754	      1356 ns/op	       0 B/op	       0 allocs/op
    
    BenchmarkAero_StaticAll               	  138068	      8120 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCloudyKitRouter_StaticAll    	   95910	     12613 ns/op	       0 B/op	       0 allocs/op
    BenchmarkDenco_StaticAll              	  235926	      5106 ns/op	       0 B/op	       0 allocs/op
    BenchmarkGoji_StaticAll               	   32330	     36144 ns/op	       0 B/op	       0 allocs/op
    BenchmarkGojiv2_StaticAll             	    3615	    321214 ns/op	  205984 B/op	    1570 allocs/op
    BenchmarkGorillaMux_StaticAll         	    1255	    806566 ns/op	  153236 B/op	    1413 allocs/op
    BenchmarkHttpRouter_StaticAll         	  111940	     10263 ns/op	       0 B/op	       0 allocs/op
(3) BenchmarkLenrouter_StaticAll          	  129529	      9154 ns/op	       0 B/op	       0 allocs/op

PASS
ok  	github.com/julienschmidt/go-http-routing-benchmark	177.100s

```    

(*) `BenchmarkLenrouter_Param20` it gets second rank but it could be first. See https://github.com/julienschmidt/go-http-routing-benchmark/pull/93
