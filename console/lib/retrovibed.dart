import 'dart:convert';
import 'dart:ffi';
import 'dart:io';
import 'dart:typed_data';
import 'package:ffi/ffi.dart' as ffi;
import 'package:flutter/widgets.dart';
import 'package:media_kit/media_kit.dart';
import 'package:retrovibed/retrovibed/gen.dart' as lib;
import 'package:retrovibed/caching.dart' as caching;
import 'package:retrovibed/design.kit/theme.defaults.dart' as theming;
import 'package:retrovibed/env.dart' as env;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/windowx.dart';
import 'package:window_manager/window_manager.dart' show TitleBarStyle;

File _defaultlib() {
  if (Platform.isMacOS) {
    // Flutter's native-assets build wraps the raw libretrovibed.dylib built
    // by hook/build.dart into a versioned .framework and embeds it in
    // Contents/Frameworks, stripping the "lib" prefix and ".dylib"
    // extension along the way (see frameworkUri() in flutter_tools'
    // native_assets_host.dart) - so the file that ends up on disk is
    // retrovibed.framework/retrovibed, not libretrovibed.dylib.
    final execDir = File(Platform.resolvedExecutable).parent.path;
    return File("$execDir/../Frameworks/retrovibed.framework/retrovibed");
  }

  return File("libretrovibed.so");
}

String _path({String name = "libretrovibed.so"}) {
  final files = () {
    final ldlibs = env.string([
      'LD_LIBRARY_PATH',
      'APPDIR_LIBRARY_PATH',
      'RETROVIBED_SHARED_NATIVE_LIBS_DIRECTORY',
    ], fallback: '../.eg.cache/dev.native.libs');
    return ldlibs.split(":").map((path) => File("${path}/${name}"));
  }();

  final found = files.firstWhere((v) {
    try {
      return v.existsSync();
    } catch (_) {
      return false;
    }
  }, orElse: _defaultlib);
  return found.path;
}

DynamicLibrary _loadLibrary() {
  if (Platform.isIOS) {
    return DynamicLibrary.process();
  }

  return DynamicLibrary.open(_path());
}

final bridge = lib.DaemonBridge(_loadLibrary());

String build_version() {
  return _convertstring(bridge.build_version());
}

String deeppool_oauth2_bearer() {
  return _convertstring(bridge.deeppool_oauth2_bearer());
}

String bearer_token() {
  return _convertstring(bridge.authn_bearer());
}

// return a bearer authorized to connect to this process's own local
// remote control listen socket. never valid outside this process.
meta.Daemon local_device() {
  final encoded = _convertstring(bridge.local_device());
  return meta.Daemon()..mergeFromProto3Json(jsonDecode(encoded));
}

// return a bearer authorized to connect to this process's own local
// remote control listen socket. never valid outside this process.
String remote_control_listen_token() {
  return _convertstring(bridge.remote_control_listen_token());
}

String bearer_token_host(String hostname) {
  return _convertstring(
    bridge.authn_bearer_host(hostname.toNativeUtf8().cast<Char>()),
  );
}

String public_key() {
  return _convertstring(bridge.public_key());
}

String username() {
  return _convertstring(bridge.username());
}

String xdg_dir_config() {
  return _convertstring(bridge.xdg_dir_config());
}

String xdg_dir_cache() {
  return _convertstring(bridge.xdg_dir_cache());
}

String xdg_dir_data() {
  return _convertstring(bridge.xdg_dir_data());
}

String xdg_dir_download() {
  return _convertstring(bridge.xdg_dir_download());
}

String xdg_relroot() {
  return _convertstring(bridge.xdg_relroot());
}

// returns an empty string on success, non empty contains the error.
String seed(String passphrase) {
  return _convertstring(bridge.seed(passphrase.toNativeUtf8().cast<Char>()));
}

// returns an empty string on success, non empty contains the error.
String unseed() {
  return _convertstring(bridge.unseed());
}

// bootstraps a local-only guest identity derived from this device's hostname
// and machine id. returns true on success.
bool guest() {
  return bridge.guest() == 0;
}

List<String> ips() {
  final String? ipaddrs = _convertstring(bridge.ips());
  final List<dynamic> res = jsonDecode(ipaddrs ?? "") ?? [];
  return res.whereType<String>().toList();
}

void setenv(String key, String value) {
  bridge.gsetenv(
    key.toNativeUtf8().cast<Char>(),
    value.toNativeUtf8().cast<Char>(),
  );
}

void logging() {
  bridge.logging();
}

// Returns true if the certificate is valid/enrolled, false otherwise.
bool validatecert(String hostname, Uint8List derBytes) {
  final hostnamePtr = hostname.toNativeUtf8().cast<Char>();
  final certPtr = ffi.calloc<Uint8>(derBytes.length);
  try {
    certPtr.asTypedList(derBytes.length).setAll(0, derBytes);
    final result = bridge.validatecert(
      hostnamePtr,
      certPtr.cast<UnsignedChar>(),
      derBytes.length,
    );
    return result == 0;
  } finally {
    ffi.calloc.free(hostnamePtr);
    ffi.calloc.free(certPtr);
  }
}

// Clears all locally trusted (TOFU) certificates, forcing re-enrollment on next connect.
// Returns true on success.
bool resetcerts() {
  return bridge.resetcerts() == 0;
}

void daemon({bool smoke = false}) {
  String args = jsonEncode([
    "daemon",
    "--no-mdns-advertise",
    "--mdns-discovery",
    "--remote-control",
  ]);
  bridge.egdaemon(args.toNativeUtf8().cast<Char>(), smoke ? 1 : 0);
}

String _convertstring(Pointer<Char> charPointer) {
  try {
    return charPointer.cast<ffi.Utf8>().toDartString();
  } finally {
    ffi.calloc.free(charPointer);
  }
}

bool metered() {
  return bridge.netmon_metered() != 0;
}

bool set_metered(bool b) {
  return bridge.netmon_set_metered(b ? 1 : 0) != 0;
}

void checkpointdb() {
  try {
    bridge.checkpointdb();
  } catch (e) {
    print("failed to checkpoint database ${e}");
  }
}

void fault(int errcode) {
  exit(errcode);
}

Future<void> run(void Function() fn) async {
  print("build version ${build_version()}");
  // env.printSystemEnv();
  print("cp 0");
  WidgetsFlutterBinding.ensureInitialized();
  print("cp 1");
  HttpOverrides.global = meta.DaemonHttpOverrides();
  print("cp 2");
  logging();
  print("cp 3");
  caching.setglobal(await caching.DirsWellKnown.xdg(await env.xdg()));
  print("cp 4");
  // checkpointing the database on initialization prevents
  // a significant number of issues due to hard shutdowns and state corruption
  // issues.
  checkpointdb();
  print("cp 5");
  if (theming.Defaults.defaults.desktop) {
    await windowx.ensureInitialized();
    // waitUntilReadyToShow force-corrects isFullScreen/isMaximized/isMinimized
    // to a known state before anything else reads them (see window_manager's
    // own implementation) - without this, early reads of that state (e.g. in
    // Full/Hamburger initState) can be unreliable, especially on Linux where
    // isFullScreen reads live GDK window state that may not be settled yet.
    await windowx.waitUntilReadyToShow();
    await Future.wait([
      windowx.setTitleBarStyle(TitleBarStyle.hidden),
      windowx.maximize(),
    ]);
  }
  print("cp 6");
  MediaKit.ensureInitialized();
  print("cp 7");

  fn();
}
