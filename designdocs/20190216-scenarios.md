# 201902116 Scenarios

**DRAFT**

This is a collection of the scenarios and what the orchestrator does in those cases.

General assumptions in here:

* IP Address Order: A, X, B, Y, C, Z, D

## Simple Startup

**Test `001-simple-start`**

* Given:
  I have no established etcd cluster nor etcd-controller group.
  I have 3 fresh nodes running etcd-controller (A,B,C).
  I no nodelist (nodelist is empty)
* When:
  I have a nodelist that includes the 3 nodes A,B,C.
* Then:
  I should have a three node etcd cluster running on A,B,C.
  I should have a three node etcd-controller group running on A,B,C.
  I should have the orchestrator running on A.

## Adding a new node (no orchestrator change)

**Test `002-adding-a-new-node`**

* Given:
  I have an established 3 node etcd cluster + etcd-controller group (A,B,C).
  I have B as the orchestrator.
  I have a nodelist that includes A,B,C.
  I have a fresh node ready to run etcd-controller (D).
* When:
  I start D.
  I update the nodelist to be A,B,C,D.
* Then:
  I should have a four node etcd cluster running on A,B,C,D.
  I should have a four node etcd-controller group running on A,B,C,D.
  I should have A as the orchestrator.

## Removing a single etcd node (no orchestrator change)

**Test `003-removing-single-etcd-node`**

* Given:
  I have an established 3 node etcd cluster + etcd-controller group.
  I have a nodelist that includes the 3 nodes A,B,C.
* When:
  I update the nodelist to include just A,B.
* Then:
  I should have a two node etcd cluster (A,B) with a 3 node etcd-controller group (A,B,C).
  Node C is in "watching" state.

## Adding a new node (orchestrator change)

* Given:
  I have an established 3 node etcd cluster + etcd-controller group (B,C,D).
  I have B as the orchestrator.
  I have a nodelist that includes B,C,D.
  I have a fresh node ready to run etcd-controller (A).
* When:
  I start A.
  I update the nodelist to be A,B,C,D.
* Then:
  I should have a four node etcd cluster running on A,B,C,D.
  I should have a four node etcd-controller group running on A,B,C,D.
  I should have A as the orchestrator.

## Removing a single etcd-controller node

* Given:
  I have an established 2 node etcd cluster (A,B), and three node etcd-controller group (A,B,C).
  I have a nodelist that includes 2 nodes (A,B).
* When:
  I shutdown node C.
* Then:
  I should have a two node etcd cluster (A,B) with a 2 node etcd-controller group (A,B).

## Moving orchestrator

* Given:
  I have an established 3 node etcd cluster (A,B,C) and three node etcd-controller group (A,B,C).
  I have a nodelist that includes the 3 nodes A,B,C.
* When:
  I update the node list to include just B,C.
* Then:
  I should have a two node etcd cluster (B,C) with a 2 node etcd-controller group (B,C).
  B should be the orchestrator.

## Failing orchestrator node / not fixed

* Given:
  I have an established 3 node etcd cluster (A,B,C) and three node etcd-controller group (A,B,C).
  I have a nodelist that includes the 3 nodes A,B,C.
* When:
  I kill node A.
* Then:
  I should have a three node etcd cluster (A,B,C) with one failed node (A).
  I should have a two node etcd-controller group (B,C).
  No orchestrator should be running.

## Failing orchestrator node / fixed

* Given:
  I have an established 3 node etcd cluster (A,B,C) and three node etcd-controller group (A,B,C).
  I have a nodelist that includes the 3 nodes A,B,C.
* When:
  I kill node A. I update the nodelist to only include B,C.
* Then:
  I should have a two node etcd cluster (B,C).
  I should have a two node etcd-controller group (B,C).
  B should be the orchestrator.

## Adding a single node to an unmanaged cluster

* Given:
  I have a running cluster outside of etcd-controller (X,Y,Z).
  I have a fresh node ready to run etcd-controller (A).
  I have a nodelist with A,X(unmanaged),Y(unmanaged),Z(unmanaged).
* When:
  I start up an instance of etcd-controller on A.
* Then:
  I have a four node etcd cluster (A,X,Y,Z) with the data carried over from (X,Y,Z).

## Stopping an unmanaged node - no change

* Given:
  I have a running etcd cluster A,B,C,X(unmanaged).
  I have a running etcd-controller group A,B,C.
  I have A as the orchestrator.
  I have a nodelist with A,B,C,X(unmanaged).
* When:
  I stop node X.
* Then:
  I have a four node etcd cluster (A,B,C,X(unmanaged)) with one failed node (X).
  I have a three node etcd-controller group (A,B,C).

## Stopping an unmanaged node missing from the nodelist - remove

* Given:
  I have a running etcd cluster A,B,C,X(unmanaged).
  I have a running etcd-controller group A,B,C.
  I have A as the orchestrator.
  I have a nodelist with A,B,C.
* When:
  I stop node X.
* Then:
  I have a three node etcd cluster (A,B,C).
  I have a three node etcd-controller group (A,B,C).

## Removing an unmanaged running node from the nodelist - no change

* Given:
  I have a running etcd cluster A,B,C,X(unmanaged).
  I have a running etcd-controller group A,B,C.
  I have A as the orchestrator.
  I have a nodelist with A,B,C,X(unmanaged).
* When:
  I update the nodelist to be A,B,C.
* Then:
  I have a four node etcd cluster (A,B,C,X(unmanaged)).
  I have a three node etcd-controller group (A,B,C).

## Removing an unmanaged stopped node from the nodelist - remove

* Given:
  I have a running etcd cluster A,B,C,X(unmanaged) with failed node X.
  I have a running etcd-controller group A,B,C.
  I have A as the orchestrator.
  I have a nodelist with A,B,C,X(unmanaged).
* When:
  I update the nodelist to be A,B,C.
* Then:
  I have a three node etcd cluster (A,B,C).
  I have a three node etcd-controller group (A,B,C).

