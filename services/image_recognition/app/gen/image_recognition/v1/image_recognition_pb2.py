"""
Protocol Buffers runtime module for image_recognition.v1 (handcrafted).

本ファイルは `buf` / `grpcio-tools` が利用できない環境でも
Python 側の生成物を更新できるよう、Descriptor を動的生成して
メッセージ/サービスクラスを構築します。

proto: schema/proto/image_recognition/v1/image_recognition.proto
- RecognizeImageRequest に optional string place_id = 3 を追加
- threshold は proto3 optional
"""

from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import descriptor_pb2 as _descriptor_pb2
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder

# @@protoc_insertion_point(imports)
_sym_db = _symbol_database.Default()


def _build_file_descriptor() -> bytes:
    fd = _descriptor_pb2.FileDescriptorProto()
    fd.name = "image_recognition/v1/image_recognition.proto"
    fd.package = "image_recognition.v1"
    fd.syntax = "proto3"

    # message HelloRequest { string name = 1; }
    m = fd.message_type.add(); m.name = "HelloRequest"
    f = m.field.add(); f.name = "name"; f.number = 1
    f.label = _descriptor.FieldDescriptor.LABEL_OPTIONAL
    f.type = _descriptor.FieldDescriptor.TYPE_STRING

    # message HelloReply { string message = 1; }
    m = fd.message_type.add(); m.name = "HelloReply"
    f = m.field.add(); f.name = "message"; f.number = 1
    f.label = _descriptor.FieldDescriptor.LABEL_OPTIONAL
    f.type = _descriptor.FieldDescriptor.TYPE_STRING

    # message RecognizeImageRequest {
    #   bytes image_data = 1;
    #   optional float threshold = 2;
    #   optional string place_id = 3;
    # }
    m = fd.message_type.add(); m.name = "RecognizeImageRequest"
    f = m.field.add(); f.name = "image_data"; f.number = 1
    f.label = _descriptor.FieldDescriptor.LABEL_OPTIONAL
    f.type = _descriptor.FieldDescriptor.TYPE_BYTES

    # proto3 optional threshold
    one = m.oneof_decl.add(); one.name = "_threshold"
    f = m.field.add(); f.name = "threshold"; f.number = 2
    f.label = _descriptor.FieldDescriptor.LABEL_OPTIONAL
    f.type = _descriptor.FieldDescriptor.TYPE_FLOAT
    f.proto3_optional = True
    f.oneof_index = 0

    # proto3 optional place_id
    one2 = m.oneof_decl.add(); one2.name = "_place_id"
    f = m.field.add(); f.name = "place_id"; f.number = 3
    f.label = _descriptor.FieldDescriptor.LABEL_OPTIONAL
    f.type = _descriptor.FieldDescriptor.TYPE_STRING
    f.proto3_optional = True
    f.oneof_index = 1

    # message RecognizeImageResponse {
    #   bool is_match = 1;
    #   float similarity_score = 2;
    #   string error_message = 3;
    # }
    m = fd.message_type.add(); m.name = "RecognizeImageResponse"
    f = m.field.add(); f.name = "is_match"; f.number = 1
    f.label = _descriptor.FieldDescriptor.LABEL_OPTIONAL
    f.type = _descriptor.FieldDescriptor.TYPE_BOOL
    f = m.field.add(); f.name = "similarity_score"; f.number = 2
    f.label = _descriptor.FieldDescriptor.LABEL_OPTIONAL
    f.type = _descriptor.FieldDescriptor.TYPE_FLOAT
    f = m.field.add(); f.name = "error_message"; f.number = 3
    f.label = _descriptor.FieldDescriptor.LABEL_OPTIONAL
    f.type = _descriptor.FieldDescriptor.TYPE_STRING

    # message HealthCheckRequest {}
    m = fd.message_type.add(); m.name = "HealthCheckRequest"

    # message HealthCheckResponse { bool healthy = 1; string status = 2; }
    m = fd.message_type.add(); m.name = "HealthCheckResponse"
    f = m.field.add(); f.name = "healthy"; f.number = 1
    f.label = _descriptor.FieldDescriptor.LABEL_OPTIONAL
    f.type = _descriptor.FieldDescriptor.TYPE_BOOL
    f = m.field.add(); f.name = "status"; f.number = 2
    f.label = _descriptor.FieldDescriptor.LABEL_OPTIONAL
    f.type = _descriptor.FieldDescriptor.TYPE_STRING

    # service ImageRecognitionService
    s = fd.service.add(); s.name = "ImageRecognitionService"
    mth = s.method.add(); mth.name = "Hello"
    mth.input_type = ".image_recognition.v1.HelloRequest"
    mth.output_type = ".image_recognition.v1.HelloReply"
    mth = s.method.add(); mth.name = "RecognizeImage"
    mth.input_type = ".image_recognition.v1.RecognizeImageRequest"
    mth.output_type = ".image_recognition.v1.RecognizeImageResponse"
    mth = s.method.add(); mth.name = "HealthCheck"
    mth.input_type = ".image_recognition.v1.HealthCheckRequest"
    mth.output_type = ".image_recognition.v1.HealthCheckResponse"

    return fd.SerializeToString()


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(_build_file_descriptor())

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'image_recognition.v1.image_recognition_pb2', _globals)

# @@protoc_insertion_point(module_scope)
