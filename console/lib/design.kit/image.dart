import 'package:flutter/material.dart' as m;

class Image extends m.StatelessWidget {
  final String current;
  final double? size;
  const Image(
    this.current, {
    super.key,
    this.size,
  });

  // attaches a silent error listener before Image.network resolves so an
  // expected load failure (e.g. missing cover art) doesn't get dumped to the
  // console even though errorBuilder already shows a fallback icon. no-op for
  // an empty url.
  static void precache(m.BuildContext context, String current, {Map<String, String>? headers}) {
    if (current == "") return;
    m.precacheImage(m.NetworkImage(current, headers: headers), context, onError: (exception, stackTrace) {});
  }

  @override
  m.Widget build(m.BuildContext context) {
    if (current == "") return m.Icon(m.Icons.image_outlined, size: size);

    precache(context, current);

    return m.Image.network(
      current,
      height: size,
      errorBuilder: (context, error, stackTrace) => m.Icon(m.Icons.image_outlined, size: size),
    );
  }
}
