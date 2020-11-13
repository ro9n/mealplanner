your code goes in this directory

NOTE The following assumptions are also listed in the code as comments and they share the same number
# Assumptions:
1. All three arguments are always passed to the main function, the program terminates otherwise
2. Data directory and it's relative ordering does not change
3. All dates are of layout yyyy-mm-dd
4. preceding time refers to same size window ending before start date
	// e.g. if start=16/1 end=30/1, preceding start=1/1 end = 15/1
5. lookup date range is inclusive i.e. [start, end], not (start, end)
6. Meal ids appearing in the meal to day map are valid
7. Meal can have 0 or more dishes
8. Filename will always represent a valid user id
9. JSON is always valid
10. active during preceding period means greater eq 5 meals

# Environment
The solution is written in Go v1.15 using gomodules on an UNIX OS, for installation please refer to https://golang.org/doc/install

# Build
Assuming you are is the `src` directory, run the following command to build
`go build -o ../dist/run .`

# Run 
Assuming you are is the `src` directory, run the following command to build
`../dist/run active 2016-09-01 2016-09-08`

# Approach
1. CLI sends query to main
2. main creates a query context based on activity level and passes in the date range
3. context is applied on data directory
  1. the apply call is delegated to underlying context: activeContext, superActiveContext, boredContext each have their own apply implementation
  2. worker receives directory, pulls file with the help of reader, then reads files and creates activity for each one
  3. for each activity, context verifies whether it complies or not and selects users who qualify
4. main receives a list of users and prints on stdout
