import 'dart:io';
import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:path/path.dart' as path;
import 'package:path_provider/path_provider.dart';

abstract class Cache {
  T put<T>(String k, T v, {Duration? ttl});
  T? get<T>(String k);
  T delete<T>(String k);
}

abstract class Codec {
  List<int> encode<T>(T v);
  T decode<T>(List<int> v);
}

class jsoncodec implements Codec {
  const jsoncodec();

  List<int> encode<T>(T v) {
    return utf8.encode(jsonEncode(v));
  }

  T decode<T>(List<int> v) {
    return jsonDecode(utf8.decode(v)) as T;
  }
}

class record<T> {
  final T data;
  final DateTime ttl;
  const record(this.data, this.ttl);
}

class disk {
  final Duration ttl;
  final Directory dir;
  final Codec codec;

  const disk(
    this.dir, {
    this.ttl = const Duration(minutes: 5),
    this.codec = const jsoncodec(),
  });

  static Future<disk> defaulted({Duration? ttl, Codec? codec}) {
    return getApplicationCacheDirectory().then(
      (v) => disk(
        v,
        ttl: ttl ?? const Duration(minutes: 5),
        codec: codec ?? const jsoncodec(),
      ),
    );
  }

  T put<T>(String k, T v, {Duration? ttl}) {
    final dst = path.join(dir.path, md5.convert(utf8.encode(k)).toString());
    File(dst).writeAsBytesSync(
      codec.encode(record(v, DateTime.now().add(ttl ?? this.ttl))),
    );
    return v;
  }

  T? get<T>(String k) {
    final dst = path.join(dir.path, md5.convert(utf8.encode(k)).toString());
    final decoded = codec.decode<record<T>>(File(dst).readAsBytesSync());
    if (decoded.ttl.isAfter(DateTime.now())) {
      return decoded.data;
    }

    File(dst).deleteSync();
    return null;
  }

  T? delete<T>(String k) {
    final dst = File(
      path.join(dir.path, md5.convert(utf8.encode(k)).toString()),
    );

    if (!dst.existsSync()) {
      return null;
    }

    final decoded = codec.decode<record<T>>(dst.readAsBytesSync());
    dst.deleteSync();

    return decoded.ttl.isAfter(DateTime.now()) ? decoded.data : null;
  }
}
