// ignore_for_file: experimental_member_use

import 'dart:convert';
import 'package:crypto/crypto.dart' as crypto;
import 'package:uuid/uuid.dart' as uuid;

uuid.UuidValue fromString(String v) =>
    v.isEmpty ? uuid.Namespace.nil.uuidValue : uuid.UuidValue.fromString(v);

String max() => uuid.Namespace.max.value;
String min() => uuid.Namespace.nil.value;

bool isMin(uuid.UuidValue v) => v == uuid.Namespace.nil.uuidValue;
bool isMax(uuid.UuidValue v) => v == uuid.Namespace.max.uuidValue;
bool isMinMax(uuid.UuidValue v) => isMin(v) || isMax(v);

T pattern<T>(String v, T min, T max, T value) {
  final _v = fromString(v);
  return switch (_v) {
    _ when isMin(_v) => min,
    _ when isMax(_v) => max,
    _ => value,
  };
}

String random() => uuid.Uuid().v4();

String withSuffix(int v) =>
    '00000000-0000-0000-0000-${v.toString().padLeft(12, '0')}';

bool prefix(String p, v) => v.startsWith(p);

String md5x(String s) => crypto.md5.convert(utf8.encode(s)).toString();
