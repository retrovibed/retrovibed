import 'package:flutter/material.dart' as m;

class Image extends m.StatelessWidget {
  final String current;
  final double? size;
  const Image(
    this.current, {
    super.key,
    this.size,
  });

  // builds an Image.network for current, attaching a silent error listener
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
    m.Widget Function(m.BuildContext context)? error,
  }) {
    if (current == "") return null;

    final onError = error ?? (context) => m.Icon(m.Icons.image_outlined, size: height);
    m.precacheImage(m.NetworkImage(current, headers: headers), context, onError: (exception, stackTrace) {});

    return m.Image.network(
      current,
      headers: headers,
      width: width,
      height: height,
      fit: fit,
      errorBuilder: (context, exception, stackTrace) => onError(context),
    );
  }

  @override
  m.Widget build(m.BuildContext context) {
    return precache(context, current, height: size) ?? m.Icon(m.Icons.image_outlined, size: size);
  }
}
