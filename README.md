# Autoimport

Given source code, class names can be found in available `.jar` files, and import statements can be generated, for Java and for Kotlin.

This currently only works for OpenJDK 8, not OpenJDK 11 and beyond.

Includes the `autoimport` utility for looking up packages, given the start of a class name.

Example use (with OpenJDK 8 installed):

    $ ./autoimport FilePe
    import java.io.*; // FilePermissionCollection
    import java.io.*; // FilePermission
    import sun.security.tools.policytool.*; // FilePerm
    import net.rubygrapefruit.platform.*; // FilePermissionException

Example use (with OpenJDK 19 and openjdk-src installed):

    $ ./autoimport -e FileSystem
    import java.io.*; // FileSystem

#### Features and limitation

* Searches directories of `.jar` files for class names.
* Given the start of the class name, searches for the matching shortest class, and also returns the import path (like `java.io.*`).
* Also searches `*/lib/src.zip` files, if found.
* Intended to be used for simple autocompletion of class names.

#### General info

* Version: 1.0.0
* License: BSD-3
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
