import 'dart:convert';
import 'dart:ffi';
import 'package:ffi/ffi.dart' as ffi;
import 'package:console/retrovibed/gen.dart' as lib;

final bridge = lib.DaemonBridge(DynamicLibrary.open("lib/retrovibed.so"));

String bearerToken() {
  return _convertstring(bridge.authn_bearer());
}

String public_key() {
  return _convertstring(bridge.public_key());
}

List<String> ips() {
  return [_convertstring(bridge.ips())];
}

void daemon() {
  String args = jsonEncode(["daemon"]);
  bridge.daemon(args.toNativeUtf8().cast<Char>());
}

String _convertstring(Pointer<Char> charPointer) {
  try {
    return charPointer.cast<ffi.Utf8>().toDartString();
  } finally {
    ffi.calloc.free(charPointer);
  }
}
