# Development Trains

## Installation

To run the project locally, you will need to have Go with at least version 1.23.1 installed. After that, you may build the source using:

```bash
make build
```

The latest binary source should
If you don't have Go installed and need to build the source, you can use Docker to build and run the project locally. You can build the Docker image using:

```bash
make build-image
```

This will create an image called `idea456:development-trains` which contains the build output of the source.

## Usage

To run the program, you must specify a file path to the input path first by using the `-i` flag:

```
./development-trains -i ./tests/sample.txt
```

If using a Docker image, you can run this example command with a file path to your local environment holding the test files:

```bash
docker run -v <path to your test folder containing test files>:/tests idea456:development-trains -i /tests/<name of test file>.txt
```

If choosing to input from a text file, ensure that it is following the format:

```
2 // number of stations
A
B

6 // number of routes
E1,A,B,10
E2,A,C,20
E3,B,D,10
E4,C,D,30
E5,D,E,10
E6,D,F,10

4 // number of packages
K1,1,A,E
K2,1,B,E
K3,1,D,E
K4,3,C,F

2 // number of trains
Q1,3,A
Q2,3,A
```

You may also choose to prompt for the input instead using the `--prompt` flag without the `-i` flag, it will prompt and output:

```bash
Enter the number of stations:
3
Enter the stations (one per line):
A
B
C
Enter the number of routes:
2
Enter the routes (one per line, format: E1,A,B,10):
E1,A,B,30
E2,B,C,10
Enter the number of packages:
1
Enter the packages (one per line, format: K1,3,A,E):
K1,5,A,C
Enter the number of trains:
1
Enter the trains (one per line, format: Q1,3,A):
Q1,6,B

```

Using either input method will return a list of moves with the specified format:

```
W=0, T=Q1, N1=B, P1=[], N2=A, P2=[]
W=30, T=Q1, N1=A, P1=[K1], N2=B, P2=[]
W=60, T=Q1, N1=B, P1=[], N2=C, P2=[K1]
```

To print out a more detailed list of moves, you may specify the `--verbose` flag which will print out the moves and steps taken by each train:

```bash
./development-trains -i ./tests/sample.txt --verbose
```

which returns the following output:

```
[0 minutes] Train Q1 moving from station B to station A

[30 minutes] Train Q1 moving from station A to station B
Carried packages:
	- K1 package with weight 5 heading to C station

[60 minutes] Train Q1 moving from station B to station C
Dropped packages:
	- K1 package with weight 5 at C station
```

To include a summary of time taken for each package to be delivered, you can specify the `--summary` flag:

```bash
./development-trains -i ./tests/sample.txt --summary
```

which returns the following output:

```
W=0, T=Q1, N1=B, P1=[], N2=A, P2=[]
W=30, T=Q1, N1=A, P1=[K1], N2=B, P2=[]
W=60, T=Q1, N1=B, P1=[], N2=C, P2=[K1]

Name Weight DeliveredAt Train
K1   5kg    60m         Q1
```
