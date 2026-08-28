import 'dart:io';
import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';
import 'package:retrovibed/retrovibed.dart' as retro;
import 'package:retrovibed/httpx.dart' as httpx;

bool boolean(String key, {bool fallback = false}) {
  try {
    return bool.parse(Platform.environment[key] ?? "");
  } catch (_) {
    return fallback;
  }
}

String lookup(Iterable<String> key, Map<String, String?> env, String fallback) {
  for (final k in key) {
    if (env.containsKey(k)) {
      return env[k] ?? '';
    }
  }
  return fallback;
}

String string(Iterable<String> key, {String fallback = ''}) => lookup(key, Platform.environment, fallback);

void printSystemEnv() {
  Platform.environment.forEach((key, value) {
    print('$key: $value');
  });
}

Future<({String configDir, String dataDir, String cacheDir, String downloadDir})> xdg() async {
  if (Platform.isLinux || Platform.isMacOS) {
    final relroot = retro.xdg_relroot();
    return (
      configDir: p.join(retro.xdg_dir_config(), relroot),
      dataDir: p.join(retro.xdg_dir_data(), relroot),
      cacheDir: p.join(retro.xdg_dir_cache(), relroot),
      downloadDir: retro.xdg_dir_download(),
    );
  }

  final dataDir = await getApplicationDocumentsDirectory();
  final cacheDir = await getApplicationCacheDirectory();
  final downloadDir = Platform.isIOS ? dataDir : await getDownloadsDirectory();
  final configDir = await getApplicationSupportDirectory();

  print("config ${configDir}");
  print("data ${dataDir}");
  print("cache ${cacheDir}");
  print("download ${downloadDir}");

  retro.setenv("XDG_CONFIG_HOME", configDir.path);
  retro.setenv("XDG_DATA_HOME", dataDir.path);
  retro.setenv("XDG_CACHE_HOME", cacheDir.path);
  retro.setenv("XDG_DOWNLOAD_DIR", downloadDir?.path ?? "");
  // On mobile the native bridge is a separate OS process env from whatever
  // launched the app (e.g. Android's Zygote/ActivityManager), so it never
  // sees RETROVIBED_META_ENDPOINT set by dev tooling the way Dart's own
  // Platform.environment lookup does. Without this, Go's Deeppool() falls
  // back to its own (build-tag-dependent) default instead of the endpoint
  // Dart is actually using, causing Go-side calls like oauth2_bearer() to
  // target the wrong host while Dart's HTTP calls hit the right one.
  retro.setenv("RETROVIBED_META_ENDPOINT", httpx.metaendpoint());

  return (
    configDir: configDir.path,
    dataDir: dataDir.path,
    cacheDir: cacheDir.path,
    downloadDir: downloadDir?.path ?? "",
  );
}

class vars {
  static const AutoIdentifyMedia = "RETROVIBED_MEDIA_AUTO_IDENTIFY";
  static const WindowManagerNativeFullScreen = "RETROVIBED_WINDOW_MANAGER_NATIVE_FULL_SCREEN";
}
