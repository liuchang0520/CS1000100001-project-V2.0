# General Overview

This project simulates a map reduce system:

Different processes simulate master and workers. Task assignment could be achieved by RPC between master and worker processes.

Implemented very basic fault tolerance mechanism by Go Channel: https://gobyexample.com/channels , in which idle worker process ports are stored.

More detailed idea/implementation is in proposal.txt.

# Steps to run this project in Shell:





