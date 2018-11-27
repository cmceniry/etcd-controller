# 20181126 Conductor RPC

** DRAFT **

This describes the baseline RPC for the conductor.
The RPC is implemented as GRPC.

Since the conductor is supposed to be bundled with the driver in the final controller, it will register with the same GRPC server then.
Currently, it'll need its own.

# RPC Calls

* GetNodeStatus:
  Asks for the current status of a specific named node

  Parameters: NodeName

  Returns: NodeInfo structure
  * Name
  * PeerURL
  * Status: Unknown, Running, Stopped, Failed