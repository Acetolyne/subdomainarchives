package main

import (
	"flag"
	"fmt"
	"os"
	"net/http"
	"io"
	"strings"
)

func main(){
	/*
	flags
	-e existing domains only
	-n nonexisting domains only
	-f full url including paths
	-r filter with regex
	api info at https://github.com/internetarchive/wayback/blob/master/wayback-cdx-server/README.md
	*/
	//Set the available flags
	//Custom error message when usage is wrong
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("EXAMPLE USAGE: subdomainarchives [-e -n] [-f] [-r {regex}] BASE URL without subdomain")
	}
	exists := flag.Bool("e", false, "only show existing urls")
	notexists := flag.Bool("n", false, "only show non-existing urls")
	full := flag.Bool("f", false, "search full urls not just subdomains")
	regex := flag.String("r", "", "use a regex to match the urls")
	//parse all the flags the user inputed
  	flag.Parse() 
	//@todo remove the printing debig only
  	fmt.Println("exists:", *exists)
	fmt.Println("not exists:", *notexists)
	fmt.Println("full:", *full)
	fmt.Println("regex:", *regex)
	//Get the url
	url := flag.Arg(0)
	if url == ""{
		flag.Usage()
	}
	fmt.Println("url:", url)

	set := make(map[string]struct{})
	//make sure that the user did not use the e flag and the n flag together, that doesnt make sense
	if (*exists == true && *notexists == true){
		fmt.Println("CAN NOT USE e AND n FLAGS TOGETHER\n")
		flag.Usage()
	}

	//Generate the url based on the flags the user passed in
	requesturl := "https://web.archive.org/cdx/search/cdx?url=*." + url
	//make the request to archive.org
	resp, err := http.Get(requesturl)
	if err != nil {
   		fmt.Println(err)
	}
	fmt.Println(resp.StatusCode)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	results := string(body)
	lines := strings.Split(results, "\n")
	for _, l := range lines {
		finalurl := ""
		//Extract the url from the line
		cur := strings.Split(l, " ")
		if (len(cur) >1) {
			//If f flag not used then split on url and remove the end
			if (*full == false) {
				fullurl := strings.Split(cur[2], url)
				finalurl = fullurl[0] + url
			}else{
				finalurl = cur[2]
			}
			//add url to set if it is not already in there
			set[finalurl] = struct{}{}
			//fmt.Println(finalurl)
			
		}
	}
	for key := range(set) {
		if (*exists == true) {
			resp, err := http.Get(key)
			if err != nil {
				fmt.Println("'"+key+"':")
				fmt.Println(err)
			}
			if (resp.StatusCode == 200) {
				fmt.Println(key)
			}
		}
		if (*notexists == true) {
			resp, err := http.Get(key)
			if err != nil {
				fmt.Println("'"+key+"':")
				fmt.Println(err)
			}
			if (resp.StatusCode != 200) {
				fmt.Println(key)
			}
		}
	}
}