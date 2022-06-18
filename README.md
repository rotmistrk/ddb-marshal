# ddb-marshal
Marshal/Unmarshal AWS Dynamo DB entries to GoLang structs

# Usage

## 1. Annotate struct with ddb tag

```go
type Entry struct {
    Name string `ddb:"name,required"`
    Value uint64 `ddb:"value,required"`
    Timestamp time.Time `ddb:"ts"`
    Items map[string]string `ddb:"items"`
}
```

Use public names (capitalized names) for the fields 

## 2. Create Marshaller

```go
marshaller := NewMarshaller()
```

## 3. Reading: Unmarshal results of ddb read API

```go
if output, err := dynamodb.GetItem(input); err == nil {
    var entry Entry
    err = marshaller.Unmarshal(&entry, output.Item)
}
```

### 3.1 Get unmarshalled fields 

These can be:
1. Checked to be empty, or
2. Preserved to persist on consequent update of the entry
3. Passed to unmarshaller again with different data structure for data of "specialization" of the entry

```go
unmarshalled, err := marshaller.GetUnmarshalledFields(&entry, output.Item)
```

## 3.3 Marshal structure for insert/update operation

```go
var input dynamodb.PutItemInput
input.Item, _ = marshaller.Marshal(&entry)
for k, v := range unmarshalled {
    input.Item[k] = v
}
```

## Field tags

Minimal support:

1. comma-separated 
2. first element is name, should follow DDB requirements
3. other entries may be any
4. if one of them is "required", there is minimal validation on the value to be present during unmarshal
5. future extensions are possible, for example HashKet/RangeKey specifications, GSI/LSI specifications 



## Field types support

Many types are supported, but far from covering any variant

Supported:
1. Primitive types (string, ints, floats, boolean)
2. time.Time as unix time (numeric) - specifically for TTL field support
3. Arrays of the primitives and time.Time
4. Byte arrays
5. Arrays of byte arrays
6. Maps with string as key and primitives or time.Time as a value


# BUGS

1. No default behavior (required/optional)
2. No transparent encryption support

# TODO

1. Separate set of classes to create table, GSI, LSI based on the tags
2. Separate set of classes to query, scan, and update table
