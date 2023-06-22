# `kafkaless`

an experiment with [go](https://go.dev/) + [service weaver](https://serviceweaver.dev/) + [msk serverless](https://docs.aws.amazon.com/msk/latest/developerguide/serverless.html)

# Getting Started

You will need to have the following installed to work with the project

1. [Go](https://go.dev/doc/install)
2. [Weaver](https://serviceweaver.dev/docs.html#installation)

## Running The App

The following assumes you are running a machine with `Make` available. If not, you can use the the `Go` toolchain and `Weaver` CLI directly. Please consult the [`Makefile`](./Makefile) for the commands that are used to compose the following `Make` recipes.

Run the following in your terminal
```
// run the app as a monolith on your local machine
make monolith.run
```

```
// run the app as a fleet of micro services on
// your local machine
make services.run

// to view the diagnostics dashboard run the following
// in a new terminal
make services.dashboard
```

