import 'dart:convert';
import 'dart:ffi';
import 'dart:io';
import 'package:ffi/ffi.dart' as ffi;
import 'package:console/retrovibed/gen.dart' as lib;

String _path() {
  final files = [File("build/nativelib/retrovibed.so")];
  final found = files.firstWhere(
    (v) => v.existsSync(),
    orElse: () => File("lib/retrovibed.so"),
  );
  return found.path;
}

final bridge = lib.DaemonBridge(DynamicLibrary.open(_path()));

String bearer_token() {
  return _convertstring(bridge.authn_bearer());
}

String bearer_token_host(String hostname) {
  return _convertstring(
    bridge.authn_bearer_host(hostname.toNativeUtf8().cast<Char>()),
  );
}

String public_key() {
  return _convertstring(bridge.public_key());
}

List<String> ips() {
  final List<dynamic> res = jsonDecode(_convertstring(bridge.ips()));
  return res.whereType<String>().toList();
}

void daemon() {
  String args = jsonEncode(["daemon", "--no-mdns"]);
  bridge.daemon(args.toNativeUtf8().cast<Char>());
}

String _convertstring(Pointer<Char> charPointer) {
  try {
    return charPointer.cast<ffi.Utf8>().toDartString();
  } finally {
    ffi.calloc.free(charPointer);
  }
}
