# Protobuf Decompiler


<p align="center">
  <img src="https://github.com/MUCZ/protobuf-decompiler/blob/main/doc/img.png">
</p>

The Protobuf Decompiler is a CLI designed to help you recover the original `*.proto` file from a protobuf binary descriptor.

In most supported languages, you can locate the raw descriptor of the original protobuf definition within the generated template codes. This descriptor usually takes the form of a binary string (check the next section for examples).

Please note that the Protobuf Decompiler cannot restore comments since they are discarded by `protoc`.

Here's how to use it:
```sh
make
./protodec examples/helloworld.pb.go
```

The output should resemble the following:
```proto
syntax = "proto3";

option java_multiple_files = true;
option java_package = "io.grpc.examples.helloworld";
option java_outer_classname = "HelloWorldProto";
option objc_class_prefix = "HLW";

package helloworld;

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply) {}
  rpc SayHelloStreamReply (HelloRequest) returns (stream HelloReply) {}
  rpc SayHelloBidiStream (stream HelloRequest) returns (stream HelloReply) {}
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```

# Protobuf Descriptor Example

This section provides examples of the program's input, specifically the protobuf descriptors. In most cases, these descriptors are encoded in the generated template code as literal binary arrays. Additionally, you can extract them from the generated binary executable if available.

Python
```py
# helloworld_pb2.py
DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x10helloworld.proto\x12\nhelloworld\"\x1c\n\x0cHelloRequest\x12\x0c\n\x04name\x18\x01 \x01(\t\"\x1d\n\nHelloReply\x12\x0f\n\x07message\x18\x01 \x01(\t2I\n\x07Greeter\x12>\n\x08SayHello\x12\x18.helloworld.HelloRequest\x1a\x16.helloworld.HelloReply\"\x00\x42\x36\n\x1bio.grpc.examples.helloworldB\x0fHelloWorldProtoP\x01\xa2\x02\x03HLWb\x06proto3')
```

Go
```go
// helloworld.pb.go
var file_examples_helloworld_helloworld_helloworld_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f,
	0x77, 0x6f, 0x72, 0x6c, 0x64, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64,
	0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0a, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x22, 0x22, 0x0a,
    ... ... 
}
```

Java
```Java
// HelloWorldProto.java
  static {
    java.lang.String[] descriptorData = {
      "\n\020helloworld.proto\022\nhelloworld\"\034\n\014HelloR" +
      "equest\022\014\n\004name\030\001 \001(\t\"\035\n\nHelloReply\022\017\n\007me" +
      "ssage\030\001 \001(\t2I\n\007Greeter\022>\n\010SayHello\022\030.hel" +
      "loworld.HelloRequest\032\026.helloworld.HelloR" +
      "eply\"\000B6\n\033io.grpc.examples.helloworldB\017H" +
      "elloWorldProtoP\001\242\002\003HLWb\006proto3"
    };
```

C++
```c++
// helloworld.pb.cc
const char descriptor_table_protodef_helloworld_2eproto[] PROTOBUF_SECTION_VARIABLE(protodesc_cold) =
  "\n\020helloworld.proto\022\nhelloworld\"\236\001\n\014Hello"
  "Request\022\014\n\004name\030\001 \001(\t\022\021\n\ttipe_data\030\002 \003(\005"
  "\022\020\n\010size_all\030\003 \001(\005\022\020\n\010size_arr\030\004 \003(\005\022\021\n\t"
  "timeEpoch\030\005 \003(\005\022\020\n\010datablob\030\006 \001(\014\022\021\n\tfla"
  "gparam\030\007 \001(\005\022\021\n\tdone_send\030\010 \001(\005\"\234\001\n\nHell"
  "oReply\022\014\n\004name\030\001 \001(\t\022\021\n\ttipe_data\030\002 \003(\005\022"
  "\020\n\010size_all\030\003 \001(\005\022\020\n\010size_arr\030\004 \003(\005\022\021\n\tt"
  "imeEpoch\030\005 \003(\005\022\020\n\010datablob\030\006 \001(\014\022\021\n\tflag"
  "param\030\007 \001(\005\022\021\n\tdone_send\030\010 \001(\0052J\n\010Greete"
  "r2\022>\n\010SayHello\022\030.helloworld.HelloRequest"
  "\032\026.helloworld.HelloReply\"\000b\006proto3"
  ;
```

C#
```c#
// Helloworld.cs
    static HelloworldReflection() {
      byte[] descriptorData = global::System.Convert.FromBase64String(
          string.Concat(
            "ChBoZWxsb3dvcmxkLnByb3RvEgpoZWxsb3dvcmxkIhwKDEhlbGxvUmVxdWVz",
            "dBIMCgRuYW1lGAEgASgJIh0KCkhlbGxvUmVwbHkSDwoHbWVzc2FnZRgBIAEo",
            "CTJJCgdHcmVldGVyEj4KCFNheUhlbGxvEhguaGVsbG93b3JsZC5IZWxsb1Jl",
            "cXVlc3QaFi5oZWxsb3dvcmxkLkhlbGxvUmVwbHkiADJQCgxNdWx0aUdyZWV0",
            "ZXISQAoIU2F5SGVsbG8SGC5oZWxsb3dvcmxkLkhlbGxvUmVxdWVzdBoWLmhl",
            "bGxvd29ybGQuSGVsbG9SZXBseSIAMAFCNgobaW8uZ3JwYy5leGFtcGxlcy5o",
            "ZWxsb3dvcmxkQg9IZWxsb1dvcmxkUHJvdG9QAaICA0hMV2IGcHJvdG8z"));
```

Rust 
```rs

// todo: add rust example
```
