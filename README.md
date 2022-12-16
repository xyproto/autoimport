# importmatcher

# A WORK IN PROGRESS

Given source code, class names can be found in available `.jar` files, and import statements can be generated, for Java and for Kotlin.

## Only tested on Arch Linux! It's a work in progress for other platforms!

Includes the `w00t` utility for looking up packages, given the start of a class name.

Example use:

```
$ time ./w00t FilePe
import java.io.*; // FilePermissionCollection
import java.io.*; // FilePermission
import sun.security.tools.policytool.*; // FilePerm
import net.rubygrapefruit.platform.*; // FilePermissionException
./w00t FilePe  0,21s user 0,04s system 324% cpu 0,076 total
```

#### Features and limitation

* Searches directories of .jar files for class names.
* Given the start of the class name, searches for the matching shortest class, and also returns the import path (like `java.io.*`).
* Intended to be used for simple autocompletion of class names.

#### General info

* Version: 0.0.1
* License: BSD-3
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
