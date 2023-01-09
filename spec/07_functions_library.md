# Functions Library

The Wars Module is deployed with a built-in set of libraries for commonly-used algorithmic pricing and reserve functions. It also includes algorithmic application logic and features, such as _Augmented Waring_. Additional functions can be added to the Library through SDK updates. This requires a formal process of governance to approve updates, to assure the integrity of these functions.

## Function Types

The following function types will be included in the standard Wars SDK Module:

* Power \(exponential\)
* Logistic \(sigmoidal\)
* Constant Product \(swapper\)

  Algorithmic Applications include:

* Alpha Wars \(Risk-adjusted waring\)
* Innovation Wars \(offers war shareholders contingent rights to future IP rights and/or revenues\)
* Impact Wars \(offers war shareholders contingent rights to success-based outcomes payments and/or rewards\)

### Exponential Function \(power\)

Function \(used as pricing function\):

![power function price](../.gitbook/assets/power1%20%281%29.png)

Integral \(used as reserve function\):

![power function reserve](../.gitbook/assets/power2%20%281%29.png)

### Logistic Function \(sigmoid\)

Function \(used as pricing function\):

![sigmoid function price](../.gitbook/assets/sigmoid1.png)

Integral \(used as reserve function\):

![sigmoid function reserve](../.gitbook/assets/sigmoid2.png)

### Augmented Waring Curves \(augmented\)

Initial reserve:

![initial reserve](../.gitbook/assets/augmented1%20%281%29.png)

Initial supply:

![initial supply](../.gitbook/assets/augmented2%20%281%29.png)

Constant power function invariant:

![constant power](../.gitbook/assets/augmented3%20%281%29.png)

Invariant function:

![invariant function](../.gitbook/assets/augmented4%20%281%29.png)

Pricing function:

![pricing function](../.gitbook/assets/augmented5%20%281%29.png)

Reserve function:

![reserve function](../.gitbook/assets/augmented6%20%281%29.png)

Ref: [https://medium.com/giveth/deep-dive-augmented-waring-curves-3f1f7c1fa751](https://medium.com/giveth/deep-dive-augmented-waring-curves-3f1f7c1fa751)

### Constant Product Function \(swapper\)

Reserve function:

![swapper function](../.gitbook/assets/swapper%20%281%29.png)

