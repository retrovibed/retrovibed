// ignore_for_file: experimental_member_use

import 'dart:convert';
import 'package:crypto/crypto.dart' as crypto;
import 'package:uuid/uuid.dart' as uuid;
import 'package:uuid/data.dart' show V7Options;

uuid.UuidValue fromString(String v) => v.isEmpty ? uuid.Namespace.nil.uuidValue : uuid.UuidValue.fromString(v);

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

// at, when given, pins the embedded timestamp instead of using the current
// time - primarily for tests that need deterministically-ordered sids
// without depending on real wall-clock timing between calls.
String v7({DateTime? at}) =>
    at == null ? uuid.Uuid().v7() : uuid.Uuid().v7(config: V7Options(at.millisecondsSinceEpoch, null));

DateTime v7timestamp(String s) {
  final bytes = fromString(s).toBytes();
  final ms = (bytes[0] << 40) |
      (bytes[1] << 32) |
      (bytes[2] << 24) |
      (bytes[3] << 16) |
      (bytes[4] << 8) |
      bytes[5];
  return DateTime.fromMillisecondsSinceEpoch(ms, isUtc: true);
}

String random() => uuid.Uuid().v4();

String withSuffix(int v) => '00000000-0000-0000-0000-${v.toString().padLeft(12, '0')}';

bool prefix(String p, v) => v.startsWith(p);

String md5x(String s) => crypto.md5.convert(utf8.encode(s)).toString();
