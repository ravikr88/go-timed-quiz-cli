# Go Timed Quiz CLI

This is a command-line quiz application written in Go. The application reads questions and answers from a CSV file, presents the questions to the user, and evaluates the answers within a specified time limit.

## Features

- Read questions and answers from a CSV file.
- Randomize the order of questions.
- Timed quiz session.
- Evaluate and summarize the user's performance.

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/ravikr88/go-timed-quiz-cli.git
    cd go-timed-quiz-cli
    ```

2. Build the application:

    ```sh
    go build -o quiz
    ```

## Usage

```sh
./quiz -file=path/to/questions.csv -random=true -time=10
```

##  A smple look of interface
 <img width="1080" alt="Screenshot 2024-06-08 at 2 42 39 AM" src="https://github.com/ravikr88/go-timed-quiz-cli/assets/135989427/be7fff8e-bc21-4703-86c5-f55c6ad62fb8">
