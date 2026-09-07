import 'dart:io';
import 'package:flutter/material.dart' as m;
import 'package:retrovibed/httpx.dart' as httpx;
import 'image.cache.dart' as imagecache;

class Image extends m.StatelessWidget {
  final String current;
  final double? size;
  const Image(
    this.current, {
    super.key,
    this.size,
  });

  // builds a disk-cached network image, attaching a silent error listener
  // beforehand so an expected load failure (e.g. missing cover art) doesn't
  // get dumped to the console even though errorBuilder already shows a
  // fallback widget. returns null when current is unset -- callers should
  // fall back with `precache(...) ?? fallback`.
  static m.Widget? precache(
    m.BuildContext context,
    String current, {
    Map<String, String>? headers,
    double? width,
    double? height,
    m.BoxFit? fit,
    m.Widget? missing,
  }) {
    if (current == "") return null;

    final pending = imagecache.cached(current, () {
      final options = (headers ?? const {}).entries.map((e) => httpx.Request.header(e.key, e.value)).toList();
      return httpx.get(Uri.parse(current), options: options).then((r) => r.bodyBytes);
    });

    return _CachedImage(
      pending: pending,
      width: width,
      height: height,
      fit: fit,
      missing: missing ?? m.Icon(m.Icons.image_outlined, size: height),
    );
  }

  @override
  m.Widget build(m.BuildContext context) {
    final missing = m.Icon(m.Icons.image_outlined, size: size);
    return precache(context, current, height: size, missing: missing) ?? missing;
  }
}

// _CachedImage renders pending's disk-cached file once it resolves, and
// missing in the meantime (or if it fails). FutureBuilder already handles
// the async lifecycle off the build path, so no bespoke State is needed.
class _CachedImage extends m.StatelessWidget {
  final Future<File> pending;
  final double? width;
  final double? height;
  final m.BoxFit? fit;
  final m.Widget missing;

  const _CachedImage({
    required this.pending,
    this.width,
    this.height,
    this.fit,
    required this.missing,
  });

  @override
  m.Widget build(m.BuildContext context) {
    return m.FutureBuilder<File>(
      future: pending,
      builder: (context, snapshot) {
        if (snapshot.hasError) {
          print("failed to load image ${snapshot.error}");
          return missing;
        }

        final file = snapshot.data;
        if (file == null) return missing;

        final provider = m.FileImage(file);
        m.precacheImage(
          provider,
          context,
          onError: (exception, stackTrace) {
            print("unable to cache image ${exception}");
          },
        );

        return m.Image(
          image: provider,
          width: width,
          height: height,
          fit: fit,
          errorBuilder: (context, exception, stackTrace) => missing,
        );
      },
    );
  }
}
