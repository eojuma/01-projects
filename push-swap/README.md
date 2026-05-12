## circularly sorting

During the sorting process, Stack A may not always be strictly sorted top to bottom.

Instead, it is often circularly sorted, meaning:

Elements are in order if you consider the stack as circular.

The smallest number may not be at the top yet, but all numbers are in relative order.

Why circularly sorted?

Final rotations (ra or rra) at the end bring the smallest number to the top, making A fully sorted.

Example

Initial stack A:

```
A (top → bottom): 1 4 7 8
B (empty)
```

Insert 6 from B:

Find the smallest number in A bigger than 6 → 7

Rotate A until 7 is at top → ra 2 times

```
A: 7 8 1 4
pa → push 6 on top:
A: 6 7 8 1 4
```


Notice that A is not strictly sorted top-to-bottom.

But it is circularly sorted:
1 4 6 7 8 if rotated appropriately.

Final rotation to place the smallest element (1) on top:

```
ra → 7 8 1 4 6
ra → 8 1 4 6 7
ra → 1 4 6 7 8
```

* Push elements from A → B until only 3 remain in A.

* Sort the remaining 3 in A with minimal moves (sa, ra, rra).

* For each element in B:

* Find the optimal insertion index in A.

* Rotate A and B minimally to align insertion.

* pa to move the element to A.

* Repeat until B is empty.

* Rotate A so that the smallest number is on top.


essentially circularly sorting is this once we have sorted stack a with all three elements we will go to stack b get the first element then loop thru stack b and find out which interger is the smallest in stack a but is bigger than our element in stack b then we will rotate the number of times the index of the item in stack a bringing the item to the top of stack a then we will pa and put the element in stack b to the top of stack a as sorted then will do this cmany times till stack b is empty then we will run rotate or rra till we get the final stack a