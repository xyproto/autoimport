package com.example.demo;

public class Example {
    public static void main(String[] args) {
        List<String> names = new ArrayList<>();
        names.add("Alice");
        names.add("Bob");

        Map<String, Integer> ageMapping = new HashMap<>();
        ageMapping.put("Alice", 30);
        ageMapping.put("Bob", 25);

        for (String name : names) {
            System.out.println(name + " is " + ageMapping.get(name) + " years old.");
        }
    }
}
