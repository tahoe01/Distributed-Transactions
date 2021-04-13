# Distributed-Transactions

## Notes

- You are welcome to switch languages or groups. If you do change groups, please update us so we can update the cluster assignments.
- You are allowed to use any programming language; however, the TAs will only help with the four supported languages: C/C++, Go, Java, and Python.
- Your implementation will be tested on the CS425 VMs. It is your responsibility to make sure that the code runs on these VMs.

## Overview

In this MP you will be implementing a distributed transaction system. You goal is to support transactions that read and write to distributed objects while ensuring full ACI(D) properties. (The D is in parentheses because you are not required to store the values in durable storage or implement crash recovery.)

## Clients, Servers, Branches, and Accounts

You are (once again) implementing a system that represents a collection of accounts and their balances. The accounts are stored in five different branches (named A, B, C, D, E).  An account is named with the identifier of the branch followed by an account name; e.g., `A.xyz` is account `xyz` stored at branch `A`. Account names will be comprised of lowercase english letters.

You will need to implement a server that represents a branch and keeps track of all accounts in that branch, and a client that executes transactions. *For each transaction*, a client randomly picks one of the five servers to communicate with. The chosen server acts as the coordinator for that transaction, and communicates with the client and all other necessary servers. The coordinator can therefore be different for different transactions. More details on transactions and their properties have been provided later in this specification.

Unlike the previous MPs, you *do not have to handle failures* and can assume that all the servers remain up for the duration of the demo. Clients may exit but you do not have to deal with clients crashing in the middle of a transaction.


## Configuration

Each server must take three arguments. The first argument identifies the branch that the server must handle. The second argument is the port number it listens on. The third argument is a configuration file. The configuration file has 5 lines, each containing the branch, hostname, and the port no. of a server. The configuration file provided to each server is the same. A sample configuration file for a cluster running on group g01 would look like this:

```
A sp21-cs425-g01-01.cs.illinois.edu 1234
B sp21-cs425-g01-02.cs.illinois.edu 1234
C sp21-cs425-g01-03.cs.illinois.edu 1234
D sp21-cs425-g01-04.cs.illinois.edu 1234
E sp21-cs425-g01-05.cs.illinois.edu 1234
```

A client takes only one argument -- the configuration file (same as above), which provides the required details for connecting to a server when processing a transaction.

## Client Interface and Transactions

At start up, the starts accepting commands from `stdin` (typed in by the user). The user will execute the following commands:

* `BEGIN`: Open a new transaction. The client must connect to a randomly selected server which will coordinate the transaction, and reply "OK" to the user.
* `DEPOSIT server.account amount`: Deposit some amount into an account. The amount will be a positive integer. (You can assume that the value of any account will never exceed 100,000,000.) The account balance should increase by the given amount. If the account was previously unreferenced, it should be created with an initial balance of `amount`. The client should reply with `OK`.

<pre><samp><kbd>DEPOSIT A.foo 10</kbd>
OK
<kbd>DEPOSIT B.bar 30</kbd>
OK
</samp></pre>

* `BALANCE server.account`: The client should display the current balance in
  the given account:

<pre><samp><kbd>BALANCE A.foo</kbd>
A.foo = 10
</samp></pre>

If a query is made to an account that has not previously received a deposit, the client should print `NOT FOUND` and abort the transaction.

* `WITHDRAW server.account amount`: Withdraw some amount from an account. The account balance should
  decrease by the withdrawn amount. The client should reply with `OK` if the operation is successful. If the
  account does not exist (i.e, has never received any deposits), the client should print `NOT FOUND` and abort
  the transaction.


<pre><samp><kbd>BEGIN
WITHDRAW C.baz 5</kbd>
NOT FOUND
</samp></pre>

* `COMMIT`: Commit the transaction, making its results visible to other transactions. The client should reply either with `COMMIT OK` or `ABORTED`, in the case that the transaction had to be aborted during the commit process.
* `ABORT`: Abort the transaction. All updates made during the transaction must be rolled back. The client should reply with `ABORTED` to confirm that the transaction was aborted.

The client will forward each command in a transaction to the randomly selected coordinating server. The coordinating server communicates with the appropriate server (corresponding to the branch in which the account is maintained) sending it the command and receiving the outcome, which it then forwards to the client. The coordinating server handles ABORT and COMMIT commands by appropriately communicating with all other servers. Each server must process each command it receives in a way that satisfies the properties discussed below.


*Notes*:
* You should ignore any commands occuring outside a transaction (other than `BEGIN`).
* You can assume that all commands are using valid format. E.g., you will not see `DEPOSIT A.foo -232` or
  `WITHDRAW B.$#@% 0`.
