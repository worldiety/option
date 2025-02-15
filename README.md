# option

[![Go Reference](https://pkg.go.dev/badge/github.com/worldiety/enum.svg)](https://pkg.go.dev/github.com/worldiety/option)


Package option provides a missing type to allow modelling optional facts in cases where the natural zero value of a type is not sufficient to express this fact.
The need for this is not a problem caused by the expectation, that other language just provides this, but because

* using pointers to express optionality, has not the same semantics as pointers - and is often not even acceptable when "explicitly" modelling  memory layouts. 
* There is no sense in using a pointer to a pointer type and omitting the option fact
* Pointers introduce more questions around the overloaded meaning of the zero value of a pointer type. 
* Pointers are also used to express a kind of ownership concerns, so that makes the water even more muddy.
* Semantics with iterators are entirely broken, because there is no iter.Seq3 (T,bool,error) you would need to resort to (*T,error) where *T is for optional cases *T which is against the semantic guidelines, that the tupel of (*T,error) means that if error is nil, *T is not. 


To shorten up that discussion, even the standard library provides a generic helper for this in form of the [sql.Null](https://pkg.go.dev/database/sql#Null) type.
However, the usage of this type cannot be generally recommended, because:

* exposing the Valid field and the wrapped T in public fields violates the concept of information hiding and
intended consistency of the author, which returns the Option. We experimented with a copy of that and have
seen a regular misconception, especially for new programmers, which started to modify the value and the
valid flag causing all sorts of confusion, which the usage of the type should have solved, actually.
Nonetheless, code using an option type makes the fact a lot better readable and when mixing pointer types
if they don't make sense, especially when mixing with ownership concerns.
* the sql package indicates a persistence concern. It is against our architecture guidelines (and most others) to import
persistent dependencies into a domain core.
* there is no transparent JSON Marshal/Unmarshal support using the natural absent or null representation
* there is no omitzero support

Why not using other libaries like samber/mo?

* some do not express it in a clear way
* have missing json support
* introduce zero values on accessors, which is the entirely wrong thing. We would not need the Option at all.
* provide no lazy functors
* unneeded dependencies