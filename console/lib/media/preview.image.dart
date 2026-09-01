import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import './media.pb.dart';
import './preview.unsupported.dart';

// the daemon already publishes a url for image mimetypes, so this costs no new transport.
// InteractiveViewer keeps the whole image on screen by default and lets the reader zoom
// into detail a fitted image loses.
class PreviewImage extends StatelessWidget {
  final Media current;

  const PreviewImage({super.key, required this.current});

  @override
  Widget build(BuildContext context) {
    final image = ds.Image.precache(
      context,
      current.image,
      headers: httpx.localheaders(current.image),
      fit: BoxFit.contain,
    );

    if (image == null) return PreviewUnsupported(current: current);

    return InteractiveViewer(
      minScale: 1.0,
      maxScale: 8.0,
      child: SizedBox.expand(child: image),
    );
  }
}
