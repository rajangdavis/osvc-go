# osvc-rest

A command line interface for using the [Oracle Service Cloud REST API](https://docs.oracle.com/en/cloud/saas/service/18b/cxsvc/toc.htm).

## Installation

You can install either by fetching with go and building the source

	$ go get github.com/rajangdavis/osvc-rest && cd $GOPATH && go build

Or downloading from the relevant architecture under the [releases link](https://github.com/rajangdavis/osvc-rest/releases)

Windows Example ([using Git Bash for Windows](https://gitforwindows.org/)): 

	$ curl -L -O https://github.com/rajangdavis/osvc-rest/releases/download/v1.0-beta/osvc-rest-windows-386 && chmod u+x ./osvc-rest-windows-386 && mv ./osvc-rest-windows-386 ./osvc-rest

Linux Example (tested on Ubuntu 16.04.3 LTS):

	$ curl -L -O https://github.com/rajangdavis/osvc-rest/releases/download/v1.0-beta/osvc-rest-linux-amd64 && chmod u+x ./osvc-rest-linux-amd64 && mv ./osvc-rest-linux-amd64 ./osvc-rest

You will be able to run the executable in the following way:

	$ ./osvc-rest <commands>

If you want to get rid of the "./", open your .bash_profile and create the following alias:
	
	alias osvc-rest="~/path/to/osvc-rest"

After you have saved and restarted your bash session, you can use the command as expected


## Basic Usage
The basic formula for this CLI is the following:
	
	$ osvc-rest <command to run> <something the command needs> <optional flags to change various settings> <some way to authenticate>

The basic commands come in the following flavors:

1. [HTTP Methods](#http-methods)
	1. For creating objects and [uploading file attachments](#uploading-file-attachments), use the [**post** command](#post)
	2. For fetching objects and [downloading file attachments](#downloading-file-attachments), use the [**get** command](#get)
	3. For updating objects and [uploading file attachments to already existing objects](#uploading-file-attachments) use the [**patch** command](#patch)
	4. For deleting objects use the [**delete** command](#delete)
	5. For looking up options, use the [**options** command](#options)
2. [Running one or more ROQL queries](#running-one-or-more-roql-queries)
3. [Running reports](#running-reports)

Here are the _spicier_ (more advanced) commands:

1. [Bulk Delete](#bulk-delete)
2. [Running multiple ROQL Queries in parallel](#running-multiple-roql-queries-in-parallel)
3. [Performing Session Authentication](#performing-session-authentication)

## Authentication:
Use the following flags to authenticate

	  -i, --interface (string) Oracle Service Cloud Interface to connect with
	  
	  Basic Authentication
	  -u, --username (string)  Username to use for basic authentication
	  -p, --password (string)  Password to use for basic authentication

	  Session Authentication
	  -s, --session (string)  Sets the session token for session authentication

	  OAuth Authentication (untested but should work)
	  -o, --oauth (string)  Sets the bearer token for OAuth authentication

## Performing Session Authentication

1. Create a custom script with the following code:

```php
<?php

// Find our position in the file tree
if (!defined('DOCROOT')) {
$docroot = get_cfg_var('doc_root');
define('DOCROOT', $docroot);
}
 
/************* Agent Authentication ***************/
 
// Set up and call the AgentAuthenticator
require_once (DOCROOT . '/include/services/AgentAuthenticator.phph');

// get username and password
$username = $_GET['username'];
$password = $_GET['password'];
 
// On failure, this includes the Access Denied page and then exits,
// preventing the rest of the page from running.
echo json_encode(AgentAuthenticator::authenticateCredentials($username,$password));

```

2. Run a curl command against the location of the custom script and pass in the account credentials that you would like to create a session with

		curl "https://$OSC_SITE.custhelp.com/cgi-bin/$OSC_CONFIG.cfg/php/custom/login_test.php?username=$OSC_ADMIN&password=$OSC_PASSWORD"


3. If successful, a JSON object with the following properties should return:

		{"acct_id":644,"session_id":"ILsRhQkUWW2ef6wOR90yMcy0Isdfsf8b02POI"}

4. To retrieve the "session_id" of the JSON object, we will set the session id to a bash variable and parse the JSON object using [jq, a popular command line tool for manipulating JSON](https://stedolan.github.io/jq/).

		export SESSION_ID=$(curl "https://$OSC_SITE.custhelp.com/cgi-bin/$OSC_CONFIG.cfg/php/custom/login_test.php?username=$OSC_ADMIN&password=$OSC_PASSWORD" | ~/Desktop/jq-win64.exe -r '.session_id') && echo $SESSION_ID

5. If successful, you wil see the "session_id" property by itself:

		ILsRhQkUWW2ef6wOR90yMcy0Isdfsf8b02POI

6. Finally, we will chain the output of the above commands and use it with osvc-rest to generate JSON:

		export SESSION_ID=$(curl "https://$OSC_SITE.custhelp.com/cgi-bin/$OSC_CONFIG.cfg/php/custom/login_test.php?username=$OSC_ADMIN&password=$OSC_PASSWORD" | ~/Desktop/jq-win64.exe -r '.session_id') && osvc-rest get "" -s $SESSION_ID -i $OSC_SITE --demosite

## HTTP Methods
All the of HTTP Methods have the following formula:
	
	$ osvc-rest <http-verb> <resource-url> (optional flags) <authentication-method>


### POST
In order to create a resource, you must use the **post** command to send JSON data to the resource of your choice

	$ osvc-rest post "opportunities" --data '{"name":"TEST"}' -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE


### GET
In order to fetch a resource, you must use the **get** command to request JSON data from the resource of your choice
	
	$ osvc-rest get "opportunities/?q=name like 'TEST'" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE
	
### PATCH	
In order to update a resource, you must use the **patch** command to send JSON data to update the resource of your choice

	$ osvc-rest patch "opportunities/5" --data '{"name":"updated NAME for TEST"}'  -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

### DELETE
In order to delete a resource, you must use the **delete** command to delete the resource of your choice
	
	$ osvc-rest delete "opportunities/5" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

### OPTIONS
To review the options of what HTTP verbs you can use against a resource, use the **options** command
	
	$ osvc-rest options "opportunities" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

## Uploading File Attachments

In order to upload a file attachment, use the --attach-file (or -f) flag to attach a file with the file location

	$ osvc-rest post "opportunities" --data '{"name":"TEST"}' --attach-file "./proof_of_purchase.jpg" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

To attach multiple files, use the --attach-file (or -f) flag for each file you wish to attach

	$ osvc-rest patch "incidents/302" -f "front_angle.png" -f "back_angle.png" -f "side_angle.png" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

## Downloading File Attachments

In order to download a file attachment from a given resource, [add "?download" to the file attachment URL](https://docs.oracle.com/en/cloud/saas/service/18b/cxsvc/c_osvc_managing_file_attachments.html#ManagingFileAttachments-07BABEF6__concept-406-3A92801C). The file will be downloaded in the same directory that the command is run in.

	$ osvc-rest get "incidents/24898/fileAttachments/253?download" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

To download all file attachmentss from a given resource, add ["?download" to the file Attachments URL](https://docs.oracle.com/en/cloud/saas/service/18b/cxsvc/c_osvc_managing_file_attachments.html#ManagingFileAttachments-07BABEF6__concept-410-3A92801F). A file called "downloadedAttachment.tgz" will be downloaded to your computer. 

	$ osvc-rest get "incidents/24898/fileAttachments?download" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

You can extract the file using [tar](https://askubuntu.com/questions/499807/how-to-unzip-tgz-file-using-the-terminal/499809#499809)
    
    $ tar -xvzf ./downloadedAttachment.tgz

## Running one or more ROQL queries
Runs one or more ROQL queries and returns parsed results
	
	Single Query Example:
	$ osvc-rest query "DESCRIBE" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE
	
	Multiple Queries Example: (Queries should be wrapped in quotes and space separated)
	$ osvc-rest query "SELECT * FROM INCIDENTS LIMIT 100" "SELECT * FROM SERVICEPRODUCTS LIMIT 100" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

## Running Reports
Runs an analytics report and returns parsed results

	Report (without filters) Example:
	$ osvc-rest report --id 176 -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

	Report (with filters and limiting) Example:
	$ osvc-rest report --id 176 --limit 10 --filters '[{"name":"search_ex","values":"returns"}]' -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

## Bulk Delete
This CLI provides a very simple interface to use the Bulk Delete feature within the latest versions of the REST API. Before you can use this feature, make sure that you have the [correct permissions set up for your profile](https://docs.oracle.com/en/cloud/saas/service/18b/cxsvc/c_osvc_bulk_delete.html#BulkDelete-10689704__concept-212-37785F91).

Here is an example of the how to use the Bulk Delete feature: 

	$ osvc-rest query "DELETE from incidents limit 1000" "DELETE from incidents limit 1000" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE --demosite -v latest -a "Testing bulk delete multiple requests"

## Running multiple ROQL Queries in parallel
Instead of running multiple queries in with 1 GET request, you can run multiple GET requests and combine the results

	$ osvc-rest query --parallel "SELECT * FROM INCIDENTS LIMIT 20000" "SELECT * FROM INCIDENTS Limit 20000 OFFSET 20000" "SELECT * FROM INCIDENTS Limit 20000 OFFSET 40000" "SELECT * FROM INCIDENTS Limit 20000 OFFSET 60000" "SELECT * FROM INCIDENTS Limit 20000 OFFSET 80000" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE -v latest -a "Fetching a ton of incidents info"


## Optional Flags:
	    --access-token (string) 	Adds an access token to ensure quality of service
	-a, --annotate (string)     	Adds a custom header that adds an annotation (CCOM version must be set to "v1.4" or "latest"); limited to 40 characters
	    --debug                 	Prints request headers for debugging
	    --demosite              	Change the domain from 'custhelp' to 'rightnowdemo'
	-e, --exclude-null          	Adds a custom header to excludes null from results
	    --next-request (int)      	Number of milliseconds before another HTTP request can be made; this is an anti-DDoS measure
	    --no-ssl-verify         	Turns off SSL verification
	    --schema                	Sets 'Accept' header to 'application/schema+json'
	    --suppress-rules        	Adds a header to suppress business rules
	-t, --utc-time              	Adds a custom header to return results using Coordinated Universal Time (UTC) format for time (Supported on November 2016+)
	-v, --version (string)      	Changes the CCOM version (default "v1.3")

## Help
	Use "osvc-rest [command] --help" for more information about a command.

## Underlying Tech and Attribution
osvc-rest is written in [golang](https://golang.org).

I started with [this initial codebase](https://github.com/dharmeshkakadia/cobra-example) and built/tested the features from there.