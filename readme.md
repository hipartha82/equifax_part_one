**##** **Golang Dashboard API**



This project demonstrates a Golang-based REST API that fetches user details and todo information concurrently from external sources, aggregates, filters, merges them, and serves a consolidated dashboard response.



**##** **Overview of the API**



Exposes an endpoint /dashboard/{userId}.

Fetches user details from https://dummyjson.com/users/{userId}.

Fetches todo information from https://dummyjson.com/todos/users/{userId}.

Processes both datasets concurrently using goroutines and sync.WaitGroup.

Formats the merged output into JSON with fields like Id, Full\_name, Status, Pending\_task\_count, Next\_urgent\_task, and Error\_warning.



\## **Requirements**



Go 1.18 or later



Internet connection (to fetch data from dummyjson.com)



**## Running the API**



 	# Clone the repository:



 	git clone https://github.com/yourusername/dashboard-api.git

 	cd dashboard-api



 	# Run the server:

 	go run main.go



 	# Access the API:

 	curl http://localhost:8080/dashboard/3





**## Project Structure**



.

├── main.go        # Golang source code for the API

└── README.md      # Documentation





**## Example Output**



&nbsp;	\*\***Request**\*\*:


 	curl http://localhost:8080/dashboard/3



&nbsp;	

&nbsp;	\*\***Response**\*\*:


 	{

 	 "id": "3",

  	 "full\_name": "John Doe",

  	 "status": "Veteran",

  	 "pending\_task\_count": "2",

  	 "next\_urgent\_task": "Finish report",

  	 "error\_warning": "null"

 	}



\## **Key Features**



**Data Transformation**: Processes raw API data into a structured summary with aggregation \& Filtering as per requirements both for User details and ToDos

**Concurrency**: Fetches user details and todos in parallel.

**Timeout Handling**: Uses context.WithTimeout to enforece a global timeout of 2 sec to avoid long waits for each concurrent tasks

**Error Handling**: Handles multiple erroneous conditions distinctly based on possible dataset possibilities \& provides meaningful error messages in the JSON output.





**## Appendix**



For easy testing purpose, another file based interface version for the program is also created maintaining the same logic and above key features and had been tested successfully.

To keep it simpler, inputs for the User details and Todos were created in the program itself and was queried for a specific user Id to perform quick validation on the core logic.



e.g. Below is the sample output when queried for user ID : 2



-------------------------------------------------------------



Details File Created

ToDo File Created



Original User file entries:

-------------------------------

1,John, Paul, 40   

2,Sameer, kumar, 56   

3,Rahul, Paul,  49  

4 , Raymond, Lee, 58   

5 , Rakesh, Chauhan, 38   



Expected output of User data from the array:

------------------------------------------------

1#John Paul#Rookie#

2#Sameer kumar#Veteran#

3#Rahul Paul#Rookie#

4#Raymond Lee#Veteran#

5#Rakesh Chauhan#Rookie#



Original Todo file entries:

-------------------------------

1,HLD to Write, Y   

1,Dev Process to build, Y   

1,Testing to certify, N   

&nbsp;2, HLD to Write, N   

2,Dev Process to build, N   

2,Testing To certify, N   

&nbsp;3, HLD to Write, Y   

3,Dev Process to build, N   

3,Testing To certify, N   

&nbsp;4, HLD to Write, Y   

4,Dev Process to build, Y   

4,Testing To certify, Y   

6, HLD to Write, N   

6,Dev Process to build, N   

6,Testing To certify, Y  



Expected output of Todo data from the array:

------------------------------------------------

1#Testing to certify#1#

2#HLD to Write#3#

3#Dev Process to build#2#

4##0#

6#HLD to Write#2#



The User ID to get the data for details and Todo is :  2





Both user details and Todo Summary are successfully received

User Details :  {\[2 Sameer kumar Veteran] true <nil>}

Todo Summary :  {\[2 HLD to Write 3] true <nil>}



The Merged output for the User details with Todo Summary is : 

&nbsp;{

&nbsp; "id": "2",

&nbsp; "full\_name": "Sameer kumar",

&nbsp; "status": "Veteran",

&nbsp; "pending\_task\_count": "3",

&nbsp; "next\_urgent\_task": "HLD to Write",

&nbsp; "error\_warning": "null"

}



=== Code Execution Successful ===

