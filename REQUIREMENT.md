## Client Interface and Transactions

At start up, the client starts accepting commands from `stdin` (typed in by the user). The user will execute the following commands:

* `BEGIN`: Open a new transaction. The client must connect to a randomly selected server which will coordinate the transaction, and reply "OK" to the user.

<pre><samp><kbd>BEGIN</kbd>
OK
</samp></pre>

* `DEPOSIT server.account amount`: Deposit some amount into an account. The amount will be a positive integer. (You can assume that the value of any account will never exceed 100,000,000.) The account balance should increase by the given amount. If the account does not exist (in other words, if it is the first time a deposit operation on the account has been issued), the account should be created with an initial balance of `amount`. The client should reply with `OK`.

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

If a query is made to an account that has not previously received a deposit, the server system must abort the transaction and inform the client. The client should then print `NOT FOUND, ABORTED`.

<pre><samp><kbd>BEGIN</kbd>
OK
<kbd>BALANCE C.baz 5</kbd>
NOT FOUND, ABORTED
</samp></pre>

* `WITHDRAW server.account amount`: Withdraw some amount from an account. The account balance should
  decrease by the withdrawn amount. The client should reply with `OK` if the operation is successful. If the
  account does not exist (i.e, has never received any deposits), the server system must abort the transaction and inform the client. The client should then print `NOT FOUND, ABORTED`.


<pre><samp><kbd>BEGIN</kbd>
OK
<kbd>WITHDRAW B.bar 5</kbd>
OK
<kbd>WITHDRAW C.baz 5</kbd>
NOT FOUND, ABORTED
</samp></pre>

* `COMMIT`: Commit the transaction, making its results visible to other transactions. The client should reply either with `COMMIT OK` or `ABORTED`, in the case that the transaction had to be aborted during the commit process.
* `ABORT`: Abort the transaction. All updates made during the transaction must be rolled back. The client should reply with `ABORTED` to confirm that the transaction was aborted.

The client will forward each command in a transaction to the randomly selected coordinating server. The coordinating server communicates with the appropriate server (corresponding to the branch in which the account is maintained) by sending it the command and receiving the outcome, which it then forwards to the client. The coordinating server handles ABORT and COMMIT commands by appropriately communicating with all other servers. Each server must process each command it receives in a way that satisfies the properties discussed below.

