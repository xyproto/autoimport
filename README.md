# importmatcher

Given source code, class names can be found in available `.jar` files, and import statements can be generated, for Java and for Kotlin.

#### Features and limitation

* Searches directories of .jar files for class names.
* Given the start of the class name, searches for the matching shortest class, and also returns the import path (like `java.io.*`).
* Intended to be used for simple autocompletion of class names.

#### General info

* Version: 0.0.1
* License: BSD-3
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
