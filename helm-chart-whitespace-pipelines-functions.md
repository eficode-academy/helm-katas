# Helm Chart Whitespace Handling, Pipelines and Functions

## Learning Goals

TODO
- Handle whitespace in helm charts

## Introduction

TODO

Since we are templating `yaml` files we have to be careful with getting the indentation of parameterized values right, so therefore the templating includes functionality for handling whitespace.

### Whitespace Handling with Helm



### Helm Functions

Helm has a number of `functions` available that enable more elaborate templating.

Functions are used in actions and usually take an argument:
```
{{ function argument }}
```
The result of applying the argument to the function will be returned by the action.

A useful and simple example of a function could be to add quotes to a string:

```
shouldBeAString: {{ quote .Values.myString }}
```

Where we assume `myString=FooBar` the result of the function will be `shouldBeAString: "FooBar"`.

### Helm Pipelines

Pipelines allows us to use output of one function as the input of another function:

```
{{ function1 | function2 }}
```

Where the result of function1 is used as the argument for function2, and the result of function2 is returned from the action.

> :bulb: Referencing a values is actually an implicit function!

We can rewrite our quoting example above with a pipeline:
```
shouldBeAString: {{ .Values.myString | quote }}
```
Which will produce the exact same result.

We can use as many functions as we want to in a pipeline.

For example if we wanted to make sure that our string only contains lower case characters, we can use the `lower` function in our pipeline:
```
shouldBeALowerCaseString: {{ .Values.myString | lower | quote }}
```

Which would first change the value of `myString=FooBar` to lowercase, and then add quotes.

The result would be: `shouldBeALowerCaseString: "foobar"`




## Exercise

### Overview

### Step-by-Step

<details>
      <summary>Steps:</summary>
</details>

## Cleanup
