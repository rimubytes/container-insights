# Cgroups Project

This project demonstrates the use of Linux cgroups (Control Groups) to manage and limit resource usage for processes. It uses the containerd/cgroups Go library to create a cgroup that restricts CPU usage for a process.

## Features

- CPU Resource Limiting: Restricts a process's CPU usage using cgroup v2.
- Process Isolation: Adds specific processes to a cgroup to limit their resource consumption.
- Simple Implementation: Demonstrates the basic usage of cgroups in Go, ideal for learning purposes.

### Requirements

- Go 1.16 or higher.
- Linux system with cgroup v2 support enabled (most modern Linux distributions).
- `stress` command-line tool for generating CPU load.

### Installation

1. Install Go Dependencies: Install the required Go modules using the following command:

```bash
go mod tidy
```

2. Install `stress` Tool: The project uses the `stress` tool to simulate CPU load. You can install it using:

```bash
sudo apt-get update
sudo apt-get install stress
```

### Running the project

1. Compile the program

```bash
go build -o cgroups main.go
```

2. Run the program: Execute the compile binary

```bash
sudo ./cgroups
```

Running the program will do the following:

- Create a new cgroup named my-cgroup.
- Limit the CPU usage of the stress process to 20% of a single CPU core.
- Run the stress tool inside the cgroup for 10 seconds.
- Expected Output: The stress tool will run with limited CPU usage due to the cgroup restrictions. You can observe the CPU usage in another terminal using a tool like htop or top to verify that it doesn't exceed the defined limit.

![cgroup execution](/assets/cgroups-example.png)