package com.example.demo

// StringList is a mutable list of Strings
class StringList() {
    private var elements: MutableList<String> = mutableListOf()

    public fun add(element: String) {
        elements.add(element)
    }

    operator fun iterator(): Iterator<String> {
        return elements.iterator()
    }
}

// ArrayList is a type alis for StringList, on purpose, to test that
// a java.util.ArrayList import will NOT be added automatically.
typealias ArrayList = StringList

fun main() {
    var names = ArrayList()
    names.add("Alice")
    names.add("Bob")

    var ageMapping = HashMap<String, Int>()
    ageMapping["Alice"] = 30
    ageMapping["Bob"] = 25

    for (name in names) {
        println("$name is ${ageMapping.get(name)} years old.")
    }
}
