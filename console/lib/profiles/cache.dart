// import 'dart:convert';
import 'dart:typed_data';
import 'package:synchronized/synchronized.dart';
import 'package:lru/lru.dart';
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/design.kit/bytesx.dart';

final LruTypedDataCache<String, Uint8List> cache = LruTypedDataCache<String, Uint8List>(
  capacity: 256,
  capacityInBytes: bytesx.MiB,
);

final Lock _lock = Lock();
Future<meta.ProfileLookupResponse> cached(
  String id,
  Future<meta.ProfileLookupResponse> Function() fetch,
) {
  return _lock.synchronized(() {
    final c = cache[id];
    return c == null
        ? fetch().then((v) {
          cache[id] = v.profile.writeToBuffer();
          return v;
        })
        : Future.value(meta.ProfileLookupResponse(profile: meta.Profile.fromBuffer(c)));
  });
}
