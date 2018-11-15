# 20181114 Simple Driver RPC

**DRAFT**

This describes the general RPC calls used 
Simple Driver implements the RPC as GRPC.

## Note about GRPC/H2 TLS requirement

Since we're already producing certificates for etcd, we will just use what is supplied in Simple Driver.

# RPC Calls

* InitializeCluster:
  Tells a node to start etcd as a new cluster by itself.
  If etcd is currently running or DataDir is not empty, this will return failure.

  Parameters: empty

  Returns: Success bool, ErrorMessage string

  Optional argument (TBD) of initializing from a snapshot.

* JoinCluster:
  Tells a node to start etcd and join an existing cluster.
  If etcd is currently running or DataDir is not empty, this will return failure.

  Parameters: Array of existing cluster node peer address:port

  Returns: Success bool, ErrorMessage string

* Start:
  Tells a node to start etcd with no discovery information.
  This expects the cluster to exist already.
  If etcd is currently running, this will still return success (with error message).

  Parameters: empty

  Returns: Success bool, ErrorMessage string

* Stop:
  Tells a node to shutdown etcd.
  If etcd is currently not running, this will still return success (with error message).

  Parameters: empty

  Returns: Success bool, ErrorMessage string
  
* Purge:
  Tells a node to delete etcd's DataDir.
  Will error if etcd is currently running.

  Parameters: empty

  Returns: Success bool, ErrorMessage string