* A transaction should see its own tentative updates; e.g., if I `DEPOSIT A.foo 10` and then call `BALANCE` on `A.foo` in the same transaction, I should see the deposited amounts.  Whether updates from other transactions are seen depends on whether those transactions are committed and isolation properties, discussed below.
* A server may, as needed, spontaneously abort a transaction while waiting for the next command from the user.
  The user will need to open a new transaction using `BEGIN`. E.g.:

<pre><samp><kbd>BEGIN
DEPOSIT A.foo 10</kbd>
OK
<kbd>DEPOSIT B.bar 30</kbd>
OK
ABORTED
</samp></pre>


## Atomicity

Transactions should execute atomically. In particular, any changes made by a transaction should be rolled
back in case of an abort (initiated either by the user or a server) and all account values should be
restored to their state before the transaction.

## Consistency

As described above, a transaction should not reference any accounts that have not yet received any deposits in a `WITHDRAW` or `BALANCE` command. An additional consistency constraint is that, *at the end of a transaction* no account balance should be negative. If, when a user specifies `COMMIT` any balances are negative, the transaction should be aborted.

<pre><samp><kbd>BEGIN 
DEPOSIT B.bar 20</kbd>
OK
<kbd>WITHDRAW B.bar 30</kbd>
OK
<kbd>COMMIT</kbd>
ABORTED
</samp></pre>

However, it is OK for accounts to have negative balances *during* the transaction, assuming those are eventually resolved:


<pre><samp><kbd>BEGIN 
DEPOSIT B.bar 20</kbd>
OK
<kbd>WITHDRAW B.bar 30</kbd>
OK
<kbd>DEPOSIT B.bar 15</kbd>
OK
<kbd>COMMIT</kbd>
COMMIT OK
</samp></pre>

## Isolation

Your system should support multiple simultaneous clients (up to 10) that execute transactions concurrently. It should guarantee the serializability of the executed transactions. This means that the results should be equivalent to a serial execution of all committed transactions. (Aborted transactions should have no impact on other transactions.) A natural choice is to use either two-phase locking or timestamped concurrency to achieve this (though these are not strict requirements).

You *must* support concurrency between transactions that do not interfere with each other. E.g., if T1 on client 1 executes `DEPOSIT A.x, BALANCE B.y` and then T2 on client 2 executes `DEPOSIT A.w, BALANCE B.z`, the transactions should both proceed without waiting for each other. In particular, using a single global lock (or one lock per server) will not satisfy the concurrency requirements of this MP. You should support read sharing as well, so `BALANCE A.x` executed by two transactions should not be considered interfering.

On the other hand, if T1 executes `DEPOSIT A.x` and T2 executes `BALANCE A.x`, you may delay the execution
of one of the transactions while waiting for the other to complete; e.g., `BALANCE A.x` in T2 may wait to return a response until T1 is committed or aborted.

## Deadlock Resolution

Your system must implement a deadlock resolution strategy. One option is deadlock detection, where the system detects a deadlock and aborts one of the transactions. You can also use concurrency control strategies that avoid deadlocks altogether (e.g. timestamped concurrency).

As discussed earlier, a client can spontaneously display `ABORTED` to the user at any point in time to indicate that the transaction has been aborted. Remember that deadlocks may span multiple servers and clients.

You should not use timeouts as your deadlock detection strategy because transactions will be executed interactively and this will therefore result in too many false positives. Moreover, your deadlock resolution strategy cannot assume that the objects requiring a lock are known apriori (when the transaction begins), as the user will input commands interactively.


**Design document**

1. A detailed explanation of your concurrency control approach. Explain how and where locks are maintained, when they are acquired, and when they are released. If you are using a lock-free strategy, explain the other data structures (timestamps, dependency lists) used in your implementation.

If your algorithm implements a strategy that does not directly follow a concurrency strategy described in the lecture or the literature, you will also need to include an argument for why your strategy ensures serial equivalence of transactions.

2. A description of how transactions are aborted and their actions are rolled back. Be sure to mention how you ensure that other transactions do not use partial results from aborted transactions.

3. Describe how you detect or prevent deadlocks.


**Code guidelines**

The client should only print the responses to the commands or the ABORTED message (as described above). It should not print any additional messages.
The client should support reading transactions from a file redirected to `stdin`. (This would require appropriately handling EOF character at the end.)
Every time a server *commits* any updates to its objects, it should print the balance of all accounts with non-zero values.  
Other than the above, the servers should not print any additional messages.

**Graphs** You do not need to perform any experiment, or plot any graphs for this MP.



## High-level Rubric

- Correct submission format and build instructions (5 points)
- Design Document (15 points)
- Functionality testing (60 points)
    - Atomicity (15 points)
    - Consistency (15 points)
    - Isolation (25 points)
    - Deadlock avoidance (5 points)




