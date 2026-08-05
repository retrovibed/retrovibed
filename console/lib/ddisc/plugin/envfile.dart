import 'dart:convert';
import 'dart:ffi';
import 'package:ffi/ffi.dart' as ffi;
import 'package:retrovibed/retrovibed.dart' as retrovibed;
import 'package:retrovibed/uuidx.dart' as uuidx;

/// Client-side view of a .env sidecar file's KEY=VALUE entries, with the
/// nearest preceding (or trailing, same-line) comment attached as help
/// text. The server treats .env files as an opaque byte blob - parsing and
/// re-serializing is entirely a client concern.
///
/// The real parsing rules live in Go (retroapi/envfile) and are exposed to
/// the console app over FFI (console/retrovibedbind's envfile_parse
/// export, bound in console/lib/retrovibed/gen.dart via ffigen); [parseEnv]
/// calls through that bridge so the rules only need to be written once.
class Variable {
  // id is the digest at the time this row was first created (parsed from
  // content, or freshly added blank); digest is always the *current*
  // digest. id == digest means the row hasn't been edited since creation -
  // used both for dirty-tracking and, via [id], as a stable Flutter widget
  // key that survives value/hint edits.
  final String id, digest, key, value, hint;

  Variable(this.key, this.value, this.hint) : id = _digestOf(key, value, hint), digest = _digestOf(key, value, hint);

  Variable._edit(this.id, this.key, this.value, this.hint) : digest = _digestOf(key, value, hint);

  static String _digestOf(String key, String value, String hint) => uuidx.md5x('$key$value$hint');

  bool get unchanged => id == digest;

  Variable copyWith({String? key, String? value, String? hint}) =>
      Variable._edit(id, key ?? this.key, value ?? this.value, hint ?? this.hint);

  @override
  bool operator ==(Object other) => other is Variable && other.key == key;

  @override
  int get hashCode => Object.hash(key, value, hint);

  @override
  String toString() => "Variable($key, $value, $hint)";
}

List<Variable> parseEnv(String content) {
  final contentPtr = content.toNativeUtf8().cast<Char>();
  final resultPtr = retrovibed.bridge.envfile_parse(contentPtr);
  ffi.calloc.free(contentPtr);

  final json = resultPtr.cast<ffi.Utf8>().toDartString();
  ffi.calloc.free(resultPtr);

  final decoded = jsonDecode(json) as List<dynamic>? ?? [];
  return decoded
      .map((e) => e as Map<String, dynamic>)
      .map((e) => Variable(e['key'] as String, e['value'] as String, e['hint'] as String))
      .toList();
}

String _quoteIfNeeded(String value) {
  if (value.isEmpty) return value;
  if (RegExp(r'[\s#"]').hasMatch(value)) {
    return '"${value.replaceAll('"', '\\"')}"';
  }
  return value;
}

/// Serializes [variables] as the complete, authoritative content of a .env
/// file - every line comes from [variables], nothing else survives. There
/// is no base string to go stale against, so this can't resurrect a
/// removed variable or clobber one it never knew about.
String serializeEnv(List<Variable> variables) {
  if (variables.isEmpty) return '';

  final lines = variables.map((v) {
    final quoted = _quoteIfNeeded(v.value);
    final comment = v.hint.isNotEmpty ? ' # ${v.hint}' : '';
    return '${v.key}=$quoted$comment';
  });
  return '${lines.join('\n')}\n';
}
