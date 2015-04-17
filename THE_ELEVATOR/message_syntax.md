Message syntax
==============

This defines the syntax for all messages to be passed between lifts.

New orders
----------
- Syntax:	"neworder",<floor>,<button>
- Example:	neworder,2,0
- Meaning:	New order at floor 2 in direction 0 (up)

Cost
----
- Syntax:	"cost",<cost>,<floor>,<button>
- Example:	cost,13,2,0
- Meaning:	My cost is 13 for order at floor 2 in direction 0 (up)

I'm alive
---------
- Syntax:	"alive"
- Example:	alive
- Meaning:	I'm alive!

Order completed
---------------
- Syntax:	"ordercomplete",<floor>,<button>
- Example:	ordercomplete,2,0
- Meaning:	I've completed the order at floor 2 in direction 0 (up)
